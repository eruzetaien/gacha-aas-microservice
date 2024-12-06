import { createBrowserRouter } from "react-router-dom";
import App from "../App";
import CharacterPage from "../pages/CharacterPage";
import GachaPullPage from "../pages/GachaPullPage";
import GachaSystemPage from "../pages/GachaSystemPage";
import UserPage from "../pages/UserPage";


export const router = createBrowserRouter([
    {
        path: "/",
        element: <App />,
    },
    {
        path:'manage', 
        element: <UserPage />,
    },
    {
        path:'gacha-system/:gachaSystemId', 
        children: [
            {path:'', element: <GachaSystemPage/>},
            {path:'character/:characterId', element: <CharacterPage/>}
        ]
    },
    {
        path:'gacha/:userId/:gachaSystemName', 
        element: <GachaPullPage/>
    },
]);

