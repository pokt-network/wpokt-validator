import {
  Hex,
  TransactionReceipt,
  createPublicClient,
  createWalletClient,
  http,
  parseAbi,
} from "viem";
import { privateKeyToAccount } from "viem/accounts";
import { goerli, hardhat, mainnet } from "viem/chains";
import { config } from "./config";

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

const account = privateKeyToAccount(`0x${config.ethereum.private_key}`);

const walletClient = createWalletClient({
  account,
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
    abi: parseAbi(["function balanceOf(address) view returns (uint256)"]),
    functionName: "balanceOf",
    args: [address],
  });

  return balance;
};

const getAddress = async (): Promise<Hex> => {
  return account.address;
};

const sendWPOKT = async (
  recipient: Hex,
  amount: bigint
): Promise<TransactionReceipt> => {
  const hash = await walletClient.writeContract({
    address: config.ethereum.wrapped_pocket_address as Hex,
    abi: parseAbi([
      "function transfer(address _to, uint256 _value) public returns (bool success)",
    ]),
    functionName: "transfer",
    args: [recipient, amount],
  });
  const receipt = await publicClient.waitForTransactionReceipt({ hash });
  return receipt;
};

const mintWPOKT = async (
  data: { recipient: Hex; amount: bigint; nonce: bigint },
  signatures: Array<Hex>
): Promise<TransactionReceipt> => {
  const hash = await walletClient.writeContract({
    address: config.ethereum.mint_controller_address as Hex,
    abi: parseAbi([
      "function mintWrappedPocket(tuple(address recipient, uint256 amount, uint256 nonce), bytes[] signatures) public",
    ]),
    functionName: "mintWrappedPocket",
    args: [data, signatures],
  });

  const receipt = await publicClient.waitForTransactionReceipt({ hash });
  return receipt;
};

const burnAndBridgeWPOKT = async (
  amount: bigint,
  poktAddress: Hex
): Promise<TransactionReceipt> => {
  const hash = await walletClient.writeContract({
    address: config.ethereum.wrapped_pocket_address as Hex,
    abi: parseAbi([
      "function burnAndBridge(uint256 amount, address poktAddress) public",
    ]),
    functionName: "burnAndBridge",
    args: [amount, poktAddress],
  });

  const receipt = await publicClient.waitForTransactionReceipt({ hash });
  return receipt;
};

export default {
  getBalance,
  getWPOKTBalance,
  getAddress,
  sendWPOKT,
  mintWPOKT,
  burnAndBridgeWPOKT,
};
