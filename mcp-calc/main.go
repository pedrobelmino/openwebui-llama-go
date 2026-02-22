package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Cria o servidor MCP
	s := server.NewMCPServer(
		"Go-Calculator",
		"1.0.0",
		server.WithLogging(),
	)

	// Define a ferramenta de Calculadora
	calculatorTool := mcp.NewTool("calculate",
		mcp.WithDescription("Realiza operacoes matematicas basicas (soma, multiplicacao)"),
		mcp.WithString("operation", mcp.Required(), mcp.Description("A operacao: 'add' ou 'multiply'")),
		mcp.WithNumber("a", mcp.Required(), mcp.Description("Primeiro numero")),
		mcp.WithNumber("b", mcp.Required(), mcp.Description("Segundo numero")),
	)

	// Registra a ferramenta e a lógica
	s.AddTool(calculatorTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		log.Printf("Recebida chamada de ferramenta. Params: %+v", request.Params.Arguments)

		// Função auxiliar para converter argumentos numéricos com segurança
		getFloat := func(args map[string]interface{}, key string) (float64, error) {
			val, ok := args[key]
			if !ok {
				return 0, fmt.Errorf("argumento '%s' faltando", key)
			}
			switch v := val.(type) {
			case float64:
				return v, nil
			case int:
				return float64(v), nil
			case float32:
				return float64(v), nil
			default:
				return 0, fmt.Errorf("argumento '%s' tem tipo invalido: %T", key, val)
			}
		}

		opArg, ok := request.Params.Arguments["operation"]
		if !ok {
			return mcp.NewToolResultError("Argumento 'operation' faltando"), nil
		}
		op, ok := opArg.(string)
		if !ok {
			return mcp.NewToolResultError(fmt.Sprintf("Argumento 'operation' deve ser string, recebeu %T", opArg)), nil
		}

		a, err := getFloat(request.Params.Arguments, "a")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		b, err := getFloat(request.Params.Arguments, "b")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		log.Printf("Executando operacao: %s com %f e %f", op, a, b)

		var result float64
		switch op {
		case "add":
			result = a + b
		case "multiply":
			result = a * b
		default:
			log.Printf("Erro: Operacao desconhecida: %s", op)
			return mcp.NewToolResultError(fmt.Sprintf("Operacao desconhecida: %s", op)), nil
		}

		log.Printf("Resultado: %f", result)
		return mcp.NewToolResultText(fmt.Sprintf("%f", result)), nil
	})

	// Configura o servidor SSE
	baseURL := os.Getenv("MCP_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
	sseServer := server.NewSSEServer(s, baseURL)

	log.Println("Servidor MCP Calculator rodando na porta 8080...")
	if err := sseServer.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
