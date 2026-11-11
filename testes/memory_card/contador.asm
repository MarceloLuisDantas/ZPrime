# Conta quantas vezes o programa foi executado
#define CONTADOR_SAVE 100 # endereço do contador no save
#define PRINT_INT16 1003

.text
    request_save *program_name # rt = save_id
    move $t1, $rt              # t1 = save_id

    li $t0, CONTADOR_SAVE
    read_save $t1, $t0         # rt = contador
    inc $rt                    # contador++
    write_save $t1, $t0, $rt   # atualiza o contador no save

    move $t0, $rt
    li $sc, PRINT_INT16
    syscall

    li $sc, 0
    syscall

.data
    program_name: .string "contador\0"
