# Deployment Guide for wPOKT Validators

This guide will walk you through the process of deploying wPOKT Validators on Google Cloud Platform (GCP).

## Steps to Deploy the Validators:

### Step 1: Create a Project on Google Cloud Platform (GCP)

1. Start by creating a project on Google Cloud Platform. This project will serve as the deployment environment for your validators.

2. Within the GCP project, enable the "Compute Engine" and "Secret Manager" services. "Compute Engine" will allow you to create and manage virtual machines, while "Secret Manager" will securely store your validator's private keys.

### Step 2: Create a MongoDB Project and Cluster

1. Create a project on MongoDB Cloud. This project will serve as the database for your validators.

2. Within the MongoDB project, create a MongoDB cluster with backups enabled. The cluster will provide the necessary storage and reliability for your validator's data.

### Step 3: Generate Validator Private Keys

1. Determine the number (N) of validators you will be deploying. For each validator, generate two sets of private keys: one for the Ethereum network and one for the Pocket network.

2. Make a note of the validator Ethereum addresses and Pocket public keys. You'll need this information later during the deployment process.

3. Generate a Pocket multisig address using the N Pocket public keys. This multisig address will also be required during the template creation.

4. Update the MintController Smart Contract on the Ethereum network with the Ethereum addresses of the N validators. The MintController Smart Contract will utilize these addresses to validate signatures from the deployed validators during the bridging process.

### Step 4: Store Secrets in Secret Manager

1. Add all the Ethereum and Pocket private keys to the Secret Manager on GCP. Ensure you securely store these keys as they are crucial for your validator's operation.

2. Also add the MongoDB URI with read-and-write permissions to the Secret Manager. This URI will be used to connect to the MongoDB cluster.

3. Note down the names of all the secrets created in Secret Manager. You will use these secret names during the deployment process.

4. Additionally, consider storing copies of the private keys in other secure places for additional redundancy and security. You might want to use hardware wallets, cold storage devices, or other secure offline storage methods to safeguard your validator's private keys.

### Step 5: Optional - Create Service Accounts and Separate Key Pairs

1. Optionally, you can create service accounts on GCP with access to specific keys. This can enhance security and limit access to certain resources.

2. If desired, you can also consider separating the Ethereum and Pocket key pairs into separate GCP projects to further isolate resources and permissions.

### Step 6: Create a VM Template on Compute Engine with Docker Image and Env Variables

1. Create a VM template on GCP's "Compute Engine" that includes the docker image for the wPOKT Validator and valid environment variables.

2. Set the default environment variables for:

    - Ethereum network configuration

    - Pocket network configuration

    - Google secret manager configuration

Refer to the sample `config.sample.yml` or `sample.env` files for reference on how to structure the environment variables.

3. In the "Advanced Options" under "Management", ensure that logging and monitoring are enabled for the VM template. This will help you monitor the performance and status of the deployed validators effectively.

4. Ensure that the service account assigned to this VM template has the necessary access to the secrets stored in Secret Manager. This access ensures that the validator instances can securely retrieve the required private keys for the Ethereum and Pocket networks.

### Step 7: Create and Start VM Instances

1. With the VM template prepared, you can now create and start N instances of the validators, where N is the number of validators you determined in Step 3.

2. For each validator, ensure you have performed the following updates:

    - Update the environment variables with the specific secret names for the Ethereum and Pocket private keys stored in Secret Manager.

    - If you chose to optionally create service accounts and separated the key pairs into separate GCP projects in Step 5, make sure to update the necessary service accounts and ensure the correct access permissions are granted to each validator instance.

3. Proceed to start the N instances of the validators with the valid environment variables and configurations.

Following these steps, you'll successfully deploy N wPOKT Validators on Google Cloud Platform. The validators will be ready to operate, securely bridging POKT tokens from the POKT network to wPOKT on the Ethereum Mainnet.
