import { parseUnits } from "viem";
import ethereum from "../util/ethereum";
import pocket from "../util/pocket";
import { expect } from "chai";

const ZERO_ADDRESS = "0000000000000000000000000000000000000000";

type MintMemo = {
  address: string;
  amount: string;
};

type MessageSend = {
  type: string;
  value: {
    from_address: string;
    to_address: string;
    amount: string;
  };
};

export const mintFlow = async () => {
  it("should transfer pokt", async () => {
    const signer = await pocket.signerPromise;
    const fromAddress = await pocket.getAddress();
    const recipientAddress = await ethereum.getAddress();
    const toAddress = ZERO_ADDRESS;
    const amount = parseUnits("1", 6);
    const fee = BigInt(10000);

    const memo: MintMemo = {
      address: recipientAddress,
      amount: amount.toString(),
    };

    const fromBeforeBalance = await pocket.getBalance(fromAddress);
    const toBeforeBalance = await pocket.getBalance(toAddress);

    const tx = await pocket.sendPOKT(
      signer,
      toAddress,
      amount.toString(),
      JSON.stringify(memo),
      fee.toString()
    );

    expect(tx).to.not.be.null;

    if (!tx) return;

    expect(tx.hash).to.be.a("string");
    expect(tx.hash).to.have.lengthOf(64);
    expect(tx.tx_result.code).to.equal(0);
    expect(tx.stdTx.memo).to.equal(JSON.stringify(memo));
    expect(tx.stdTx.fee[0].amount).to.equal(fee.toString());

    const msg: MessageSend = tx.stdTx.msg as MessageSend;
    expect(msg.type).to.equal("pos/Send");
    expect(msg.value.from_address.toLowerCase()).to.equal(
      fromAddress.toLowerCase()
    );
    expect(msg.value.to_address.toLowerCase()).to.equal(
      toAddress.toLowerCase()
    );
    expect(msg.value.amount).to.equal(amount.toString());

    const fromAfterBalance = await pocket.getBalance(fromAddress);
    const toAfterBalance = await pocket.getBalance(toAddress);

    expect(fromAfterBalance).to.equal(fromBeforeBalance - amount - fee);
    expect(toAfterBalance).to.equal(toBeforeBalance + amount);
  });
};
