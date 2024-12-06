import React from "react";
import { toast } from "react-toastify";
import { handleRequest } from "../utils/api";

interface GachaSystemCreateModalProps {
  isOpen: boolean;
  onClose: () => void;
  onRefresh: () => void;
}

const GachaSystemCreateModal: React.FC<GachaSystemCreateModalProps> = ({
  isOpen,
  onClose,
  onRefresh,
}) => {

  const [name, setName] = React.useState("");

  if (!isOpen) return null;

  const handleCreate = async () => {
    if (name) {
      try {
        const createEndpoint = `${import.meta.env.VITE_GACHA_MASTER_URL}/create`;
        const response = await handleRequest<{message:string}>(createEndpoint, "POST", {
          name: name,
        });

        if (response.code === 200) {
          toast.success("Gacha system created successfully!");
          onClose();
          onRefresh();
        } else {
          toast.error(`Create failed: ${response.data.message}`);
        }
      } catch (error) {
        toast.error("An error occurred. Please try again.");
      }
    }
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white p-6 rounded-md shadow-md w-80">
        <h2 className="text-xl font-bold mb-4">Create Gacha System</h2>
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
        <div className="flex justify-end space-x-2">
          <button
            onClick={onClose}
            className="bg-gray-500 hover:bg-gray-700 text-white py-2 px-4 rounded"
          >
            Close
          </button>
          <button
            onClick={handleCreate}
            className="bg-blue-500 hover:bg-blue-700 text-white py-2 px-4 rounded"
          >
            Save
          </button>
        </div>
      </div>
    </div>
  );
};

export default GachaSystemCreateModal;
