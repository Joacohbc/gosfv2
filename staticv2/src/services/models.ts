export type cFile = {
    id: number;
    filename: string;
    name: string;
    extension: string;
    contentType: string;
    url: string;
    sharedUrl: string;
    createdAt: string;
    updatedAt: string;
    parentId: number;
    children: cFile[];
    savedLocal: boolean;
}

export type User = {
    id: number;
    icon: string;
    username: string;
}

const emptyUser: User = {
    id: 0,
    icon: '',
    username: '',
};

const emptyFile: cFile = {
    id: 0,
    filename: '',
    name: '',
    extension: '',
    contentType: '',
    url: '',
    sharedUrl: '',
    createdAt: '',
    updatedAt: '',
    parentId: 0,
    children: [],
    savedLocal: false,
};

const emptyFileList = [ emptyFile ];
const emptyUserList = [ emptyUser ];
export { emptyFile, emptyUser, emptyFileList, emptyUserList }