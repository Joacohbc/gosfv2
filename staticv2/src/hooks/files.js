import { useCallback, useContext } from "react";
import getContentType from "../utils/content-types";
import AuthContext from "../context/auth-context";

/**
 * Un Custom Hook para poder obtener extraer de un archivo (que viene del backend) la información que se necesita.
 */
export const useGetInfo = () => {
    const { baseUrl, addTokenParam } = useContext(AuthContext);

    // Función por defecto para realizar cambios adicionales a los datos
    const defaultAdditionalChanges = (data) => data;

    return {
        /**
        * Función que extrae información de un archivo (que viene del backend).
        * @param {Object} fileRawData - Los datos del archivo (que vienen del backend).
        * @param {Boolean} withToken - Si se debe agregar el token a la url del archivo.
        * @param {Function} additionalChanges - Una función que realiza cambios adicionales a los datos.
        * @returns {Object} Un objeto con la información del archivo (id, filename, name, extension, contentType, url, sharedUrl).
        */
        getFilenameInfo: useCallback((fileRawData, withToken = false, additionalChanges = defaultAdditionalChanges) => {
            const filenameArray = fileRawData.filename.split('.');
            const fileUrl = `${baseUrl}/api/files/${fileRawData.id}`;
            const sharedFileUrl = `${baseUrl}/api/files/share/${fileRawData.id}`;
            
            return additionalChanges({
                id: fileRawData.id,
                filename: fileRawData?.filename,
                name: filenameArray.slice(0, -1).join('.'),
                extension: filenameArray.length > 1 ? filenameArray.pop() : '',
                contentType: getContentType(fileRawData?.filename),
                url: withToken ? addTokenParam(fileUrl) : fileUrl,
                sharedUrl: withToken ? addTokenParam(sharedFileUrl) : sharedFileUrl,
            });
        }, [ baseUrl, addTokenParam ]),

        /**
         * Función que extrae información de un usuario (que viene del backend).
         * @param {Object} userRawData - Los datos del usuario (que vienen del backend).
         * @param {Boolean} withToken - Si se debe agregar el token a la url del icono del usuario.
         * @param {Function} additionalChanges - Una función que realiza cambios adicionales a los datos.
         * @returns {Object} Un objeto con la información del usuario (id, icon, username).
        */
        getUserInfo: useCallback((userRawData, withToken = false, additionalChanges = defaultAdditionalChanges) => { 
            const iconUrl = `${baseUrl}/api/users/icon/${userRawData.id}`;
            return additionalChanges({
                id: userRawData.id,
                icon: withToken ? addTokenParam(iconUrl) : iconUrl,
                username: userRawData.username,
            });
        }, [ baseUrl, addTokenParam ])
    }
}

/**
 * Un Custom Hook para consumir la API del Backend relacionada con los archivos.
 */
export const useFiles = () => {
    const { cAxios } = useContext(AuthContext);

    /**
     * Función que obtiene información de un archivo.
     * @param {String} fileId - El id del archivo.
     * @returns {Object} Un objeto con la información del archivo.
     */
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

    /**
     * Función que obtiene información de un archivo compartido con el usuario actual.
     * @param {String} fileId - El id del archivo.
     * @returns {Object} Un objeto con la información del archivo.
     */
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

    /**
     * Función que obtiene los archivos del usuario actual.
     * @returns {Array} Un arreglo con los archivos del usuario actual.
     */
    const getFiles = useCallback(async () => {
        if(!cAxios) return [];

        try {
            const res = await cAxios.get('/api/files');
            const data = res.data || [];
            localStorage.setItem('files', JSON.stringify(data));
            return data;
        } catch(err) {

            // Si no hay conexión al servidor y hay archivos en el localStorage
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

    /**
     * Función que actualiza un archivo.
     * @param {String} fileId - El id del archivo.
     * @param {Object} fileData - Los datos del archivo a actualizar.
     */
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

    /**
     * Función que elimina un archivo.
     * @param {String} fileId - El id del archivo.
     * @param {Boolean} force - Si se debe eliminar el archivo de la base de datos.
     * @returns {Object} Un objeto con la información del archivo eliminado.
     */
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

    /**
     * Función que agrega un usuario a un archivo.
     * @param {String} fileId - El id del archivo.
     * @param {String} userId - El id del usuario.
     * @returns {Object} Un objeto con la información del archivo.
     */
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

    /**
     * Función que elimina un usuario de un archivo.
     * @param {String} fileId - El id del archivo.
     * @param {String} userId - El id del usuario.
     * @returns {Object} Un objeto con la información del archivo.
     */
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

    /**
     * Función que sube un archivo.
     * @param {Object} files - Los archivos a subir.
     * @returns {Object} Un objeto con la información del archivo.
     */
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