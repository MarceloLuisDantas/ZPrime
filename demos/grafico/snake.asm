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


_draw_player:
    lb $t0, 17($zero)  # t1 = player y
    # player y *= 60 | valor na vram conrespondente a linha
    multi $t0, $t0, 60 
    
    lb $t1, 16($zero)  # t0 = player x
    # t0 = linha + coluna | valor na vram a [y][x] 
    add $t0, $t0, $t1 

    # 17 = [00000000][0001][0001] sem character, fonte branca, fundo branco
    li $t1, 17 

    svr $t1, 0($t0) # draw player

    return

_move:
    lb $t0, 18($zero) # t0 = direction
    lb $t2, 16($zero) # t2 = x
    lb $t3, 17($zero) # t3 = y

    # if player.direction == up { }
    li $t1, 0
    bne $t0, $t1, *else_up
        # if $t0 > 0 { player.y -= 1 }
        ble $t3, $zero, *set_dead
            # player y -= 1
            dec $t3           
            sb $t3, 17($zero) 
            return
    else_up:

    # if player.direction == right { }
    li $t1, 1
    bne $t0, $t1, *else_right
        # if player.x < 59 { player.x += 1 }
        li $t1, 59
        bge $t2, $t1, *set_dead
            # player x += 1
            inc $t2           
            sb $t2, 16($zero) 
            return
    else_right:

    # if player.direction == down { }
    li $t1, 2
    bne $t0, $t1, *else_down
        li $t1, 59
        bge $t3, $t1, *set_dead
            inc $t3           # player y += 1
            sb $t3, 17($zero) # salve player y
            return
    else_down:

    # if player.direction == left { }
    li $t1, 3
    bne $t0, $t1, *else_left
        # if player.x > 0 { player.x -= 1 }
        beq $t2, $zero, *set_dead
            # player x -= 1
            dec $t2           
            sb $t2, 16($zero)
            return
    else_left:

    set_dead:
        sb $zero, 19($zero)
        return

_spawn_fuit:
    rand $t0 # random number 

_main:
    # ram[16] = player x
    li $t0, 5 # player x
    sb $t0, 16($zero) 

    # ram[17] = player y
    li $t0, 5 # player y
    sb $t0, 17($zero) 

    # ram[18] = player direction
    li $t0, 1 # player direction
    sb $t0, 18($zero)

    # ram[19] = player alive
    li $t0, 1 # player alive
    sb $t0, 19($zero)

    # ram[20] = player length
    li $t0, 1 # player length
    sb $t0, 20($zero)

    # ram[21] = fruit
    li $t0, 0 # fruit don't exit
    sb $t0, 21($zero) 
    
    game_loop:
        _start_logic:
            _startd_moviment_logic:
                # if player_alive { move }
                lb $t0, 19($zero)
                beq $t0, $zero, *game_over
                    lb $t0, 0($zero) # key up
                    beq $t0, $zero, *else_set_up
                        li $t0, 0
                        sb $t0, 18($zero)
                        j *go_move
                    else_set_up:

                    lb $t0, 1($zero) # key down
                    beq $t0, $zero, *else_set_down
                        li $t0, 2
                        sb $t0, 18($zero)
                        j *go_move
                    else_set_down:

                    lb $t0, 2($zero) # key left
                    beq $t0, $zero, *else_set_left
                        li $t0, 3
                        sb $t0, 18($zero)
                        j *go_move
                    else_set_left:

                    lb $t0, 3($zero) # key right
                    beq $t0, $zero, *else_set_right
                        li $t0, 1
                        sb $t0, 18($zero)
                        j *go_move
                    else_set_right:

                    go_move:
                        jal *_move

            _end_movimento_logic:

            # _start_spawn_fuit_logic:
            #     # if fuit not exit { spawn_fuit }
            #     lb $t0, 21($zero)
            #     bne $t0, $zero, *_end_spawn_fuit_logic
            #         jal *_spawn_fuit
            #         inc $t0
            #         sb $t0, 21($zero)
            # _end_spawn_fuit_logic:


                j *_end_logic

            game_over:
                j *_end_logic
                    
            
        _end_logic:

        _start_drawing:
            jal *_clear_scream
            jal *_draw_player

            # render frame
            li $sc, 100
            syscall
        _end_drawing:

        j *game_loop
    


.data
