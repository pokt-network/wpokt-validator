import {
  Account,
  Hex,
  TransactionReceipt,
  Transport,
  WalletClient,
  createPublicClient,
  createWalletClient,
  decodeEventLog,
  encodeEventTopics,
  http,
  parseAbi,
  parseUnits,
} from "viem";
import { generatePrivateKey, privateKeyToAccount } from "viem/accounts";
import { Chain, goerli, hardhat, mainnet } from "viem/chains";
import { config } from "./config";
import { MintData } from "../types";
import {MINT_CONTROLLER_ABI, WRAPPED_POCKET_ABI} from "./abis";

const ETH_CHAIN = (() => {
  switch (config.ethereum.chain_id) {
    case "1":
      return mainnet;
    case "5":
      return goerli;
    case "31337":
    default:
      return hardhat;
  }
})();

const defaultWalletClient: WalletClient<Transport, Chain, Account> =
  createWalletClient({
    account: privateKeyToAccount(`0x${config.ethereum.private_key}`),
    chain: ETH_CHAIN,
    transport: http(),
  });

const publicClient = createPublicClient({
  chain: ETH_CHAIN,
  transport: http(),
});

const getBalance = async (address: Hex): Promise<bigint> => {
  const balance = await publicClient.getBalance({
    address,
  });
  return balance;
};

const getWPOKTBalance = async (address: Hex): Promise<bigint> => {
  const balance = await publicClient.readContract({
    address: config.ethereum.wrapped_pocket_address as Hex,
    abi: WRAPPED_POCKET_ABI,
    functionName: "balanceOf",
    args: [address],
  });

  return balance as bigint;
};

const sendETH = async (
  wallet: WalletClient<Transport, Chain, Account>,
  recipient: Hex,
  amount: bigint
): Promise<TransactionReceipt> => {
  const hash = await wallet.sendTransaction({
    to: recipient,
    value: amount,
  });
  const receipt = await publicClient.waitForTransactionReceipt({ hash });
  return receipt;
};

const sendWPOKT = async (
  wallet: WalletClient<Transport, Chain, Account>,
  recipient: Hex,
  amount: bigint
): Promise<TransactionReceipt> => {
  const hash = await wallet.writeContract({
    address: config.ethereum.wrapped_pocket_address as Hex,
    abi: WRAPPED_POCKET_ABI,
    functionName: "transfer",
    args: [recipient, amount],
  });
  const receipt = await publicClient.waitForTransactionReceipt({ hash });
  return receipt;
};

const walletPromise: Promise<WalletClient<Transport, Chain, Account>> =
  (async () => {
    const pKey = generatePrivateKey();
    const walletClient = createWalletClient({
      account: privateKeyToAccount(pKey),
      chain: ETH_CHAIN,
      transport: http(),
    });

    await sendETH(
      defaultWalletClient,
      walletClient.account.address,
      parseUnits("1000", 18)
    );

    return walletClient;
  })();

const getAddress = async (): Promise<Hex> => {
  const wallet = await walletPromise;
  return wallet.account.address;
};

const mintWPOKT = async (
  wallet: WalletClient<Transport, Chain, Account>,
  data: MintData,
  signatures: string[]
): Promise<TransactionReceipt> => {
  const hash = await wallet.writeContract({
    address: config.ethereum.mint_controller_address as Hex,
    abi: MINT_CONTROLLER_ABI,
    functionName: "mintWrappedPocket",
    args: [data, signatures],
  });

  const receipt = await publicClient.waitForTransactionReceipt({ hash });
  return receipt;
};

type MintedEvent = {
  recipient: Hex;
  amount: bigint;
  nonce: bigint;
};

const findMintedEvent = (receipt: TransactionReceipt): MintedEvent | null => {
  const eventTops = encodeEventTopics({
    abi: WRAPPED_POCKET_ABI,
    eventName: "Minted",
  });

  const event = receipt.logs.find((log) => log.topics[0] === eventTops[0]);

  if (!event) {
    return null;
  }

  const decodedLog = decodeEventLog({
    abi: WRAPPED_POCKET_ABI,
    eventName: "Minted",
    data: event.data,
    topics: event.topics,
  });

  return decodedLog.args as MintedEvent;
};

const burnAndBridgeWPOKT = async (
  wallet: WalletClient<Transport, Chain, Account>,
  amount: bigint,
  poktAddress: Hex
): Promise<TransactionReceipt> => {
  const hash = await wallet.writeContract({
    address: config.ethereum.wrapped_pocket_address as Hex,
    abi: WRAPPED_POCKET_ABI,
    functionName: "burnAndBridge",
    args: [amount, poktAddress],
  });

  const receipt = await publicClient.waitForTransactionReceipt({ hash });
  return receipt;
};

export default {
  walletPromise,
  CHAIN: ETH_CHAIN,
  getBalance,
  getWPOKTBalance,
  getAddress,
  sendWPOKT,
  mintWPOKT,
  findMintedEvent,
  burnAndBridgeWPOKT,
};
