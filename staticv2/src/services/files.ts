import { User, cFile } from "./models.ts";
import getAuthBasic from "./utils.ts";
import getContentTypeByFileName from "../utils/content-types";
import { getCacheService } from "./cache.ts";

interface FilesAPI {
    getFileInfo: (fileId: string) => Promise<cFile>;
    getShareFile: (fileId: string) => Promise<cFile>;
    getShareFileInfo(fileId: string): Promise<cFile>;
    getFiles: () => Promise<Array<cFile>>;
    updateFile: (fileId: string, fileData: any) => Promise<{ data: cFile, message: string }>;
    deleteFile: (fileId: string, force: boolean) => Promise<{ data: cFile, message: string }>;
    deleteFiles: (fileIds: string[], force: boolean) => Promise<{ data: cFile[], message: string }>;
    addUserToFile: (fileId: string, userId: string) => Promise<{ data: cFile, message: string }>;
    removeUserFromFile: (fileId: string, userId: string) => Promise<{ data: cFile, message: string }>;
    uploadFile: (files: cFile[]) => Promise<{ data: cFile[], message: string }>;
    getFilenameInfo: (fileRawData: cFile, withToken: boolean, additionalChanges: (data: cFile) => cFile) => cFile;
    getUserInfo: (userRawData: User, withToken: boolean, additionalChanges: (data: User) => User) => User;
    addTokenParam: (url: string) => string;
};

const getFileService = (baseUrlInput: string, tokenInput: string) : FilesAPI => {
    const { addTokenParam, cAxios, baseUrl } = getAuthBasic(baseUrlInput, tokenInput);
    const { setCacheFiles, addCacheFiles, removeCacheFile, updateCacheFile } = getCacheService();

    return {
        addTokenParam,
        async getFileInfo(fileId: string): Promise<cFile> {
            try {
                const res = await cAxios.get(`/api/files/${fileId}/info`);
                return res.data;
            } catch (err : any) {
                throw new Error(err.response.data.message);
            }
        },    
        async getShareFile(fileId: string): Promise<cFile> {
            try {
                const res = await cAxios.get(`/api/files/share/${fileId}`);
                return res.data ?? {};
            } catch (err : any) {
                throw new Error(err.response.data.message);
            }
        },
        async getShareFileInfo(fileId: string): Promise<cFile> {
            try {
                const res = await cAxios.get(`/api/files/share/${fileId}/info`);
                return res.data ?? {};
            } catch (err : any) {
                throw new Error(err.response.data.message);
            }
        },        
        async getFiles(): Promise<cFile[]> {
            try {                
                const res = await cAxios.get('/api/files');
                setCacheFiles(res.data ?? []);
                return res.data ?? [];
            } catch (err : any) {
                throw new Error(err.response.data.message);
            }
        },        
        async updateFile(fileId: string, fileData: cFile): Promise<{ data: cFile; message: string; }> {
            try {
                const res = await cAxios.put(`/api/files/${fileId}`, fileData);
                updateCacheFile(res.data.id, res.data);
                return {
                    data: res.data,
                    message: `${res.data.filename} (${res.data.id}) updated successfully`,
                };
            } catch (err : any) {
                throw new Error(err.response.data.message);
            }
        },        
        async deleteFile(fileId: string, force: boolean): Promise<{ data: cFile; message: string; }> {
            try {
                const res = await cAxios.delete(`/api/files/${fileId}${force ? '?force=yes' : ''}`);
                removeCacheFile(res.data.id);
                return {
                    data: res.data,
                    message: `${res.data.filename} (${res.data.id}) deleted successfully`,
                };
            } catch (err : any) {
                throw new Error(err.response.data.message);
            }
        },
        async deleteFiles(fileIds: string[], force: boolean): Promise<{ data: cFile[]; message: string; }> {
            try {
                const res = await cAxios.delete(`/api/files${force ? '?force=yes' : ''}`, {
                    data: fileIds,
                });

                console.log(res);
                
                res.data.forEach((file: cFile) => removeCacheFile(file.id));
                return {
                    data: res.data,
                    message: `${res.data.length} files deleted successfully`,
                };
            } catch (err : any) {
                throw new Error(err.response.data.message);
            }
        },   
        async addUserToFile(fileId: string, userId: string): Promise<{ data: cFile; message: string; }> {
            try {
                const res = await cAxios.post(`/api/files/share/${fileId}/user/${userId}`);
                res.data.sharedWith ||= [];
                return {
                    data: res.data,
                    message: 'User added successfully',
                };
            } catch (err : any) {
                throw new Error(err.response.data.message);
            }
        },
        async removeUserFromFile(fileId: string, userId: string): Promise<{ data: cFile; message: string; }> {
            try {
                const res = await cAxios.delete(`/api/files/share/${fileId}/user/${userId}`);
                res.data.sharedWith ||= [];
                return {
                    data: res.data,
                    message: 'User removed successfully',
                };
            } catch (err : any) {
                throw new Error(err.response.data.message);
            }
        },
        async uploadFile(files: cFile[]): Promise<{ data: cFile[]; message: string; }> {
            try {
                const res = await cAxios.post('/api/files', files);
                addCacheFiles(res.data);
                return {
                    data: res.data,
                    message: `${res.data.length} files uploaded successfully`,
                };
            } catch (err : any) {
                throw new Error(err.response.data.message);
            }
        },
        getFilenameInfo(fileRawData: cFile, withToken: boolean, additionalChanges: (data: cFile) => cFile): cFile {
            if(additionalChanges === undefined) additionalChanges = (data: cFile) => data;

            const filenameArray = fileRawData.filename.split('.');
            const fileUrl = `${baseUrl}/api/files/${fileRawData.id}`;
            const sharedFileUrl = `${baseUrl}/api/files/share/${fileRawData.id}`;
        
            return additionalChanges({
                id: fileRawData.id,
                filename: fileRawData?.filename,
                name: filenameArray.slice(0, -1).join('.'),
                extension: filenameArray.pop() ?? '',
                contentType: getContentTypeByFileName(fileRawData?.filename),
                url: withToken ? addTokenParam(fileUrl) : fileUrl,
                sharedUrl: withToken ? addTokenParam(sharedFileUrl) : sharedFileUrl,
                createdAt: fileRawData?.createdAt,
                updatedAt: fileRawData?.updatedAt,
                parentId: fileRawData?.parentId,
                children: fileRawData?.children || [],
                savedLocal: false
            });
        },
        getUserInfo(userRawData: User, withToken: boolean, additionalChanges: (data: User) => User): User {
            if(additionalChanges === undefined) additionalChanges = (data: User) => data;

            const iconUrl = `${baseUrl}/api/users/icon/${userRawData.id}`;
            return additionalChanges({
                id: userRawData.id,
                icon: withToken ? addTokenParam(iconUrl) : iconUrl,
                username: userRawData.username,
            });
        }
    }
}

export default getFileService;