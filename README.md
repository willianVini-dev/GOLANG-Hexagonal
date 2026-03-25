# GOLANG Hexagonal

Projeto de estudos em Go para aprender Arquitetura Hexagonal (Ports and Adapters), com uma API HTTP que consulta noticias na NewsAPI.

## Objetivo

Este repositorio foi criado para fins de estudo sobre:

- separacao de responsabilidades por camadas
- uso de interfaces (ports) para desacoplamento
- adaptadores de entrada (HTTP) e saida (cliente HTTP externo)
- validacao de entrada e padrao de erros REST

## Arquitetura

Estrutura principal:

```text
.
|-- main.go
|-- adapter/
|   |-- input/
|   |   |-- routes/
|   |   |-- controller/
|   |   `-- model/request/
|   `-- output/
|       |-- news_http/
|       `-- model/response/
|-- application/
|   |-- domain/
|   |-- port/input/
|   |-- port/output/
|   `-- service/
`-- configuration/
    |-- env/
    |-- logger/
    |-- rest_err/
    `-- validation/
```

Resumo das camadas:

- adapter/input: recebe HTTP, valida request e chama o caso de uso
- application: regras de negocio (service + domain + interfaces)
- adapter/output: chama servicos externos (NewsAPI)
- configuration: logger, env, validacao e erros padronizados

Para ver fluxogramas detalhados, consulte [FLUXOGRAMAS_ARQUITETURA.md](FLUXOGRAMAS_ARQUITETURA.md).

## Tecnologias

- Go
- Gin
- Resty
- Zap
- go-playground/validator
- godotenv

## Como rodar localmente

### 1. Clonar o projeto

```bash
git clone https://github.com/willianVini-dev/GOLANG-Hexagonal.git
cd GOLANG-Hexagonal
```

### 2. Configurar variaveis de ambiente

Copie [.env.example](.env.example) para .env e preencha a chave da NewsAPI:

```env
NEWS_API_KEY=your_news_api_key_here
```

### 3. Instalar dependencias

```bash
go mod tidy
```

### 4. Executar

```bash
go run main.go
```

Servidor padrao: http://localhost:8080

## Endpoint

### GET /news

Query params:

- subject (obrigatorio, minimo 2 caracteres)
- from (obrigatorio, formato YYYY-MM-DD)

Exemplo:

```bash
curl "http://localhost:8080/news?subject=technology&from=2026-03-24"
```

Resposta de sucesso (exemplo):

```json
{
  "Status": "ok",
  "TotalResults": 10,
  "Articles": [
    {
      "Source": {
        "Id": null,
        "Name": "Example"
      },
      "Author": "Author name",
      "Title": "News title",
      "Description": "Description...",
      "UrlToImage": "https://...",
      "PublishedAt": "2026-03-24T10:00:00Z",
      "Content": "Content..."
    }
  ]
}
```

Resposta de erro de validacao (exemplo):

```json
{
  "message": "Some fields are invalid",
  "code": 400,
  "error": "bad_request",
  "causes": [
    {
      "field": "Subject",
      "message": "Subject is a required field"
    }
  ]
}
```

## Observacoes

- Projeto com foco didatico (arquitetura e organizacao de camadas)
- Nomes/campos de resposta seguem o dominio atual do projeto
- Evolucoes naturais: testes unitarios, mocks para ports, e tratamento mais completo de erros da API externa

## Licenca

Uso educacional / estudos.
