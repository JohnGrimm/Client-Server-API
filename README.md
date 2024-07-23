# Client Server API ( FullCycle )

### Client
  - Solicita a cotação do dólar com timeout de 300ms
  - Salva a cotação em um txt
    
### Server
  - Cria um servidor na porta 8080 com uma rota: /cotação
  - Realiza uma requisição a uma API externa com timeout de 200ms
  - Salva no banco de dados ( SQLite3 ) usando timeout de 10ms
  - Retorna a cotação em JSON


-------------------------------------------


## Starting Server

#### Inicia o servidor na porta: 8080
```sh
go run server/server.go
```


## Starting Client

#### Inicia o client
```sh
go run client/client.go
```
