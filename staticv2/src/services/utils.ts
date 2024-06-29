import axios, { AxiosInstance } from "axios";


type AuthBasic = {
    baseUrl: string;
    token: string;
    cAxios: AxiosInstance;
    addTokenParam: (url: string) => string;
};

const getAuthBasic = (baseUrlInput : string, tokenInput : string) : AuthBasic => {
    const baseUrl = baseUrlInput;
    const token = tokenInput;
    
    const cAxios = axios.create({ baseURL: baseUrl });

    const addTokenParam = (url: string): string => {
        const urlObj = new URL(url);
        urlObj.searchParams.append('api-token', token);
        return urlObj.toString();
    };

    return {
        baseUrl,
        token,
        cAxios,
        addTokenParam,
    }
}

export default getAuthBasic;