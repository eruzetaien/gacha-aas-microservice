import React from 'react';
import { LogoutButton } from '../components/LogoutButton';
import TextLink from '../components/TextLink';

const NotFound: React.FC = () => {
  return (
    <div className="min-h-screen w-screen flex items-center justify-center bg-gradient-to-br from-blue-600 to-green-400">
      <LogoutButton />
      <div className="bg-white bg-opacity-20 backdrop-blur-md rounded-lg px-10 py-12 shadow-lg max-w-md w-full text-center">
        <h1 className="text-6xl font-bold text-white mb-4">404</h1>
        <p className="text-lg text-white mb-6">Oops! Not Found Coy.</p>
        <TextLink text="Go back to home" link='/manage'/>
      </div>
    </div>
  );
};

export default NotFound;
