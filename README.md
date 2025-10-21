<div align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="./assets/logo-dark.png" />
    <source media="(prefers-color-scheme: light)" srcset="./assets/logo.png" />
    <img alt="Project Logo" src="./assets/logo.png" width="100" />
  </picture>

  <h3>A lightweight, Go-based open source document-oriented JSON database.</h3>
</div>

---

**deo** ([/'dioh/](https://ipa-reader.com/?text=%CB%88dioh), pronounced `dee-oh`) is a Go-based, open source document-oriented database that stores data as msgpack-compressed JSON files.
It’s entirely accessed through a RESTful API, making it simple to integrate and easy to deploy as a single binary.

```
go install github.com/deo/deo@latest
```

---

## cURL Examples

### Create a document
```bash
curl -X POST http://localhost:6741/api/dbs/my_db/collections/my_collection/documents \
  -H "Content-Type: application/json" \
  -d '{"name": "John Doe", "age": 30}'
```

### Read a document

```bash
curl http://localhost:6741/api/dbs/my_db/collections/my_collection/documents/<document_id>
```

### Update a document

```bash
curl -X PUT http://localhost:6741/api/dbs/my_db/collections/my_collection/documents/<document_id> \
  -H "Content-Type: application/json" \
  -d '{"name": "John Doe", "age": 31}'
```

### Delete a document

```bash
curl -X DELETE http://localhost:6741/api/dbs/my_db/collections/my_collection/documents/<document_id>
```

---

## Architecture

### Storage Layout

```
$HOME/.deo/{database_name}/{collection_name}/*.msgpack
│      │              │                   │
│      │              │                   └── Documents (*.msgpack)
│      │              └── Collection
│      └── Deo Database
└── User Home Directory
```

### RESTful API Structure

```
localhost:6741   /api   /dbs   /{database_name}   /collections   /{collection_name}   /documents
│                │      │      │                  │              │                    │
│                │      │      │                  │              │                    └── Documents endpoint
│                │      │      │                  │              └── Collection name slug
│                │      │      │                  └── Collections endpoint
│                │      │      └── Database name slug
│                │      └── Database endpoint
│                └── API root
└── Hostname + Port
```

---

## Using Docker

Images:
- [`myferr/deo`](https://hub.docker.com/r/myferr/deo)
- [`ghcr.io/myferr/deo`](https://github.com/myferr/deo/pkgs/container/deo)

[Tags](https://github.com/myferr/deo/tags)

### Command line
To run Deo in a Docker container, use the following command:

```bash
docker run -d -p 6741:6741 --name deo ghcr.io/myferr/deo:latest
```

> [!TIP]
> `myferr/deo:latest` may not exist, so either don't specify a tag or use `ghcr.io/myferr/deo:latest`

### Compose
To run Deo in a Docker Compose setup, use the following configuration:

```yaml
services:
  deo:
    image: ghcr.io/myferr/deo:latest
    ports:
      - "6741:6741"
```

and then start the container:

```bash
docker-compose up -dv
```

## Studio
**deo** has a "studio" web UI to manage databases, collections, and documents which is useful for testing.
Accessible at [`http://localhost:6741/studio`](http://localhost:6741/studio)

### Demonstration

https://github.com/user-attachments/assets/bf61fd84-8498-4ca6-9376-10ef5dc9d32c

---

# API Documentation

**deo** relies upon a RESTful application programming interface (API) to interact with databases, collections, and documents. The same route can be called in different methods to perform various operations so to explain the API in a more detailed manner, read below for documentation.

> [!NOTE]
> Don't want to interact with RESTful APIs? [**deo** has libraries!](https://github.com/myferr/deo#libraries)

## Base URL
`http://localhost:6741/api`

## Error Responses
All API errors return a JSON object with `success: false` and a `message` field.

```json
{
  "success": false,
  "message": "Error description"
}
```

## Endpoints

### 1. Database Operations

#### 1.1 Create Database
*   **Endpoint:** `/api/dbs`
*   **Method:** `POST`
*   **Description:** Creates a new database.
*   **Request Body:**
    ```json
    {
      "db_name": "string"
    }
    ```
    *   `db_name` (string, required): The name of the database to create.
*   **Success Response:**
    *   **Code:** `201 Created`
    ```json
    {
      "success": true,
      "message": "Database created successfully"
    }
    ```
*   **Error Responses:**
    *   `400 Bad Request`: Invalid request body (e.g., `db_name` missing).
    *   `500 Internal Server Error`: Failed to create database (e.g., file system error).

#### 1.2 List Databases
*   **Endpoint:** `/api/dbs`
*   **Method:** `GET`
*   **Description:** Lists all available databases.
*   **Success Response:**
    *   **Code:** `200 OK`
    ```json
    {
      "success": true,
      "data": ["db1", "db2"]
    }
    ```
    *   `data` (array of strings): An array of database names.
*   **Error Responses:**
    *   `500 Internal Server Error`: Failed to list databases.

#### 1.3 Delete Database
*   **Endpoint:** `/api/dbs/:db_name`
*   **Method:** `DELETE`
*   **Description:** Deletes a specific database and all its contents (collections and documents).
*   **Path Parameters:**
    *   `db_name` (string, required): The name of the database to delete.
*   **Success Response:**
    *   **Code:** `200 OK`
    ```json
    {
      "success": true,
      "message": "Database deleted successfully"
    }
    ```
*   **Error Responses:**
    *   `500 Internal Server Error`: Failed to delete database.

### 2. Collection Operations

#### 2.1 Create Collection
*   **Endpoint:** `/api/dbs/:db_name/collections`
*   **Method:** `POST`
*   **Description:** Creates a new collection within a specified database.
*   **Path Parameters:**
    *   `db_name` (string, required): The name of the database.
*   **Request Body:**
    ```json
    {
      "collection_name": "string"
    }
    ```
    *   `collection_name` (string, required): The name of the collection to create.
*   **Success Response:**
    *   **Code:** `201 Created`
    ```json
    {
      "success": true,
      "message": "Collection created successfully"
    }
    ```
*   **Error Responses:**
    *   `400 Bad Request`: Invalid request body (e.g., `collection_name` missing).
    *   `500 Internal Server Error`: Failed to create collection.

#### 2.2 List Collections
*   **Endpoint:** `/api/dbs/:db_name/collections`
*   **Method:** `GET`
*   **Description:** Lists all collections within a specified database.
*   **Path Parameters:**
    *   `db_name` (string, required): The name of the database.
*   **Success Response:**
    *   **Code:** `200 OK`
    ```json
    {
      "success": true,
      "data": ["collection1", "collection2"]
    }
    ```
    *   `data` (array of strings): An array of collection names.
*   **Error Responses:**
    *   `500 Internal Server Error`: Failed to list collections.

#### 2.3 Delete Collection
*   **Endpoint:** `/api/dbs/:db_name/collections/:collection_name`
*   **Method:** `DELETE`
*   **Description:** Deletes a specific collection and all its documents within a database.
*   **Path Parameters:**
    *   `db_name` (string, required): The name of the database.
    *   `collection_name` (string, required): The name of the collection to delete.
*   **Success Response:**
    *   **Code:** `200 OK`
    ```json
    {
      "success": true,
      "message": "Collection deleted successfully"
    }
    ```
*   **Error Responses:**
    *   `500 Internal Server Error`: Failed to delete collection.

### 3. Document Operations

#### 3.1 Create Document
*   **Endpoint:** `/api/dbs/:db_name/collections/:collection_name/documents`
*   **Method:** `POST`
*   **Description:** Creates a new document within a specified collection. A unique `_id` (UUID) is automatically generated for the document.
*   **Path Parameters:**
    *   `db_name` (string, required): The name of the database.
    *   `collection_name` (string, required): The name of the collection.
*   **Request Body:**
    *   Any valid JSON object. The `_id` field will be added automatically.
    ```json
    {
      "field1": "value1",
      "field2": "value2"
    }
    ```
*   **Success Response:**
    *   **Code:** `201 Created`
    ```json
    {
      "success": true,
      "message": "Document created successfully",
      "data": {
        "_id": "uuid-string",
        "field1": "value1",
        "field2": "value2"
      }
    }
    ```
    *   `data` (object): The created document, including its generated `_id`.
*   **Error Responses:**
    *   `400 Bad Request`: Invalid JSON in request body.
    *   `500 Internal Server Error`: Failed to create document.

#### 3.2 List Documents
*   **Endpoint:** `/api/dbs/:db_name/collections/:collection_name/documents`
*   **Method:** `GET`
*   **Description:** Lists all documents within a specified collection.
*   **Path Parameters:**
    *   `db_name` (string, required): The name of the database.
    *   `collection_name` (string, required): The name of the collection.
*   **Success Response:**
    *   **Code:** `200 OK`
    ```json
    {
      "success": true,
      "data": [
        { "_id": "doc1-id", "fieldA": "valueA" },
        { "_id": "doc2-id", "fieldB": "valueB" }
      ]
    }
    ```
    *   `data` (array of objects): An array of documents.
*   **Error Responses:**
    *   `500 Internal Server Error`: Failed to list documents.

#### 3.3 Read Document
*   **Endpoint:** `/api/dbs/:db_name/collections/:collection_name/documents/:document_id`
*   **Method:** `GET`
*   **Description:** Retrieves a specific document by its ID.
*   **Path Parameters:**
    *   `db_name` (string, required): The name of the database.
    *   `collection_name` (string, required): The name of the collection.
    *   `document_id` (string, required): The ID of the document to retrieve.
*   **Success Response:**
    *   **Code:** `200 OK`
    ```json
    {
      "success": true,
      "data": {
        "_id": "document-id",
        "field1": "value1",
        "field2": "value2"
      }
    }
    ```
    *   `data` (object): The retrieved document.
*   **Error Responses:**
    *   `404 Not Found`: Document not found.
    *   `500 Internal Server Error`: Error loading document.

#### 3.4 Update Document
*   **Endpoint:** `/api/dbs/:db_name/collections/:collection_name/documents/:document_id`
*   **Method:** `PUT`
*   **Description:** Updates an existing document with the provided JSON data. The entire document is replaced.
*   **Path Parameters:**
    *   `db_name` (string, required): The name of the database.
    *   `collection_name` (string, required): The name of the collection.
    *   `document_id` (string, required): The ID of the document to update.
*   **Request Body:**
    *   Any valid JSON object representing the new content of the document.
    ```json
    {
      "newField": "newValue",
      "anotherField": 123
    }
    ```
*   **Success Response:**
    *   **Code:** `200 OK`
    ```json
    {
      "success": true,
      "message": "Document updated successfully",
      "data": {
        "newField": "newValue",
        "anotherField": 123
      }
    }
    ```
    *   `data` (object): The updated document.
*   **Error Responses:**
    *   `400 Bad Request`: Invalid JSON in request body.
    *   `500 Internal Server Error`: Failed to update document.

#### 3.5 Delete Document
*   **Endpoint:** `/api/dbs/:db_name/collections/:collection_name/documents/:document_id`
*   **Method:** `DELETE`
*   **Description:** Deletes a specific document from a collection.
*   **Path Parameters:**
    *   `db_name` (string, required): The name of the database.
    *   `collection_name` (string, required): The name of the collection.
    *   `document_id` (string, required): The ID of the document to delete.
*   **Success Response:**
    *   **Code:** `200 OK`
    ```json
    {
      "success": true,
      "message": "Document deleted successfully"
    }
    ```
*   **Error Responses:**
    *   `404 Not Found`: Document not found.
    *   `500 Internal Server Error`: Failed to delete document.

---

## Libraries

**deo** has official libraries to provide a cleaner developer experience, letting you access the database through clean code and not web requests.

[All official libraries are provided here.](https://github.com/myferr/deo/tree/main/packages)

* JavaScript / TypeScript
  * [deo-db](https://npmjs.com/package/deo-db)
