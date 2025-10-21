import { DeoClient } from "./index";

describe("DeoClient", () => {
  const client = new DeoClient("http://localhost:6741");
  const dbName = "test_db";
  const collectionName = "test_collection";

  beforeAll(async () => {
    // Ensure the database and collection are clean before tests
    await client.deleteDatabase(dbName).catch(() => {});
    await client.createDatabase(dbName);
    await client.dbs[dbName].createCollection(collectionName);
  });

  afterAll(async () => {
    // Clean up the database after all tests
    await client.deleteDatabase(dbName);
  });

  it("should create and list databases", async () => {
    const dbs = await client.listDatabases();
    expect(dbs.success).toBe(true);
    expect(dbs.data).toContain(dbName);
  });

  it("should create and list collections", async () => {
    const collections = await client.dbs[dbName].listCollections();
    expect(collections.success).toBe(true);
    expect(collections.data).toContain(collectionName);
  });

  it("should create, read, update, and delete a document", async () => {
    const doc = { name: "Test Document", value: 123 };
    const createdDocResponse = await client.dbs[dbName].collections[collectionName].createDocument(doc);
    expect(createdDocResponse.success).toBe(true);
    expect(createdDocResponse.data).toHaveProperty("_id");
    expect(createdDocResponse.data?.name).toBe("Test Document");

    const documentId = createdDocResponse.data!._id;

    const readDocResponse = await client.dbs[dbName].collections[collectionName].readDocument(documentId);
    expect(readDocResponse.success).toBe(true);
    expect(readDocResponse.data?._id).toBe(documentId);
    expect(readDocResponse.data?.name).toBe("Test Document");

    const updatedDoc = { name: "Updated Document", value: 456 };
    const updatedDocResponse = await client.dbs[dbName].collections[collectionName].updateDocument(documentId, updatedDoc);
    expect(updatedDocResponse.success).toBe(true);
    expect(updatedDocResponse.data?.name).toBe("Updated Document");

    const deletedDocResponse = await client.dbs[dbName].collections[collectionName].deleteDocument(documentId);
    expect(deletedDocResponse.success).toBe(true);

    const listDocsResponse = await client.dbs[dbName].collections[collectionName].listDocuments();
    expect(listDocsResponse.success).toBe(true);
    expect(listDocsResponse.data?.some((d) => d._id === documentId)).toBeFalsy();
  });

  it("should list documents with query parameters", async () => {
    const queryCollectionName = "query_test_collection";
    const collection = client.dbs[dbName].collections[queryCollectionName];
    await client.dbs[dbName].createCollection(queryCollectionName);

    await collection.createDocument({ name: "C doc", value: 30, group: "X" });
    await collection.createDocument({ name: "A doc", value: 10, group: "Y" });
    await collection.createDocument({ name: "B doc", value: 20, group: "X" });

    // Test filtering
    const filteredDocs = await collection.listDocuments({ filters: { group: "X" } });
    expect(filteredDocs.success).toBe(true);
    expect(filteredDocs.data).toHaveLength(2);
    expect(filteredDocs.data?.every(d => d.group === "X")).toBe(true);

    // Test sorting
    const sortedDocs = await collection.listDocuments({ sortBy: "value", order: "asc" });
    expect(sortedDocs.success).toBe(true);
    expect(sortedDocs.data).toHaveLength(3);
    expect(sortedDocs.data?.map(d => d.value)).toEqual([10, 20, 30]);

    // Test pagination
    const paginatedDocs = await collection.listDocuments({ sortBy: "value", order: "asc", limit: 1, offset: 1 });
    expect(paginatedDocs.success).toBe(true);
    expect(paginatedDocs.data).toHaveLength(1);
    expect(paginatedDocs.data![0].value).toBe(20);

    // Test combined query
    const combined = await collection.listDocuments({ filters: { group: "X" }, sortBy: "value", order: "desc", limit: 1 });
    expect(combined.success).toBe(true);
    expect(combined.data).toHaveLength(1);
    expect(combined.data![0].name).toBe("C doc");

    // Cleanup
    await client.dbs[dbName].deleteCollection(queryCollectionName);
  });
});