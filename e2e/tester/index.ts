import { formatUnits } from "viem";
import { findHealthChecks } from "./util/mongodb";
import pocket from "./util/pocket";
import ethereum from "./util/ethereum";
import { mintFlow } from "./flows/mint";

const init = async () => {
  const healths = await findHealthChecks();

  console.log("Number of validators:", healths.length);

  const pocketAddress = await pocket.getAddress();

  console.log("Pocket address:", pocketAddress);

  console.log(
    "Pocket network:",
    pocket.CHAIN.network,
    "at",
    pocket.CHAIN.rpcUrl
  );

  const pocketBalance = await pocket.getBalance(pocketAddress);

  console.log("Pocket balance:", formatUnits(pocketBalance, 6), "POKT");

  const ethAddress = await ethereum.getAddress();

  console.log("Ethereum address:", ethAddress);

  console.log(
    "Ethereum network:",
    ethereum.CHAIN.network,
    "at",
    ethereum.CHAIN.rpcUrls.default.http[0]
  );

  const ethBalance = await ethereum.getBalance(ethAddress);

  console.log("Ethereum balance:", formatUnits(ethBalance, 18), "ETH");
};

before(async () => {
  await init();
  console.log("\n");
});

describe("E2E tests", async () => {
  describe("Mint flow", mintFlow);
});
