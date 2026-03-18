# Registradores
| Nome      | Uso                 |
|-----------|---------------------|
| $zero     | Sempre retorna zero |
| $t0 ~ $t5 | Uso temporario, suponha que sempre vão ser alterados ao chamar funções |
| $rt       | Valor de retorno das funções |
| $sc       | Setar qual syscall deseja ser chamada |
| $sp       | Stack Point |
| $sf       | Stack Frame |
| $gp       | Indica onde o espaço de data começa na ROM |
| $ir       | Guarda a instrução a ser executada |
| $pc       | Program counter |   
| RAM 0x01  | KEY UP |
| RAM 0x00  | KEY UP |
| RAM 0x01  | KEY DOWN |
| RAM 0x02  | KEY LEFT |
| RAM 0x03  | KEY RIGHT |
| RAM 0x04  | KEY SPACE |
| RAM 0x05  | KEY ENTER |
| RAM 0x06  | KEY BACKSPACE |
| RAM 0x07  | KEY W |
| RAM 0x08  | KEY A |
| RAM 0x09  | KEY S |
| RAM 0x0A  | KEY D |
| RAM 0x0B  | KEY Q |
| RAM 0x0C  | KEY E |
| RAM 0x0D  | KEY I |
| RAM 0x0E  | KEY O |
| RAM 0x0F  | KEY P |

# Instruções 
## Instruções Aritméticas
- `add dest, src1, src2` — Soma src1 e src2, armazena em dest
- `addi dest, src, value` — Soma src e valor imediato, armazena em dest
- `addu dest, src1, src2` — Soma unsigned src1 e src2, armazena em dest
- `addui dest, src, value` — Soma unsigned src e valor imediato, armazena em dest
- `sub dest, src1, src2` — Subtrai src2 de src1, armazena em dest
- `subi dest, src, value` — Subtrai valor imediato de src, armazena em dest
- `subu dest, src1, src2` — Subtrai unsigned src2 de src1, armazena em dest
- `subui dest, src, value` — Subtrai unsigned valor imediato de src, armazena em dest
- `mult dest, src1, src2` — Multiplica src1 por src2, armazena em dest
- `multi dest, src, value` — Multiplica src por valor imediato, armazena em dest
- `multu dest, src1, src2` — Multiplica unsigned src1 por src2, armazena em dest
- `multui dest, src, value` — Multiplica unsigned src por valor imediato, armazena em dest
- `div dest, src1, src2` — Divide src1 por src2, armazena em dest
- `divi dest, src, value` — Divide src por valor imediato, armazena em dest
- `divu dest, src1, src2` — Divide unsigned src1 por src2, armazena em dest
- `divui dest, src, value` — Divide unsigned src por valor imediato, armazena em dest

## Instruções Lógicas
- `and dest, src1, src2` — AND bit a bit entre src1 e src2
- `andi dest, src, value` — AND bit a bit entre src e valor imediato
- `or dest, src1, src2` — OR bit a bit entre src1 e src2
- `ori dest, src, value` — OR bit a bit entre src e valor imediato

## Comparação
- `slt dest, src1, src2` — dest = 1 se src1 < src2, senão 0
- `slti dest, src, value` — dest = 1 se src < value, senão 0
- `sltu dest, src1, src2` — unsigned, dest = 1 se src1 < src2
- `sltui dest, src, value` — unsigned, dest = 1 se src < value

## Movimentação e Carga
- `move dest, src` — Copia valor de src para dest
- `li dest, value` — Carrega valor imediato em dest
- `la dest, value` — Carrega endereço imediato em dest

## Saltos e Controle
- `j point` — Salta para endereço point
- `jal point` — Salta para endereço point, salvando retorno
- `jr src` — Salta para endereço em src
- `return` — Retorna para endereço salvo
- `beq src1, src2, point` — Salta se src1 == src2
- `bne src1, src2, point` — Salta se src1 != src2
- `bgt src1, src2, point` — Salta se src1 > src2
- `bge src1, src2, point` — Salta se src1 >= src2
- `blt src1, src2, point` — Salta se src1 < src2
- `ble src1, src2, point` — Salta se src1 <= src2

## Incremento e Decremento
- `inc dest` — Incrementa dest
- `dec dest` — Decrementa dest

## Random
- `rand dest` — Gera número aleatório e armazena em dest

## Shift
- `sll dest, src, value` — Shift left lógico
- `srl dest, src, value` — Shift right lógico

## RAM
- `sb src, offset, reg_index` — Salva byte de src na RAM
- `sw src, offset, reg_index` — Salva word de src na RAM
- `lb dest, offset, reg_index` — Carrega byte da RAM para dest
- `lw dest, offset, reg_index` — Carrega word da RAM para dest

## VRAM (Gráficos)
- `svr src, offset, reg_index` — Salva valor de src na VRAM
- `lvr dest, offset, reg_index` — Carrega valor da VRAM para dest

## ROM
- `lrb dest, offset, point` — Carrega byte da ROM para dest
- `lrw dest, offset, point` — Carrega word da ROM para dest

## Entrada (Teclado)
Os registradores de RAM de 0x0000 a 0x000F representam teclas:
- 0x0000: KEY UP
- 0x0001: KEY DOWN
- 0x0002: KEY LEFT
- 0x0003: KEY RIGHT
- 0x0004: KEY SPACE
- 0x0005: KEY ENTER
- 0x0006: KEY BACKSPACE
- 0x0007: KEY W
- 0x0008: KEY A
- 0x0009: KEY S
- 0x000A: KEY D
- 0x000B: KEY Q
- 0x000C: KEY E
- 0x000D: KEY I
- 0x000E: KEY O
- 0x000F: KEY P


