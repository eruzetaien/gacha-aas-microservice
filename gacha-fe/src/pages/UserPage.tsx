import { useEffect, useState } from "react";
import { useNavigate } from "react-router";
import { toast } from "react-toastify";
import GachaSystemCreateModal from "../components/GachaSystemModal";
import { DeletIcon } from "../components/Icon";
import { LogoutButton } from "../components/LogoutButton";
import TextLink from "../components/TextLink";
import { GachaSystem } from "../types/gachaSystemType";
import { handleRequest } from "../utils/api";

const UserPage: React.FC = () => {
  const [gachaSystems, setGachaSystems] = useState<GachaSystem[]>([]);

  const [refreshToggle, setRefreshToggle] = useState(false);
  const [isModalOpen, setIsModalOpen] = useState(false);

  const navigate = useNavigate();

  const fetchData = async () => {
    const endpoint = `${import.meta.env.VITE_GACHA_MASTER_URL}/all`;
    try {
      const response = await handleRequest<GachaSystem[]>(endpoint, "GET");
    
      if (response.code === 200) {
        setGachaSystems(response.data);
      } else {
        toast.error(`Error: ${response.status}`);
      }
      } catch (error) {
        if (error instanceof Error) {
          if (error.message === "401") {
            navigate("/");
            toast.error("You are not authorized to view the page.");
          }
        } else {
          toast.error("An error occurred. Please try again : " + error);
        }
      }
    }

    useEffect(() => {
      fetchData();
    }, [refreshToggle]); 

  const handleDelete = async (id: number) => {
    const confirmation = window.confirm("Are you sure you want to delete this Gacha system?");
    if (confirmation) {
      const deleteGachaEndpoint = `${import.meta.env.VITE_GACHA_MASTER_URL}/id/${id}`;
      try {
        const response = await handleRequest<{message:string}>(deleteGachaEndpoint, "DELETE");
      
        if (response.code === 200) {
          toast.success("Delete gacha system success");
          setRefreshToggle(!refreshToggle)
        } else {
          toast.error(`Authentication failed: ${response.data.message}`);
        }
      } catch (error) {
          toast.error("An error occurred. Please try again : " + error);
      }
    }
  };

  

  return (
    <div className="min-h-screen w-screen flex items-center justify-center bg-gradient-to-br from-blue-600 to-green-400">
      <LogoutButton />
      <div className="bg-white bg-opacity-20 backdrop-blur-md rounded-lg px-6 py-10 shadow-lg max-w-xl w-full"> 
        <h1 className="text-3xl font-bold text-white text-center mb-6">
          Gacha Systems
        </h1>
        <table className="min-w-full table-auto text-left text-white">
        <thead>
            <tr className="border-b flex justify-between">
            <div>
              <th className="px-4 py-2">#</th>
              <th className="px-4 py-2">Name</th>
            </div>
            <th className="px-4 py-2">
              <button
                  onClick={() => setIsModalOpen(true)}
                  className="bg-transparent hover:bg-white text-white hover:bg-opacity-50 outline outline-white outline-1 rounded px-2 py-1 flex items-center justify-center">
                  Create
                  </button>
            </th>
            </tr>
          </thead>
          <tbody>
          {gachaSystems && gachaSystems.length > 0 ? (
            gachaSystems.map((system, index) => (
              <tr key={system.id} className="border-b flex justify-between items-center">
                <div>
                  <td className="px-4 py-2">{index + 1}</td>
                  <td className="px-4 py-2">
                    <TextLink text={system.name} link={`/gacha-system/${system.id}`} />
                  </td>
                </div>
                <td className="px-4 py-2">
                  <button
                    onClick={() => handleDelete(system.id)}
                    className="bg-red-500 hover:bg-red-700 rounded p-2 flex items-center justify-center"
                  >
                    <DeletIcon />
                  </button>
                </td>
              </tr>
            ))
          ) : (
            <tr>
              <td colSpan={3} className="text-center py-4 text-white">
                No data available.
              </td>
            </tr>
          )}
        </tbody>
        </table>
      </div>
    {/* Modal */}
      <GachaSystemCreateModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onRefresh={() => setRefreshToggle(!refreshToggle)}
      />     
    </div>
  );
};

export default UserPage;
