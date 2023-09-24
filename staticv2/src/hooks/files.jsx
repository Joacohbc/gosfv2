import { useCallback, useContext } from "react";
import getContentType from "../utils/content-types";
import AuthContext from "../context/auth-context";

export const useGetInfo = () => {
    const { baseUrl, addTokenParam } = useContext(AuthContext);

    return {
        getFileInfo: useCallback((fileRawData, withToken = false, addicionalChanges = (file) => file) => {
            const filenameArray = fileRawData.filename.split('.');
            const fileUrl = `${baseUrl}/api/files/${fileRawData.id}`;
            const sharedFileUrl = `${baseUrl}/api/files/share/${fileRawData.id}`;
            
            return addicionalChanges({
                id: fileRawData.id,
                filename: fileRawData?.filename,
                name: filenameArray.slice(0, -1).join('.'),
                extesion: filenameArray.length > 1 ? filenameArray.pop() : '',
                contentType: getContentType(fileRawData?.filename),
                url: withToken ? addTokenParam(fileUrl) : fileUrl,
                sharedUrl: withToken ? addTokenParam(sharedFileUrl) : sharedFileUrl,
            });
        }, [ baseUrl, addTokenParam ]),
        getUserInfo: useCallback((userRawData, withToken = false, addicionalChanges = (user) => user) => { 
            const iconUrl = `${baseUrl}/api/users/icon/${userRawData.id}`;
            return addicionalChanges({
                id: userRawData.id,
                icon: withToken ? addTokenParam(iconUrl) : iconUrl,
                username: userRawData.username,
            });
        }, [ baseUrl, addTokenParam ])
    }
}

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

export const useFiles = () => {
    const { cAxios } = useContext(AuthContext);

    const getFileInfo = useCallback(async (fileId) => {
        if(!fileId) throw new Error('File id is required')
        if(!cAxios) return {};

        try {
            const res = await cAxios.get(`/api/files/${fileId}/info`);
            return res.data || {};
        } catch(err) {
            throw new Error(err.response.data.message);
        }
    }, [ cAxios ]);

    const getShareFileInfo = useCallback(async (fileId) => {
        if(!fileId) throw new Error('File id is required')
        if(!cAxios) return {};

        try {
            const res = await cAxios.get(`/api/files/share/${fileId}/info`);
            return res.data || {};
        } catch(err) {
            throw new Error(err.response.data.message);
        }
    }, [ cAxios ]);

    const getFiles = useCallback(async () => {
        if(!cAxios) return [];

        try {
            const res = await cAxios.get('/api/files');
            return res.data || [];
        } catch(err) {
            throw new Error(err.response.data.message);
        }
    }, [ cAxios ]);

    const updateFile = useCallback(async (fileId, fileData) => {
        if(!fileId) throw new Error('File id is required')
        if(!cAxios) return '';

        try {
            const res = await cAxios.put(`/api/files/${fileId}`, fileData);

            return {
                data: res.data,
                message: `${res.data.filename} (${res.data.id}) updated successfully`,
            }
        } catch(err) {
            throw new Error(err.response.data.message);
        }
    }, [ cAxios ]);

    const deleteFile = useCallback(async (fileId, force) => {
        if(!fileId) throw new Error('File id is required')
        if(!cAxios) return '';

        try {
            const res = await cAxios.delete(`/api/files/${fileId}${force ? '?force=yes' : ''}`);
            return {
                data: res.data,
                message: `${res.data.filename} (${res.data.id}) deleted successfully`,
            };
        } catch(err) {
            throw new Error(err.response.data.message);
        }
    }, [ cAxios ]);

    const addUserToFile = useCallback(async (fileId, userId) => {
        if(!fileId) throw new Error('File id is required')
        if(!userId) throw new Error('User id is required')
        if(!cAxios) return '';
        
        try {
            const res = await cAxios.post(`/api/files/share/${fileId}/user/${userId}`);
            res.data.sharedWith ||= [];
            return {
                data: res.data,
                message: 'User added successfully',
            }
        } catch(err) {
            throw new Error(err.response.data.message);
        }
    }, [ cAxios ]);

    const removeUserFromFile = useCallback(async (fileId, userId) => {
        if(!fileId) throw new Error('File id is required')
        if(!userId) throw new Error('User id is required')
        if(!cAxios) return '';

        try {
            const res = await cAxios.delete(`/api/files/share/${fileId}/user/${userId}`);
            res.data.sharedWith ||= [];
            return {
                data: res.data,
                message: 'User removed successfully',
            };
        } catch(err) {
            throw new Error(err.response.data.message);
        }
    }, [ cAxios ]);

    const uploadFile = useCallback(async (files) => {
        try {
            const res = await cAxios.post('/api/files', files);
            return {
                data: res.data,
                message: `${res.data.length} files uploaded successfully`,
            }
        } catch(err) {
            throw new Error(err.response.data.message);
        }
    }, [ cAxios ]);
    
    return {
        getFileInfo, 
        getShareFileInfo,
        getFiles,
        updateFile,
        deleteFile,
        addUserToFile,
        removeUserFromFile,
        uploadFile
    }   
}

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