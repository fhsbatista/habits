# Spec: Task List

## Descrição
Listagem unificada de tarefas `em_andamento` e hábitos pendentes do dia. Acessível via `habits list` ou `habits ls`.

## Comando
```
habits list
habits ls
```

## Layout

A listagem é dividida em três seções:

### Seção 1 — Hábitos de hoje
Hábitos que devem ser performados no dia atual e ainda não foram. Exibe nome e pilar vinculado.

### Seção 2 — Executar
Tarefas cuja próxima ação é `executar`, **ou** que ainda não tiveram nenhum `update`.

### Seção 3 — Aguardar
Tarefas cuja próxima ação é `aguardar`.

### Ordenação
Dentro de cada seção: tarefas com **mais tempo sem update aparecem primeiro** (ou seja, as atualizadas mais recentemente ficam por último).

### Colunas
| Coluna           | Descrição                                                         |
|------------------|-------------------------------------------------------------------|
| ID               | Identificador da tarefa                                           |
| Título           | Título da tarefa                                                  |
| Descrição        | Descrição da tarefa (truncada se longa)                           |
| Sem update       | Tempo desde o último update (ou desde criação se nunca atualizada)|

### Formato do tempo sem update
Exibido na unidade mais adequada, sem arredondamento para cima:
- `X minuto(s)` — até 59 minutos
- `X hora(s)` — até 23 horas
- `X dia(s)` — a partir de 1 dia

**Exemplo visual:**
```
[ HÁBITOS DE HOJE ]
ID       Nome                Pilar
-------- ------------------- ----------
x9y8z7   Meditação           Saúde
m1n2o3   Leitura             Carreira

[ EXECUTAR ]
ID       Título              Descrição                  Sem update
-------- ------------------- -------------------------- ----------
a1b2c3   Refatorar auth      Separar middleware JWT     5 dias
d4e5f6   Responder e-mail    Cliente X aguardando       2 horas
g7h8i9   Revisar PR          -                          10 minutos

[ AGUARDAR ]
ID       Título              Descrição                  Sem update
-------- ------------------- -------------------------- ----------
j1k2l3   Deploy produção     Aguardando janela de maint 3 dias
```

## Regras de negócio
- Tarefas `concluida` e `arquivada` não aparecem.
- Tarefas sem nenhum `update` aparecem na seção **Executar**.
- Tempo sem update calculado a partir de `updatedAt` se existir, senão `createdAt`.
- Hábitos aparecem na seção **Hábitos de hoje** apenas se forem devidos no dia atual e ainda não tiverem sido performados.
- Se não houver hábitos pendentes, a seção **Hábitos de hoje** é omitida.

## Camadas

| Responsabilidade                                       | Camada       |
|--------------------------------------------------------|--------------|
| Filtrar tarefas em andamento                           | Domínio      |
| Separar por tipo de próxima ação                       | Domínio      |
| Ordenar por tempo sem update                           | Domínio      |
| Calcular tempo sem update                              | Domínio      |
| Determinar hábitos pendentes do dia                    | Domínio      |
| Formatar tempo em minutos/horas/dias                   | Apresentação |
| Renderizar tabela no terminal                          | Apresentação |
| Omitir seção vazia de hábitos                          | Apresentação |
| Buscar tarefas e hábitos no SQLite                     | Infra        |
