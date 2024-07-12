import { User, cFile } from "./models";

const LS_KEYS = {
    files: 'files',
    numberOfFiles: 'numberOfFiles',
    icon: 'icon',
    user: 'user'
};

type LocalStorageItem<T> = {
    value: T;
    timestamp: Date;
};

interface CacheAPI {
    getCacheFiles: () => LocalStorageItem<Array<cFile>>;
    setCacheFiles: (files: Array<cFile>) => void;
    addCacheFiles: (files: Array<cFile>) => void;
    removeCacheFile: (fileId: number) => void;
    updateCacheFile: (fileId: number, fileData: cFile) => void;
    getCacheNumberOfFiles: () => LocalStorageItem<Number>;
    setCacheNumberOfFiles: (numberOfFiles: number) => void;

    setCacheUser: (user: User) => void;
    getCacheUser: () => LocalStorageItem<User>;
    setCacheIcon: (icon: string) => void;
    getCacheIcon: () => LocalStorageItem<string>;
    
    clean: () => void;
}

export const getCacheService = () : CacheAPI => {
    return {
        getCacheFiles: getFiles,
        setCacheFiles: setFiles,
        getCacheNumberOfFiles: getNumberOfFiles,
        setCacheNumberOfFiles: setNumberOfFiles,
        addCacheFiles: addFiles,
        removeCacheFile: removeFile,
        updateCacheFile: updateFile,
        setCacheUser: setUser,
        getCacheUser: getUser,
        setCacheIcon: setIcon,
        getCacheIcon: getIcon,
        clean
    };
}

const setLocalStorage = (key: string, value: any) => {
    localStorage.setItem(key, JSON.stringify({
        value,
        timestamp: new Date(),
    }));
}

const getLocalStorage = (key) : LocalStorageItem<any> => {
    const item = localStorage.getItem(key);
    if (!item) {
        return { value: null, timestamp: new Date() };
    }
    const { value, timestamp } = JSON.parse(item);

    const parsedTimestamp = new Date(Date.parse(timestamp));
    return { value, timestamp: parsedTimestamp };
};

// Get the files from the local storage
const getFiles = () : LocalStorageItem<Array<cFile>> => {
    return getLocalStorage(LS_KEYS.files);
};

// Set the files in the local storage
const setFiles = (files: Array<cFile>) => {
    setLocalStorage(LS_KEYS.files, files);
    setNumberOfFiles(files.length);
};

// Get the number of files from the local storage
const getNumberOfFiles = () : LocalStorageItem<Number> => {
    return getLocalStorage(LS_KEYS.numberOfFiles);
};

// Set the number of files in the local storage (is called automatically when the files are set)
const setNumberOfFiles = (numberOfFiles: number) => {
    setLocalStorage(LS_KEYS.numberOfFiles, numberOfFiles);
};

const addFiles = (file: Array<cFile>) => {
    const { value: files } = getFiles();
    setFiles([...files, ...file]);
}

const removeFile = (fileId: number) => {
    const { value: files } = getFiles();
    const newFiles = files.filter(file => file.id !== fileId);
    setFiles(newFiles);
}

const updateFile = (fileId: number, fileData: cFile) => {
    const { value: files } = getFiles();
    const newFiles = files.map(file => {
        if (file.id === fileId) {
            return fileData;
        }
        return file;
    });
    setFiles(newFiles);
}

const setUser = (user: User) => {
    setLocalStorage(LS_KEYS.user, user);
}

const getUser = () : LocalStorageItem<User> => {
    return getLocalStorage(LS_KEYS.user);
}

const setIcon = (icon: string) => {
    setLocalStorage(LS_KEYS.icon, icon);
}

const getIcon = () : LocalStorageItem<string> => {
    return getLocalStorage(LS_KEYS.icon);
}


const clean = () => {
    localStorage.removeItem(LS_KEYS.files);
    localStorage.removeItem(LS_KEYS.numberOfFiles);
}