#define HP 16($zero)    # endereço do hp na RAM
#define CURA 10
#define CURA_DUPLA 20
#define PRINT_INT16 1003

.text
    li $t0, CURA
    sb $t0, HP

    lb $t0, HP
    addi $t0, $t0, CURA_DUPLA
    sb $t0, HP

    li $sc, PRINT_INT16
    lb $t0, HP
    syscall # 30

    li $sc, 0
    syscall

.data
