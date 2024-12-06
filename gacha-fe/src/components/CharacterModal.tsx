import { useEffect, useState } from "react";
import { toast } from "react-toastify";
import { Character } from "../types/characterType";
import { Rarity } from "../types/rarityType";
import { handleRequest } from "../utils/api";

interface InputCharacterModalProps {
  isOpen: boolean;
  rarities: Rarity[];
  character: Character;
  gachaSystemId: number;
  onClose: () => void;
  onRefresh: () => void;
}

const InputCharacterModal: React.FC<InputCharacterModalProps> = ({
  isOpen,
  rarities,
  character,
  gachaSystemId,
  onClose,
  onRefresh,
}) => {
  const [id, setId] = useState(character.id);
  const [name, setName] = useState(character.name);
  const [rarityId, setRarityId] = useState(character.rarityId);
  const [image, setImage] = useState<File | null>(null);

  useEffect(() => {
    setId(character.id);
    setName(character.name);
    setRarityId(character.rarityId);
    setImage(null); // Reset image since it's a new character
  }, [character]);

  if (!isOpen || !character) return null;
  
  const validateImage = (file: File): boolean => {
    return file.type === "image/png" && file.size <= 2 * 1024 * 1024; // 2MB
  };

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      const file = e.target.files[0];
      if (validateImage(file)) {
        setImage(file);
      } else {
        setImage(null);
        alert("Please select a valid PNG file smaller than 2MB.");
      }
    }
  };

  const resetInput = () => {
    setId(-1);
    setName("");
    setRarityId(rarities[0].id);
    setImage(null);
  }


  const validateCharacter = (formData: FormData) => {
    if (!formData.get("name")) {
      toast.error("Name is required.");
      return false;
    }

    if (!formData.get("image")) { 
      toast.error("Image is required.");
      return false;
    }    

    return true;
  }

  const handleCreate = async (formData: FormData) => {
    if (formData  && validateCharacter(formData)) {
      try {
        const createEndpoint = `${import.meta.env.VITE_GACHA_MASTER_URL}/character/create`;
        const response = await handleRequest<Character>(createEndpoint, "POST", formData);

        if (response.code === 200) {
          toast.success("Character added successfully!");
          resetInput();
          onClose();
          onRefresh();
        } else {
          toast.error(`Update character failed: ${response.data.message}`);
        }
      } catch (error) {
        toast.error("An error occurred. Please try again.");
      }
    }
  };

  const handleUpdate = async (formData: FormData) => {
    if (formData) {
      try {
        const createEndpoint = `${import.meta.env.VITE_GACHA_MASTER_URL}/character/update`;
        const response = await handleRequest<Character>(createEndpoint, "PATCH", formData);

        if (response.code === 200) {
          toast.success("Character updated successfully!");
          onClose();
          onRefresh();
        } else {
          toast.error(`Update character failed: ${response.data.message}`);
        }
      } catch (error) {
        toast.error("An error occurred. Please try again.");
      }
    }
  };

  const handleSaveClick = () => {
    const formData = new FormData();
    if (id !== -1) formData.append("id", id.toString());
    formData.append("name", name);
    formData.append("gachaSystemId", gachaSystemId.toString());
    formData.append("rarityId", rarityId === -1? rarities[0].id.toString(): rarityId.toString());

    if (image) formData.append("image", image);

    if (id === -1) {
      handleCreate(formData);
    } else {
      handleUpdate(formData);
    }
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white p-6 rounded-md shadow-md w-96">
        <h2 className="text-xl font-bold mb-4">{character.id !== -1 ? "Edit Character" : "Create Character"}</h2>
        <div className="mb-4">
          <label htmlFor="name" className="block text-sm font-medium text-gray-700">Name</label>
          <input
            type="text"
            id="name"
            name="name"
            value={name}
            onChange={(e) => setName(e.target.value)}
            className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2"
          />
        </div>
        <div className="mb-4">
          <label htmlFor="rarityId" className="block text-sm font-medium text-gray-700">Rarity</label>
          <select
            id="rarityId"
            name="rarityId"
            value={rarityId === -1 ? rarities[0].id : rarityId}
            onChange={(e) => setRarityId(parseInt(e.target.value))}
            className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2"
          >
            {rarities.map((rarity) => (
              <option key={rarity.id} value={rarity.id}>
                {rarity.name}
              </option>
            ))}
          </select>
        </div>
        <div className="mb-4">
        <label htmlFor="image" className="block text-sm font-medium text-gray-700">
          Image
        </label>
          <input
            type="file"
            id="image"
            name="image"
            accept="image/png"
            onChange={handleFileChange}
            className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2"
          />
          <p className="text-sm text-gray-500 mt-1">
            Only PNG files with a maximum size of 2MB are allowed.
          </p>
          {image && !validateImage(image) && (
            <p className="text-red-500 text-sm mt-1">
              Invalid file. Please select a PNG file smaller than 2MB.
            </p>
          )}
        </div>

        <div className="flex justify-end space-x-2">
          <button
            onClick={onClose}
            className="bg-gray-500 hover:bg-gray-700 text-white py-2 px-4 rounded"
          >
            Close
          </button>
          <button
            onClick={handleSaveClick}
            className="bg-blue-500 hover:bg-blue-700 text-white py-2 px-4 rounded"
          >
            Save
          </button>
        </div>
      </div>
    </div>
  );
};

export default InputCharacterModal;
