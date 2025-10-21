// src/index.ts
class DeoError extends Error {
  constructor(message) {
    super(message);
    this.name = "DeoError";
  }
}
async function _fetch(url, options) {
  const response = await fetch(url, options);
  const clonedResponse = response.clone();
  let data;
  try {
    const json = await response.json();
    if (typeof json !== "object" || json === null || !("success" in json)) {
      throw new DeoError(`Invalid JSON structure received from ${url}: ${clonedResponse.status} ${clonedResponse.statusText}. Raw response: ${await clonedResponse.text()}`);
    }
    data = json;
  } catch (error) {
    throw new DeoError(`Failed to parse JSON response from ${url}: ${clonedResponse.status} ${clonedResponse.statusText}. Raw response: ${await clonedResponse.text()}`);
  }
  if (!data.success) {
    throw new DeoError(data.message || "An unknown error occurred.");
  }
  return data;
}

class Collection {
  baseUrl;
  constructor(baseUrl) {
    this.baseUrl = baseUrl;
  }
  async createDocument(document) {
    return _fetch(`${this.baseUrl}/documents`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify(document)
    });
  }
  async listDocuments(options = {}) {
    const params = new URLSearchParams;
    if (options.filters) {
      for (const key in options.filters) {
        if (Object.prototype.hasOwnProperty.call(options.filters, key)) {
          params.append(`filter[${key}]`, options.filters[key]);
        }
      }
    }
    if (options.sortBy) {
      params.append("sort_by", options.sortBy);
      if (options.order) {
        params.append("order", options.order);
      }
    }
    if (options.limit !== undefined) {
      params.append("limit", String(options.limit));
    }
    if (options.offset !== undefined) {
      params.append("offset", String(options.offset));
    }
    const queryString = params.toString();
    const url = queryString ? `${this.baseUrl}/documents?${queryString}` : `${this.baseUrl}/documents`;
    return _fetch(url);
  }
  async readDocument(documentId) {
    return _fetch(`${this.baseUrl}/documents/${documentId}`);
  }
  async updateDocument(documentId, document) {
    return _fetch(`${this.baseUrl}/documents/${documentId}`, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify(document)
    });
  }
  async deleteDocument(documentId) {
    return _fetch(`${this.baseUrl}/documents/${documentId}`, {
      method: "DELETE"
    });
  }
}

class Database {
  baseUrl;
  collections;
  constructor(baseUrl) {
    this.baseUrl = baseUrl;
    this.collections = new Proxy({}, {
      get: (target, name) => {
        if (!target[name]) {
          target[name] = new Collection(`${this.baseUrl}/collections/${name}`);
        }
        return target[name];
      }
    });
  }
  async createCollection(collectionName) {
    return _fetch(`${this.baseUrl}/collections`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({ collection_name: collectionName })
    });
  }
  async listCollections() {
    return _fetch(`${this.baseUrl}/collections`);
  }
  async deleteCollection(collectionName) {
    return _fetch(`${this.baseUrl}/collections/${collectionName}`, {
      method: "DELETE"
    });
  }
}

class DeoClient {
  baseUrl;
  dbs;
  constructor(host = "http://localhost:6741") {
    this.baseUrl = `${host}/api`;
    this.dbs = new Proxy({}, {
      get: (target, name) => {
        if (!target[name]) {
          target[name] = new Database(`${this.baseUrl}/dbs/${name}`);
        }
        return target[name];
      }
    });
  }
  async createDatabase(dbName) {
    return _fetch(`${this.baseUrl}/dbs`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({ db_name: dbName })
    });
  }
  async listDatabases() {
    return _fetch(`${this.baseUrl}/dbs`);
  }
  async deleteDatabase(dbName) {
    return _fetch(`${this.baseUrl}/dbs/${dbName}`, {
      method: "DELETE"
    });
  }
}
export {
  DeoClient
};
