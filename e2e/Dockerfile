FROM golang:1.19 as base

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# copy the source code
COPY app ./app
COPY eth ./eth
COPY pokt ./pokt
COPY models ./models
COPY e2e ./e2e
COPY main.go ./
COPY defaults.yml ./

# build
RUN CGO_ENABLED=0 GOOS=linux go build -o /validator

ENV LOG_LEVEL debug

# run
CMD ["/validator", "--config", "/app/defaults.yml"]
