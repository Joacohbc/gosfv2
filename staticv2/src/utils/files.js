import getContentType from "./content-types";

export function getFileInfo(fileRawData, addicionalChanges = (file) => file) {
    const filenameArray = fileRawData.filename.split('.');
    return addicionalChanges({
        id: fileRawData.id,
        filename: fileRawData?.filename || '',
        name: filenameArray.slice(0, -1).join('.'),
        extesion: filenameArray.length > 1 ? filenameArray.pop() : '',
        contentType: getContentType(fileRawData?.filename || ''),
        // url: `${window.location.origin}/api/files/${fileRawData.id}`,
        url: `${'http://localhost:3000'}/api/files/${fileRawData.id}`,
        // sharedUrl: `${window.location.origin}/share/${fileRawData.id}`,
        sharedUrl: `${'http://localhost:3000'}/api/files/share/${fileRawData.id}`,
    })
}

export function getUserInfo(userRawData, addicionalChanges = (user) => user) { 
    return addicionalChanges({
        id: userRawData.id,
        // icon: `${window.location.origin}/api/users/icon/${userRawData.id}` || '',
        icon: `${'http://localhost:3000'}/api/users/icon/${userRawData.id}` || '',
        username: userRawData.username || '',
    })
}