
# Processador de Webhooks - GO

Projeto tutorial para testar a execução paralela de HTTP requests no Go em altas cargas, utilizando `goroutines`.

Este processador de webhooks recebe uma requisição para enviar uma certa mensagem para um endpoint de destino correspondente, chama uma `goroutine` para realizar o POST, e retorna imediatamente para o cliente.

Por ser apenas um teste de paralelismo, não foi implementado o controle de entrega das mensagens (e retentativas, etc).

## Instalação

Utilizando docker:
`````
git clone https://github.com/dmkorb/webhook-processor-go.git
docker build -t webhook-processor-go .
docker run --publish 8000:8000 --name webhook-go --rm webhook-processor-go
`````

## Utilização

O processor expõe um endpoint para o envio de requisições, `/webhooks/message`. O `body` deve contem um objeto com dois campos, `user_id` e `data`.

O projeto contém um db dummy com dois `user_id`'s, "1" e "2".

Para testar, pode-se mandar um POST http para este endpoint: 
`````
curl  -d '{ "user_id": "1","data": "Test message to be sent to user 1!"}' http://localhost:8000/webhooks/message
`````
O processador retornará imediatamente contendo o endpoint de destino na mensagem de retorno.

## Teste de carga

Para testar o paralelismo, podemos utilizar alguma ferramenta de teste de carga, como por exemplo o [loadtest](https://www.npmjs.com/package/loadtest).

Este comando abaixo envia 2000 mensagens, com 10 usuários paralelos e uma carga de 100 mensagens por segundo:
```
loadtest http://localhost:8000/webhooks/message -m POST -P '{ "user_id": "1","data": "This is a test message!"}' -T 'application/json' --rps 100 -n 2000 -c 10
```
Abaixo um exemplo de resultado deste teste rodando localmente, com um tempo de resposta médio de 7.3ms para 2000 requests:
```
INFO Requests: 0 (0%), requests per second: 0, mean latency: 0 ms
INFO Requests: 443 (22%), requests per second: 89, mean latency: 7.2 ms
INFO Requests: 943 (47%), requests per second: 100, mean latency: 7.2 ms
INFO Requests: 1443 (72%), requests per second: 100, mean latency: 7.2 ms
INFO Requests: 1942 (97%), requests per second: 100, mean latency: 7.1 ms
INFO 
INFO Target URL:          http://localhost:8000/webhooks/message
INFO Max requests:        2000
INFO Concurrency level:   10
INFO Agent:               none
INFO Requests per second: 100
INFO 
INFO Completed requests:  2000
INFO Total errors:        0
INFO Total time:          20.564933608 s
INFO Requests per second: 97
INFO Mean latency:        7.3 ms
INFO 
INFO Percentage of the requests served within a certain time
INFO   50%      6 ms
INFO   90%      10 ms
INFO   95%      12 ms
INFO   99%      16 ms
INFO  100%      37 ms (longest request)
```
