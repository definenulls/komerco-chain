FROM golang:1.18-alpine AS builder

WORKDIR /komerco-chain

ADD go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o komerco-chain main.go

FROM alpine:latest AS runner

RUN apk --no-cache add ca-certificates jq

WORKDIR /komerco-chain

COPY --from=builder /komerco-chain/komerco-chain ./
COPY ./docker/local/komerco-chain.sh ./
COPY ./core-contracts/artifacts ./core-contracts/artifacts

# Expose json-rpc, libp2p and grpc ports
EXPOSE 8545 9632 1478 5001

ENTRYPOINT ["./komerco-chain.sh"]
