# Fluxogramas da Arquitetura Hexagonal

---

## O que é Arquitetura Hexagonal? (Explicação simples)

Imagine uma **cebola com camadas**. A parte mais interna (o miolo) é a regra de negócio — o que seu sistema realmente faz. As camadas de fora são responsáveis por **falar com o mundo externo** (receber chamadas HTTP, chamar APIs, banco de dados, etc.).

A regra de ouro é: **o miolo nunca sabe quem está do lado de fora**. Ele só conhece contratos (interfaces). Isso permite trocar um banco de dados, uma API externa ou um framework HTTP sem mexer na lógica de negócio.

```
    ┌──────────────────────────────────────────────────────────┐
    │  MUNDO EXTERNO                                           │
    │  (Usuário HTTP, News API, etc.)                          │
    │                                                          │
    │   ┌────────────────────────────────────────────────┐     │
    │   │  ADAPTERS (Tradutores)                         │     │
    │   │  Convertem o mundo externo para a linguagem    │     │
    │   │  interna da aplicação e vice-versa             │     │
    │   │                                                │     │
    │   │   ┌──────────────────────────────────────┐     │     │
    │   │   │  PORTS (Contratos / Interfaces)      │     │     │
    │   │   │  Definem o que pode entrar e sair    │     │     │
    │   │   │  sem expor os detalhes internos      │     │     │
    │   │   │                                      │     │     │
    │   │   │   ┌──────────────────────────────┐   │     │     │
    │   │   │   │  CORE (Miolo — Negócio Puro) │   │     │     │
    │   │   │   │  Service + Domain            │   │     │     │
    │   │   │   │  Não depende de nada externo │   │     │     │
    │   │   │   └──────────────────────────────┘   │     │     │
    │   │   └──────────────────────────────────────┘     │     │
    │   └────────────────────────────────────────────────┘     │
    └──────────────────────────────────────────────────────────┘
```

---

## Mapa do Projeto — Arquivo por Arquivo

Antes dos diagramas, veja onde cada arquivo vive e o que é responsável por fazer:

```
Hexagonal/
│
├── main.go                              ← Ponto de partida. Liga tudo.
│
├── adapter/                             ← TRADUTORES (falam com o mundo externo)
│   ├── input/                           ← Entrada: recebe chamadas HTTP
│   │   ├── routes/routes.go             ← Define quais URLs existem e monta dependências
│   │   ├── controller/news_controller.go← Recebe a requisição, valida, chama regra de negócio
│   │   └── model/request/new_request.go ← Molde dos dados que o usuário precisa enviar
│   │
│   └── output/                          ← Saída: faz chamadas para fora
│       ├── news_http/news_api_client.go ← Chama a API externa newsapi.org
│       └── model/response/              ← Molde da resposta que a API externa devolve
│
├── application/                         ← MIOLO (regra de negócio pura)
│   ├── domain/news_domain.go            ← As estruturas de dados do negócio
│   ├── port/
│   │   ├── input/news_use_case.go       ← Contrato: o que a camada de entrada pode pedir
│   │   └── output/news_port.go          ← Contrato: o que a camada de saída precisa fazer
│   └── service/news_service.go          ← Implementa a regra de negócio
│
├── configuration/                       ← UTILITÁRIOS (usados por qualquer camada)
│   ├── env/get_env.go                   ← Lê variáveis de ambiente (.env)
│   ├── logger/logger.go                 ← Sistema de log estruturado
│   ├── rest_err/rest_err.go             ← Padrão de resposta de erro HTTP
│   └── validation/validation.go         ← Traduz erros de validação para mensagens amigáveis
│
└── .env                                 ← Configurações sensíveis (token da API, etc.)
```

---
## Diagrama 1 — Visão Geral das Camadas

> **Para iniciantes:** Pense nisso como um funil. A requisição do usuário entra pelo lado esquerdo, passa pelas camadas, busca dado real numa API externa e a resposta volta. Cada caixa tem uma responsabilidade bem separada.

```mermaid
flowchart LR
    user["👤 Usuário\nfaz GET /news"]

    subgraph ai["ADAPTER INPUT\nadapter/input/\n(Porta de entrada HTTP)"]
        routes["routes.go\nRegistra a rota GET /news\ne monta todas as dependências"]
        controller["news_controller.go\nRecebe a requisição,\nvalida os dados,\nchama a regra de negócio"]
        req["new_request.go\nMolde dos dados esperados:\n- subject (texto, min 2 chars)\n- from (data AAAA-MM-DD)"]
    end

    subgraph app["APPLICATION\napplication/\n(Miolo — Regra de Negócio)"]
        iport["port/input/news_use_case.go\nCONTRATO de entrada:\ndefine o que o controller\npode pedir ao service"]
        svc["service/news_service.go\nIMPLEMENTA a regra:\nrecebe dados, chama\na porta de saída"]
        oport["port/output/news_port.go\nCONTRATO de saída:\ndefine o que o adapter\nprecisa fornecer"]
        dom["domain/news_domain.go\nESTRUTURAS:\nNewsDomain, Article,\nNewsRequestDomain"]
    end

    subgraph ao["ADAPTER OUTPUT\nadapter/output/\n(Porta de saída externa)"]
        apiclient["news_api_client.go\nFaz a chamada HTTP real\npara newsapi.org/v2/everything\ncom resty"]
        resp["model/response/\nMolde da resposta\nda API externa:\nNewsClientResponse"]
    end

    subgraph cfg["CONFIGURATION\nconfiguration/\n(Utilitários — usados por todos)"]
        log["logger.go\nLog estruturado\nem JSON (zap)"]
        val["validation.go\nTraduz erros técnicos\nem mensagens amigáveis"]
        rerr["rest_err.go\nPadrão de erro HTTP:\n400, 404, 500"]
        env["get_env.go\nLê NEWS_API_KEY\ndo arquivo .env"]
    end

    newsapi["🌐 newsapi.org\nAPI externa de notícias"]

    user --> routes
    routes --> controller
    controller --> req
    req -->|dados OK| iport
    req -->|dados inválidos| val
    val --> rerr
    rerr -->|JSON erro 400| user

    iport --> svc
    svc --> log
    svc --> oport
    oport --> apiclient
    apiclient --> env
    apiclient --> resp
    apiclient --> newsapi
    newsapi -->|JSON com artigos| apiclient
    apiclient --> dom
    dom -->|NewsDomain preenchido| svc
    svc -->|retorna domínio| controller
    controller -->|JSON 200| user

    classDef blueBox fill:#4A90E2,stroke:#2E5C8A,color:#fff
    classDef greenBox fill:#27AE60,stroke:#1A6B3A,color:#fff
    classDef purpleBox fill:#8E44AD,stroke:#5B2C6F,color:#fff
    classDef orangeBox fill:#E67E22,stroke:#A04000,color:#fff
    classDef grayBox fill:#7F8C8D,stroke:#566573,color:#fff

    class ai blueBox
    class app greenBox
    class ao purpleBox
    class cfg orangeBox
    class newsapi grayBox
```

---

## Diagrama 2 — Como as Dependências são Montadas (routes.go)

> **Para iniciantes:** Antes de qualquer requisição chegar, o `routes.go` já montou a "cadeia de dependências" como peças de Lego encaixadas. Cada peça só conhece o contrato da próxima, nunca a implementação direta.

```mermaid
flowchart TD
    start["main.go chama\nroutes.InitRoutes(router)"]

    step1["1. Cria newsApiClient\nnews_http.NewNewsApiClient()\n\nConfigura resty com\nbaseURL: newsapi.org/v2"]

    step2["2. Cria newsService\nservice.NewNewsService(newsApiClient)\n\nInjecta newsApiClient dentro\ndo service via interface NewsPort\n\nService não sabe que é resty,\nsó sabe que implementa NewsPort"]

    step3["3. Cria newsController\ncontroller.NewNewsController(newsService)\n\nInjecta newsService dentro\ndo controller via interface NewsUseCase\n\nController não sabe que é newsService,\nsó sabe que implementa NewsUseCase"]

    step4["4. Registra a rota\nr.GET('/news', newsController.GetNews)\n\nGin vai chamar GetNews\nquando chegar GET /news"]

    explain["Resultado: a cadeia está montada\nnewsApiClient → newsService → newsController\nmas cada um só conhece o CONTRATO do próximo"]

    start --> step1 --> step2 --> step3 --> step4 --> explain
```

---

## Diagrama 3 — Fluxo Completo de uma Requisição (Caminho Feliz)

> **Para iniciantes:** Este é o passo a passo exato de O QUE ACONTECE no código quando você chama `GET /news?subject=tech&from=2026-03-24` e tudo dá certo.

```mermaid
sequenceDiagram
    actor User as Usuário
    participant Ctrl as news_controller.go
    participant Svc as news_service.go
    participant Client as news_api_client.go
    participant Env as get_env.go
    participant Log as logger.go
    participant API as newsapi.org

    User->>Ctrl: GET /news?subject=tech&from=2026-03-24

    Note over Ctrl: ShouldBindQuery(&NewsRequest)<br/>Gin converte query string em struct:<br/>NewsRequest{Subject:"tech", From: time.Time(2026-03-24)}

    Ctrl->>Log: Info("Getting news...")

    Note over Ctrl: Monta NewsRequestDomain{<br/>  Subject: "tech",<br/>  From: "2026-03-24"<br/>}

    Ctrl->>Svc: GetNewsService(NewsRequestDomain)

    Svc->>Log: Info("init get news service, subject=tech, from=2026-03-24")

    Svc->>Client: GetNewsPort(NewsRequestDomain)

    Client->>Env: GetNewsTokenApi()
    Env-->>Client: "c8aa5658b7e14da..." (lê do .env)

    Note over Client: Monta query params:<br/> q=tech<br/> from=2026-03-24<br/> apiKey=c8aa5658...

    Client->>API: GET newsapi.org/v2/everything?q=tech&from=...&apiKey=...
    API-->>Client: JSON com Status, TotalResults, Articles[]

    Note over Client: Mapeia resposta para domínio:<br/>NewsDomain{<br/>  Status: "ok",<br/>  TotalResults: 35,<br/>  Articles: [...]<br/>}

    Client-->>Svc: NewsDomain, nil (sem erro)
    Svc-->>Ctrl: NewsDomain, nil

    Ctrl->>User: HTTP 200 + JSON NewsDomain
```

---

## Diagrama 4 — Fluxo de Erro (Validação Falha)

> **Para iniciantes:** O que acontece quando o usuário manda dados errados? Por exemplo, `subject=a` (menos de 2 caracteres) ou sem a data `from`.

```mermaid
sequenceDiagram
    actor User as Usuário
    participant Ctrl as news_controller.go
    participant Valid as validation.go
    participant Log as logger.go

    User->>Ctrl: GET /news?subject=a
    Note over Ctrl: ShouldBindQuery falha<br/>porque subject tem só 1 char<br/>(regra: min=2)

    Ctrl->>Log: Error("Error trying to bind query:", err)
    Ctrl->>Valid: ValidateUseError(err)

    Note over Valid: Detecta que é validator.ValidationErrors<br/>Itera cada campo com problema:<br/>- field: "Subject"<br/>- message: "Subject must be at least 2 chars"<br/>Monta lista de Causes[]

    Valid-->>Ctrl: RestErr{Message:"Some fields are invalid", Causes:[...]}

    Ctrl->>User: HTTP 400 + JSON de erro

    Note over User: Resposta recebida:<br/>{<br/>  "message": "Some fields are invalid",<br/>  "causes": [<br/>    {"field":"Subject","message":"..."}<br/>  ]<br/>}
```

---

## Diagrama 5 — O Papel dos Contratos (Interfaces / Ports)

> **Para iniciantes:** As interfaces são como "tomadas padronizadas". O controller não sabe se está falando com `newsService` ou outro service qualquer — só sabe que o "plugue" (interface `NewsUseCase`) encaixa. O mesmo vale para o service com o API client.

```mermaid
flowchart LR
    subgraph controller_box["news_controller.go"]
        ctrl_dep["Depende de:\ninput.NewsUseCase\n(interface)"]
    end

    subgraph input_port_box["port/input/news_use_case.go"]
        iface1["interface NewsUseCase\n─────────────────\nGetNewsService(\n  domain.NewsRequestDomain\n) (*NewsDomain, *RestErr)"]
    end

    subgraph service_box["service/news_service.go"]
        svc_impl["newsService\nIMPLEMENTA NewsUseCase\n─────────────────\nTambém depende de:\noutput.NewsPort\n(interface)"]
    end

    subgraph output_port_box["port/output/news_port.go"]
        iface2["interface NewsPort\n─────────────────\nGetNewsPort(\n  domain.NewsRequestDomain\n) (*NewsDomain, *RestErr)"]
    end

    subgraph client_box["news_api_client.go"]
        client_impl["newsApiClient\nIMPLEMENTA NewsPort\n─────────────────\nChama newsapi.org\ncom resty HTTP client"]
    end

    ctrl_dep -->|satisfeita por| svc_impl
    ctrl_dep -- define contrato --> iface1
    svc_impl -- define contrato --> iface2
    svc_impl -->|depende de| iface2
    iface2 -->|satisfeita por| client_impl

    classDef portStyle fill:#E67E22,stroke:#A04000,color:#fff
    classDef implStyle fill:#27AE60,stroke:#1A6B3A,color:#fff
    classDef ctrlStyle fill:#4A90E2,stroke:#2E5C8A,color:#fff

    class input_port_box,output_port_box portStyle
    class service_box,client_box implStyle
    class controller_box ctrlStyle
```

---

## Diagrama 6 — Estrutura do Erro REST (rest_err.go)

> **Para iniciantes:** Toda vez que algo dá errado, a API responde com um JSON padronizado. Este diagrama mostra quais tipos de erro existem e quando cada um é usado.

```mermaid
flowchart TD
    err_struct["RestErr\n─────────────────────\nmessage: string  → mensagem legível\ncode: int        → código HTTP (400, 500...)\nerror: string    → chave do erro\ncauses: Causes[] → detalhes dos campos"]

    subgraph factories["Funções que criam erros"]
        f1["NewBadRequestError(msg)\n→ code: 400\n→ error: 'bad_request'\nUSO: erro simples de entrada"]
        f2["NewBadRequestValidationError(msg, causes)\n→ error: 'bad_request'\n→ inclui lista de campos inválidos\nUSO: falha de validação de campo"]
        f3["NewInternalServerError(msg)\n→ code: 500\n→ error: 'internal_server_error'\nUSO: erro ao chamar API externa"]
        f4["NewNotFoundError(msg)\n→ code: 404\n→ error: 'not_found'\nUSO: recurso não existe"]
    end

    subgraph causes_struct["Causes (detalhes do campo)"]
        c["field: string   → nome do campo com problema\nmessage: string → mensagem traduzida"]
    end

    err_struct --- f1
    err_struct --- f2
    err_struct --- f3
    err_struct --- f4
    f2 --> causes_struct
```

---

## Diagrama 7 — Sistema de Log (logger.go)

> **Para iniciantes:** O logger grava tudo que acontece na aplicação em formato JSON para facilitar monitoramento. Você configura o destino e o nível pelo `.env`.

```mermaid
flowchart LR
    env_log["Variáveis no .env\n─────────────────\nLOG_OUTPUT=stdout\nLOG_LEVEL=info"]

    subgraph logger_init["init() — roda automaticamente ao iniciar"]
        parse["Lê LOG_OUTPUT e LOG_LEVEL\ndo ambiente"]
        build["Monta zap.Config:\n- encoding: JSON\n- campos: level, time, message"]
        create["Cria logger global\nlog *zap.Logger"]
    end

    subgraph funcs["Funções disponíveis"]
        info_fn["logger.Info(msg, ...fields)\nNível: INFO\nUso: eventos normais"]
        error_fn["logger.Error(msg, err, ...fields)\nNível: INFO (com campo 'error')\nUso: erros ocorridos"]
    end

    subgraph output_ex["Exemplo de saída JSON"]
        json_out["{\"level\":\"info\",\n\"time\":\"2026-03-24T10:00:00\",\n\"message\":\"Getting news...\"}"]
    end

    env_log --> parse --> build --> create
    create --> info_fn
    create --> error_fn
    info_fn --> json_out
    error_fn --> json_out
```

---

## Diagrama 8 — Resumo Visual das Camadas (para colar na cabeça)

```mermaid
flowchart TB
    subgraph world["MUNDO EXTERNO"]
        user2["Usuário HTTP"]
        ext_api["newsapi.org"]
    end

    subgraph adapters["ADAPTERS — Conversores entre mundo externo e aplicação\nnunca têm regra de negócio"]
        in_ad["INPUT ADAPTER\nroutes + controller + request model\nConverte HTTP → domínio interno"]
        out_ad["OUTPUT ADAPTER\nnews_api_client + response model\nConverte domínio interno → HTTP externo"]
    end

    subgraph ports["PORTS — Contratos (interfaces)\ndefinem o que pode entrar/sair sem expor detalhes"]
        in_p["INPUT PORT\nNewsUseCase interface\nO que o controller pode pedir"]
        out_p["OUTPUT PORT\nNewsPort interface\nO que o adapter de saída deve fornecer"]
    end

    subgraph core["CORE — Coração da aplicação\nnão depende de framework, banco ou HTTP"]
        service2["newsService\nOrquestra: recebe pedido,\npede dado externo,\ndevolve domínio"]
        domain2["NewsDomain / Article\nNewsRequestDomain\nAs estruturas de dado do negócio"]
    end

    subgraph cross["CONFIGURATION — Utilitários transversais\nqualquer camada pode usar"]
        cfg_items["logger · validation · rest_err · env"]
    end

    user2 -->|HTTP| in_ad
    in_ad -->|interface| in_p
    in_p --> service2
    service2 --> domain2
    service2 -->|interface| out_p
    out_p --> out_ad
    out_ad -->|HTTP| ext_api
    ext_api -->|resposta| out_ad
    out_ad --> domain2
    domain2 --> service2
    service2 --> in_p
    in_p --> in_ad
    in_ad -->|HTTP| user2

    core -.->|usa| cross
    adapters -.->|usa| cross

    classDef worldStyle fill:#7F8C8D,stroke:#566573,color:#fff
    classDef adapterStyle fill:#4A90E2,stroke:#2E5C8A,color:#fff
    classDef portStyle fill:#E67E22,stroke:#A04000,color:#fff
    classDef coreStyle fill:#27AE60,stroke:#1A6B3A,color:#fff
    classDef crossStyle fill:#C0392B,stroke:#922B21,color:#fff

    class world worldStyle
    class adapters adapterStyle
    class ports portStyle
    class core coreStyle
    class cross crossStyle
```
