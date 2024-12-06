import React from "react";
import { Rarity } from "../types/RarityType";

interface InputRarityModalProps {
  isOpen: boolean;
  rarity: Rarity | null;
  setRarity: (rarity: Rarity) => void;
  onClose: () => void;
  onUpdate: () => void;
  onCreate: () => void;
}

const InputRarityModal: React.FC<InputRarityModalProps> = ({
  isOpen,
  rarity,
  setRarity,
  onClose,
  onUpdate,
  onCreate,
}) => {
  if (!isOpen || !rarity) return null;

  const handleRarityInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;

    if (name === "chance" && (parseInt(value) <= 0 || parseInt(value) > 100)) {
      return;
    }

    if (rarity) {
      setRarity({ ...rarity, [name]: value });
    }
  };

  const handleSaveClick = () => {
    if (rarity.id !== -1){
      onUpdate();
    }
    else{
      onCreate();
    }
  }
    

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white p-6 rounded-md shadow-md w-80">
        <h2 className="text-xl font-bold mb-4">Edit Rarity</h2>
        <div className="mb-4">
          <label htmlFor="name" className="block text-sm font-medium text-gray-700">Name</label>
          <input
            type="text"
            id="name"
            name="name"
            value={rarity.name}
            onChange={handleRarityInputChange}
            className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2"
          />
        </div>
        <div className="mb-4">
          <label htmlFor="chance" className="block text-sm font-medium text-gray-700">Chance (%)</label>
          <input
            type="number"
            id="chance"
            name="chance"
            value={rarity.chance}
            onChange={handleRarityInputChange}
            step="0.01"
            className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2"
          />
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

export default InputRarityModal;
