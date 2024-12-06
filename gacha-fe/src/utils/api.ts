import { ApiResponse } from "../types/fetchTypes";

export const handleRequest = async <T> (
    url: string,
    method: string,
    body: object | FormData | string = {},
    headers: Record<string, string> = {}
  ): Promise<ApiResponse<T>> => {
    const jwtToken = localStorage.getItem("jwt_token");
  
      // Jika ada JWT token, tambahkan ke headers Authorization
      if (jwtToken) {
        headers["Authorization"] = `Bearer ${jwtToken}`;
      }
  
      // Tentukan apakah body adalah FormData
      const isFormData = body instanceof FormData;
  
      // Jika body bukan FormData dan bukan string, anggap itu objek JSON
      if (!isFormData && typeof body !== "string") {
        body = JSON.stringify(body); // Mengubah objek menjadi JSON string
        headers["Content-Type"] = "application/json"; // Set header untuk JSON
      }
  
      // Kirim permintaan
      const response = await fetch(url, {
        method,
        headers: {
          ...headers, // Pastikan header dikirim, termasuk header Authorization jika ada JWT
        },
        ...(method !== "GET" && { body: body as BodyInit }), // Hanya tambahkan body jika bukan GET
      }); 
      
      try {
        const result = await response.json(); // Parse response JSON        
        return result;

      } catch (error) {
        throw new Error(response.status.toString() || "Something went wrong");
      }
  };
  