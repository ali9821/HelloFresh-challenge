FROM golang:1.18-alpine AS build_base

RUN apk add --no-cache git

WORKDIR /tmp/server

ADD . .

RUN go env -w GO111MODULE=on
RUN go mod download

# Unit test
#RUN CGO_ENABLED=0 go test -v

# Build the Go app
RUN go build -o ./out/server cmd/server/main.go

#start runner
FROM alpine:3.9
RUN apk add ca-certificates
WORKDIR /app/server
COPY --from=build_base /tmp/server/out/server .

COPY --from=build_base /tmp/server/config.yml .

CMD ["./server"]