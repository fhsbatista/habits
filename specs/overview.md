# Habits CLI — Visão Geral

## Problema
Ferramenta de linha de comando para cadastrar e acompanhar hábitos e tarefas do dia a dia, com rastreamento de tempo de execução e visualização em timeline.

## Conceitos centrais

### Pilar
Categoria que agrupa hábitos. Exemplos: Casa, Finanças, Relacionamentos, Carreira.

### Habit
Um hábito recorrente vinculado a um Pilar.

### Task
Uma tarefa pontual com título e descrição.

### Session
Registro de execução com início e fim. É o mecanismo de rastreamento de tempo.
Pode estar vinculada a um Habit ou a uma Task (sessionável).
A timeline é construída a partir das Sessions.

## Armazenamento
SQLite local.

## Interface
CLI com subcomandos. Exemplos de uso:
- `habits pillar add <nome>`
- `habits habit add <nome> --pillar <pilar>`
- `habits task add <título> --desc <descrição>`
- `habits start habit <id>`
- `habits stop`
- `habits log` (timeline do dia atual)
- `habits log --date 2026-05-30`

## Expansão futura
A abstração de Session e timeline foi desenhada para suportar outros tipos de entidades além de Habit e Task.
