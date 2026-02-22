FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY . .
# Baixa dependências e garante que o go.mod/go.sum estejam corretos com base no código fonte
RUN go mod tidy
RUN go build -o mcp-server main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/mcp-server .
EXPOSE 8080
CMD ["./mcp-server"]
