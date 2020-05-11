
# Processador de Webhooks - GO

Projeto tutorial para testar a execução paralela de HTTP requests no Go em altas cargas, utilizando `goroutines`.
Este processador de webhooks recebe uma requisição para enviar uma certa mensagem para um endpoint de destino correspondente, chama uma `goroutine` para realizar o POST, e retorna imediatamente para o cliente.
