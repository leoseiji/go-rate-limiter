# go-rate-limiter

Comando docker-compose: docker-compose up -d

Comando para rodar a aplicação: go run main.go

Configuração arquivo na raiz: .env

* WEB_SERVER_PORT: Porta do servidor
* ENABLE_RATE_LIMIT_BY_IP: true para ativar
* ENABLE_RATE_LIMIT_BY_TOKEN: true para ativar
* MAX_REQUESTS_BY_IP: Requests por IP por segundo
* BLOCK_DURATION_IP: Duração do bloqueio por ip em segundos
* BLOCK_DURATION_TOKEN: Duração do bloqueio por token em segundos
* TOKEN_LIMIT_LIST: Lista dos tokens e limite por segundo.
* STORAGE_TYPE: Tipo da base.

Arquivos de exemplo para rodar os requests:
* api/get_token.http
* api/get.http
