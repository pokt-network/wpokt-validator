import { formatUnits, parseUnits } from "viem";
import { databasePromise } from "./util/mongodb";
import pocket from "./util/pocket";
import ethereum from "./util/ethereum";

const main = async () => {
  const db = await databasePromise;

  const collection = db.collection("healthchecks");

  const result = await collection.find({}).toArray();

  console.log("Number of healthchecks: ", result.length);

  const pocketAddress = await pocket.getAddress();

  console.log("Pocket address: ", pocketAddress);

  const pocketBalance = await pocket.getBalance(pocketAddress);

  console.log("Pocket balance: ", formatUnits(pocketBalance, 6), " POKT");

  // const user = "00104055c00bed7c983a48aac7dc6335d7c607a7";

  // console.log(
  //   "User Balance: ",
  //   formatUnits(await getBalance(user), 6),
  //   " POKT"
  // );
  // const txRes = await sendPOKT(user, parseUnits("100", 6).toString(), "test");
  // console.log(txRes);
  
  const ethAddress = await ethereum.getAddress();

  console.log("Ethereum address: ", ethAddress);

  const ethBalance = await ethereum.getBalance(ethAddress);

  console.log("Ethereum balance: ", formatUnits(ethBalance, 18), " ETH");

  

};

main()
  .catch((err) => {
    console.error(err);
  })
  .finally(() => {
    process.exit(0);
  });
