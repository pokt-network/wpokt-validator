# wPOKT Validator

The wPOKT Validator is a validator node that facilitates the bridging of POKT tokens from the POKT network to wPOKT on the Ethereum Mainnet. It achieves this by monitoring transactions to a specific vault address on the POKT network. Upon receiving transactions with a valid memo, it creates a signed transaction on the Ethereum Mainnet, which is then stored in the MongoDB database. Users can access the signed transaction data through the provided UI and submit it on the Ethereum Mainnet to mint wPOKT. Additionally, the process supports burning wPOKT tokens on the Mainnet.

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
    $ go run main.go config.testnet.yml
    ```

    c. Running in a Docker container:

    - Set the environment variables in the environment or use a file.
    - Execute the following command in the project directory with the `docker-compose.yml` file:

    ```bash
    $ docker-compose up
    ```
