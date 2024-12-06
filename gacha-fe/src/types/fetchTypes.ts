export type FetchState<T> = {
    response: ApiResponse<T> | null;
    loading: boolean;
    error: string | null;
};

export type ApiResponse<T> = {
    code: number;
    status: string;
    data: T; // Data bisa berupa apa saja, seperti objek atau array
  };
  