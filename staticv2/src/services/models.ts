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