import React from "react";

interface ButtonProps {
  children: React.ReactNode; // Untuk konten di dalam tombol
  onClick?: () => void;      // Fungsi opsional untuk menangani klik
}

const Button: React.FC<ButtonProps> = ({ children, onClick }) => {
  return (
    <button
      onClick={onClick}
      className="w-full bg-gradient-to-r from-blue-400 to-green-600 text-white font-bold py-3 rounded-full"
    >
      {children}
    </button>
  );
};

export default Button;
