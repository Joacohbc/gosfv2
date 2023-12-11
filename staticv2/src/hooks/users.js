import { useCallback, useContext } from "react";
import AuthContext from "../context/auth-context";

export const useUsers = () => {
    const { cAxios, addTokenParam, baseUrl } = useContext(AuthContext);


    const getMyIconURL = useCallback(() => {
        const url = new URL(addTokenParam(baseUrl + '/api/users/me/icon'));
        url.searchParams.append('random', new Date().getTime());
        return url.toString();
    }, [ addTokenParam, baseUrl ]);

    const getMyInfo = useCallback(async () => {
        if(!cAxios) return '';
        
        try {
            const res = await cAxios.get('api/users/me');
            return res.data || {};
        } catch(err) {
            throw new Error(err.response.data.message);
        }
    }, [ cAxios ]);

    const updateUser = useCallback(async (newUsername) => {
        if(!newUsername) throw new Error('New username is required')
        if(!cAxios) return '';
        
        try {
            const res = await cAxios.put('/api/users/me', { username: newUsername });
            return {
                data: res.data,
                message: 'Username updated successfully',
            }
        } catch(err) {
            throw new Error(err.response.data.message);
        }
    }, [ cAxios ]);

    const changePassword = useCallback(async (oldPassword, newPassword) => {
        if(!oldPassword) throw new Error('Old password is required')
        if(!newPassword) throw new Error('New password is required')
        if(!cAxios) return '';

        try {
            await cAxios.put('/api/users/me/password', { old_password: oldPassword, new_password: newPassword });
            return {
                message: 'Password updated successfully',
            }
        } catch(err) {
            throw new Error(err.response.data.message);
        }
    }, [ cAxios ]);
    
    const uploadIcon = useCallback(async (iconFile) => {
        if(!iconFile) throw new Error('Icon file is required')
        if(!cAxios) return '';

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
        } catch(err) {
            throw new Error(err.response.data.message);
        }
    }, [ cAxios ]);

    const deleteIcon = useCallback(async () => {
        if(!cAxios) return '';
        
        try {
            await cAxios.delete('/api/users/me/icon');
            return {
                message: 'Icon deleted successfully',
            }
        } catch(err) {
            throw new Error(err.response.data.message);
        }
    }, [ cAxios ]);

    const deleteAccount = useCallback(async () => {
        if(!cAxios) return '';

        try {
            await cAxios.delete('/api/users/me');
        } catch(err) {
            throw new Error(err.response.data.message);
        }
    }, [ cAxios ]);

    return {
        getMyIconURL,
        getMyInfo,
        updateUser,
        changePassword,
        uploadIcon,
        deleteIcon,
        deleteAccount
    }
}