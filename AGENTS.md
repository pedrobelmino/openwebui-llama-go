# AGENTS.md — openwebui-llama-go

Stack local de LLM com tool calling: Open WebUI como interface, Ollama como motor, um servidor MCP
em Go (`mcp-calc`) expondo `add`/`multiply` via SSE, e o MCPO como proxy que traduz MCP -> OpenAPI
para o Open WebUI enxergar a ferramenta. Tudo orquestrado por um único `docker-compose.yml`.

## Stack (lida de go.mod / Dockerfile / compose)

- Go 1.23 (`go 1.23` no `mcp-calc/go.mod`)
- `github.com/mark3labs/mcp-go` v0.8.2
- Build: `golang:1.23-alpine` -> `alpine:latest`
- `ollama/ollama:latest` (porta 11434)
- `ghcr.io/open-webui/open-webui:main` (8080)
- `ghcr.io/open-webui/mcpo:main` (8000, exposto em 8002)
- Nenhuma dependência de banco; estado em volumes Docker (`ollama_storage`, `open-webui_storage`)

## Comandos

- subir tudo: `docker compose up -d --build`
- derrubar: `docker compose down`  (com `-v` apaga modelos e dados do WebUI)
- baixar o modelo: `docker exec -it ollama ollama run llama3`  (sai com `/bye`)
- rebuild só da ferramenta: `docker compose up -d --build mcp-calc`
- build do binário: `cd mcp-calc && go build -o mcp-server main.go`
- test: `go test ./...`  (não existe nenhum teste no projeto)
- lint: não existe no projeto (`gofmt`/`go vet` disponíveis, sem config)

## Estrutura

- `docker-compose.yml` — 4 serviços: `ollama`, `open-webui`, `mcp-calc`, `mcpo`
- `mcp-calc/main.go` — servidor MCP + ferramenta `calculate`
- `mcp-calc/Dockerfile` — multi-stage; roda `go mod tidy` + `go build` no build
- `mcp-calc/go.mod` — módulo `mcp-calc`
- `prints/` — screenshots do passo a passo (não é código)

## Convenções observadas

- Um único `main.go`: ferramenta, handler e servidor no mesmo arquivo. Não dividir sem motivo.
- Nome da ferramenta e dos argumentos são contrato visível ao LLM: `calculate`,
  `operation` (`add`|`multiply`), `a`, `b`.
- Erro de execução volta como `mcp.NewToolResultError(...)`, nunca `panic`/exit.
- Conversão numérica defensiva em `getFloat` (aceita float64, int, float32; senão devolve erro tipado).
- Logs com `log.Printf` em português, incluindo os argumentos recebidos.

## Regras

- Segredo/credencial só por variável de ambiente. O `.gitignore` já cobre `.env` e `.env.*`
  (`.env.example` é permitido) — não comitar `.env`.
- Não alterar o contrato da tool (nome, `operation`, `a`, `b`, descrição) sem reimportar o
  `openapi.json` no Open WebUI: o LLM vê a descrição antiga e chama errado.
- `MCP_BASE_URL` é o que define a URL base anunciada pelo SSE; mantê-la coerente com o nome do
  serviço no compose (`http://mcp-calc:8080`).

## Pegadinhas reais

- Portas fixas e concorridas: `8080` (Open WebUI), `8081` (mcp-calc), `8002` (mcpo), `11434` (Ollama).
  O `8080` colide com qualquer Spring Boot local — pare um antes de subir o outro.
- Ordem é manual e importa: subir os containers, puxar o `llama3` no Ollama, e só então criar a Tool
  no Open WebUI com o JSON de `http://localhost:8002/openapi.json`. Dentro do JSON, a URL precisa
  apontar para `http://mcpo:8000` (nome do serviço), não para `localhost`.
- `open-webui` declara `depends_on` para `ollama` e `mcp-calc`, mas não para `mcpo`: o WebUI pode
  subir antes do proxy ficar pronto.
- `mcpo` precisa de `--server-type sse` e da URL `http://mcp-calc:8080/sse`; sem isso o proxy não
  encontra a ferramenta e o `openapi.json` sai vazio.
- Sem GPU por padrão (bloco `deploy` comentado no compose). Modelo maior que o Llama 3 em CPU estoura
  timeout no tool calling.
