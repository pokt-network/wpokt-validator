import { Db, MongoClient } from "mongodb";

import { config } from "./config";
import {
  CollectionHealthChecks,
  CollectionMints,
  Health,
  Mint,
} from "../types";

const createDatabasePromise = async (): Promise<Db> => {
  const client = new MongoClient(config.mongodb.uri, {});

  await client.connect();
  return client.db(config.mongodb.database);
};

export const databasePromise: Promise<Db> = createDatabasePromise();

export const findHealthChecks = async (): Promise<Health[]> => {
  const db = await databasePromise;
  return db
    .collection(CollectionHealthChecks)
    .find({
      updated_at: {
        $gte: new Date(Date.now() - 10000),
      },
    })
    .toArray() as Promise<Health[]>;
};

export const findMint = async (txHash: string): Promise<Mint | null> => {
  const db = await databasePromise;
  return db
    .collection(CollectionMints)
    .findOne({
      transaction_hash: txHash.toLowerCase(),
    }) as Promise<Mint | null>;
};
