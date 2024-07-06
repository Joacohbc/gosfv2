import { User } from "./models";
import getAuthBasic from "./utils";

interface UsersAPI {
    getMyIconURL: (updated) => string;
    getMyInfo: () => Promise<User>;
    updateUser: (newUsername: string) => Promise<{ data: User, message: string }>;
    changePassword: (oldPassword: string, newPassword: string) => Promise<{ message: string }>;
    uploadIcon: (iconFile: File) => Promise<{ message: string }>;
    deleteIcon: () => Promise<{ message: string }>;
    deleteAccount: () => Promise<void>;
};

const getUserService = (baseUrlInput: string, tokenInput: string) : UsersAPI => {
    const { cAxios, baseUrl } = getAuthBasic(baseUrlInput, tokenInput);

    const rawToUser = (rawData: any) : User => {
        return {
            id: rawData.id,
            icon: rawData.icon,
            username: rawData.username,
        };
    }

    return {
        getMyIconURL: (updated = false) => {
            const url = new URL(`${baseUrl}/api/users/me/icon`);
            if(updated) url.searchParams.append('random', new Date().getTime().toString());
            return url.toString();
        },
        getMyInfo: async () => {
            try {
                const res = await cAxios.get('/api/users/me');
                return rawToUser(res.data);
            } catch (err : any) {
                throw new Error(err.response.data.message);
            }
        },
        updateUser: async (newUsername: string) => {
            if(!newUsername) throw new Error('New username is required')
            try {
                const res = await cAxios.put('/api/users/me', { username: newUsername });
                return {
                    data: rawToUser(res.data),
                    message: 'Username updated successfully',
                }
            } catch(err : any) {
                throw new Error(err.response.data.message);
            }
        },
        changePassword: async (oldPassword: string, newPassword: string) => {
            if(!oldPassword) throw new Error('Old password is required')
            if(!newPassword) throw new Error('New password is required')
            try {
                await cAxios.put('/api/users/me/password', { old_password: oldPassword, new_password: newPassword });
                return {
                    message: 'Password updated successfully',
                }
            } catch(err) {
                throw new Error(err.response.data.message);
            }
        },
        uploadIcon: async (iconFile: File) => {
            if(!iconFile) throw new Error('Icon file is required')
            try {
                const formData = new FormData();
                formData.append('icon', iconFile);
                await cAxios.post('/api/users/me/icon', formData, {
                    headers: {
                        'Content-Type': 'multipart/form-data'
                    }
                });
                return {
                    message: 'Icon updated successfully',
                }
            } catch(err : any) {
                throw new Error(err.response.data.message);
            }
        },
        deleteIcon: async () => {
            try {
                await cAxios.delete('/api/users/me/icon');
                return {
                    message: 'Icon deleted successfully',
                }
            } catch(err : any) {
                throw new Error(err.response.data.message);
            }
        },
        deleteAccount: async () => {
            try {
                await cAxios.delete('/api/users/me');
            } catch(err : any) {
                throw new Error(err.response.data.message);
            }
        }

    }
}

export default getUserService;