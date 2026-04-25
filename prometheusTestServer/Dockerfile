FROM golang:alpine AS builder

WORKDIR /app
# We don't have go.mod in prometheusTestServer, wait, does it use the root go.mod?
# The root go.mod has the prometheus client. But we only need to build main.go.
# We can just copy the whole project or just the test server.
# Let's copy the whole root so we have go.mod, then build inside prometheusTestServer.

COPY . .
WORKDIR /app/prometheusTestServer
RUN CGO_ENABLED=0 GOOS=linux go build -o testserver main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/prometheusTestServer/testserver .
EXPOSE 2112
CMD ["./testserver"]
