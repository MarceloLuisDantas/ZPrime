.text
# Kernel 
#   Syscalls
#       0 - module
#       1 - exit
#       2 - print int16
#       3 - print uint16
#       4 - print int8
#       5 - print uint8
#       6 - print char
#       7 - print string
#       8 - println int16
#       9 - println uint16
#      10 - println int8
#      11 - println uint8
#      12 - println char
#      13 - println string

    j *_main

# __ZPRIME_MOD(number, mod)
#               sp+2,   sp
__ZPRIME_MOD:
    lw $t0, 0($sp)          # t0 = mod
    lw $t1, 2($sp)          # t1 = number

    slt $t2, $t0, $zero     # t2 = mod < 0 (mod negativo)
    slt $t3, $t1, $zero     # t3 = num < 0 (valor negativo)

    # if (number > 0) { __ZPRIME_MOD_VALUE_POSITIVE } else { __ZPRIME_MOD_VALUE_NEGATIVE }
    bne $t3, $zero, *__ZPRIME_MOD_VALUE_NEGATIVE
        __ZPRIME_MOD_VALUE_POSITIVE:
            # if (mod > 0) { __ZPRIME_MOD_VALUE_POSITIVE_MOD_NEGATIVE } else { __ZPRIME_MOD_VALUE_POSITIVE_MOD_NEGATIVE }
            bne $t2, $zero, *__ZPRIME_MOD_VALUE_POSITIVE_MOD_NEGATIVE
                __ZPRIME_MOD_VALUE_POSITIVE_MOD_POSITIVE:
                    __ZPRIME_MOD_VALUE_POSITIVE_MOD_POSITIVE_LOOP:     
                        # if t1 < t0; { break } 
                        blt $t1, $t0, *__ZPRIME_END_MOD
                        sub $t1, $t1, $t0   # t1 -= t0 
                        j *__ZPRIME_MOD_VALUE_POSITIVE_MOD_POSITIVE_LOOP

                __ZPRIME_MOD_VALUE_POSITIVE_MOD_NEGATIVE:
        __ZPRIME_MOD_VALUE_NEGATIVE:
            # if (mod > 0) { __ZPRIME_MOD_VALUE_NEGATIVE_MOD_POSITIVE } else { __ZPRIME_MOD_VALUE_NEGATIVE_MOD_NEGATIVE }
            bne $t2, $zero, *__ZPRIME_MOD_VALUE_NEGATIVE_MOD_NEGATIVE
                __ZPRIME_MOD_VALUE_NEGATIVE_MOD_POSITIVE:

                __ZPRIME_MOD_VALUE_NEGATIVE_MOD_NEGATIVE:

    __ZPRIME_END_MOD:

    move $rt, $t1
    return


# __ZPRIME_MODU(number, mod)
#               sp+2,   sp
__ZPRIME_MODU:
    lw $t0, 0($sp)          # t0 = mod
    lw $t1, 2($sp)          # t1 = number

    slt $t2, $t0, $zero     # t2 = mod < 0 (mod negativo)
    # if (mod > 0) { __ZPRIME_MODU_NEGATIVE_LOOP } else { __ZPRIME_MODU_POSITIVE_LOOP }
    bne $t2, $zero, *__ZPRIME_MODU_NEGATIVE_LOOP
        __ZPRIME_MODU_POSITIVE_LOOP:     
            sltu $t2, $t1, $t0 # t2 = t1 < t0
            bgt $t2, $zero, *__ZPRIME_END_MODU # if ( t2 > 0 ) { __ZPRIME_END_MODU }
            sub $t1, $t1, $t0  # t1 -= t0 
            j *__ZPRIME_MODU_POSITIVE_LOOP
        __ZPRIME_MODU_NEGATIVE_LOOP:
            sltu $t2, $t1, $t0 # t2 = t1 < t0
            bgt $t2, $zero, *__ZPRIME_END_MODU # if ( t2 > 0 ) { __ZPRIME_END_MODU }
            sub $t1, $t1, $t0  # t1 -= t0 
            j *__ZPRIME_MODU_NEGATIVE_LOOP


    __ZPRIME_END_MODU:

    move $rt, $t1
    return

# __ZPRIME_PRINT_U16(number, column_vram, line_vram)
#                    sp+4    sp+2         sp
__ZPRIME_PRINT_U16:
    # Calcula onde na VRAM deve começar a escrever
    lw $t1, 0($sp)      # t1 = line_vram
    addi $sp, $sp, 2    # pop line_vram
    multi $t1, $t1, 60  # t1 = line * 60

    lw $t2, 0($sp)      # t2 = column_vram
    addi $sp, $sp, 2    # pop column_vram
    add $t4, $t1, $t2   # t4 = line * 60 + column | Index incial na VRAM

    lw $t0, 0($sp)      # t0 = number
    addi $sp, $sp, 2    # pop number

    li $t1, 0
    # t0 = numero a ser printado
    # t1 = 0
    # t2 = index na VRAM 

    # push 0 (zero indica que todos os valores foram desempilhados no momento de printar)
    addi $sp, $sp, -2
    sw $zero, 0($sp)

    # while (t0 != 0)
    __ZPRIME_WHILE_STACK_NUMBERS_U16:
        # push ra
        addi $sp, $sp, -2
        sw $ra, 0($sp)
        
        # __ZPRIME_MODU(number, mod)
        #               sp+2,   sp    
        addi $sp, $sp, -4   # cria 2 epaços na stack
        li $t1, 10
        sw $t1, 0($sp)      # push 10 (modulo)
        sw $t0, 2($sp)      # push number    

        jal *__ZPRIME_MODU

        lw $t0, 2($sp)
        addi $sp, $sp, 4    # pop values

        # pop ra
        lw $ra, 0($sp)
        addi $sp, $sp, 2

        move $t1, $rt       # t1 = number % 10
        divui $t0, $t0, 10  # number /= 10

        # Monta o valor para VRAM
        ori $t1, $t1, 48    # t1 or 00110000 = 0011xxxx | valor em ASCII
        sll $t1, $t1, 8     # t1 = 0011xxxx00000000
        ori $t1, $t1, 16    # t1 or 00010000 = 0011xxxx00010000
        # 16 = [0001][0000]
        # [0001] - fonte branca
        # [0000] - fundo preto
        
        # push num
        addi $sp, $sp, -2 
        sw $t1, 0($sp)

        # if t0 <= 0 { break }
        ble $t0, $zero, *__ZPRIME_END_WHILE_STACK_NUMBERS_U16

        j *__ZPRIME_WHILE_STACK_NUMBERS_U16
    __ZPRIME_END_WHILE_STACK_NUMBERS_U16:

    __ZPRIME_UNSTACK_PRINT_NUMBERS:
        # pop
        lw $t0, 0($sp)
        addi $sp, $sp, 2

        # if (t0 == 0) { break }
        beq $t0, $zero, *__ZPRIME_END_UNSTACK_PRINT_NUMBERS

        # print
        svr $t0, 0($t4)
        inc $t4

        j *__ZPRIME_UNSTACK_PRINT_NUMBERS
    __ZPRIME_END_UNSTACK_PRINT_NUMBERS:

    return

__ZPRIME_TEST_PRINT_U16:
    addi $sp, $sp, -2   # adiciona 1 espaço na stack
    sw $ra, 0($sp)      # push ra

    addi $sp, $sp, -6   # Adiciona 3 espaços na stack       
    
    li $t0, 0
    sw $t0, 0($sp)      # line_vram 0
    sw $t0, 2($sp)      # column_vram 0

    li $t0, 1
    sw $t0, 4($sp)      # number

    # print_16(number, column_vram, line_vram)
    #          sp+4    sp+2         sp
    jal *__ZPRIME_PRINT_U16
    
    li $t0, 1
    sw $t0, 0($sp)
    li $t0, 0
    sw $t0, 2($sp)  
    li $t0, 32767
    sw $t0, 4($sp)
    jal *__ZPRIME_PRINT_U16

    li $t0, 2
    sw $t0, 0($sp)
    li $t0, 0
    sw $t0, 2($sp)
    li $t0, 32768
    sw $t0, 4($sp)
    jal *__ZPRIME_PRINT_U16
   
    li $t0, 3
    sw $t0, 0($sp)
    li $t0, 0
    sw $t0, 2($sp)
    li $t0, 65535
    sw $t0, 4($sp)
    jal *__ZPRIME_PRINT_U16

    addi $sp, $sp, 6 # pop values
 
    lw $ra, 0($sp)
    addi $sp, $sp, 2 # pop ra
    return


_main:
    jal *__ZPRIME_TEST_PRINT_U16
    
    li $sc, 100 # render frame
    loop:
        syscall
        j *loop

     