# GCP MIG Runbook (wPOKT Validator)

This runbook documents how to run `wpokt-validator` on Google Compute Engine using a Managed Instance Group (MIG), replacing the deprecated Compute Engine container startup agent flow.

## Why this exists

Google deprecated the legacy container startup agent (`gce-container-declaration`) on VMs. This runbook uses:

- a normal VM instance template
- a metadata startup script
- Docker + `--restart=always`
- per-instance metadata overrides for validator-specific keys/config

## Current operating model

- One zonal MIG in `us-central1-a`
- Seven instances (fixed size)
- Shared base config in the instance template
- Per-validator secrets and RPC endpoints in MIG per-instance config

## Files used

- Startup script: `scripts/gcp/startup-validator-from-metadata.sh`

The script reads instance metadata keys, writes `/etc/wpokt-validator.env`, then runs:

- `docker pull <IMAGE>`
- `docker rm -f <APP_NAME> || true`
- `docker run -d --name <APP_NAME> --restart=always --env-file /etc/wpokt-validator.env <IMAGE>`

## Prerequisites

- Docker image already published (multi-arch recommended): `docker.io/<DOCKERHUB_OR_REGISTRY>/<IMAGE>:<TAG>`
- Service account with required access:
  - Secret Manager read access (`roles/secretmanager.secretAccessor`)
  - KMS permissions needed by signing flow
- GCP APIs enabled (Compute Engine, Secret Manager, KMS)

## 1) Create instance template

This creates a template without `gce-container-declaration` and injects shared metadata + startup script.

```bash
gcloud compute instance-templates create validator-mainnet-shannon-mig-template-v1 \
  --machine-type=n2-custom-2-4096 \
  --boot-disk-size=10GB \
  --boot-disk-type=pd-balanced \
  --image-family=cos-stable \
  --image-project=cos-cloud \
  --service-account=<SERVICE_ACCOUNT_EMAIL> \
  --scopes=https://www.googleapis.com/auth/cloud-platform \
  --metadata=google-logging-enabled=true,google-monitoring-enabled=true,block-project-ssh-keys=true,IMAGE=docker.io/<DOCKERHUB_OR_REGISTRY>/<IMAGE>:<TAG>,MONGODB_DATABASE=<MONGODB_DATABASE>,GOOGLE_SECRET_MANAGER_ENABLED=true,GOOGLE_MONGO_SECRET_NAME=projects/<PROJECT_NUMBER_OR_ID>/secrets/<MONGO_URI_SECRET>/versions/latest,POKT_START_HEIGHT=<POKT_START_HEIGHT>,ETH_START_BLOCK_NUMBER=<ETH_START_BLOCK_NUMBER>,LOG_LEVEL=info \
  --metadata-from-file=startup-script=scripts/gcp/startup-validator-from-metadata.sh
```

## 2) Create MIG

```bash
gcloud compute instance-groups managed create validator-mainnet-shannon-mig \
  --zone=us-central1-a \
  --template=validator-mainnet-shannon-mig-template-v1 \
  --size=0 \
  --base-instance-name=validator-mainnet-shannon-mig
```

## 3) Create named instances with per-instance metadata

Create each validator with unique signer + RPC config. Example for `-01`:

```bash
gcloud compute instance-groups managed create-instance validator-mainnet-shannon-mig \
  --zone=us-central1-a \
  --instance=validator-mainnet-shannon-mig-01 \
  --stateful-metadata=APP_NAME=validator-mainnet-shannon-mig-01,GOOGLE_ETH_SECRET_NAME=projects/<PROJECT_NUMBER_OR_ID>/secrets/<ETH_PRIVATE_KEY_SECRET>/versions/latest,POKT_GCP_KMS_KEY_NAME=projects/<PROJECT_ID>/locations/global/keyRings/<KEYRING>/cryptoKeys/<KEY>/cryptoKeyVersions/<VERSION>,ETH_RPC_URL=https://<ETH_RPC_ENDPOINT>,POKT_RPC_URL=https://<POKT_RPC_ENDPOINT>
```

Repeat for `-02` ... `-07` with each validator's own values.

## 4) Verify deployment

Check MIG state:

```bash
gcloud compute instance-groups managed list-instances validator-mainnet-shannon-mig --zone=us-central1-a
```

Check per-instance metadata overrides:

```bash
gcloud compute instance-groups managed instance-configs list validator-mainnet-shannon-mig --zone=us-central1-a
```

Check startup script progress:

```bash
gcloud compute instances get-serial-port-output validator-mainnet-shannon-mig-01 --zone=us-central1-a --port=1
```

Check container status/logs:

```bash
gcloud compute ssh validator-mainnet-shannon-mig-01 --zone=us-central1-a --command="sudo docker ps -a"
gcloud compute ssh validator-mainnet-shannon-mig-01 --zone=us-central1-a --command="sudo docker logs --tail=100 validator-mainnet-shannon-mig-01"
```

## 5) Updating one instance safely (common operation)

When only one validator needs a config change (for example RPC URL), update just that instance config and recreate only that instance.

```bash
gcloud compute instance-groups managed instance-configs update validator-mainnet-shannon-mig \
  --zone=us-central1-a \
  --instance=validator-mainnet-shannon-mig-01 \
  --stateful-metadata=ETH_RPC_URL=https://<ETH_RPC_ENDPOINT>

gcloud compute instance-groups managed recreate-instances validator-mainnet-shannon-mig \
  --zone=us-central1-a \
  --instances=validator-mainnet-shannon-mig-01
```

## Troubleshooting notes

- `Exceeded the quota usage` at startup:
  - Usually RPC endpoint limit/rate quota issue. Update `ETH_RPC_URL` for affected instance.

- Container repeatedly restarts:
  - `docker logs <APP_NAME>` for fatal errors.
  - Confirm metadata values are present in `instance-configs list`.

- SSH host key warning after instance recreate:
  - Expected when instance was recreated.
  - Remove stale key and reconnect:

```bash
ssh-keygen -R compute.<host-id> -f ~/.ssh/google_compute_known_hosts
```

- `curl: (22) ... 404` in early startup-script logs:
  - This can appear transiently during metadata reads. If script exits `0` and container starts, it is non-blocking.

## Cutover guidance

- Do not run old standalone validator VMs with the same keypairs at the same time as MIG validators.
- Keep exactly one active instance per signer/keypair.
- After validation, keep old `validator-mainnet-shannon-new-*` instances terminated/deleted.
