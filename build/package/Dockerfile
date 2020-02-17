FROM golang:1.13-alpine as build_base

# Install some dependencies needed to build the project
RUN apk add bash ca-certificates git gcc g++ libc-dev
WORKDIR $GOPATH/src/petstore-service


COPY go.mod .
COPY go.sum .

ENV GO111MODULE=on

RUN go mod download

FROM build_base AS server_builder

COPY . .

# Compile the project
RUN cd cmd/server && CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go install -a -tags netgo -ldflags '-w -extldflags "-static"'

# Use project from scratch
FROM alpine:3
RUN apk add --no-cache ca-certificates

COPY --from=server_builder /go/bin/server /app/cmd/server/
COPY --from=server_builder /go/src/petstore-service/config/config.yaml /app/config/
COPY --from=server_builder /go/src/petstore-service/cert/firestore/firestore_key.json /app/cert/firestore/
WORKDIR /app/cmd/server

CMD ./server