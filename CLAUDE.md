# Habits CLI — Guia para Agentes de IA

## Visão geral

CLI em Go para rastreamento de hábitos e tarefas com visualização de timeline. Consulte `specs/overview.md` para entender os conceitos centrais do projeto.

## Linguagem e dependências

- Go 1.24+
- Sem frameworks CLI externos (stdlib apenas)
- SQLite para persistência

## Arquitetura

O projeto segue Clean Architecture. As camadas são:

- **`internal/domain`** — entidades e interfaces de repositório. Sem dependências externas. Contém apenas regras de negócio puras.
- **`internal/usecase`** — casos de uso que orquestram o domínio.
- **`internal/repository`** — implementações dos repositórios (SQLite).
- **`internal/cli`** — camada de apresentação. Trata input/output do terminal.
- **`cmd/habits`** — entrypoint da aplicação.

### Regra de ouro da arquitetura

Antes de colocar lógica numa camada, pergunte: *"isso faria sentido numa interface diferente (web, mobile)?"*
- Sim → é regra de negócio → vai no **domínio**
- Não, é específico do terminal → vai na **apresentação**
- É acesso a banco/rede → vai na **infra/repository**

Dependências sempre apontam para dentro: `cli → usecase → domain ← repository`.

## Specs

Cada feature tem uma spec em `specs/`. Antes de implementar qualquer coisa, leia a spec correspondente. Se a spec não existir, crie-a e aguarde aprovação antes de implementar.

Arquivos disponíveis:
- `specs/overview.md` — visão geral e conceitos
- `specs/pillar.md` — cadastro de pilares
- `specs/habit.md` — cadastro de hábitos com frequência e cor
- `specs/task.md` — cadastro de tarefas com status, próxima ação e eventos
- `specs/session.md` — check-in/check-out e timeline visual
- `specs/habit-grid.md` — grid de consistência de hábitos
- `specs/task-list.md` — listagem unificada de tarefas e hábitos pendentes

## Fluxo obrigatório de desenvolvimento (TDD)

Toda implementação de feature deve seguir este fluxo **sem pular etapas**:

1. **Leia a spec** da feature em `specs/`
2. **Escreva o teste** para o comportamento a ser implementado
3. **Aguarde confirmação** do usuário de que o teste faz sentido
4. **Execute o teste** e confirme que ele **falha** (`go test ./...`)
5. **Implemente** o código mínimo para o teste passar
6. **Execute o teste** novamente e confirme que ele **passa**
7. Só então a feature é considerada completa

Nunca implemente código antes de ter um teste falhando aprovado.

## Convenções de código

- Nomes em inglês no código, português nas mensagens exibidas ao usuário
- Sem comentários óbvios — apenas quando o "porquê" não é evidente pelo código
- IDs são `int64` (gerados pelo SQLite via autoincrement)
- Erros sempre propagados com contexto: `fmt.Errorf("criar hábito: %w", err)`
- Interfaces definidas no domínio, implementadas na infra
