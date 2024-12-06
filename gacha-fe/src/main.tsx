import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { RouterProvider } from 'react-router';
import { ToastContainer } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';
import './index.css';
import { router } from './routes/router.tsx';


createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <RouterProvider router={router} /> 
    <ToastContainer />      
  </StrictMode>,
);