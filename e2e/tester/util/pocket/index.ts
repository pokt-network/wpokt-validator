import { JsonRpcProvider } from "./provider";
import { KeyManager } from "@pokt-foundation/pocketjs-signer";
import { TransactionBuilder } from "@pokt-foundation/pocketjs-transaction-builder";
import { config } from "../config";
import { TransactionResponse } from "@pokt-foundation/pocketjs-types";

const provider = new JsonRpcProvider({
  rpcUrl: config.pocket.rpc_url,
});

const getBalance = async (address: string): Promise<bigint> => {
  const balance = await provider.getBalance(address);
  return balance;
};

const getAddress = async (): Promise<string> => {
  const signer = await signerPromise;
  return signer.getAddress();
};

const signerPromise = KeyManager.fromPrivateKey(config.pocket.private_key);

const sendPOKT = async (
  address: string,
  amount: string,
  memo: string
): Promise<TransactionResponse> => {
  const signer = await signerPromise;
  const transactionBuilder = new TransactionBuilder({
    provider,
    signer,
    chainID: config.pocket.chain_id,
  });

  const txMsg = transactionBuilder.send({
    fromAddress: signer.getAddress(),
    toAddress: address,
    amount,
  });

  const txresponse = await transactionBuilder.submit({
    memo,
    txMsg,
  });

  return txresponse;
};

export default {
  getBalance,
  getAddress,
  sendPOKT,
};
