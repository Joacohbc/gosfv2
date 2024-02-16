import { useContext, useMemo } from "react";
import AuthContext from "../context/auth-context";
import getNoteService from "../services/notes";

export const useNotes = () => {
    const { baseUrl, token } = useContext(AuthContext);
    return useMemo(() => getNoteService(baseUrl, token), [ baseUrl, token ]);
}