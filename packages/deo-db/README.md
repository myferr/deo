# deo-db

A powerful and lightweight TypeScript client for interacting with the DeoDB remote database service. This library provides a convenient and type-safe way to manage databases, collections, and documents via a RESTful API.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
  - [Initialization](#initialization)
  - [Database Operations](#database-operations)
  - [Collection Operations](#collection-operations)
  - [Document Operations (CRUD)](#document-operations-crud)
  - [Error Handling](#error-handling)
- [API Reference](#api-reference)
  - [`DeoClient`](#deoclient)
    - [`new DeoClient(host?: string)`](#new-deoclienthost-string)
    - [`client.createDatabase(dbName: string): Promise<DeoResponse>`](#clientcreatedatasedbname-string-promisedeoresonse)
    - [`client.listDatabases(): Promise<DeoResponse<string[]>>`](#clientlistdatabases-promisedeoresponsestring)
    - [`client.deleteDatabase(dbName: string): Promise<DeoResponse>`](#clientdeletedatasedbname-string-promisedeoresonse-1)
    - [`client.dbs[dbName]: Database`](#clientdbsdbname-database)
  - [`Database`](#database)
    - [`db.createCollection(collectionName: string): Promise<DeoResponse>`](#dbcreatecollectioncollectionname-string-promisedeoresonse)
    - [`db.listCollections(): Promise<DeoResponse<string[]>>`](#dblistcollections-promisedeoresponsestring)
    - [`db.deleteCollection(collectionName: string): Promise<DeoResponse>`](#dbdeletecollectioncollectionname-string-promisedeoresonse-1)
    - [`db.collections[collectionName]: Collection`](#dbcollectionscollectionname-collection)
  - [`Collection`](#collection)
    - [`collection.createDocument<T extends object>(document: T): Promise<DeoResponse<Document & T>>`](#collectioncreatedocumentt-extends-objectdocument-t-promisedeoresonsedocument--t)
    - [`collection.listDocuments(): Promise<DeoResponse<Document[]>>`](#collectionlistdocuments-promisedeoresonsedocument)
    - [`collection.readDocument(documentId: string): Promise<DeoResponse<Document>>`](#collectionreaddocumentdocumentid-string-promisedeoresonsedocument)
    - [`collection.updateDocument<T extends object>(documentId: string, document: T): Promise<DeoResponse<Document & T>>`](#collectionupdatedocumentt-extends-objectdocumentid-string-document-t-promisedeoresonsedocument--t)
    - [`collection.deleteDocument(documentId: string): Promise<DeoResponse<any>>`](#collectiondeletedocumentdocumentid-string-promisedeoresonseany)
  - [`DeoResponse<T>`](#deoresponset)
  - [`Document`](#document)
  - [`DeoError`](#deoerror)
- [Project Structure](#project-structure)
- [Development](#development)
  - [Running Tests](#running-tests)
  - [Building the Project](#building-the-project)
- [Contributing](#contributing)
- [License](#license)

## Features

-   **Remote Database Interaction:** Connects to a DeoDB server to manage data.
-   **Hierarchical Data Model:** Supports databases, collections, and documents.
-   **Full CRUD Operations:** Create, Read, Update, and Delete documents within collections.
-   **Type-Safe:** Built with TypeScript for robust and predictable code.
-   **Familiar API:** Intuitive methods for common database operations.
-   **Error Handling:** Provides a custom `DeoError` for clear error reporting.

## Installation

To install `deo-db` in your Bun project, run the following command:

```bash
bun add deo-db
```

## Usage

### Initialization

First, import `DeoClient` and initialize it, optionally specifying the host of your DeoDB server.

```typescript
import { DeoClient } from 'deo-db';

// Connect to the default host (http://localhost:6741)
const client = new DeoClient();

console.log('DeoClient initialized.');
```

### Database Operations

Manage your databases on the DeoDB server.

```typescript
import { DeoClient } from 'deo-db';

const client = new DeoClient();

const dbName = 'myNewDatabase';

async function manageDatabases() {
  // Create a new database
  const createDbResponse = await client.createDatabase(dbName);
  if (createDbResponse.success) {
    console.log(`Database '${dbName}' created successfully.`);
  } else {
    console.error(`Failed to create database: ${createDbResponse.message}`);
  }

  // List all databases
  const listDbsResponse = await client.listDatabases();
  if (listDbsResponse.success && listDbsResponse.data) {
    console.log('Available Databases:', listDbsResponse.data);
  } else {
    console.error(`Failed to list databases: ${listDbsResponse.message}`);
  }

  // Access a specific database instance (dynamically created)
  const myDb = client.dbs[dbName];
  console.log(`Accessed database instance for '${dbName}'.`);

  // In a real application, you might delete the database later
  // const deleteDbResponse = await client.deleteDatabase(dbName);
  // if (deleteDbResponse.success) {
  //   console.log(`Database '${dbName}' deleted successfully.`);
  // } else {
  //   console.error(`Failed to delete database: ${deleteDbResponse.message}`);
  // }
}

manageDatabases();
```

### Collection Operations

Once you have a database, you can manage collections within it.

```typescript
import { DeoClient } from 'deo-db';

const client = new DeoClient();
const dbName = 'myNewDatabase'; // Ensure this database exists
const collectionName = 'users';

async function manageCollections() {
  const myDb = client.dbs[dbName];

  // Create a new collection
  const createCollectionResponse = await myDb.createCollection(collectionName);
  if (createCollectionResponse.success) {
    console.log(`Collection '${collectionName}' created in '${dbName}'.`);
  } else {
    console.error(`Failed to create collection: ${createCollectionResponse.message}`);
  }

  // List all collections in the database
  const listCollectionsResponse = await myDb.listCollections();
  if (listCollectionsResponse.success && listCollectionsResponse.data) {
    console.log(`Collections in '${dbName}':`, listCollectionsResponse.data);
  } else {
    console.error(`Failed to list collections: ${listCollectionsResponse.message}`);
  }

  // Access a specific collection instance (dynamically created)
  const usersCollection = myDb.collections[collectionName];
  console.log(`Accessed collection instance for '${collectionName}'.`);

  // In a real application, you might delete the collection later
  // const deleteCollectionResponse = await myDb.deleteCollection(collectionName);
  // if (deleteCollectionResponse.success) {
  //   console.log(`Collection '${collectionName}' deleted successfully.`);
  // } else {
  //   console.error(`Failed to delete collection: ${deleteCollectionResponse.message}`);
  // }
}

manageCollections();
```

### Document Operations (CRUD)

Perform CRUD operations on documents within a specific collection.

```typescript
import { DeoClient } from 'deo-db';

const client = new DeoClient();
const dbName = 'myNewDatabase';
const collectionName = 'users';

interface UserDocument {
  name: string;
  email: string;
  age?: number;
}

async function manageDocuments() {
  const usersCollection = client.dbs[dbName].collections[collectionName];

  // Create a document
  const newUser: UserDocument = { name: 'Alice Smith', email: 'alice@example.com' };
  const createDocResponse = await usersCollection.createDocument(newUser);
  let userId: string | undefined;

  if (createDocResponse.success && createDocResponse.data) {
    userId = createDocResponse.data._id;
    console.log('New user created:', createDocResponse.data);
  } else {
    console.error(`Failed to create document: ${createDocResponse.message}`);
    return;
  }

  if (!userId) return;

  // Read a document
  const readDocResponse = await usersCollection.readDocument(userId);
  if (readDocResponse.success && readDocResponse.data) {
    console.log('Read user:', readDocResponse.data);
  } else {
    console.error(`Failed to read document: ${readDocResponse.message}`);
  }

  // Update a document
  const updatedUser: Partial<UserDocument> = { age: 30 };
  const updateDocResponse = await usersCollection.updateDocument(userId, updatedUser);
  if (updateDocResponse.success && updateDocResponse.data) {
    console.log('Updated user:', updateDocResponse.data);
  } else {
    console.error(`Failed to update document: ${updateDocResponse.message}`);
  }

  // List all documents in the collection
  const listDocsResponse = await usersCollection.listDocuments();
  if (listDocsResponse.success && listDocsResponse.data) {
    console.log('All users:', listDocsResponse.data);
  } else {
    console.error(`Failed to list documents: ${listDocsResponse.message}`);
  }

  // List documents with options
  const filteredAndSortedDocs = await usersCollection.listDocuments({
    filters: { age: 30 },
    sortBy: 'name',
    order: 'asc',
    limit: 10,
  });

  if (filteredAndSortedDocs.success) {
    console.log('Filtered and sorted users:', filteredAndSortedDocs.data);
  } else {
    console.error(`Failed to list documents with options: ${filteredAndSortedDocs.message}`);
  }

  // Delete a document
  const deleteDocResponse = await usersCollection.deleteDocument(userId);
  if (deleteDocResponse.success) {
    console.log(`User with ID '${userId}' deleted.`);
  } else {
    console.error(`Failed to delete document: ${deleteDocResponse.message}`);
  }
}

manageDocuments();
```

### Error Handling

The `DeoClient` methods return `DeoResponse` objects. If `success` is `false`, the `message` field will contain details about the error. Additionally, network errors or invalid JSON responses will throw a `DeoError`.

```typescript
import { DeoClient, DeoError } from 'deo-db';

const client = new DeoClient('http://non-existent-host:1234'); // An invalid host to demonstrate error handling

async function demonstrateErrorHandling() {
  try {
    const response = await client.listDatabases();
    if (!response.success) {
      console.error('API Error:', response.message);
    } else {
      console.log('Success:', response.data);
    }
  } catch (error) {
    if (error instanceof DeoError) {
      console.error('DeoDB Client Error:', error.message);
    } else {
      console.error('Unexpected Error:', error);
    }
  }
}

demonstrateErrorHandling();
```

## API Reference

### `DeoClient`

The main class for interacting with the DeoDB server.

#### `new DeoClient(host?: string)`

Creates a new instance of the `DeoClient`.

-   `host` (optional `string`): The base URL of the DeoDB server. Defaults to `http://localhost:6741`.

#### `client.createDatabase(dbName: string): Promise<DeoResponse>`

Creates a new database on the server.

-   `dbName` (`string`): The name of the database to create.
-   Returns `Promise<DeoResponse>`: A promise that resolves with the API response.

#### `client.listDatabases(): Promise<DeoResponse<string[]>>`

Lists the names of all databases available on the server.

-   Returns `Promise<DeoResponse<string[]>>`: A promise that resolves with the API response containing an array of database names.

#### `client.deleteDatabase(dbName: string): Promise<DeoResponse>`

Deletes a database from the server.

-   `dbName` (`string`): The name of the database to delete.
-   Returns `Promise<DeoResponse>`: A promise that resolves with the API response.

#### `client.dbs[dbName]: Database`

Dynamically accesses a `Database` instance for the specified `dbName`. This property is a `Proxy` that creates `Database` objects on demand.

-   `dbName` (`string`): The name of the database to access.
-   Returns `Database`: An instance of the `Database` class.

### `Database`

Represents a specific database on the DeoDB server, providing methods to manage its collections.

#### `db.createCollection(collectionName: string): Promise<DeoResponse>`

Creates a new collection within this database.

-   `collectionName` (`string`): The name of the collection to create.
-   Returns `Promise<DeoResponse>`: A promise that resolves with the API response.

#### `db.listCollections(): Promise<DeoResponse<string[]>>`

Lists the names of all collections within this database.

-   Returns `Promise<DeoResponse<string[]>>`: A promise that resolves with the API response containing an array of collection names.

#### `db.deleteCollection(collectionName: string): Promise<DeoResponse>`

Deletes a collection from this database.

-   `collectionName` (`string`): The name of the collection to delete.
-   Returns `Promise<DeoResponse>`: A promise that resolves with the API response.

#### `db.collections[collectionName]: Collection`

Dynamically accesses a `Collection` instance for the specified `collectionName` within this database. This property is a `Proxy` that creates `Collection` objects on demand.

-   `collectionName` (`string`): The name of the collection to access.
-   Returns `Collection`: An instance of the `Collection` class.

### `Collection`

Represents a specific collection within a database, providing methods for document CRUD operations.

#### `collection.createDocument<T extends object>(document: T): Promise<DeoResponse<Document & T>>`

Creates a new document in this collection.

-   `document` (`T`): The document object to create. It will be assigned an `_id` by the server.
-   Returns `Promise<DeoResponse<Document & T>>`: A promise that resolves with the API response, including the created document with its `_id`.

#### `collection.listDocuments(options?: ListDocumentsOptions): Promise<DeoResponse<Document[]>>`

Lists all documents in this collection, with optional filtering, sorting, and pagination.

-   `options` (optional `ListDocumentsOptions`): An object to specify query parameters.
    -   `filters` (optional `object`): A key-value object for filtering documents. Example: `{ field: 'value' }`.
    -   `sortBy` (optional `string`): The field to sort by.
    -   `order` (optional `'asc' | 'desc'`): The sort order. Defaults to `'asc'`.
    -   `limit` (optional `number`): The maximum number of documents to return.
    -   `offset` (optional `number`): The number of documents to skip.
-   Returns `Promise<DeoResponse<Document[]>>`: A promise that resolves with the API response containing an array of documents.

#### `collection.readDocument(documentId: string): Promise<DeoResponse<Document>>`

Retrieves a single document by its ID from this collection.

-   `documentId` (`string`): The `_id` of the document to retrieve.
-   Returns `Promise<DeoResponse<Document>>`: A promise that resolves with the API response containing the document.

#### `collection.updateDocument<T extends object>(documentId: string, document: T): Promise<DeoResponse<Document & T>>`

Updates an existing document in this collection.

-   `documentId` (`string`): The `_id` of the document to update.
-   `document` (`T`): The partial or full document object with updated fields.
-   Returns `Promise<DeoResponse<Document & T>>`: A promise that resolves with the API response, including the updated document.

#### `collection.deleteDocument(documentId: string): Promise<DeoResponse<any>>`

Deletes a document by its ID from this collection.

-   `documentId` (`string`): The `_id` of the document to delete.
-   Returns `Promise<DeoResponse<any>>`: A promise that resolves with the API response.

### `DeoResponse<T>`

Interface for the standard response structure from the DeoDB API.

```typescript
interface DeoResponse<T = any> {
  success: boolean;
  message?: string;
  data?: T;
}
```

-   `success` (`boolean`): Indicates if the API request was successful.
-   `message` (optional `string`): Provides an error or informational message.
-   `data` (optional `T`): The payload returned by the API, if any.

### `Document`

Interface for a document stored in DeoDB, which always includes an `_id`.

```typescript
interface Document {
  _id: string;
  [key: string]: any;
}
```

### `DeoError`

Custom error class thrown by `DeoClient` for API-related issues or network errors.

## Project Structure

```
.
├── .gitignore
├── bun.lockb
├── package.json
├── README.md
├── tsconfig.json
├── dist/                 # Compiled JavaScript and TypeScript declaration files
│   ├── index.d.ts
│   └── index.js
└── src/
    ├── index.test.ts     # Unit tests for the DeoClient, covering API interactions
    └── index.ts          # Main source code for the DeoClient, Database, and Collection classes
```

-   `src/index.ts`: Contains the core `DeoClient`, `Database`, and `Collection` classes, responsible for making HTTP requests to the DeoDB server.
-   `src/index.test.ts`: Comprehensive test suite for `DeoClient`, ensuring all API interactions work as expected against a running DeoDB server.
-   `package.json`: Defines project metadata, scripts, and dependencies.
-   `tsconfig.json`: TypeScript configuration file.
-   `bun.lockb`: Bun's lockfile for deterministic dependency resolution.
-   `dist/`: Output directory for compiled JavaScript and TypeScript declaration files.

## Development

### Running Tests

To ensure everything is working correctly, you can run the test suite using Bun. Note that these tests require a DeoDB server to be running at `http://localhost:6741` (or the host specified in `DeoClient` constructor).

```bash
bun test
```

### Building the Project

To compile the TypeScript source code into JavaScript and generate declaration files, use the build script:

```bash
bun run build
```

This will output the compiled files into the `dist/` directory.

## Contributing

Contributions are welcome! If you find a bug or have a feature request, please open an issue. If you'd like to contribute code, please fork the repository and submit a pull request.

## License

This project is licensed under the MIT License. See the `LICENSE` file (if present) for details.
