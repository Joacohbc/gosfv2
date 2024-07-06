export type cFile = {
    id: number;
    filename: string;
    name?: string;
    extension?: string;
    contentType?: string;
    url?: string;
    sharedUrl?: string;
    sharedWith: User[];
    shared: boolean;
    createdAt: string;
    updatedAt: string;
    owner_id: number;
    isDir: boolean;
    parentId: number;
    children: cFile[];
    savedLocal?: boolean;
}

export type User = {
    id: number;
    icon: string;
    username: string;
}

export type Note = {
    id?: number;
    title?: string;
    content: string;
    createdAt?: string;
    updatedAt?: string;
    parentId?: number;
    savedLocal?: boolean;
}
