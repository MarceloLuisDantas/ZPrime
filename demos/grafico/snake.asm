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

_draw_border:
    li $t5, 5 # [00000000][0000][0101] '║' fonte cinza fundo preto
    li $t0, 0 # border.y = 0
    li $t1, 60 # last pixel first line
    loop_draw_border:
        beq $t0, $t1, *end_draw_border
            svr $t5, 0($t0) 
            multi $t0, $t0, 60
            svr $t5, 0($t0) 
            addi $t0, $t0, 59
            svr $t5, 0($t0) 
            subi $t0, $t0, 59
            divi $t0, $t0, 60
            addi $t0, $t0, 3540
            svr $t5, 0($t0) 
            subi $t0, $t0, 3540
            inc $t0
            j *loop_draw_border
    end_draw_border:
    return

_draw_player:
    # 17 = [00000000][0000][0001] fundo branco
    li $t4, 1

    li $t2, 25        # t2 = first segment index
    lb $t3, 20($zero) # t3 = player lenght

    _draw_segments:
        # if len == 0 { break }
        beq $t3, $zero, *_end_draw_segments    
            lb $t0, 1($t2) # t0 = segment.y
            # li $sc, 1004
            # syscall

            # player y *= 60 | valor na vram conrespondente a linha
            multi $t0, $t0, 60 
            lb $t1, 0($t2) # t1 = segment.x
            # t0 = linha + coluna | valor na vram a [y][x] 
            add $t0, $t0, $t1 
            svr $t4, 0($t0) # draw segment
            dec $t3 # lenght -= 1
            addi $t2, $t2, 2
            j *_draw_segments
    _end_draw_segments:

    # 17 = [00000000][0000][0100] fundo verde
    li $t4, 4

    _draw_head:
        lb $t0, 17($zero)  # t0 = head.y
        # player y *= 60 | valor na vram conrespondente a linha
        multi $t0, $t0, 60 
        lb $t1, 16($zero)  # t1 = head.x
        # t0 = linha + coluna | valor na vram a [y][x] 
        add $t0, $t0, $t1 
        svr $t4, 0($t0) # draw player
    _end_draw_head:

    return

_draw_fruit:
    lw $t0, 22($zero) # fruit position
    li $t1, 3         # [00000000][0000][0011] fundo vermelho

    svr $t1, 0($t0)
    return

_move:
    _move_segments:
        lb $t0, 16($zero) # t0 = player.x
        lb $t1, 17($zero) # t1 = player.y
        lb $t4, 20($zero) # t4 = len
        li $t5, 25 # first segment position

        loop_move_segments:
            # if len == 0 { break }
            beq $t4, $zero, *end_move_segments

            # save the x and y of the current segment in t2 and t0
            lb $t2, 0($t5) # t2 = *(t5)
            lb $t3, 1($t5) # t3 = *(t5+1)          

            # set current segment x and y to the last segment x and y                        
            sb $t0, 0($t5) # current_segment.x = prev_seg.x
            sb $t1, 1($t5) # current_segment.y = prev_seg.y

            addi $t5, $t5, 2 # set pivot to the next segment
            dec $t4          # len -= 1

            move $t0, $t2 
            move $t1, $t3
            j *loop_move_segments
        end_move_segments:
    _end_move_segments:
    
    lb $t0, 18($zero) # t0 = direction
    lb $t2, 16($zero) # t2 = x
    lb $t3, 17($zero) # t3 = y

    # if player.direction == up { }
    li $t1, 0
    bne $t0, $t1, *else_up
        # if $t0 > 1 { player.y -= 1 }
        li $t1, 1
        ble $t3, $t1, *set_dead
            # player y -= 1
            dec $t3           
            sb $t3, 17($zero) 
            return
    else_up:

    # if player.direction == right { }
    li $t1, 1
    bne $t0, $t1, *else_right
        # if player.x < 58 { player.x += 1 }
        li $t1, 58
        bge $t2, $t1, *set_dead
            # player x += 1
            inc $t2           
            sb $t2, 16($zero) 
            return
    else_right:

    # if player.direction == down { }
    li $t1, 2
    bne $t0, $t1, *else_down
        li $t1, 58
        bge $t3, $t1, *set_dead
            inc $t3           # player y += 1
            sb $t3, 17($zero) # salve player y
            return
    else_down:

    # if player.direction == left { }
    li $t1, 3
    bne $t0, $t1, *else_left
        # if player.x > 0 { player.x -= 1 }
        li $t1, 1
        beq $t2, $t1, *set_dead
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
    
    mod_loop:
        # if t0 < 3600 { break }
        sltui $t1, $t0, 3600 # t1 = t0 < 3600
        bne $t1, $zero, *mod_end
            subi $t0, $t0, 3600
            j *mod_loop
    mod_end:

    # fruit position = t1
    sw $t0, 22($zero)

    # fruit = alive
    li $t0, 1
    sb $t0, 21($zero) 
    return

_new_segment:
    lb $t0, 20($zero)    # t0 = lenght    
    li $t1, 25           # first segment on ram
    li $t2, 2            # size of each segment
    mult $t2, $t2, $t0   # size * lenght,  
    add $t1, $t1, $t2    # new segment position
    sb $zero, 0($t1) # new_segment.x = 0
    sb $zero, 1($t1) # new_segment.y = 0
    return

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
    li $t0, 0 # player length
    sb $t0, 20($zero)

    # ram[21] = fruit
    li $t0, 0 # fruit don't exit
    sb $t0, 21($zero) 

    # ram[22][23] = fruit position
    li $t0, 0 
    sw $t0, 22($zero) 

    # ram[24] = wait frames
    li $t0, 2
    sb $t0, 24($zero) 

    # ram[25..535] = segments x and y (255)
    li $t0, 255
    
    
    game_loop:
        _start_logic:
            _startd_moviment_logic:

                # if player_alive { move }
                lb $t0, 19($zero)
                beq $t0, $zero, *game_over
                    lb $t1, 18($zero) # player poistion 
                    
                    lb $t0, 0($zero) # key up
                    beq $t0, $zero, *else_set_up
                        # if direction != down { set down }
                        li $t2, 2
                        beq $t1, $t2, *else_set_up
                            li $t0, 0
                            sb $t0, 18($zero)
                            j *go_move
                    else_set_up:

                    lb $t0, 1($zero) # key down
                    beq $t0, $zero, *else_set_down
                        # if direction != up { set down }
                        li $t2, 0
                        beq $t1, $t2, *else_set_down
                            li $t0, 2
                            sb $t0, 18($zero)
                            j *go_move
                    else_set_down:

                    lb $t0, 2($zero) # key left
                    beq $t0, $zero, *else_set_left
                        # if direction != right { set down }
                        li $t2, 1
                        beq $t1, $t2, *else_set_left
                            li $t0, 3
                            sb $t0, 18($zero)
                            j *go_move
                    else_set_left:

                    lb $t0, 3($zero) # key right
                    beq $t0, $zero, *else_set_right
                        # if direction != left { set down }
                        li $t2, 3
                        beq $t1, $t2, *else_set_right
                            li $t0, 1
                            sb $t0, 18($zero)
                            j *go_move
                    else_set_right:

                    go_move:
                        lb $t0, 24($zero)
                        bne $t0, $zero, *_waiting_frames
                            li $t0, 2
                            sb $t0, 24($zero)
                            jal *_move
                            j *_end_logic
                        _waiting_frames:
                            lb $t0, 24($zero)
                            dec $t0
                            sb $t0, 24($zero)

            _end_movimento_logic:

            _start_fruit_colision:
                lb $t0, 17($zero)  # t0 = player.y                
                multi $t0, $t0, 60 # line in the vram
                lb $t1, 16($zero)  # t1 = player.x
                add $t0, $t0, $t1  # t0 = player position in the vram
                lw $t1, 22($zero)  # t1 = fruit position
                # if fruit.position == player.position { }
                bne $t0, $t1, *_end_fruit_colision
                    lb $t0, 20($zero) # t0 = player.lenght
                    inc $t0           # t0 += 1
                    li $t1, 255       # t1 = 255 (max lenght)
                    # if player.lenght == 255 { win }
                    bne $t0, $t1, *not_win
                        sb $zero, 19($zero) # player win
                        j *game_over
                    not_win:
                        sb $t0, 20($zero) # player.lenght = t0
                        sb $zero, 21($zero) # despawn fruit
                        jal *_new_segment
            _end_fruit_colision:

            _start_spawn_fuit_logic:
                # if fuit_exit == 0 { spawn_fuit }
                lb $t0, 21($zero)
                bne $t0, $zero, *_end_spawn_fuit_logic
                    jal *_spawn_fuit
                    inc $t0
                    sb $t0, 21($zero)
            _end_spawn_fuit_logic:

            game_over:
                j *_end_logic
                    
            
        _end_logic:

        _start_drawing:
            jal *_clear_scream
            jal *_draw_fruit
            jal *_draw_player
            jal *_draw_border

            # render frame
            li $sc, 100
            syscall
        _end_drawing:

        j *game_loop
    


.data
