interface DeoResponse<T = any> {
    success: boolean;
    message?: string;
    data?: T;
}
interface Document {
    _id: string;
    [key: string]: any;
}
interface ListDocumentsOptions {
    filters?: {
        [key: string]: string;
    };
    sortBy?: string;
    order?: "asc" | "desc";
    limit?: number;
    offset?: number;
}
declare class Collection {
    private baseUrl;
    constructor(baseUrl: string);
    createDocument<T extends object>(document: T): Promise<DeoResponse<Document & T>>;
    listDocuments(options?: ListDocumentsOptions): Promise<DeoResponse<Document[]>>;
    readDocument(documentId: string): Promise<DeoResponse<Document>>;
    updateDocument<T extends object>(documentId: string, document: T): Promise<DeoResponse<Document & T>>;
    deleteDocument(documentId: string): Promise<DeoResponse<any>>;
}
declare class Database {
    private baseUrl;
    collections: {
        [key: string]: Collection;
    };
    constructor(baseUrl: string);
    createCollection(collectionName: string): Promise<DeoResponse>;
    listCollections(): Promise<DeoResponse<string[]>>;
    deleteCollection(collectionName: string): Promise<DeoResponse>;
}
export declare class DeoClient {
    private baseUrl;
    dbs: {
        [key: string]: Database;
    };
    constructor(host?: string);
    createDatabase(dbName: string): Promise<DeoResponse>;
    listDatabases(): Promise<DeoResponse<string[]>>;
    deleteDatabase(dbName: string): Promise<DeoResponse>;
}
export {};
