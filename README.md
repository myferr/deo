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
docker-compose up -d
```
