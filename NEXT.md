# Next steps — controller web

Contexto: hoje só existem os pacotes de domínio `mailbox` e `tape`, sem
servidor HTTP, persistência, camada de serviço, tipo `User` ou testes. O
objetivo deste passo é construir a fatia mais fina de um controller HTTP que
exercite o ciclo de vida completo (criar mailbox → agendar abertura → gravar
tape → listar tapes → marcar tape como tocado), só para tornar óbvias as
lacunas reais que faltam para uma primeira versão utilizável.

Decisões já tomadas:

- Roteamento: `net/http` da stdlib (Go 1.22+ `ServeMux`, padrões tipo
  `"POST /mailboxes/{id}/tapes"`), sem novas dependências.
- Persistência: só as **interfaces** de repositório por enquanto, sem
  implementação concreta.
- Escopo: fluxo completo (mailbox + tape).

## Checklist

### 1. `tape/tape.go`

- [x] Adicionar `func (t *Tape) MarkPlayed(playedAt time.Time)` (seta
      `played = true`, `playedAt = &playedAt`) — simétrico ao
      `SetOpeningTime` de `Mailbox`.
- [x] Adicionar getters: `MailboxID()`, `RecordedBy()`, `RecordedAt()`,
      `Played()`, `PlayedAt()`.
- [x] Adicionar `type Repository interface` com `Save(ctx, *Tape) error`,
      `FindByID(ctx, uuid.UUID) (*Tape, error)`,
      `ListByMailboxID(ctx, uuid.UUID) ([]*Tape, error)`.
- [x] Definir `ErrNotFound` sentinel para o repositório sinalizar "não
      encontrado".

### 2. `mailbox/mailbox.go`

- [x] Adicionar getters: `OwnerID()`, `CreatedAt()`, `LastListenedAt()`,
      `OpensAt()`.
- [x] Adicionar `type Repository interface` com `Save(ctx, *Mailbox) error`,
      `FindByID(ctx, uuid.UUID) (*Mailbox, error)`.
- [x] Definir `ErrNotFound` sentinel.

### 3. `internal/api/` (pacote novo do controller)

- [ ] `router.go`: `NewRouter(mailboxes mailbox.Repository, tapes
  tape.Repository) *http.ServeMux` registrando as rotas abaixo.
- [ ] `mailbox_handlers.go`: handlers de mailbox.
- [ ] `tape_handlers.go`: handlers de tape.
- [ ] `dto.go`: structs de request/response com `json:"..."` + conversão
      de/para os tipos de domínio via os getters do passo 1/2 (nunca acessar
      campo privado diretamente).
- [ ] `respond.go`: helpers `writeJSON(w, status, v)` e `writeError(w,
  status, err)`.

Rotas:

| Método  | Caminho                        | Ação                                     |
| ------- | ------------------------------ | ---------------------------------------- |
| `POST`  | `/mailboxes`                   | cria mailbox (`ownerId`, `opensAt?`)     |
| `GET`   | `/mailboxes/{id}`              | busca mailbox                            |
| `PATCH` | `/mailboxes/{id}/opening-time` | `SetOpeningTime`                         |
| `POST`  | `/mailboxes/{id}/tapes`        | grava tape (valida que o mailbox existe) |
| `GET`   | `/mailboxes/{id}/tapes`        | lista tapes do mailbox                   |
| `POST`  | `/tapes/{id}/played`           | `MarkPlayed`                             |

- [ ] Mapear erros para status: 400 (JSON/UUID inválido), 404
      (`ErrNotFound`), 500 (qualquer outro erro do repositório).

### 4. `cmd/server/main.go`

- [ ] Montar `internal/api.NewRouter(nil, nil)` (repositórios `nil` de
      propósito, ainda sem implementação) e subir com
      `http.ListenAndServe(":8080", router)`.
- [ ] Deixar um `// TODO:` indicando que os repositórios precisam de
      implementação real antes de qualquer request funcionar de fato —
      chamadas que tocam o repositório vão panicar com nil pointer, e isso é
      esperado neste passo.

## Fora de escopo por enquanto (mas já são os próximos passos óbvios depois deste)

- [ ] Implementação concreta de `mailbox.Repository` / `tape.Repository`
      (em memória ou banco real).
- [ ] Tipo `User` (hoje `ownerID`/`recordedBy` são só UUIDs soltos).
- [ ] Autenticação/autorização.
- [ ] Testes automatizados.
- [ ] Commitar as mudanças pendentes (`tape.go` modificado, `mailbox/` não
      rastreado).

## Testes unitários (podem ser feitos já, antes do controller)

Atenção: como a maioria dos campos de `Tape`/`Mailbox` é não exportada, testes
fora do pacote (`tape_test`, `mailbox_test`) só enxergam `ID`. Para checar os
outros campos, o teste precisa estar no próprio pacote (`package tape`,
`package mailbox`) — ou esperar pelos getters já listados acima.

### `tape/tape_test.go` — `NewTape`

- [ ] Campos (`mailboxID`, `recordedBy`, `recordedAt`) são atribuídos
      corretamente a partir dos argumentos.
- [ ] `played` começa `false`.
- [ ] `playedAt` começa `nil`.
- [ ] `ID` é um UUID não-zero (`!= uuid.Nil`).
- [ ] Duas chamadas seguidas geram `ID`s diferentes e crescentes (propriedade
      do UUIDv7).

### `mailbox/mailbox_test.go` — `NewMailbox`

- [ ] Campos (`ownerID`, `createdAt`, `opensAt`) são atribuídos corretamente.
- [ ] `lastListenedAt` começa `nil`.
- [ ] `opensAt` aceita `nil` e também um ponteiro preenchido na construção.
- [ ] `ID` não-zero e crescente entre chamadas, mesma checagem do tape.

### `mailbox/mailbox_test.go` — `SetOpeningTime`

- [ ] Definir um horário novo atualiza `opensAt`.
- [ ] Chamar com `nil` limpa o campo.
- [ ] Chamar duas vezes seguidas — a segunda sobrescreve a primeira.

### Fora de escopo por enquanto

- [ ] Caminho de erro do `uuid.NewV7()` — não testável sem injetar a fonte de
      aleatoriedade, não vale o esforço agora.
- [ ] Testes de `MarkPlayed` (`tape`) — só depois de implementar o método na
      seção 1 acima.

## Verificação ao terminar

- [ ] `go build ./...` compila sem erros.
- [ ] `go run ./cmd/server` sobe o servidor na porta 8080.
- [ ] `curl -X POST localhost:8080/mailboxes -d '{"ownerId":"<uuid>"}'`
      retorna 500 (nil pointer) hoje, mas confirma que roteamento e parsing
      de JSON estão corretos mesmo sem persistência real.
