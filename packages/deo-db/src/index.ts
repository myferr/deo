interface DeoResponse<T = any> {
  success: boolean;
  message?: string;
  data?: T;
}

interface CreateDatabaseRequest {
  db_name: string;
}

interface CreateCollectionRequest {
  collection_name: string;
}

interface Document {
  _id: string;
  [key: string]: any;
}

class DeoError extends Error {
  constructor(message: string) {
    super(message);
    this.name = "DeoError";
  }
}

async function _fetch<T>(
  url: string,
  options?: RequestInit,
): Promise<DeoResponse<T>> {
  const response = await fetch(url, options);
  const clonedResponse = response.clone(); // Clone the response
  let data: DeoResponse<T>;

  try {
    const json = await response.json();
    if (typeof json !== "object" || json === null || !("success" in json)) {
      throw new DeoError(
        `Invalid JSON structure received from ${url}: ${clonedResponse.status} ${clonedResponse.statusText}. Raw response: ${await clonedResponse.text()}`,
      );
    }
    data = json as DeoResponse<T>;
  } catch (error) {
    throw new DeoError(
      `Failed to parse JSON response from ${url}: ${clonedResponse.status} ${clonedResponse.statusText}. Raw response: ${await clonedResponse.text()}`,
    );
  }

  if (!data.success) {
    throw new DeoError(data.message || "An unknown error occurred.");
  }

  return data;
}

class Collection {
  private baseUrl: string;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
  }

  async createDocument<T extends object>(
    document: T,
  ): Promise<DeoResponse<Document & T>> {
    return _fetch<Document & T>(`${this.baseUrl}/documents`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(document),
    });
  }

  async listDocuments(): Promise<DeoResponse<Document[]>> {
    return _fetch<Document[]>(`${this.baseUrl}/documents`);
  }

  async readDocument(documentId: string): Promise<DeoResponse<Document>> {
    return _fetch<Document>(`${this.baseUrl}/documents/${documentId}`);
  }

  async updateDocument<T extends object>(
    documentId: string,
    document: T,
  ): Promise<DeoResponse<Document & T>> {
    return _fetch<Document & T>(`${this.baseUrl}/documents/${documentId}`, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(document),
    });
  }

  async deleteDocument(documentId: string): Promise<DeoResponse<any>> {
    return _fetch<any>(`${this.baseUrl}/documents/${documentId}`, {
      method: "DELETE",
    });
  }
}

class Database {
  private baseUrl: string;
  public collections: { [key: string]: Collection };

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
    this.collections = new Proxy(
      {},
      {
        get: (target: Record<string, Collection>, name: string) => {
          if (!target[name]) {
            target[name] = new Collection(
              `${this.baseUrl}/collections/${name}`,
            );
          }
          return target[name];
        },
      },
    );
  }

  async createCollection(collectionName: string): Promise<DeoResponse> {
    return _fetch(`${this.baseUrl}/collections`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ collection_name: collectionName }),
    });
  }

  async listCollections(): Promise<DeoResponse<string[]>> {
    return _fetch<string[]>(`${this.baseUrl}/collections`);
  }

  async deleteCollection(collectionName: string): Promise<DeoResponse> {
    return _fetch(`${this.baseUrl}/collections/${collectionName}`, {
      method: "DELETE",
    });
  }
}

export class DeoClient {
  private baseUrl: string;
  public dbs: { [key: string]: Database };

  constructor(host: string = "http://localhost:6741") {
    this.baseUrl = `${host}/api`;
    this.dbs = new Proxy(
      {},
      {
        get: (target: Record<string, Database>, name: string) => {
          if (!target[name]) {
            target[name] = new Database(`${this.baseUrl}/dbs/${name}`);
          }
          return target[name];
        },
      },
    );
  }

  async createDatabase(dbName: string): Promise<DeoResponse> {
    return _fetch(`${this.baseUrl}/dbs`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ db_name: dbName }),
    });
  }

  async listDatabases(): Promise<DeoResponse<string[]>> {
    return _fetch<string[]>(`${this.baseUrl}/dbs`);
  }

  async deleteDatabase(dbName: string): Promise<DeoResponse> {
    return _fetch(`${this.baseUrl}/dbs/${dbName}`, {
      method: "DELETE",
    });
  }
}
