import { JsonRpcProvider } from "./pocket/provider";
import { KeyManager } from "@pokt-foundation/pocketjs-signer";
import { TransactionBuilder } from "@pokt-foundation/pocketjs-transaction-builder";
import { config } from "./config";
import { TransactionResponse } from "@pokt-foundation/pocketjs-types";

// Instantiate a provider for querying information on the chain!
export const provider = new JsonRpcProvider({
  rpcUrl: config.pocket.rpc_url,
});

export const getBalance = async (address: string): Promise<bigint> => {
  const balance = await provider.getBalance(address);
  return balance;
};

export const getAddress = async (): Promise<string> => {
  const signer = await signerPromise;
  return signer.getAddress();
};

export const signerPromise = KeyManager.fromPrivateKey(
  config.pocket.private_key
);

export const sendPOKT = async (
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
