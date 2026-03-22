# 🌦️ Clima OTEL

Sistema distribuído em Go composto por dois microsserviços com
rastreamento distribuído utilizando OpenTelemetry e Zipkin.

------------------------------------------------------------------------

## 🧱 Arquitetura

    Client
      ↓
    Service A (Input)
      ↓
    Service B (Orquestração)
      ↓
    ViaCEP API
      ↓
    WeatherAPI

### Observabilidade

    Service A
       ↓
    OTEL Collector
       ↓
    Zipkin

    Service B
       ↓
    OTEL Collector
       ↓
    Zipkin

------------------------------------------------------------------------

## 🚀 Tecnologias

-   Go
-   Docker / Docker Compose
-   OpenTelemetry (OTEL)
-   Zipkin
-   ViaCEP API
-   WeatherAPI

------------------------------------------------------------------------

## 📦 Como executar

Na raiz do projeto:

``` bash
docker compose up --build
```

------------------------------------------------------------------------

## 🔎 Como testar

``` bash
curl --location 'http://localhost:8080/weather' \
--header 'Content-Type: application/json' \
--data '{"cep":"29902555"}'
```

------------------------------------------------------------------------

## ✅ Resposta de sucesso (200)

``` json
{
  "city": "São Paulo",
  "temp_C": 28.5,
  "temp_F": 83.3,
  "temp_K": 301.65
}
```

------------------------------------------------------------------------

## ❌ Respostas de erro

### CEP inválido

-   Status: `422`
-   Body:

```{=html}
<!-- -->
```
    invalid zipcode

### CEP não encontrado

-   Status: `404`
-   Body:

```{=html}
<!-- -->
```
    can not find zipcode

------------------------------------------------------------------------

## 📊 Observabilidade (Zipkin)

Acesse:

http://localhost:9411

### Como visualizar traces

1.  Execute uma requisição no Service A
2.  Abra o Zipkin
3.  Clique em **Run Query**
4.  Selecione um trace

### Exemplo de trace esperado

    service-a-weather
      ↓
    HTTP call service-b
      ↓
    service-b-weather
      ↓
    weather.execute
      ├─ viacep.lookup
      └─ weatherapi.lookup

------------------------------------------------------------------------

## 🧠 Funcionalidades

-   Validação de CEP (8 dígitos)
-   Consulta de cidade via ViaCEP
-   Consulta de temperatura via WeatherAPI
-   Conversão de temperatura:
    -   Celsius
    -   Fahrenheit
    -   Kelvin
-   Tratamento de erros (422 e 404)
-   Tracing distribuído com OpenTelemetry
-   Spans manuais para APIs externas

------------------------------------------------------------------------

## 🐳 Serviços do Docker Compose

-   service-a → porta 8080
-   service-b → porta 8081
-   otel-collector → portas 4317/4318
-   zipkin → porta 9411

------------------------------------------------------------------------

## 📌 Observações

-   O tracing distribuído é propagado via `context.Context`
-   O Service A atua como gateway/proxy
-   O Service B concentra a lógica de negócio
-   OTEL Collector centraliza o envio de traces para o Zipkin

------------------------------------------------------------------------

## 📁 Estrutura do Projeto

    clima-otel/
      service-a/
      service-b/
      infra/
        otel-collector/
          config.yaml
      docker-compose.yaml
      README.md
