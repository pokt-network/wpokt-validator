FROM golang:1.19 as base

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY app ./app
COPY ethereum ./ethereum
COPY pocket ./pocket
COPY models ./models
COPY main.go ./
COPY config.testnet.yml ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /validator-service

# Set environment variables
ENV ETH_PRIVATE_KEY ${ETH_PRIVATE_KEY}
ENV ETH_RPC_URL ${ETH_RPC_URL}
ENV ETH_CHAIN_ID ${ETH_CHAIN_ID}
ENV ETH_START_BLOCK_NUMBER ${ETH_START_BLOCK_NUMBER}

ENV POKT_PRIVATE_KEY ${POKT_PRIVATE_KEY}
ENV POKT_RPC_URL ${POKT_RPC_URL}
ENV POKT_CHAIN_ID ${POKT_CHAIN_ID}
ENV POKT_START_HEIGHT ${POKT_START_HEIGHT}

ENV MONGODB_URI ${MONGODB_URI}
ENV MONGODB_DATABASE ${MONGODB_DATABASE}

# Run
CMD ["/validator-service", "/app/config.testnet.yml"]
