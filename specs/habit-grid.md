# Spec: Habit Grid

## Descrição
Visualização de consistência dos hábitos ao longo do tempo. Exibida abaixo da timeline no comando `habits log`.

## Layout

- Eixo horizontal: dias
- Eixo vertical: título dos hábitos
- Junção linha/coluna: `■` colorido se o hábito foi executado no dia, vazio se não foi
- Cada hábito tem uma cor associada (escolhida no cadastro)

**Exemplo visual:**
```
          01  02  03  04  05  06  07 ...
Meditação  ■           ■   ■   ■
Exercício      ■   ■       ■
Leitura    ■   ■           ■   ■
```

## Parâmetros

| Parâmetro       | Descrição                                          | Padrão         |
|-----------------|----------------------------------------------------|----------------|
| (nenhum)        | Exibe todos os dias do mês atual                   | mês corrente   |
| `--days <n>`    | Exibe os últimos N dias (ex: `--days 60`)          | —              |
| `--month <mes>` | Exibe todos os dias do mês especificado (YYYY-MM)  | —              |

## Cores dos hábitos

Definidas no cadastro do hábito. Cores disponíveis: `azul`, `verde`, `laranja`, `vermelho`, `amarelo`.

Se não informada, uma cor é escolhida aleatoriamente dentre as disponíveis.

### Seleção de cor no cadastro
Ao criar um hábito, um prompt interativo é exibido:
```
Escolha uma cor [azul, verde, laranja, vermelho, amarelo] (Enter para aleatório):
```

## Regras de negócio
- Um hábito é considerado "executado no dia" se existir ao menos uma Session finalizada vinculada a ele naquele dia.
- A cor é um atributo do Habit, independente da interface de uso.
- A lógica de determinar quais dias um hábito foi executado pertence ao domínio.

## Camadas

| Responsabilidade                                      | Camada       |
|-------------------------------------------------------|--------------|
| Determinar se hábito foi executado num dia            | Domínio      |
| Buscar Sessions por hábito e intervalo de datas       | Domínio      |
| Escolha e validação de cores disponíveis              | Domínio      |
| Cor aleatória quando não informada                    | Domínio      |
| Renderização do grid com caracteres e cores ANSI      | Apresentação |
| Prompt interativo de seleção de cor                   | Apresentação |
| Persistência das Sessions e Habits no SQLite          | Infra        |
