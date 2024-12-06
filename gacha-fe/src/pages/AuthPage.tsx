import { useState } from 'react';
import { useNavigate } from 'react-router';
import { toast } from 'react-toastify';
import Button from '../components/GreenishButton';
import { AuthResponse } from '../types/authResponseType';
import { handleRequest } from '../utils/api';

export default function AuthPage() {

  const [isLogin, setIsLogin] = useState(true)

  const toggleForm = () => setIsLogin(!isLogin)

  const [name, setName] = useState<string>("");
  const [username, setUsername] = useState<string>("");
  const [password, setPassword] = useState<string>("");

  const navigate = useNavigate();

  const handleAuth = async (endpoint: string, body: object) => {
    try {
      const response = await handleRequest<AuthResponse>(endpoint, "POST", body);

      if (response.code === 200 && response.data.userToken) {
        // Save JWT token to localStorage
        localStorage.setItem("jwt_token", response.data.userToken);
        localStorage.setItem("user_id", response.data.id.toString());
        localStorage.setItem("user_name", response.data.name);

        // Navigate to home page
        navigate("/manage");
      } else {
        toast.error(`Authentication failed: ${response.data.message}`);
      }
    } catch (error) {
      toast.error("An error occurred. Please try again.");
    }
  };

  // Login handler
  const handleLogin = () => {
    const endpoint = `${import.meta.env.VITE_GACHA_AUTH_URL}/login`;
    handleAuth(endpoint, { username, password });
  };

  // Register handler
  const handleRegister = () => {
    const endpoint = `${import.meta.env.VITE_GACHA_AUTH_URL}/register`;
    handleAuth(endpoint, { name, username, password });
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Enter") {
      if (isLogin) {
        handleLogin();
      } else {
        handleRegister();
      }
    }
  };

  return (
      <div className="min-h-screen w-screen flex items-center justify-center bg-gradient-to-br from-blue-600 to-green-400">
        <div className="bg-white bg-opacity-20 w-1/4 backdrop-blur-lg text-center  rounded-3xl shadow-xl overflow-hidden">
            <div className="px-10 py-10">
              <h2 className="text-4xl font-bold text-center text-white mb-6">
              {isLogin ? "Welcome Back!" : "Ready to create?"}
              </h2>
              <div className="space-y-2 mt-10 mb-4">
                {!isLogin && (
                  <div className='flex flex-col'>
                  <label className="text-white text-left">Name</label>
                  <input
                    type="text"
                    className="rounded-md outline-none mt-1 p-2"
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                  />
                </div>
                )}
                <div className='flex flex-col'>
                  <label className="text-white text-left">Username</label>
                  <input
                    type="username"
                    className="rounded-md outline-none mt-1 p-2"
                    value={username}
                    onChange={(e) => setUsername(e.target.value)}
                  />
                </div>
                <div className='flex flex-col'>
                  <label className="text-white text-left">Password</label>
                  <input
                    type="password"
                    className="rounded-md outline-none mt-1 p-2"
                    value={password}
                    onKeyDown={handleKeyDown} // Tangani tombol Enter
                    onChange={(e) => setPassword(e.target.value)}
                  />  
                  {!isLogin && (
                    <p className="text-sm text-white  mt-1">
                      Password must be at least 8 characters long.
                    </p>
                  )}
                </div>
              </div>
              {isLogin? (<Button onClick={handleLogin}>Login</Button>) : (<Button onClick={handleRegister}>Register</Button>) }
              
              <p className="mt-4 text-center text-white">
                <a
                  href="#"
                  onClick={toggleForm}
                  className="text-white hover:underline hover:text-white"
                >
                  {isLogin ? "Need an account? Sign up" : "Already have an account? Log in"}
                </a>
              </p>
            </div>
          </div>  
      </div>
  )
}

