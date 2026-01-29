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

# Instruções Aritmeticas
add
addu
addi
addui
sub
subu
subi
subui
mult
multu
multi
multui
div
divu
divi
divui
and
andi
or
ori
slt
sltu
slti
sltui
move
li
la
j
jr
jal
return
beq
bne
bgt
bge
blt
ble
inc
dec
rand
lw
lb
sw
sb
lrw
lrb
lvr
svr
sll
srl
