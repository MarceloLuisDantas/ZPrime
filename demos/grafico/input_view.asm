.text
    j *_main

_clear_scream:
    li $t2, 3600 # pixels na tela 

    # 0 = [00000000][0000][0000] tudo preto
    li $t1, 0

    li $t0, 0
    loop_cleaning:
        beq $t0, $t2, *cleaning_end
        svr $t1, 0($t0)
        inc $t0
        j *loop_cleaning
    cleaning_end:

    return

_draw_controler_border:
    li $t0, 32 # largura
    li $t1, 1154 # start
    li $t5, 1 # fundo branco
    loop_draw_top_bottom:
        svr $t5, 0($t1)
        svr $t5, 840($t1)
        inc $t1
        dec $t0
        bne $t0, $zero, *loop_draw_top_bottom

    li $t0, 15 # altura
    li $t1, 1153 # start
    loop_draw_left_right:
        svr $t5, 0($t1)
        svr $t5, 32($t1)
        addi $t1, $t1, 60
        dec $t0
        bne $t0, $zero, *loop_draw_left_right
    return

_draw_inputs:
    li $t5, 5 # fundo cinza, buttom off
    _draw_up:
        li $t0, 1340
        lb $t1, 0($zero)
        lb $t2, 7($zero)
        or $t1, $t1, $t2
        beq $t1, $zero, *_draw_up_off
            li $t5, 4# fundo verde, buttom on
        _draw_up_off:
        svr $t5, 0($t0)
        svr $t5, 59($t0)
        svr $t5, 60($t0)
        svr $t5, 61($t0)

    li $t5, 5 # fundo cinza, buttom off
    _draw_down:
        li $t0, 1760
        lb $t1, 1($zero)
        lb $t2, 9($zero)
        or $t1, $t1, $t2
        beq $t1, $zero, *_draw_down_off
            li $t5, 4 # fundo verde, buttom on
        _draw_down_off:
        svr $t5, 0($t0)
        svr $t5, -1($t0)
        svr $t5, 1($t0)
        svr $t5, 60($t0)

    li $t5, 5 # fundo cinza, buttom off
    _draw_left:
        li $t0, 1517
        lb $t1, 2($zero)
        lb $t2, 8($zero)
        or $t1, $t1, $t2
        beq $t1, $zero, *_draw_left_off
            li $t5, 4 # fundo verde, buttom on
        _draw_left_off:
        svr $t5, 0($t0)
        svr $t5, 60($t0)
        svr $t5, 59($t0)
        svr $t5, 120($t0)

    li $t5, 5 # fundo cinza, buttom off
    _draw_right:
        li $t0, 1523
        lb $t1, 3($zero)
        lb $t2, 10($zero)
        or $t1, $t1, $t2
        beq $t1, $zero, *_draw_right_off
            li $t5, 4 # fundo verde, buttom on
        _draw_right_off:
        svr $t5, 0($t0)
        svr $t5, 60($t0)
        svr $t5, 61($t0)
        svr $t5, 120($t0)

    li $t5, 5 # fundo cinza, buttom off
    _draw_select:
        li $t0, 1467
        lb $t1, 6($zero)
        beq $t1, $zero, *_draw_select_off
            li $t5, 4 # fundo verde, buttom on
        _draw_select_off:
        svr $t5, 0($t0)
        svr $t5, 1($t0)

    li $t5, 5 # fundo cinza, buttom off
    _draw_enter:
        li $t0, 1470
        lb $t1, 5($zero)
        beq $t1, $zero, *_draw_enter_off
            li $t5, 4 # fundo verde, buttom on
        _draw_enter_off:
        svr $t5, 0($t0)
        svr $t5, 1($t0)

    li $t5, 5 # fundo cinza, buttom off
    _draw_i:
        li $t0, 1715
        lb $t1, 13($zero)
        beq $t1, $zero, *_draw_i_off
            li $t5, 4 # fundo verde, buttom on
        _draw_i_off:
        svr $t5, 0($t0)
        svr $t5, 60($t0)
        svr $t5, 61($t0)
        svr $t5, 59($t0)
        svr $t5, 120($t0)

    li $t5, 5 # fundo cinza, buttom off
    _draw_o:
        li $t0, 1538
        lb $t1, 14($zero)
        beq $t1, $zero, *_draw_o_off
            li $t5, 4 # fundo verde, buttom on
        _draw_o_off:
        svr $t5, 0($t0)
        svr $t5, 60($t0)
        svr $t5, 61($t0)
        svr $t5, 59($t0)
        svr $t5, 120($t0)

    li $t5, 5 # fundo cinza, buttom off
    _draw_p:
        li $t0, 1361
        lb $t1, 15($zero)
        beq $t1, $zero, *_draw_p_off
            li $t5, 4 # fundo verde, buttom on
        _draw_p_off:
        svr $t5, 0($t0)
        svr $t5, 60($t0)
        svr $t5, 61($t0)
        svr $t5, 59($t0)
        svr $t5, 120($t0)

    li $t5, 5 # fundo cinza, buttom off
    _draw_space:
        li $t0, 1648
        lb $t1, 4($zero)
        beq $t1, $zero, *_draw_space_off
            li $t5, 4 # fundo verde, buttom on
        _draw_space_off:
        svr $t5, 0($t0)
        svr $t5, 1($t0)
        svr $t5, 2($t0)

    li $t5, 3 # fundo vermelho, buttom off
    _draw_q:
        li $t0, 1096
        lb $t1, 11($zero)
        beq $t1, $zero, *_draw_q_off
            li $t5, 4 # fundo verde, buttom on
        _draw_q_off:
        svr $t5, 0($t0)
        svr $t5, 1($t0)
        svr $t5, 2($t0)
        svr $t5, 3($t0)
        svr $t5, 4($t0)

    li $t5, 3 # fundo vermelho, buttom off
    _draw_e:
        li $t0, 1119
        lb $t1, 12($zero)
        beq $t1, $zero, *_draw_e_off
            li $t5, 4 # fundo verde, buttom on
        _draw_e_off:
        svr $t5, 0($t0)
        svr $t5, 1($t0)
        svr $t5, 2($t0)
        svr $t5, 3($t0)
        svr $t5, 4($t0)


    return

_main:
    game_loop:
        jal *_clear_scream
        jal *_draw_controler_border
        jal *_draw_inputs


        # render
        li $sc, 100
        syscall

    j *game_loop

.data