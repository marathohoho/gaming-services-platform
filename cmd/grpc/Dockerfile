# build app
FROM golang:1.17.8-alpine3.15 as base_build

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . .
RUN go build -v -o grpc-sc cmd/grpc/main.go

# create runtime
FROM alpine:3.15.0

WORKDIR /app

RUN apk add ca-certificates

COPY --from=base_build /app/grpc-sc .

EXPOSE 5103

ENTRYPOINT ["/app/grpc-sc"]