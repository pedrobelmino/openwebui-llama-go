# Openweb UI LLama GO(mcpserver)

Este projeto demonstra a integração entre **LLMs (Large Language Models)** rodando localmente e **ferramentas externas** (Tools) utilizando o protocolo **MCP (Model Context Protocol)**.

Através da interface amigável do **Open WebUI**, o usuário pode interagir com o modelo **Llama 3** (via Ollama) usando linguagem natural. O modelo é capaz de identificar quando precisa realizar cálculos matemáticos e, transparentemente, aciona uma ferramenta de calculadora implementada em **Go**, rodando em um container Docker separado.

## 🚀 Funcionalidades

*   **Chat com IA Local:** Execução do modelo Llama 3 totalmente offline e privado.
*   **Chamada de Ferramentas (Tool Calling):** O modelo entende comandos como "quanto é 50 vezes 30?" e delega a execução para o código Go.
*   **Protocolo MCP:** Utilização do padrão *Model Context Protocol* para padronizar a comunicação entre a IA e as ferramentas.
*   **Interface Web Moderna:** Uso do Open WebUI para gerenciar chats, modelos e ferramentas.

## 🏗️ Arquitetura

O projeto é composto por quatro serviços principais orquestrados via Docker Compose:

1.  **Ollama:** O motor de inferência que executa o modelo de linguagem (Llama 3).
2.  **Open WebUI:** A interface gráfica que o usuário acessa. Ela atua como o cliente MCP, gerenciando o chat e a invocação das ferramentas.
3.  **MCPO (MCP OpenAPI Proxy):** Um proxy que converte as definições do servidor MCP para o formato OpenAPI, permitindo que o Open WebUI entenda e consuma as ferramentas facilmente.
4.  **MCP-Calc (Go Server):** O servidor MCP implementado em **Go**. Ele contém a lógica de negócio real (a calculadora) e expõe as funções `add` e `multiply` via SSE (Server-Sent Events).

## 🛠️ Tecnologias e Dependências

*   **Docker & Docker Compose:** Para orquestração e isolamento dos serviços.
*   **Go (Golang):** Linguagem utilizada para criar o servidor MCP de alta performance (`mcp-calc`).
*   **Python:** Utilizado internamente pelo Open WebUI e MCPO.
*   **Ollama:** Plataforma para rodar LLMs localmente.
*   **Llama 3:** O modelo de linguagem utilizado (pode ser substituído por outros compatíveis com tool calling).
*   **Open WebUI:** Interface de chat extensível.
*   **Model Context Protocol (MCP):** Padrão aberto para conectar assistentes de IA a sistemas onde os dados vivem.

## 📦 Como Rodar

### Pré-requisitos
*   Docker e Docker Compose instalados.
*   Git.

### Passo a Passo

1.  **Clone o repositório:**
    ```bash
    git clone <seu-repo-url>
    cd ollama-mcp-project
    ```

2.  **Suba os containers:**
    ```bash
    docker-compose up -d --build
    ```

3.  **Baixe o modelo no Ollama:**
    Acesse o container do Ollama e baixe o modelo (caso ainda não tenha):
    ```bash
    docker exec -it ollama ollama run llama3
    ```
    *(Após o download e o prompt aparecer, você pode sair com `/bye`)*

4.  **Configure a Ferramenta no Open WebUI:**
    *   Acesse `http://localhost:8080`.
    *   Crie uma conta administrativa (os dados ficam salvos localmente no volume).
    *   Vá em **Workspace** -> **Tools** -> **Create Tool**.
    *   Obtenha a especificação OpenAPI em: `http://localhost:8002/openapi.json`.
    *   Cole o JSON na definição da ferramenta no Open WebUI.
    *   **Importante:** Verifique se a URL no JSON aponta para o container interno (ex: `http://mcpo:8000`) e não para `localhost`.

5.  **Use a Ferramenta:**
    *   Inicie um novo chat.
    *   Ative a ferramenta criada.
    *   Pergunte: *"Calcule a soma de 123 e 456"*.

## 📂 Estrutura de Pastas

```
ollama-mcp-project/
├── docker-compose.yml  # Definição dos serviços
├── mcp-calc/           # Código fonte da ferramenta em Go
│   ├── main.go         # Servidor MCP e lógica da calculadora
│   ├── Dockerfile      # Build do container Go
│   └── go.mod          # Dependências Go
└── prints/             # Imagens e screenshots
```

## 🔧 Desenvolvimento

Para modificar a lógica da calculadora:
1.  Edite o arquivo `mcp-calc/main.go`.
2.  Reconstrua o container:
    ```bash
    docker-compose up -d --build mcp-calc
    ```
