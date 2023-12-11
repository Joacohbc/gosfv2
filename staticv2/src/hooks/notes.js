import { useCallback, useContext } from "react";
import AuthContext from "../context/auth-context";

export const useNotes = () => {
    const { cAxios } = useContext(AuthContext);

    const getNote = useCallback(async () => {
        if(!cAxios) return [];

        try {
            const res = await cAxios.get('/api/notes');
            return res.data || [];
        } catch(err) {
            throw new Error(err.response.data.message);
        }
    }
    , [ cAxios ]);

    const setNote = useCallback(async (note) => {
        if(!cAxios) return '';

        try {
            const res = await cAxios.post('/api/notes', { content: note });
            return {
                data: res.data,
                message: 'Note saved successfully',
            }
        } catch(err) {
            throw new Error(err.response.data.message);
        }
    }
    , [ cAxios ]);

    return {
        getNote,
        setNote
    }
}