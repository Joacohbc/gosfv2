import { useCallback, useContext } from "react";
import getContentType from "../utils/content-types";
import AuthContext from "../context/auth-context";

/**
 * Custom hook for getting file and user information.
 * @returns {Object} An object containing the following functions:
 *   - getFilenameInfo: A function that retrieves entire filename information (extension, url, content).
 *   - getUserInfo: A function that retrieves entire user information.
 */
export const useGetInfo = () => {
    const { baseUrl, addTokenParam } = useContext(AuthContext);

    // Do nothing with data by default
    const defaultAddicionalChanges = (data) => data;

    return {
        getFilenameInfo: useCallback((fileRawData, withToken = false, addicionalChanges = defaultAddicionalChanges) => {
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
        
        getUserInfo: useCallback((userRawData, withToken = false, addicionalChanges = defaultAddicionalChanges) => { 
            const iconUrl = `${baseUrl}/api/users/icon/${userRawData.id}`;
            return addicionalChanges({
                id: userRawData.id,
                icon: withToken ? addTokenParam(iconUrl) : iconUrl,
                username: userRawData.username,
            });
        }, [ baseUrl, addTokenParam ])
    }
}

/**
 * Custom hook for handling file operations.
 * @returns {Object} An object containing various file-related functions.
 */
export const useFiles = () => {
    const { cAxios } = useContext(AuthContext);

    // Get total information about a file (if the user is the owner)
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

    // Get total information about a file (if is shared with the user)
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

    // Get all files
    const getFiles = useCallback(async () => {
        if(!cAxios) return [];

        try {
            const res = await cAxios.get('/api/files');
            const data = res.data || [];
            localStorage.setItem('files', JSON.stringify(data));
            return data;
        } catch(err) {

            // Si no hay conexion al servidor y hay archivos en el localStorage
            // se retornan los archivos del localStorage
            if(err.response === undefined) {
                const files = localStorage.getItem('files');
                return files ? JSON.parse(files) : [];
            }

            // Si el error es 401 (Unauthorized), 403 (Forbidden) o 400 (Bad Request)
            if(err.response.status === 401 || err.response.status === 403 || err.response.status === 400) {
                localStorage.removeItem('files');
                return [];
            }

            throw new Error(err.response ? err.response.data.message : err);
        }
    }, [ cAxios ]);

    // const getBlob = useCallback(async (fileId) => {
    //     if(!fileId) throw new Error('File id is required')
    //     if(!cAxios) return '';

    //     try {
    //         const res = await cAxios.get(`/api/files/${fileId}`);
    //         return res.data || '';
    //     } catch(err) {
    //         throw new Error(err.response.data.message);
    //     }
    // }, [ cAxios ]);

    // Update a file
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

    // Delete a file
    const deleteFile = useCallback(async (fileId, force) => {
        if(!fileId) throw new Error('File id is required')
        if(!cAxios) return {};

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

    // Share a file with a user
    const addUserToFile = useCallback(async (fileId, userId) => {
        if(!fileId) throw new Error('File id is required')
        if(!userId) throw new Error('User id is required')
        if(!cAxios) return {};
        
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

    // Remove a user from a file
    const removeUserFromFile = useCallback(async (fileId, userId) => {
        if(!fileId) throw new Error('File id is required')
        if(!userId) throw new Error('User id is required')
        if(!cAxios) return {};

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

    // Upload a file
    const uploadFile = useCallback(async (files) => {
        if(!files) throw new Error('Files are required')
        if(!cAxios) return {};

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