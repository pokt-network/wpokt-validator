import yaml from "js-yaml";
import fs from "fs";
import { ChainID } from "@pokt-foundation/pocketjs-transaction-builder";

const CONFIG_PATH = process.env.CONFIG_PATH || "./config.local.yml";

export type Config = {
  mongodb: {
    uri: string;
    database: string;
  };
  ethereum: {
    start_block_number: number;
    private_key: string;
    rpc_url: string;
    chain_id: string;
    rpc_timeout_ms: number;
    wrapped_pocket_address: string;
    mint_controller_address: string;
    validator_addresses: string[];
  };
  pocket: {
    start_height: number;
    private_key: string;
    rpc_url: string;
    chain_id: ChainID;
    rpc_timeout_ms: number;
    tx_fee: number;
    vault_address: string;
    multisig_public_keys: string[];
  };
};

export const config = yaml.load(fs.readFileSync(CONFIG_PATH, "utf8")) as Config;
