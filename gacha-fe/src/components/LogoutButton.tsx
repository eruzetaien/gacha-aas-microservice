import { useNavigate } from "react-router";
import { toast } from "react-toastify";

export const LogoutButton = () => {

    const navigate = useNavigate();

    const logout = () => {
        localStorage.removeItem("jwt_token");
        localStorage.removeItem("user_id");
        localStorage.removeItem("user_name");

        navigate("/");
        toast.success("You have been logged out.");
    }

    return(
        <button className='fixed right-5 top-5 outline-1 hover:bg-white hover:bg-opacity-50 bg-transparent rounded outline outline-white text-white' onClick={logout}>
            Logout
        </button>
    );
}