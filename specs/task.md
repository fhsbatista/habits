# Spec: Task

## Descrição
Uma tarefa pontual com título e descrição. Assim como um Habit, pode ter Sessions vinculadas.
Segue o modelo GTD simplificado: toda tarefa em andamento possui uma "próxima ação" que indica o que fazer com ela.

## Campos
| Campo       | Tipo         | Obrigatório |
|-------------|--------------|-------------|
| id          | int64        | sim         |
| title       | string       | sim         |
| description | string       | não         |
| status      | enum         | sim         |
| createdAt   | time         | sim         |
| updatedAt   | time         | não         |

### Status
- `em_andamento` — tarefa ativa, aparece na listagem
- `concluida` — tarefa finalizada
- `arquivada` — tarefa arquivada sem conclusão

### Próxima Ação
Existe apenas após o primeiro `update`. Vinculada à tarefa enquanto estiver `em_andamento`.

| Campo  | Tipo   | Obrigatório |
|--------|--------|-------------|
| tipo   | enum   | sim         |
| detail | string | não         |

Tipos:
- `executar` — há uma ação imediata a tomar
- `aguardar` — aguardando algo externo acontecer

### Evento
Registro do que aconteceu com a tarefa num dado momento. Criado a cada `update`.

| Campo     | Tipo   | Obrigatório |
|-----------|--------|-------------|
| id        | int64  | sim         |
| taskID    | int64  | sim         |
| what      | string | sim         |
| createdAt | time   | sim         |

## Comandos CLI

### Adicionar
```
habits task add <título>
habits task add <título> --desc <descrição>
```
Tarefa nasce com status `em_andamento` e sem próxima ação.

### Atualizar (update)
```
habits task update <id>
```
Abre um prompt interativo:
1. "O que aconteceu?" → registra um Evento com `createdAt = now`
2. "Próxima ação [executar/aguardar]:" → atualiza a próxima ação da tarefa
Atualiza `updatedAt` da tarefa.

### Concluir
```
habits task done <id>
```

### Arquivar
```
habits task archive <id>
```

### Remover
```
habits task remove <id>
```

## Regras de negócio
- Título deve ser único.
- Próxima ação só existe após o primeiro `update`.
- Apenas tarefas `em_andamento` aparecem na listagem padrão.
- `updatedAt` é atualizado a cada `update`; se nunca houve update, `updatedAt` é nulo e o tempo sem update é calculado a partir de `createdAt`.
- A lógica de calcular "tempo sem update" pertence ao domínio.

## Camadas

| Responsabilidade                              | Camada       |
|-----------------------------------------------|--------------|
| Calcular tempo sem update                     | Domínio      |
| Validar transições de status                  | Domínio      |
| Determinar próxima ação vigente               | Domínio      |
| Prompt interativo de update                   | Apresentação |
| Formatação do tempo (minutos, horas, dias)    | Apresentação |
| Persistência de Task, Evento no SQLite        | Infra        |
