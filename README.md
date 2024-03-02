# League of Legends Elo Police

Bot do Whatsapp para informar atualizações nas filas ranqueada dos players cadastrados.

## Como funciona?

Cadastre o player pelo endpoint `/player`. Olhe `docs/bruno` para referência.

Sempre quando hover uma mudança nos stats na fila ranqueada Solo/Duo vai ser enviado uma mensagem no grupo de Whatsapp com os dados do player.

Para utilizar o bot é preciso antes registrar um número de Whatsapp. Executando `cmd/register/main.go` retorna o código que via QR Code é possível registrar igual Whatsapp Web.

## Como rodar?

Primeiro criar `.env`. Use o `.env.example` como referência.

- Para registar número do Whatsapp:
    ```
    make register
    ```

- Para rodar aplicação:
    ```
    make run
    ```

- Para de-registar número do Whatsapp:
    ```
    make deregister
    ```

- Para buildar binários:
    ```
    make build
    ```

- Para buildar binários para `armv6`:
    ```
    make build-armv6
    ```