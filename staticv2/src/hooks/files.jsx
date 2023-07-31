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