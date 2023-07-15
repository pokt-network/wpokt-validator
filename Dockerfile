FROM golang:1.19

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY app ./app
COPY ethereum ./ethereum
COPY pocket ./pocket
COPY main.go ./
COPY config.testnet.yml ./
COPY scripts ./scripts

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /validator-service

# Install python3
RUN apt-get update && apt-get install -y python3 python3-pip python3-setuptools

ENV ETH_PRIVATE_KEY ${ETH_PRIVATE_KEY}
ENV ETH_RPC_URL ${ETH_RPC_URL}
ENV ETH_CHAIN_ID ${ETH_CHAIN_ID}

ENV POKT_PRIVATE_KEY ${POKT_PRIVATE_KEY}
ENV POKT_RPC_URL ${POKT_RPC_URL}
ENV POKT_CHAIN_ID ${POKT_CHAIN_ID}

ENV MONGODB_URI ${MONGODB_URI}
ENV MONGODB_DATABASE ${MONGODB_DATABASE}

# Setup Config file using python script
RUN python3 scripts/setup_config.py config.testnet.yml config.yml

# Run
CMD ["/validator-service", "config.yml"]
