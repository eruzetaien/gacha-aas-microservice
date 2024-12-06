import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router";
import { toast } from "react-toastify";
import { CharacterGacha } from "../types/characterType";
import { GachaSystemEndpoint } from "../types/gachaSystemType";
import { handleRequest } from "../utils/api";
import { NoCharacterPage } from "./NoCharacterPage";
import NotFound from "./NotFoundPage";

const ANIMATION_START = 0;
const ANIMATION_FINISH = 1;
const ANIMATION_NONE = 2;

const GachaPullPage: React.FC = () => {
  let { gachaSystemName } = useParams<{ gachaSystemName: string | undefined }>();

  if (!gachaSystemName) {
    return <p className="text-red-500 text-center">Error: Gacha System name is required.</p>;
  }

  const [gachaEndpoint, setGachaEndpoint] = useState("");
  const [character, setCharacter] = useState<CharacterGacha|null>(null);

  const [gachaToggle, setGachaToggle] = useState(false);

  const [isNotFound, setIsNotFound] = useState(false);

  const [totalCharacters, setTotalCharacters] = useState(1);
  const [animationStat, setAnimationStat] = useState(ANIMATION_NONE);

  const navigate = useNavigate();

  const startAnimation = () => {
    setAnimationStat(ANIMATION_START);
    setTimeout(() => { finishAnimation()}, 4000);
  }

  const finishAnimation = () => {
    setAnimationStat(ANIMATION_FINISH);
    setTimeout(() => {
      setAnimationStat(ANIMATION_NONE);
      setGachaToggle(false);
    }, 2000);
  }

  const fetchGachaEndpoint = async () => {
    const userId = localStorage.getItem("user_id");
    if (!gachaSystemName.trim()) {
      return
    }

    const endpoint = `${import.meta.env.VITE_GACHA_MASTER_URL}/userId/${userId}?name=${encodeURIComponent(gachaSystemName)}`;
    try {
      const response = await handleRequest<GachaSystemEndpoint>(endpoint, "GET");
      if (response.code === 200) {
        setGachaEndpoint(response.data.endpoint);
        setTotalCharacters(response.data.totalCharacters);
      } else {
        if (response.code === 404) {
          setIsNotFound(true);
          return
        }
        toast.error(`Error: ${response.data.message}`);
      }
    } catch (error) {
      if (error instanceof Error) {
        if (error.message === "401") {
          navigate("/");
          toast.error("You are not authorized to view this page.");
        }
      } else {
        toast.error("An error occurred. Please try again :" + error);
      }
    }
  };


  const pullCharacter = async () => {
    console.log(gachaEndpoint);
    if (!gachaEndpoint) {
      return
    }
    startAnimation();

    const endpoint = gachaEndpoint;
    try {
      const response = await handleRequest<CharacterGacha>(endpoint, "GET");
      if (response.code === 200) {
        setCharacter(response.data);
      } else {
        if (response.code === 404) {
          setIsNotFound(true);
          return
        }
        toast.error(`Error: ${response.data.message}`);
      }
    } catch (error) {
      if (error instanceof Error) {
        if (error.message === "401") {
          navigate("/");
          toast.error("You are not authorized to view this page.");
        }
      } else {
        toast.error("An error occurred. Please try again :" + error);
      }
    }
  };

  useEffect(() => {
    if (!gachaEndpoint) {
      fetchGachaEndpoint();
    }
    if (gachaToggle === true) {
      pullCharacter();
    }
  }, [gachaToggle]);

  if (isNotFound) {
    return <NotFound />;
  }

  if (totalCharacters <= 0) {
    return <NoCharacterPage/>;
  }


  return (
    <div className="relative min-h-screen w-screen flex items-center justify-center bg-gradient-to-t from-slate-800 to-cyan-900">
      
      
      {character && (
        <div className="h-full relative">
          <img
            src={`${character.imageUrl}?v=${new Date().getMinutes()}`}
            alt={character.name}
            className={`h-screen transition-opacity duration-500 ${
              animationStat == ANIMATION_START ? "opacity-0" : "opacity-100"
            }`}
          />
        </div>
        )
      }


      { animationStat != ANIMATION_START && character && (
        <div className="fixed bottom-20 flex justify-center w-full scale-effect ">
          <div className="flex flex-col space-y-4 bg-white bg-opacity-20 text-white backdrop-blur-md px-10 pt-10 pb-6 shadow-lg w-2/4 ">
            <h1 className="text-center text-3xl md:text-5xl font-sans font-bold leading-loose">
              {character.name}
            </h1>

            <p className="text-center text-lg font-semibold">
              {character.rarity}
            </p>
          </div>
        </div>
        )
      }


      {animationStat == ANIMATION_NONE && (
          <div className='fixed bottom-5 flex justify-center w-1/12 fade-in-effect ' >
              <button className='w-full bg-gradient-to-tr from-violet-800 to-violet-500 text-white outline outline-2 outline-violet-300 hover:scale-110 transition-transform duration-300' onClick={() => setGachaToggle(true)}>
              Pull
              </button>
          </div>
        )}


      {animationStat == ANIMATION_START && (
        <>
              <div className="white-zoom-in"> </div>
              <div className="black-zoom-out"> </div>
        </>
      )}

      {animationStat == ANIMATION_FINISH && (
        <>
              <div className="white-fade-out"> </div>
        </>
      )}

    </div>
  );
};

export default GachaPullPage;



