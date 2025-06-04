-include .env

all : clean install test build

.PHONY: dev
dev : dev_one

.PHONY: dev_one
dev_one :; SIGNER_MNEMONIC="infant apart enroll relief kangaroo patch awesome wagon trap feature armor approve" go run . --yaml ./defaults/config.local.yml

.PHONY: dev_two
dev_two :; SIGNER_MNEMONIC="shy smile praise educate custom fashion gun enjoy zero powder garden second" go run . --yaml ./defaults/config.local.yml

.PHONY: dev_three
dev_three :; SIGNER_MNEMONIC="wink giant track dwarf visa feed visual drip play grant royal noise" go run . --yaml ./defaults/config.local.yml

.PHONY: beta
beta :; go run . --config ./config.beta.yml

.PHONY: clean
clean : clean_tmp_data
	go clean && go mod tidy

.PHONY: clean_tmp_data
clean_tmp_data :; if [ -d "/tmp/data" ]; then sudo rm -rf /tmp/data; fi

.PHONY: install
install :; go mod download && go mod verify

.PHONY: format
format :; go fmt ./...

.PHONY: lint
lint :; golangci-lint run

.PHONY: test
test :; go test -v ./...

.PHONY: test_coverage
test_coverage :; bash ./coverage.sh 

.PHONY: test_coverage_html
open_test_coverage :; bash ./coverage.sh && open ./coverage.html

.PHONY: build
build :; go build -o wpokt-validator .

.PHONY: docker_build
docker_build :; docker buildx build . -t dan13ram/wpokt-validator:v0.0.2 --file ./docker/Dockerfile

.PHONY: docker_build_local_gcp
docker_build_local_gcp :; docker buildx build . -t dan13ram/wpokt-validator-local-gcp:v0.0.2 --file ./docker/Dockerfile-gcp

.PHONY: docker_push
docker_push :; docker push dan13ram/wpokt-validator:v0.0.2

.PHONY: docker_dev
docker_dev : docker_one

.PHONY: docker_one
docker_one :; SIGNER_MNEMONIC="infant apart enroll relief kangaroo patch awesome wagon trap feature armor approve" YAML_FILE=/app/defaults/config.local.yml docker compose -f docker/docker-compose.yml up --force-recreate

.PHONY: docker_two
docker_two :; SIGNER_MNEMONIC="shy smile praise educate custom fashion gun enjoy zero powder garden second" YAML_FILE=/app/defaults/config.local.yml docker compose -f docker/docker-compose.yml up --force-recreate

.PHONY: docker_three
docker_three :; SIGNER_MNEMONIC="wink giant track dwarf visa feed visual drip play grant royal noise" YAML_FILE=/app/defaults/config.local.yml docker compose -f docker/docker-compose.yml up --force-recreate

.PHONY: localnet_up
localnet_up :; docker compose -f e2e/docker-compose.yml up --force-recreate

.PHONY: localnet_up_gcp
localnet_up_gcp :; GCP_KEY_ONE=${GCP_KEY_ONE} GCP_KEY_TWO=${GCP_KEY_TWO} GCP_KEY_THREE=${GCP_KEY_THREE} docker compose -f e2e/docker-compose-gcp.yml up --force-recreate

.PHONY: prompt_user
prompt_user :
	@echo "Are you sure? [y/N] " && read ans && [ $${ans:-N} = y ]

.PHONY: docker_wipe
docker_wipe : prompt_user ## [WARNING] Remove all the docker containers, images and volumes.
	docker ps -a -q | xargs -r -I {} docker stop {}
	docker ps -a -q | xargs -r -I {} docker rm {}
	docker images -q | xargs -r -I {} docker rmi {}
	docker volume ls -q | xargs -r -I {} docker volume rm {}

.PHONY: e2e_test
e2e_test :; cd e2e && yarn install && yarn test

.PHONY: e2e_test_gcp
e2e_test_gcp :; cd e2e && yarn install && CONFIG_PATH="../defaults/config.local.gcp.yml" yarn test

.PHONY: generate_keys
generate_keys :; go run scripts/generate_keys/main.go --mnemonic "${mnemonic}"

.PHONY: generate_multisig
generate_multisig :; go run scripts/generate_multisig/main.go --publickeys "${publickeys}" --threshold ${threshold}

.PHONY: gcp_kms
gcp_kms :; GOOGLE_APPLICATION_CREDENTIALS=${GOOGLE_APPLICATION_CREDENTIALS} GCP_KMS_KEY_NAME=${GCP_KMS_KEY_NAME} go run scripts/gcp_kms/main.go

.PHONY: gcp_dev_one
gcp_dev_one :; GOOGLE_APPLICATION_CREDENTIALS=${GOOGLE_APPLICATION_CREDENTIALS} SIGNER_GCP_KMS_KEY_NAME=${GCP_KEY_ONE} go run . --yaml ./defaults/config.local.gcp.yml

.PHONY: gcp_dev_two
gcp_dev_two :; GOOGLE_APPLICATION_CREDENTIALS=${GOOGLE_APPLICATION_CREDENTIALS} SIGNER_GCP_KMS_KEY_NAME=${GCP_KEY_TWO} go run . --yaml ./defaults/config.local.gcp.yml

.PHONY: gcp_dev_three
gcp_dev_three :; GOOGLE_APPLICATION_CREDENTIALS=${GOOGLE_APPLICATION_CREDENTIALS} SIGNER_GCP_KMS_KEY_NAME=${GCP_KEY_THREE} go run . --yaml ./defaults/config.local.gcp.yml

.PHONY: gcp_dev
gcp_dev : gcp_dev_one
