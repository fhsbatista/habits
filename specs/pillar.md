# Spec: Pilar

## Descrição
Categoria que agrupa hábitos. Serve para organizar os hábitos por área da vida.

## Campos
| Campo     | Tipo   | Obrigatório |
|-----------|--------|-------------|
| id        | int64  | sim         |
| name      | string | sim         |
| createdAt | time   | sim         |

## Comandos CLI

### Adicionar
```
habits pillar add <nome>
```

### Listar
```
habits pillar list
```

### Remover
```
habits pillar remove <id>
```

## Regras
- Nome deve ser único.
- Não é possível remover um pilar que possui hábitos vinculados.
