import { useContext, useMemo } from "react";
import AuthContext from "../context/auth-context";
import getFileService from "../services/files";

/**
 * Un Custom Hook para poder obtener extraer de un archivo (que viene del backend) la información que se necesita.
 */
export const useGetInfo = () => {
    const { baseUrl, token } = useContext(AuthContext);
    return useMemo(() => getFileService(baseUrl, token), [baseUrl, token]);
}

/**
 * Un Custom Hook para consumir la API del Backend relacionada con los archivos.
 */
export const useFiles = () => {
    const { baseUrl, token } = useContext(AuthContext);
    return useMemo(() => getFileService(baseUrl, token), [baseUrl, token]);
}