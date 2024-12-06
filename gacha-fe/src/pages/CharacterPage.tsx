import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router";
import { toast } from "react-toastify";
import InputCharacterModal from "../components/CharacterModal";
import { DeletIcon, EditIcon } from "../components/Icon";
import { LogoutButton } from "../components/LogoutButton";
import { Character } from "../types/characterType";
import { Rarity } from "../types/rarityType";
import { handleRequest } from "../utils/api";
import NotFound from "./NotFoundPage";

const CharacterPage: React.FC = () => {
  let { gachaSystemId } = useParams<{ gachaSystemId: string | undefined }>();
  let { characterId } = useParams<{ characterId: string | undefined }>();

  if (!gachaSystemId || !characterId) {
    return <p className="text-red-500 text-center">Error: Gacha System and character ID are required.</p>;
  }

  const emptyCharacter = {
    id: parseInt(characterId),
    name: "",
    imageUrl: "",
    rarityId: -1,
  };

  const [character, setCharacter] = useState<Character>(emptyCharacter);
  const [rarities, setRarities] = useState<Rarity[]>([]);

  const [refreshToggle, setRefreshToggle] = useState(false);
  const [isModalOpen, setIsModalOpen] = useState(false);

  const [rarityMap, setRarityMap] = useState<Record<number, string>>({});

  const [isNotFound, setIsNotFound] = useState(false);

  const navigate = useNavigate();

  const fetchCharacterData = async () => {
    const endpoint = `${import.meta.env.VITE_GACHA_MASTER_URL}/id/${gachaSystemId}/character/${characterId}`;
    try {
      const response = await handleRequest<Character>(endpoint, "GET");
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

  const fetchRaritiesData = async () => {
    const endpoint = `${import.meta.env.VITE_GACHA_MASTER_URL}/id/${gachaSystemId}/rarity/all`;
    try {
      const response = await handleRequest<Rarity[]>(endpoint, "GET");
      if (response.code === 200) {
        const newMap = response.data.reduce<Record<number, string>>((acc, rarity) => {
          acc[rarity.id] = rarity.name;
          return acc;
        }, {});
        setRarityMap(newMap);
        setRarities(response.data);
      } else {
        toast.error(`Error: ${response.status}`);
      }
    } catch (error) {
      toast.error("An error occurred. Please try again." + error);
    }
  };

  const fetchData = async () => {
    fetchCharacterData();
    fetchRaritiesData();
  };

  useEffect(() => {
    console.log("fetching data");
    fetchData();
    console.log(character);
  }, [refreshToggle]);

  const handleDelete = async (id: number) => {
    const confirmation = window.confirm("Are you sure you want to delete this Gacha system?");
    if (confirmation) {
      const endpoint = `${import.meta.env.VITE_GACHA_MASTER_URL}/id/${id}/character/${characterId}`;
      try {
        const response = await handleRequest<{ message: string }>(endpoint, "DELETE");
        if (response.code === 200) {
          toast.success("Delete gacha system success");
          navigate(`/gacha-system/${gachaSystemId}`);
        } else {
          toast.error(`Delete character failed: ${response.data.message}`);
        }
      } catch (error) {
        toast.error("An error occurred. Please try again.");
      }
    }
  };

  if (isNotFound) {
    return <NotFound />;
  }

  return (
    <div className="relative min-h-screen w-screen flex items-center justify-center bg-gradient-to-t from-slate-800 to-cyan-900">
      <LogoutButton />

      {character && (
        <div className="h-full">
        <img
          src={`${character.imageUrl}?v=${new Date().getTime()}`}
          alt={character.name}
          className="w-full h-screen"
        />
      </div>
      )}
      
      {/* Konten modal dan tombol */}
      <div className="fixed bottom-8 flex justify-center w-full">
        <button
          onClick={() => handleDelete(parseInt(gachaSystemId))}
          className="bg-red-500 hover:bg-red-500 rounded-r-none bg-opacity-70 backdrop-blur-md  p-2 flex items-center justify-center"
          >
          <DeletIcon />
        </button> 

        <div className="flex flex-col space-y-4 bg-white bg-opacity-20 backdrop-blur-md px-10 pt-10 pb-6 shadow-lg w-2/4">
          <h1 className="text-center text-3xl md:text-5xl font-sans font-bold leading-loose text-white">
            {character.name}
          </h1>

          <p className="text-center text-lg text-white font-semibold">
            {rarityMap[character.rarityId]}
          </p>
        </div>

        <button
          onClick={() => setIsModalOpen(true)}
          className="bg-blue-500 hover:bg-blue-500 rounded-l-none bg-opacity-70 backdrop-blur-md p-2 flex items-center justify-center"
          >
          <EditIcon />
        </button>


      </div>

      

      <InputCharacterModal
        isOpen={isModalOpen}
        rarities={rarities}
        character={character}
        gachaSystemId={parseInt(gachaSystemId)}
        onClose={() => setIsModalOpen(false)}
        onRefresh={() => setRefreshToggle(!refreshToggle)}
      />
    </div>
  );
};

export default CharacterPage;

