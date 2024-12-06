export type AuthResponse = {
    id: number;
    name: string;
    userToken: string;
    message? : string
};

export const isAuthResponse = (response: any): response is AuthResponse => {
    return response && typeof response.userToken === "string";
  };
  