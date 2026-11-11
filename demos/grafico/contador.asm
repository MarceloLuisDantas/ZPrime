# Mostra na tela quantas vezes o programa foi aberto, usando o memory card

# Contador de 16 bits, guardado em 2 bytes no save
#define CONTADOR_ALTO 0
#define CONTADOR_BAIXO 1

# Posições na VRAM (linha * 60 + coluna)
#define POS_TITULO 1578        # linha 26, coluna 18
#define POS_NUMERO 1710        # linha 28, coluna 30 (centro do numero)
#define POS_VEZES 1827         # linha 30, coluna 27
#define POS_SEM_CARD 2002      # linha 33, coluna 22

# [0000][0000] cor da fonte e cor do fundo
#define COR_TEXTO 16           # [0001][0000] fonte branca, fundo preto
#define COR_NUMERO 112         # [0111][0000] fonte amarela, fundo preto
#define COR_AVISO 48           # [0011][0000] fonte vermelha, fundo preto

#define RENDER 100

.text
    j *_main

# Desenha uma string terminada em \0
# t5 = endereço da string na ROM
# t1 = posição na VRAM
# t3 = cor
_draw_string:
    lrb $t0, 0($t5)
    beq $t0, $zero, *end_draw_string
        sll $t0, $t0, 8         # [xxxxxxxx][00000000]
        or $t0, $t0, $t3        # [xxxxxxxx][cor]
        svr $t0, 0($t1)
        inc $t5
        inc $t1
        j *_draw_string
    end_draw_string:
    return

# Desenha um numero sem sinal centralizado
# t0 = numero
# t1 = posição na VRAM do centro do numero
_draw_number:
    # t3 = quantidade de digitos
    move $t2, $t0
    li $t3, 0
    loop_count_digits:
        inc $t3
        divui $t2, $t2, 10
        bne $t2, $zero, *loop_count_digits

    # t1 = posição do ultimo digito = centro - digitos / 2 + digitos - 1
    divui $t2, $t3, 2
    sub $t1, $t1, $t2
    add $t1, $t1, $t3
    dec $t1

    # Desenha os digitos da direita para a esquerda
    loop_draw_digits:
        divui $t2, $t0, 10      # t2 = numero / 10
        multui $t4, $t2, 10
        subu $t4, $t0, $t4      # t4 = numero % 10

        addi $t4, $t4, 48       # digito em ASCII
        sll $t4, $t4, 8         # [digito][00000000]
        ori $t4, $t4, COR_NUMERO
        svr $t4, 0($t1)

        dec $t1
        move $t0, $t2           # numero /= 10
        bne $t0, $zero, *loop_draw_digits

    return

_main:
    request_save *program_name  # rt = save_id, 0 se não houver memory card
    move $t5, $rt               # t5 = save_id
    sb $t5, 16($zero)           # ram[16] = save_id

    # t1 = contador
    li $t0, CONTADOR_ALTO
    read_save $t5, $t0
    sll $t1, $rt, 8
    li $t0, CONTADOR_BAIXO
    read_save $t5, $t0
    or $t1, $t1, $rt

    inc $t1                     # contador++
    sw $t1, 18($zero)           # ram[17..18] = contador

    # Salva o contador
    li $t0, CONTADOR_ALTO
    srl $t2, $t1, 8
    write_save $t5, $t0, $t2
    li $t0, CONTADOR_BAIXO
    write_save $t5, $t0, $t1    # salva só o byte menos significativo

    # Desenha a tela, a VRAM mantem o conteudo entre os frames
    la $t5, *titulo
    li $t1, POS_TITULO
    li $t3, COR_TEXTO
    jal *_draw_string

    lw $t0, 18($zero)
    li $t1, POS_NUMERO
    jal *_draw_number

    la $t5, *vezes
    li $t1, POS_VEZES
    li $t3, COR_TEXTO
    jal *_draw_string

    # Sem memory card o contador não é salvo
    lb $t0, 16($zero)
    bne $t0, $zero, *game_loop
        la $t5, *sem_card
        li $t1, POS_SEM_CARD
        li $t3, COR_AVISO
        jal *_draw_string

    game_loop:
        li $sc, RENDER
        syscall
        j *game_loop

.data
    program_name: .string "contador_grafico\0"
    titulo: .string "VOCE ABRIU ESTE PROGRAMA\0"
    vezes: .string "VEZES\0"
    sem_card: .string "SEM MEMORY CARD\0"
