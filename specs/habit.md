# Spec: Habit

## Descrição
Um hábito recorrente que o usuário deseja rastrear. Está vinculado a um Pilar e possui uma frequência semanal que define em quais dias da semana deve ser performado.

## Campos
| Campo     | Tipo     | Obrigatório |
|-----------|----------|-------------|
| id        | int64    | sim         |
| name      | string   | sim         |
| pillarID  | int64    | sim         |
| color     | enum     | sim         |
| frequency | []weekday| sim         |
| createdAt | time     | sim         |

### Frequência
Lista de dias da semana em que o hábito deve ser performado.
Dias: `dom`, `seg`, `ter`, `qua`, `qui`, `sex`, `sab`.

Exemplos:
- Todo dia: `["dom","seg","ter","qua","qui","sex","sab"]`
- Dias úteis: `["seg","ter","qua","qui","sex"]`
- Só segunda e quinta: `["seg","qui"]`

### Cores disponíveis
`azul`, `verde`, `laranja`, `vermelho`, `amarelo`

Se não informada, uma cor é escolhida aleatoriamente.

## Comandos CLI

### Adicionar
```
habits habit add <nome> --pillar <pillar-id> --days seg,qua,sex
habits habit add <nome> --pillar <pillar-id> --days todos
```
Durante o cadastro, exibe prompt interativo para seleção de cor:
```
Escolha uma cor [azul, verde, laranja, vermelho, amarelo] (Enter para aleatório):
```

### Listar
```
habits habit list
habits habit list --pillar <pillar-id>
```

### Remover
```
habits habit remove <id>
```

## Regras de negócio
- O pilar referenciado deve existir.
- Nome deve ser único dentro do mesmo pilar.
- Um hábito "é devido hoje" se o dia da semana atual estiver na sua frequência.
- Um hábito "foi performado hoje" se existir ao menos uma Session finalizada vinculada a ele com `startedAt` no dia atual.
- Um hábito "está pendente hoje" se for devido hoje e não tiver sido performado hoje.

## Camadas

| Responsabilidade                                    | Camada       |
|-----------------------------------------------------|--------------|
| Determinar se hábito é devido num dado dia          | Domínio      |
| Determinar se hábito foi performado num dado dia    | Domínio      |
| Determinar se hábito está pendente num dado dia     | Domínio      |
| Validação e seleção aleatória de cor                | Domínio      |
| Prompt interativo de seleção de cor                 | Apresentação |
| Renderização colorida no terminal                   | Apresentação |
| Persistência no SQLite                              | Infra        |
