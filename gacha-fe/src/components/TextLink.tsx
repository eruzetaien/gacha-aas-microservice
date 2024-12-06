import React from "react";
import { Link } from "react-router-dom";

interface TextLinkProps {
  text: string;
  link: string;
}

const TextLink: React.FC<TextLinkProps> = ({ text, link }) => {
  return (
    <Link
      to={link}
      className="text-white hover:text-blue-500 hover:underline break-all"
    >
      {text}
    </Link>
  );
};

export default TextLink;