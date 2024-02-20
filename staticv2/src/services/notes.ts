import { Note } from "./models";
import getAuthBasic from "./utils.ts";

interface NoteAPI {
    getNote() : Promise<Note>
    setNote(note : string) : Promise<{ data: Note, message: string}>
};

const getNoteService = (baseUrlInput: string, tokenInput: string) : NoteAPI => {
    const { cAxios } = getAuthBasic(baseUrlInput, tokenInput);

    return {
        getNote: async () => {
            try {
                const res = await cAxios.get('/api/notes');
                return res.data;
            } catch(err) {
                throw new Error(err.response.data.message);
            }
        },
        setNote: async (note : string) => {
            try {
                const res = await cAxios.post('/api/notes', { content: note });
                return {
                    data: res.data,
                    message: 'Note saved successfully',
                }
            } catch(err) {
                throw new Error(err.response.data.message);
            }
        },
    }
}

export default getNoteService;