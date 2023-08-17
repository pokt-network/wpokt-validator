import { parseUnits } from "viem";
import ethereum from "../util/ethereum";
import pocket from "../util/pocket";
import { expect } from "chai";
import { findInvalidMint, findMint } from "../util/mongodb";
import { MintMemo, Status } from "../types";
import { sleep, debug } from "../util/helpers";

type MessageSend = {
  type: string;
  value: {
    from_address: string;
    to_address: string;
    amount: string;
  };
};

export const mintFlow = async () => {
  it("should mint for send tx to vault with valid memo", async () => {
    const signer = await pocket.signerPromise;
    const fromAddress = await pocket.getAddress();
    const recipientAddress = await ethereum.getAddress();
    const toAddress = pocket.VAULT_ADDRESS;
    const amount = parseUnits("1", 6);
    const fee = BigInt(10000);

    const memo: MintMemo = {
      address: recipientAddress,
      chain_id: ethereum.CHAIN.id.toString(),
    };

    const fromBeforeBalance = await pocket.getBalance(fromAddress);
    const toBeforeBalance = await pocket.getBalance(toAddress);

    debug("Sending transaction...");
    const sendTx = await pocket.sendPOKT(
      signer,
      toAddress,
      amount.toString(),
      JSON.stringify(memo),
      fee.toString()
    );

    expect(sendTx).to.not.be.null;

    if (!sendTx) return;
    debug("Transaction sent: ", sendTx.hash);

    expect(sendTx.hash).to.be.a("string");
    expect(sendTx.hash).to.have.lengthOf(64);
    expect(sendTx.tx_result.code).to.equal(0);
    expect(sendTx.stdTx.memo).to.equal(JSON.stringify(memo));
    expect(sendTx.stdTx.fee[0].amount).to.equal(fee.toString());

    const msg: MessageSend = sendTx.stdTx.msg as MessageSend;
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

    debug("Waiting for mint to be created...");
    await sleep(2000);

    let mint = await findMint(sendTx.hash);

    expect(mint).to.not.be.null;

    if (!mint) return;
    debug("Mint created");

    expect(mint.height.toString()).to.equal(sendTx.height.toString());
    expect(mint.sender_address.toLowerCase()).to.equal(
      fromAddress.toLowerCase()
    );
    expect(mint.recipient_address.toLowerCase()).to.equal(
      recipientAddress.toLowerCase()
    );
    expect(mint.amount.toString()).to.equal(amount.toString());
    expect(mint.status).to.equal(Status.PENDING);

    await sleep(12000);

    mint = await findMint(sendTx.hash);

    expect(mint).to.not.be.null;

    if (!mint) return;

    expect(mint.status).to.be.oneOf([Status.CONFIRMED, Status.SIGNED]);
    debug("Mint confirmed");

    await sleep(3000);

    mint = await findMint(sendTx.hash);

    expect(mint).to.not.be.null;

    if (!mint) return;

    expect(mint.status).to.equal(Status.SIGNED);
    debug("Mint signed");

    expect(mint.signers.length).to.equal(3);
    expect(mint.signatures.length).to.equal(3);
    expect(mint.nonce.toString()).to.equal("1");
    expect(mint.data).to.not.be.null;

    if (!mint.data) return;

    const beforeWPOKTBalance = await ethereum.getWPOKTBalance(recipientAddress);

    debug("Minting WPOKT...");
    const mintTx = await ethereum.mintWPOKT(
      await ethereum.walletPromise,
      mint.data,
      mint.signatures
    );

    expect(mintTx).to.not.be.null;

    if (!mintTx) return;
    debug("WPOKT minted: ", mintTx.transactionHash);

    expect(mintTx.transactionHash).to.be.a("string");

    const mintedEvent = ethereum.findMintedEvent(mintTx);

    expect(mintedEvent).to.not.be.null;

    if (!mintedEvent) return;

    expect(mintedEvent.nonce.toString()).to.equal(mint.nonce.toString());
    expect(mintedEvent.recipient.toLowerCase()).to.equal(
      recipientAddress.toLowerCase()
    );
    expect(mintedEvent.amount.toString()).to.equal(amount.toString());

    const afterWPOKTBalance = await ethereum.getWPOKTBalance(recipientAddress);

    expect(afterWPOKTBalance).to.equal(beforeWPOKTBalance + amount);

    await sleep(2000);

    mint = await findMint(sendTx.hash);

    expect(mint).to.not.be.null;

    if (!mint) return;

    expect(mint.status).to.equal(Status.SUCCESS);

    expect(mint.mint_tx_hash.toLowerCase()).to.equal(
      mintTx.transactionHash.toLowerCase()
    );
    debug("Mint success");
  });

  it("should fail mint to invalidMint for failed send tx to vault", async () => {
    const signer = await pocket.signerPromise;
    const fromAddress = await pocket.getAddress();
    const toAddress = pocket.VAULT_ADDRESS;
    const amount = parseUnits("1000000000", 6);

    debug("Sending transaction...");
    const sendTx = await pocket.sendPOKT(signer, toAddress, amount.toString());

    expect(sendTx).to.not.be.null;

    if (!sendTx) return;
    debug("Transaction sent: ", sendTx.hash);

    expect(sendTx.hash).to.be.a("string");
    expect(sendTx.hash).to.have.lengthOf(64);
    expect(sendTx.tx_result.code).to.equal(10);

    debug("Waiting for invalidMint to be created...");
    await sleep(2000);

    let invalidMint = await findInvalidMint(sendTx.hash);

    expect(invalidMint).to.not.be.null;

    if (!invalidMint) return;
    debug("InvalidMint created");

    expect(invalidMint.height.toString()).to.equal(sendTx.height.toString());
    expect(invalidMint.sender_address.toLowerCase()).to.equal(
      fromAddress.toLowerCase()
    );
    expect(invalidMint.amount.toString()).to.equal(amount.toString());
    expect(invalidMint.status).to.equal(Status.FAILED);

    debug("InvalidMint failed");
  });

  it("should return amount for send tx to vault with invalid memo", async () => {
    const signer = await pocket.signerPromise;
    const fromAddress = await pocket.getAddress();
    const toAddress = pocket.VAULT_ADDRESS;
    const amount = parseUnits("1", 6);
    const fee = BigInt(10000);

    const beforeBalance = await pocket.getBalance(fromAddress);

    const memo = "not a json";

    debug("Sending transaction...");
    const sendTx = await pocket.sendPOKT(
      signer,
      toAddress,
      amount.toString(),
      memo,
      fee.toString()
    );

    expect(sendTx).to.not.be.null;

    if (!sendTx) return;
    debug("Transaction sent: ", sendTx.hash);

    expect(sendTx.hash).to.be.a("string");
    expect(sendTx.hash).to.have.lengthOf(64);
    expect(sendTx.tx_result.code).to.equal(0);
    expect(sendTx.stdTx.memo.toString()).to.equal(memo.toString());
    expect(sendTx.stdTx.fee[0].amount).to.equal(fee.toString());

    const msg: MessageSend = sendTx.stdTx.msg as MessageSend;
    expect(msg.type).to.equal("pos/Send");
    expect(msg.value.from_address.toLowerCase()).to.equal(
      fromAddress.toLowerCase()
    );
    expect(msg.value.to_address.toLowerCase()).to.equal(
      toAddress.toLowerCase()
    );
    expect(msg.value.amount).to.equal(amount.toString());

    debug("Waiting for invalidMint to be created...");
    await sleep(2000);

    let invalidMint = await findInvalidMint(sendTx.hash);

    expect(invalidMint).to.not.be.null;

    if (!invalidMint) return;
    debug("InvalidMint created");

    expect(invalidMint.height.toString()).to.equal(sendTx.height.toString());
    expect(invalidMint.sender_address.toLowerCase()).to.equal(
      fromAddress.toLowerCase()
    );
    expect(invalidMint.amount.toString()).to.equal(amount.toString());
    expect(invalidMint.status).to.equal(Status.PENDING);

    await sleep(12000);

    invalidMint = await findInvalidMint(sendTx.hash);

    expect(invalidMint).to.not.be.null;

    if (!invalidMint) return;

    expect(invalidMint.status).to.be.oneOf([
      Status.CONFIRMED,
      Status.SIGNED,
      Status.SUMBITTED,
      Status.SUCCESS,
    ]);
    debug("InvalidMint confirmed");

    await sleep(3000);

    invalidMint = await findInvalidMint(sendTx.hash);

    expect(invalidMint).to.not.be.null;

    if (!invalidMint) return;

    expect(invalidMint.status).to.be.oneOf([
      Status.SIGNED,
      Status.SUMBITTED,
      Status.SUCCESS,
    ]);
    debug("InvalidMint signed");

    expect(invalidMint.signers.length).to.equal(3);
    expect(invalidMint.return_tx).to.not.be.null;

    await sleep(3000);

    invalidMint = await findInvalidMint(sendTx.hash);

    expect(invalidMint).to.not.be.null;

    if (!invalidMint) return;

    expect(invalidMint.status).to.equal(Status.SUCCESS);

    expect(invalidMint.return_tx_hash).to.not.be.null;

    const returnTx = await pocket.getTransction(invalidMint.return_tx_hash);

    expect(returnTx).to.not.be.null;

    if (!returnTx) return;

    expect(returnTx.tx_result.code).to.equal(0);

    const returnMsg: MessageSend = returnTx.stdTx.msg as MessageSend;

    expect(returnMsg.type).to.equal("pos/Send");
    expect(returnMsg.value.from_address.toLowerCase()).to.equal(
      pocket.VAULT_ADDRESS.toLowerCase()
    );
    expect(returnMsg.value.to_address.toLowerCase()).to.equal(
      fromAddress.toLowerCase()
    );
    expect(returnMsg.value.amount).to.equal((amount - fee).toString());

    const afterBalance = await pocket.getBalance(fromAddress);

    expect(afterBalance).to.equal(beforeBalance - BigInt(2) * fee);

    debug("InvalidMint success");
  });
};
