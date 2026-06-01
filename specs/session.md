# Spec: Session

## Descrição
Registro de execução de um Habit ou Task. Armazena o momento de início (check-in) e fim (check-out), permitindo calcular o tempo gasto e construir a timeline do dia.

## Campos
| Campo        | Tipo      | Obrigatório |
|--------------|-----------|-------------|
| id           | int64     | sim         |
| refType      | enum      | sim         |
| refID        | int64     | sim         |
| startedAt    | time      | sim         |
| finishedAt   | time      | não         |

`refType` pode ser: `habit` ou `task`.

## Comandos CLI

### Iniciar (check-in)
```
habits start habit <id>
habits start task <id>
```
Cria uma Session com `startedAt = now`. Só pode haver uma Session aberta por vez.

### Finalizar (check-out)
```
habits stop
```
Preenche `finishedAt = now` na Session aberta.

### Timeline
```
habits log
habits log --date 2026-05-30
```

Exibe uma timeline visual do dia no terminal.

**Layout:**
- Eixo horizontal: tempo (00:00 → 23:59)
- Eixo vertical: cada linha é uma Session
- À direita de cada linha: título do Habit ou Task vinculado
- Cabeçalho com marcações de hora (00, 06, 12, 18, 24)

**Escala automática:** a largura da timeline se ajusta à largura atual do terminal.
Reserva ~20 caracteres para o título à direita. O restante é dividido por 1440 minutos para determinar quantos minutos cada `█` representa. A escala é exibida no cabeçalho (ex: `escala: 1█ = 15min`).

**Cores:**
- Hábitos: `█` em verde
- Tarefas: `█` em roxo
- Títulos: branco

**Session em andamento** (sem `finishedAt`): o bloco se estende até `now`.

**Exemplo visual:**
```
00       06       12       18       24     escala: 1█ = 15min
|        |        |        |        |
                  ████                     Meditação
                      ████████             Estudar Go
                              ██           Pagar contas
```

## Regras
- Só uma Session pode estar aberta (sem `finishedAt`) por vez.
- `stop` falha se não houver Session aberta.
- `start` falha se já houver Session aberta.
- `finishedAt` deve ser posterior a `startedAt`.
