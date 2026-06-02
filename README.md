# habits

CLI para rastreamento de hábitos e tarefas com visualização de timeline e grid de consistência.

## Requisitos

- Go 1.25+

## Build

```bash
go build -o habits ./cmd/habits
```

Para instalar globalmente:

```bash
sudo go build -o /usr/local/bin/habits ./cmd/habits
```

## Banco de dados

### Local (padrão)

O banco SQLite é criado automaticamente em `~/.habits.db`. Para usar outro caminho:

```bash
export HABITS_DB=/caminho/para/habits.db
```

### Turso (cloud)

Para usar um banco na nuvem via [Turso](https://turso.tech), defina as variáveis de ambiente abaixo. O banco local é ignorado quando `TURSO_URL` está definido.

```bash
export TURSO_URL=libsql://seu-banco.turso.io
export TURSO_TOKEN=seu-token
```

Para obter a URL e o token:

```bash
# Instalar CLI do Turso
curl -sSfL https://get.tur.so/install.sh | bash

# Login e criação do banco
turso auth login
turso db create habits
turso db show habits      # exibe a URL
turso db tokens create habits  # gera o token
```

## Comandos

### Pilares

Pilares são categorias que agrupam hábitos (ex: Saúde, Finanças, Carreira).

```bash
habits pillar add <nome>
habits pillar list
```

### Hábitos

```bash
# Criar hábito (solicita cor interativamente)
habits habit add <nome> --pillar <id> --days <dias>

# Dias: dom,seg,ter,qua,qui,sex,sab  ou  todos
habits habit add Meditação --pillar 1 --days todos
habits habit add Exercício --pillar 1 --days seg,qua,sex

# Listar
habits habit list
habits habit list --pillar <id>

# Remover
habits habit remove <id>
```

**Cores disponíveis:** `azul`, `verde`, `laranja`, `vermelho`, `amarelo` (ou Enter para aleatório).

### Sessões (check-in / check-out)

Registra o tempo de execução de um hábito ou tarefa.

```bash
# Iniciar sessão
habits start habit <id>
habits start task <id>

# Finalizar sessão em andamento
habits stop
```

Só pode haver uma sessão aberta por vez.

### Lista do dia

Exibe hábitos pendentes para hoje e tarefas em andamento.

```bash
habits list
habits ls        # alias
```

### Timeline e grid

Exibe a timeline do dia com as sessões registradas e o grid de consistência dos hábitos.

```bash
# Timeline de hoje + grid dos últimos 7 dias (padrão)
habits log

# Timeline de outra data
habits log --date 2026-05-30

# Grid dos últimos N dias
habits log --days 30

# Grid de um mês específico
habits log --month 2026-05
```

O grid mostra `■` colorido nos dias em que cada hábito foi executado (ao menos uma sessão finalizada). À direita de cada linha aparece a porcentagem de execução no período, com cor em gradiente: verde para os mais consistentes, vermelho para os menos.

## Exemplo de uso

```bash
# Configurar estrutura
habits pillar add Saúde
habits habit add Meditação --pillar 1 --days todos
habits habit add Exercício --pillar 1 --days seg,qua,sex

# Registrar execução
habits start habit 1
habits stop

# Visualizar
habits log
habits list
```
