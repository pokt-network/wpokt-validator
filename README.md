# wPOKT Validator

The wPOKT Validator is a validator node that facilitates the bridging of POKT tokens from the POKT network to wPOKT on the Ethereum Mainnet. 

## How It Works

The wPOKT Validator comprises seven parallel services that enable the bridging of POKT tokens from the POKT network to wPOKT on the Ethereum Mainnet. Each service operates on an interval specified in the configuration. Here's an overview of their roles:

1. **Mint Monitor:**
   Monitors the Pocket network for transactions to the vault address. It validates transaction memos, inserting both valid `mint` and `invalid mint` transactions into the database.

2. **Mint Signer:**
   Handles pending and confirmed `mint` transactions. It signs confirmed transactions and updates the database accordingly.

3. **Mint Executor:**
   Monitors the Ethereum network for `mint` events and marks mints as successful in the database.

4. **Burn Monitor:**
   Monitors the Ethereum network for `burn` events and records them in the database.

5. **Burn Signer:**
   Handles pending and confirmed `burn` and `invalid mint` transactions. It signs the transactions and updates the status.

6. **Burn Executor:**
   Submits signed `burn` and `invalid mint` transactions to the Pocket network and updates the database upon success.

7. **Health:**
   Periodically reports the health status of the Golang service and sub-services to the database.

Through these services, the wPOKT Validator bridges POKT tokens to wPOKT, providing a secure and efficient validation process for the entire ecosystem.

## Installation

No specific installation steps are required. Users should have Golang installed locally and access to a MongoDB instance, either running locally or remotely, that they can attach to.

## Usage

1. **Configuration:**
   Users need a configuration file (template available as `config.testnet.yml`) containing the required fields. They can choose from the following options:

    a. Edit the `config.testnet.yml` file directly and add the necessary information.

    b. Set the configuration options as environment variables:

    - `ETH_PRIVATE_KEY`
    - `ETH_RPC_URL`
    - `ETH_CHAIN_ID`
    - `ETH_START_BLOCK_NUMBER`
    - `POKT_PRIVATE_KEY`
    - `POKT_RPC_URL`
    - `POKT_CHAIN_ID`
    - `POKT_START_HEIGHT`
    - `MONGODB_URI`
    - `MONGODB_DATABASE`

2. **Run the application:**
   Users can execute the program in the following ways:

    a. Using environment variables:

    ```bash
    $ ETH_PRIVATE_KEY="your_eth_private_key" ETH_RPC_URL="your_eth_rpc_url" ... go run main.go config.testnet.yml
    ```

    b. Using a `.env` file:

    - Create a `.env` file and add the environment variables in the format `VARIABLE_NAME=VALUE`.
    - Run the app:

    ```bash
    $ go run main.go config.testnet.yml .env
    ```

    c. Running in a Docker container:

    - Set the environment variables in the environment or use a file.
    - Execute the following command in the project directory with the `docker-compose.yml` file:

    ```bash
    $ docker-compose up
    ```

## Valid Memo

The validator node requires transactions on the POKT network to include a valid memo in the format of a JSON string. The memo should have the following structure:

```json
{ "address": "0xC9F2D9adfa6C24ce0D5a999F2BA3c6b06E36F75E", "chain_id": "5" }
```

-   `address`: The recipient address on the Ethereum network.
-   `chain_id`: The chain ID of the Ethereum network (represented as a string).

Transactions with memos not conforming to this format will not be processed by the validator.

## Docker Image

The wPOKT Validator is also available as a Docker image hosted on Docker Hub. You can run the validator in a Docker container using the following command:

```bash
docker run -d --env-file .env docker.io/dan13ram/wpokt-validator:latest
```

Ensure you have set the required environment variables in the `.env` file or directly in the command above.

## License

This project is licensed under the MIT License.
