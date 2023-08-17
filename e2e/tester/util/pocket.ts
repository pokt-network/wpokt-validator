import { JsonRpcProvider } from "@pokt-foundation/pocketjs-provider";
import { AbstractSigner, KeyManager } from "@pokt-foundation/pocketjs-signer";
import { TransactionBuilder } from "@pokt-foundation/pocketjs-transaction-builder";
import { config } from "./config";
import { Transaction } from "@pokt-foundation/pocketjs-types";
import { sleep } from "./helpers";
import { parseUnits } from "viem";

const CHAIN = {
  network: config.pocket.chain_id,
  rpcUrl: config.pocket.rpc_url,
};

const VAULT_ADDRESS = config.pocket.vault_address;

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

const POLL_INTERVAL = 1000;

const pollForTransaction = async (
  txHash: string
): Promise<Transaction | null> => {
  let polls = 0;
  while (polls < 5) {
    try {
      const tx = await provider.getTransaction(txHash);
      return tx;
    } catch {
      // do nothing
    } finally {
      await sleep(POLL_INTERVAL);
    }
  }
  return null;
};

const sendPOKT = async (
  signer: AbstractSigner,
  recipient: string = VAULT_ADDRESS,
  amount: string,
  memo: string = "",
  fee: string = "10000"
): Promise<Transaction | null> => {
  const transactionBuilder = new TransactionBuilder({
    provider,
    signer,
    chainID: config.pocket.chain_id,
  });

  const txMsg = transactionBuilder.send({
    fromAddress: signer.getAddress(),
    toAddress: recipient,
    amount,
  });

  const { txHash } = await transactionBuilder.submit({
    memo,
    txMsg,
    fee,
  });

  const tx = await pollForTransaction(txHash);

  return tx;
};

const signerPromise = (async (): Promise<AbstractSigner> => {
  const defaultSigner = await KeyManager.fromPrivateKey(
    config.pocket.private_key
  );
  const newSigner = await KeyManager.createRandom();
  await sendPOKT(
    defaultSigner,
    newSigner.getAddress(),
    parseUnits("1000", 6).toString(),
    "init"
  );
  return newSigner;
})();

export default {
  CHAIN,
  VAULT_ADDRESS,
  getBalance,
  getAddress,
  sendPOKT,
  signerPromise,
};
