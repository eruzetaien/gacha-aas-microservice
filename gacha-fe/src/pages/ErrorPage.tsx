import { useNavigate } from "react-router";

const ErrorPage: React.FC<{ message?: string }> = ({ message }) => {
  const navigate = useNavigate();

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-100">
      <div className="bg-white shadow-lg rounded-lg p-10 text-center max-w-md w-full">
        <h1 className="text-4xl font-bold text-red-500 mb-4">Oops!</h1>
        <p className="text-gray-700 text-lg mb-6">
          {message || "Something went wrong. Please try again later."}
        </p>
        <div className="flex justify-center gap-4">
          <button
            onClick={() => navigate("/")}
            className="bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded"
          >
            Go to Home
          </button>
        </div>
      </div>
    </div>
  );
};

export default ErrorPage;
