FROM golang:1.12 as builder
WORKDIR /go/src/github.com/containous/whoamitcp
COPY . .
RUN CGO_ENABLED=0 go build ./

# Create a minimal container to run a Golang static binary
FROM scratch
COPY --from=builder /go/src/github.com/containous/whoamitcp/whoamitcp .
ENTRYPOINT ["/whoamitcp"]
EXPOSE 8080
