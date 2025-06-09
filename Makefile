-include .env

all : clean install test build

export GOOGLE_APPLICATION_CREDENTIALS=./credentials.json
export LOG_LEVEL=debug

.PHONY: beta-one
beta-one :; go run . --config ./config/config.beta.one.yml

.PHONY: beta-two
beta-two :; go run . --config ./config/config.beta.two.yml

.PHONY: beta-three
beta-three :; go run . --config ./config/config.beta.three.yml

.PHONY: main
main :; go run . --config ./config/defaults.mainnet.yml

.PHONY: clean
clean : go clean && go mod tidy

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

.PHONY: open_test_coverage
open_test_coverage :; bash ./coverage.sh && open ./coverage.html

.PHONY: build
build :; go build -o wpokt-validator .

.PHONY: docker_build_beta
docker_build_beta :; docker buildx build . -t dan13ram/wpokt-validator-beta:v0.2.4 --file ./docker/Dockerfile.beta

.PHONY: docker_push_beta
docker_push_beta :; docker push dan13ram/wpokt-validator-beta:v0.2.4

.PHONY: docker_build
docker_build :; docker buildx build . -t dan13ram/wpokt-validator:v0.2.4 --file ./docker/Dockerfile.mainnet

.PHONY: docker_push
docker_push :; docker push dan13ram/wpokt-validator:v0.2.4

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
