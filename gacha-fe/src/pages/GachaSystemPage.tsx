import { useEffect, useState } from "react";
import { Link, useNavigate, useParams } from "react-router";
import { toast } from "react-toastify";
import InputCharacterModal from "../components/CharacterModal";
import { DeletIcon, EditIcon, RedirectIcon } from "../components/Icon";
import { LogoutButton } from "../components/LogoutButton";
import InputRarityModal from "../components/RarityModal";
import { Character } from "../types/characterType";
import { GachaSytemDetail } from "../types/gachaSystemType";
import { Rarity } from "../types/rarityType";
import { handleRequest } from "../utils/api";
import NotFound from "./NotFoundPage";

const responseExample = {
  code: 200,
  status: "OK",
  data: {
    name: "string",
    imageUrl: "string",
    rarity: "string",
  },
};

const emptyCharacter = {
  id: -1,
  name: "",
  imageUrl: "",
  rarityId: -1,
};

const emptyRarity = { id: -1, name: "", chance: 0 }


const GachaSystemPage: React.FC = () => {
  let { gachaSystemId } = useParams<{ gachaSystemId: string | undefined }>();


  // Menangani kasus ketika gachaSystemId tidak tersedia
  if (!gachaSystemId) {
    return <p className="text-red-500 text-center">Error: Gacha System ID is required.</p>;
  }

  // Mengambil data Gacha System berdasarkan ID

  const [characters, setCharacters] = useState<Character[]>([]);
  const [gachaSystemName, setGachaSystemName] = useState<string>("");
  const [endpoint, setEndpoint] = useState<string>("");
  const [rarities, setRarities] = useState<Rarity[]>([]);
  const [rarityMap, setRarityMap] = useState<Record<number, string>>({});

  const [refreshToggle, setRefreshToggle] = useState(false);
  const [inputRarity, setInputRarity] = useState<Rarity | null>(null);
  const [inputCharacter, setInputCharacter] = useState<Character>(emptyCharacter);

  const [isRarityModalOpen, setIsRarityModalOpen] = useState(false);
  const [isCharacterModalOpen, setIsCharacterModalOpen] = useState(false);

  const [isNotFound, setIsNotFound] = useState(false);

  const navigate = useNavigate();

  const fetchData = async () => {
    const endpoint = `${import.meta.env.VITE_GACHA_MASTER_URL}/id/${gachaSystemId}`;
    try {
      const response = await handleRequest<GachaSytemDetail>(endpoint, "GET");
    
      if (response.code === 200) {

        setGachaSystemName(response.data.name);
        setRarities(response.data.rarities);
        setCharacters(response.data.characters);
        setEndpoint(response.data.endpoint);

        const newMap = response.data.rarities.reduce<Record<number, string>>((acc, rarity) => {
          acc[rarity.id] = rarity.name;
          return acc;
        }, {});
        setRarityMap(newMap);

        console.log(response)
  
      } else if (response.code === 404) {
        setIsNotFound(true);
        return
      } else if (response.code === 401) {
        navigate("/");
        toast.error("You are not authorized to view the page.");
        return
      }
      else {
        toast.error(`Error : ${response.data.message}`);
      }
    } catch (error) {
      if (error instanceof Error) {
        if (error.message === "401") {
          navigate("/");
          toast.error("You are not authorized to view the page.");
        }
      } else {
        toast.error("An error occurred. Please try again." + error);
      }
    }
  }

  useEffect(() => {
    fetchData();
  }, [refreshToggle]); 

  // Mendapatkan data Gacha System dan detailnya

  const handleRarityDelete = async (id: number) => {
    const confirmation = window.confirm("Are you sure you want to delete this Gacha rarity? All characters with this rarity will also be deleted.");
    if (confirmation) {
      const deleteRarityEndpoint = `${import.meta.env.VITE_GACHA_MASTER_URL}/id/${gachaSystemId}/rarity/${id}`;
      try {
        const response = await handleRequest<{message:string}>(deleteRarityEndpoint, "DELETE");
      
        if (response.code === 200) {
          toast.success("Delete rarity success");
          setRefreshToggle(!refreshToggle)

        } else {
          toast.error(`Authentication failed: ${response.data.message}`);
        }
      } catch (error) {
        toast.error("An error occurred. Please try again.");
      }
    }
  };

  const openRarityModal = (id:number) => {
    setIsCharacterModalOpen(false);

    if (id === -1) {
      setInputRarity(emptyRarity);
      setIsRarityModalOpen(true);
      return;
    }

    const rarity = rarities.find((r) => r.id === id);
    if (rarity) {
      setInputRarity(rarity);
      setIsRarityModalOpen(true);
    } 
  }


  const openCharacterModal = (id:number) => {
    setIsRarityModalOpen(false);

    console.log(rarities)

    if (!rarities || rarities.length === 0) {
      toast.error("You must add at least one rarity before creating a character.");
      return;
    }

    if (id === -1) {
      setInputCharacter(emptyCharacter);
      setIsCharacterModalOpen(true);
      return;
    }

    const character = characters.find((r) => r.id === id);
    if (character) {
      setInputCharacter(character);
      setIsCharacterModalOpen(true);
    } 
  }


  const handleRarityModal = (id: number) => {
    if (id === -1) {
      setInputRarity({ id: -1, name: "", chance: 0 });
      setIsRarityModalOpen(true);
      return;
    }

    const rarity = rarities.find((r) => r.id === id);
    if (rarity) {
      setInputRarity({ id: rarity.id, name: rarity.name, chance: rarity.chance });
      setIsRarityModalOpen(true);
    } 
  };

  const handleRarityUpdate = async () => {
    if (inputRarity) {
      try {
        const editEndpoint = `${import.meta.env.VITE_GACHA_MASTER_URL}/rarity/update`;
        const response = await handleRequest<{message:string}>(editEndpoint, "PUT", {
          id: inputRarity.id,
          gachaSystemId: parseInt(gachaSystemId),
          name: inputRarity.name,
          chance: parseFloat(inputRarity.chance.toString()),
        });


        if (response.code === 200) {
          toast.success("Rarity updated successfully!");
          setIsRarityModalOpen(false);
          setRefreshToggle(!refreshToggle);
        } else {
          toast.error(`Update failed: ${response.data.message}`);
        }
      } catch (error) {
        toast.error("An error occurred. Please try again.");
      }
    }
  };

  const handleRarityCreate = async () => {
    if (inputRarity) {
      if (inputRarity.chance <= 0 || inputRarity.chance > 100) {
        toast.error("Chance must be greater than 0 and less than or equal to 100.");
        return;
            }

      try {
        const editEndpoint = `${import.meta.env.VITE_GACHA_MASTER_URL}/rarity/create`;
        const response = await handleRequest<{message:string}>(editEndpoint, "POST", {
          gachaSystemId: parseInt(gachaSystemId),
          name: inputRarity.name,
          chance: parseFloat(inputRarity.chance.toString()),
        });

        console.log(response)
        if (response.code === 200) {
          toast.success("Rarity created successfully!");
          setIsRarityModalOpen(false);
          setRefreshToggle(!refreshToggle);
        } else {
          console.log(response)
          toast.error(`Create failed: ${response.data.message}`);
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
    <div className="min-h-screen w-screen bg-gradient-to-br from-blue-600 to-green-400 py-8">
      <LogoutButton />
      <h1 className="text-3xl font-bold text-white text-center mb-6">{gachaSystemName}</h1>

      <div className="flex flex-col md:flex-row items-start justify-start">
        <div className="w-full md:w-1/2 p-6">
          <div className="bg-white bg-opacity-20 backdrop-blur-md rounded-lg shadow-lg p-6 flex flex-col">

            <Link to={`/gacha/${localStorage.getItem("user_id")}/${gachaSystemName}`} target="_blank" rel="noopener noreferrer"  className="w-full flex space-x-1 justify-center bg-gradient-to-r from-blue-600 to-green-400 text-white p-3 rounded hover:text-white hover:scale-105 mb-6">
              <p>Go to Gacha Page</p>  
              <RedirectIcon/>
            </Link>

            {/* Rarities */}
            <div className="mb-6">
              <div className="flex justify-between  items-center mb-2">
                <h2 className="text-xl font-semibold text-white">Rarities</h2> 
                <button
                  onClick={() => handleRarityModal(-1)}
                  className="bg-transparent hover:bg-white text-white hover:bg-opacity-50 outline outline-white outline-1 rounded px-2 py-1 flex items-center justify-center">
                  Create
                  </button>
              </div>
              <ul className="space-y-4">
              { rarities && rarities.length > 0 ? (
                rarities.map((rarity) => (
                  <li
                    key={rarity.id}
                    className="w-full flex justify-between items-center space-x-2"
                  >
                    <div className="w-full flex justify-between items-center bg-white bg-opacity-80 text-gray-800 rounded-lg shadow px-4 py-2">
                      <h3 className="text-lg font-medium">{rarity.name}</h3>
                      <p className="text-sm">Chance: {rarity.chance.toFixed(2)}%</p>
                    </div>

                    <div className="flex space-x-2">
                      <button
                        onClick={() => openRarityModal(rarity.id)}
                        className="bg-blue-500 hover:bg-blue-700 rounded p-2 flex items-center justify-center"
                      >
                        <EditIcon />
                      </button>

                      <button
                        onClick={() => handleRarityDelete(rarity.id)}
                        className="bg-red-500 hover:bg-red-700 rounded p-2 flex items-center justify-center"
                      >
                       <DeletIcon />
                      </button>
                    </div>
                  </li>
                ))
              ) : (
                <p className="text-white text-sm">No rarities available.</p>
              )}
            </ul>

            </div>

            {/* Characters */}
            <div className="mb-6 mt-3">
              <div className="flex justify-between  items-center mb-2">
                <h2 className="text-xl font-semibold text-white">Characters</h2> 
                <button
                  onClick={() => openCharacterModal(-1)}
                  className="bg-transparent hover:bg-white text-white hover:bg-opacity-50 outline outline-white outline-1 rounded px-2 py-1 flex items-center justify-center">
                  Add
                  </button>
              </div>
              <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
              {characters && characters.length > 0 ? (
                characters.map((character) => (
                  <Link
                    to={`/gacha-system/${gachaSystemId}/character/${character.id}`}
                    key={character.id}
                    className="bg-white bg-opacity-80 p-4 rounded-lg shadow text-center flex flex-col justify-between"
                  >
                    <h3 className="text-lg font-medium">{character.name}</h3>
                    <p className="text-sm text-gray-600">Rarity: {rarityMap[character.rarityId]}</p>
                  </Link>
                ))
              ) : (
                <p className="text-white text-sm">No characters available.</p>
              )}
            </div>

            </div>
          </div>
        </div>

        <div className="w-full md:w-1/2 p-6 flex">
          <div className="bg-white bg-opacity-20 backdrop-blur-md rounded-lg shadow-lg p-4 w-full">
              <h2 className="text-2xl font-semibold text-white mb-4">REST API</h2>
              <div className="flex items-center space-x-2 bg-orange-300 bg-opacity-40 p-2 rounded">
              <p className="bg-orange-500 bg-opacity-75 text-white font-semibold rounded px-2 py-1">GET</p>
                <p className="font-semibold text-white text break-words max-w-full py-1 text-ellipsis overflow-hidden">
                  {endpoint}
                </p>
              </div>

              <div className="mt-6">
              <h2 className="font-semibold text-white mb-2">Response</h2>
              <pre className="bg-gray-800 text-white p-4 rounded-md">

                {JSON.stringify(responseExample, null, 2)}
              </pre>
            </div>
          </div>
          
        </div>
      </div>

      {/* Modal */}
      <InputRarityModal
        isOpen={isRarityModalOpen}
        rarity={inputRarity}
        setRarity={setInputRarity}
        onClose={() => setIsRarityModalOpen(false)}
        onUpdate={handleRarityUpdate}
        onCreate={handleRarityCreate}
      />    

      <InputCharacterModal
        isOpen={isCharacterModalOpen}
        rarities={rarities}
        character={inputCharacter}
        gachaSystemId={parseInt(gachaSystemId)}
        onClose={() => setIsCharacterModalOpen(false)}
        onRefresh={() => setRefreshToggle(!refreshToggle)}
      />     
    </div>
  );
};

export default GachaSystemPage;
