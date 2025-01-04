# League of Legends Elo Police

Bot do Whatsapp para informar atualizações nas filas ranqueada dos players cadastrados.

## Como funciona?

Cadastre o player pelo endpoint `/player`. Olhe `docs/bruno` para referência.

Sempre quando uma partida de Normal Game, Flex Queue, Solo Queue e ARAM terminar vai ser enviado uma mensagem no grupo de Whatsapp com os dados do player e da partida.

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

- Para buildar binários:
  ```
  make build
  ```
