import { Db, MongoClient } from "mongodb";

import { config } from "./config";

const createDatabasePromise = async (): Promise<Db> => {
  const client = new MongoClient(config.mongodb.uri, {});

  await client.connect();
  return client.db(config.mongodb.database);
};

export const databasePromise: Promise<Db> = createDatabasePromise();
