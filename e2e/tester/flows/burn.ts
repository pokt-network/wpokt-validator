import { parseUnits } from "viem";
import ethereum from "../util/ethereum";
import pocket from "../util/pocket";
import { expect } from "chai";
import { findBurn } from "../util/mongodb";
import { Status } from "../types";
import { sleep, debug } from "../util/helpers";

type MessageSend = {
  type: string;
  value: {
    from_address: string;
    to_address: string;
    amount: string;
  };
};

export const burnFlow = async () => {
  it("should burn and return amount from vault", async () => {
    debug("\nTesting -- should burn and return amount from vault");

    const fromAddress = await ethereum.getAddress();
    const recipientAddress = await pocket.getAddress();
    const amount = parseUnits("1", 6);
    const fee = BigInt(10000);

    const recipientBeforeBalance = await pocket.getBalance(recipientAddress);
    const fromBeforeBalance = await ethereum.getWPOKTBalance(fromAddress);

    debug("Sending transaction...");
    const burnTx = await ethereum.burnAndBridgeWPOKT(
      await ethereum.walletPromise,
      amount,
      recipientAddress
    );

    expect(burnTx).to.not.be.null;

    if (!burnTx) return;
    debug("Transaction sent: ", burnTx.transactionHash);

    const burnEvent = ethereum.findBurnAndBridgeEvent(burnTx);

    expect(burnEvent).to.not.be.null;

    if (!burnEvent) return;

    expect(burnEvent.amount).to.equal(amount);
    expect(burnEvent.poktAddress.toLowerCase()).to.equal(
      `0x${recipientAddress.toLowerCase()}`
    );
    expect(burnEvent.from.toLowerCase()).to.equal(fromAddress.toLowerCase());

    debug("Waiting for burn to be created...");
    await sleep(1000);

    let burn = await findBurn(burnTx.transactionHash);

    expect(burn).to.not.be.null;

    if (!burn) return;
    debug("Burn created");

    expect(burn.block_number.toString()).to.equal(
      burnTx.blockNumber.toString()
    );
    expect(burn.sender_address.toLowerCase()).to.equal(
      fromAddress.toLowerCase()
    );
    expect(burn.recipient_address.toLowerCase()).to.equal(
      recipientAddress.toLowerCase()
    );
    expect(burn.amount.toString()).to.equal(amount.toString());
    expect(burn.status).to.equal(Status.PENDING);

    await sleep(10000);

    burn = await findBurn(burnTx.transactionHash);

    expect(burn).to.not.be.null;

    if (!burn) return;

    expect(burn.status).to.be.oneOf([
      Status.CONFIRMED,
      Status.SIGNED,
      Status.SUMBITTED,
      Status.SUCCESS,
    ]);

    debug("Burn confirmed");

    await sleep(3000);

    burn = await findBurn(burnTx.transactionHash);

    expect(burn).to.not.be.null;

    if (!burn) return;

    expect(burn.status).to.be.oneOf([
      Status.SIGNED,
      Status.SUMBITTED,
      Status.SUCCESS,
    ]);
    debug("Burn signed");

    expect(burn.signers.length).to.equal(3);
    expect(burn.return_tx).to.not.be.null;

    await sleep(3000);

    burn = await findBurn(burnTx.transactionHash);

    expect(burn).to.not.be.null;

    if (!burn) return;

    expect(burn.status).to.equal(Status.SUCCESS);

    expect(burn.return_tx_hash).to.not.be.null;

    const returnTx = await pocket.getTransction(burn.return_tx_hash);

    expect(returnTx).to.not.be.null;

    if (!returnTx) return;

    expect(returnTx.tx_result.code).to.equal(0);

    const returnMsg: MessageSend = returnTx.stdTx.msg as MessageSend;

    expect(returnMsg.type).to.equal("pos/Send");
    expect(returnMsg.value.from_address.toLowerCase()).to.equal(
      pocket.VAULT_ADDRESS.toLowerCase()
    );
    expect(returnMsg.value.to_address.toLowerCase()).to.equal(
      recipientAddress.toLowerCase()
    );
    expect(returnMsg.value.amount).to.equal((amount - fee).toString());

    const recipientAfterBalance = await pocket.getBalance(recipientAddress);
    const fromAfterBalance = await ethereum.getWPOKTBalance(fromAddress);

    expect(recipientAfterBalance).to.equal(
      recipientBeforeBalance + amount - fee
    );
    expect(fromAfterBalance).to.equal(fromBeforeBalance - amount);

    debug("Burn success");
  });
};
