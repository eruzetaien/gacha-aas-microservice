export const NoCharacterPage = () => {
    return (
      <div className="flex flex-col w-screen items-center justify-center h-screen bg-gradient-to-br from-blue-600 to-green-400">
        <div className="bg-white bg-opacity-20 backdrop-blur-md rounded-lg px-6 py-10 shadow-lg max-w-lg text-center">
          <h1 className="text-2xl font-bold text-white mb-4">
            No Characters Available
          </h1>
          <p className="text-white text-sm">
            There are currently no characters available to pull in this gacha system.
          </p>
          <p className="text-white text-sm mt-2">
            Please check back later or contact the creator.
          </p>
        </div>
      </div>
    );
  };
  