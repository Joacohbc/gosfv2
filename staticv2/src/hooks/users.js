import { useContext, useMemo } from "react";
import AuthContext from "../context/auth-context";
import getUserService from "../services/user";

export const useUsers = () => {
    const { token, baseUrl} = useContext(AuthContext);
    return useMemo(() => getUserService(baseUrl, token), [baseUrl, token]);
}