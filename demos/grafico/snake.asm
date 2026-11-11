# Teclas (RAM)
#define KEY_UP 0($zero)
#define KEY_DOWN 1($zero)
#define KEY_LEFT 2($zero)
#define KEY_RIGHT 3($zero)
#define KEY_ENTER 5($zero)

# Estado do jogo (RAM)
#define PLAYER_X 16($zero)
#define PLAYER_Y 17($zero)
#define PLAYER_DIRECTION 18($zero)
#define PLAYER_ALIVE 19($zero)
#define PLAYER_LENGTH 20($zero)
#define FRUIT_EXISTS 21($zero)
# Word, ocupa os bytes 21 e 22. Escrever a posição da fruta também escreve em
# FRUIT_EXISTS o byte menos significativo da posição.
#define FRUIT_POSITION 22($zero)
#define WAIT_FRAMES 24($zero)
#define GAME_STATE 536($zero)
#define PAUSE_DELAY 537($zero)
#define SAVE_ID 538($zero)
#define SCORE_SAVED 539($zero)   # 1 se o score do game over atual já foi salvo
#define HIGHSCORES 540           # RAM 540..549, copia da lista de highscores do save

# Segmentos do corpo (RAM 25..534), x e y de cada segmento
#define SEGMENTS 25
#define SEGMENT_X 0
#define SEGMENT_Y 1
#define SEGMENT_SIZE 2
#define MAX_LENGTH 255

# Direções
#define DIR_UP 0
#define DIR_RIGHT 1
#define DIR_DOWN 2
#define DIR_LEFT 3

# Estados do jogo
#define STATE_RUNNING 0
#define STATE_START_MENU 1
#define STATE_PAUSE_MENU 2

#define START_X 5
#define START_Y 5
#define MOVE_DELAY 2         # frames de espera entre cada movimento
#define PAUSE_COOLDOWN 15    # frames antes do enter pausar/despausar de novo

# Tela
#define SCREEN_WIDTH 60
#define SCREEN_PIXELS 3600
#define LAST_COLUMN 59
#define LAST_LINE_START 3540 # linha 59, coluna 0
#define FIRST_INSIDE 1       # primeira linha/coluna dentro da borda
#define LAST_INSIDE 58       # ultima linha/coluna dentro da borda

# Deslocamento na VRAM até o pixel vizinho
#define PIXEL_UP -60
#define PIXEL_RIGHT 1
#define PIXEL_DOWN 60
#define PIXEL_LEFT -1

# Cores, [00000000][fonte 4][fundo 4]
#define COLOR_EMPTY 0        # tudo preto
#define COLOR_BODY 1         # fundo branco
#define COLOR_GAME_OVER 2    # fundo azul
#define COLOR_FRUIT 3        # fundo vermelho
#define COLOR_HEAD 4         # fundo verde
#define COLOR_BORDER 5       # fundo cinza
#define COLOR_MENU_TEXT 32   # [0010][0000] fonte azul e fundo preto

# Textos
#define CHAR_SPACE 32
#define MENU_POSITION 1643   # linha 27, coluna 23
#define MENU_TEXT_LENGTH 14
#define GAME_OVER_POSITION 183 # linha 3, coluna 3
#define GAME_OVER_WIDTH 23
#define GAME_OVER_LINES 11

# Highscores, lista ordenada do maior para o menor
#define HIGHSCORE_COUNT 10
#define LAST_HIGHSCORE 9
#define SAVE_HIGHSCORES 0        # endereço da lista no save
#define HIGHSCORES_TITLE_POSITION 985 # linha 16, coluna 25
#define HIGHSCORES_FIRST_LINE 18
#define RANK_COLUMN 26           # coluna do ultimo digito da posição
#define DOT_COLUMN 27
#define SCORE_COLUMN 32          # coluna do ultimo digito do score
#define COLOR_TEXT 16            # [0001][0000] fonte branca e fundo preto
#define CHAR_ZERO 48
#define CHAR_DOT 46

#define RENDER 100

.text
    j *_main

_clear_scream:
    li $t2, SCREEN_PIXELS

    li $t1, COLOR_EMPTY

    li $t0, 0
    loop_cleaning:
        beq $t0, $t2, *cleaning_end
        svr $t1, 0($t0)
        inc $t0
        j *loop_cleaning
    cleaning_end:

    return

_draw_border:
    li $t5, COLOR_BORDER
    li $t0, 0 # border.y = 0
    li $t1, SCREEN_WIDTH
    loop_draw_border:
        beq $t0, $t1, *end_draw_border
            svr $t5, 0($t0)                 # linha 0, coluna t0
            multi $t0, $t0, SCREEN_WIDTH
            svr $t5, 0($t0)                 # linha t0, coluna 0
            addi $t0, $t0, LAST_COLUMN
            svr $t5, 0($t0)                 # linha t0, ultima coluna
            subi $t0, $t0, LAST_COLUMN
            divi $t0, $t0, SCREEN_WIDTH
            addi $t0, $t0, LAST_LINE_START
            svr $t5, 0($t0)                 # ultima linha, coluna t0
            subi $t0, $t0, LAST_LINE_START
            inc $t0
            j *loop_draw_border
    end_draw_border:
    return

_draw_player:
    li $t4, COLOR_BODY

    li $t2, SEGMENTS       # t2 = first segment index
    lb $t3, PLAYER_LENGTH  # t3 = player lenght

    _draw_segments:
        # if len == 0 { break }
        beq $t3, $zero, *_end_draw_segments
            lb $t0, SEGMENT_Y($t2) # t0 = segment.y

            # player y *= 60 | valor na vram conrespondente a linha
            multi $t0, $t0, SCREEN_WIDTH
            lb $t1, SEGMENT_X($t2) # t1 = segment.x
            # t0 = linha + coluna | valor na vram a [y][x]
            add $t0, $t0, $t1
            svr $t4, 0($t0) # draw segment
            dec $t3 # lenght -= 1
            addi $t2, $t2, SEGMENT_SIZE
            j *_draw_segments
    _end_draw_segments:

    li $t4, COLOR_HEAD

    _draw_head:
        lb $t0, PLAYER_Y  # t0 = head.y
        # player y *= 60 | valor na vram conrespondente a linha
        multi $t0, $t0, SCREEN_WIDTH
        lb $t1, PLAYER_X  # t1 = head.x
        # t0 = linha + coluna | valor na vram a [y][x]
        add $t0, $t0, $t1
        svr $t4, 0($t0) # draw player
    _end_draw_head:

    return

_draw_start_menu:
    la $t5, *press_start
    li $t1, MENU_POSITION
    li $t2, MENU_TEXT_LENGTH

    loop_draw_start_menu:
        beq $t2, $zero, *end_draw_start_menu
            lrb $t0, 0($t5)
            sll $t0, $t0, 8  # [xxxxxxxx][00000000]
            ori $t0, $t0, COLOR_MENU_TEXT
            svr $t0, 0($t1)
            inc $t5
            inc $t1
            dec $t2

            j *loop_draw_start_menu
    end_draw_start_menu:

    return

_draw_pause_menu:
    la $t5, *resume_game
    li $t1, MENU_POSITION
    li $t2, MENU_TEXT_LENGTH

    loop_draw_pause_menu:
        beq $t2, $zero, *end_draw_pause_menu
            lrb $t0, 0($t5)
            sll $t0, $t0, 8  # [xxxxxxxx][00000000]
            ori $t0, $t0, COLOR_MENU_TEXT
            svr $t0, 0($t1)
            inc $t5
            inc $t1
            dec $t2

            j *loop_draw_pause_menu
    end_draw_pause_menu:

    return

_draw_game_over:
    la $t5, *game_over_text_1
    li $t1, GAME_OVER_POSITION
    li $t2, GAME_OVER_WIDTH
    li $t3, GAME_OVER_LINES

    loop_draw_game_over:
        beq $t3, $zero, *end_draw_game_over
            loop_draw_game_over_line:
                beq $t2, $zero, *end_draw_game_over_line
                    lrb $t0, 0($t5)
                    li $t4, CHAR_SPACE
                    beq $t0, $t4, *draw_space
                        li $t0, COLOR_GAME_OVER
                        svr $t0, 0($t1)
                    draw_space:

                    inc $t5
                    inc $t1
                    dec $t2
                    j *loop_draw_game_over_line
            end_draw_game_over_line:

            li $t2, GAME_OVER_WIDTH
            subi $t1, $t1, GAME_OVER_WIDTH # back start
            addi $t1, $t1, SCREEN_WIDTH    # next line
            dec $t3
            j *loop_draw_game_over
    end_draw_game_over:

    return

_draw_fruit:
    lw $t0, FRUIT_POSITION
    li $t1, COLOR_FRUIT

    svr $t1, 0($t0)
    return

_move:
    _move_segments:
        lb $t0, PLAYER_X      # t0 = player.x
        lb $t1, PLAYER_Y      # t1 = player.y
        lb $t4, PLAYER_LENGTH # t4 = len
        li $t5, SEGMENTS      # first segment position

        loop_move_segments:
            # if len == 0 { break }
            beq $t4, $zero, *end_move_segments

            # save the x and y of the current segment in t2 and t0
            lb $t2, SEGMENT_X($t5) # t2 = *(t5)
            lb $t3, SEGMENT_Y($t5) # t3 = *(t5+1)

            # set current segment x and y to the last segment x and y
            sb $t0, SEGMENT_X($t5) # current_segment.x = prev_seg.x
            sb $t1, SEGMENT_Y($t5) # current_segment.y = prev_seg.y

            addi $t5, $t5, SEGMENT_SIZE # set pivot to the next segment
            dec $t4                     # len -= 1

            move $t0, $t2
            move $t1, $t3
            j *loop_move_segments
        end_move_segments:
    _end_move_segments:

    lb $t0, PLAYER_DIRECTION # t0 = direction
    lb $t2, PLAYER_X         # t2 = x
    lb $t3, PLAYER_Y         # t3 = y

    li $t5, SCREEN_WIDTH # t5 = position
    mult $t5, $t5, $t3
    add $t5, $t5, $t2

    # if player.direction == up { }
    li $t1, DIR_UP
    bne $t0, $t1, *else_up
        # if player.y > 1 { player.y -= 1 }
        li $t1, FIRST_INSIDE
        ble $t3, $t1, *set_dead
            li $t1, COLOR_BODY
            lvr $t0, PIXEL_UP($t5)
            # detecta morte por contato proprio
            beq $t0, $t1, *set_dead
            # player y -= 1
            dec $t3
            sb $t3, PLAYER_Y
            return
    else_up:

    # if player.direction == right { }
    li $t1, DIR_RIGHT
    bne $t0, $t1, *else_right
        # if player.x < 58 { player.x += 1 }
        li $t1, LAST_INSIDE
        bge $t2, $t1, *set_dead
            li $t1, COLOR_BODY
            lvr $t0, PIXEL_RIGHT($t5)
            # detecta morte por contato proprio
            beq $t0, $t1, *set_dead
            # player x += 1
            inc $t2
            sb $t2, PLAYER_X
            return
    else_right:

    # if player.direction == down { }
    li $t1, DIR_DOWN
    bne $t0, $t1, *else_down
        # if player.y < 58 { player.y += 1 }
        li $t1, LAST_INSIDE
        bge $t3, $t1, *set_dead
            li $t1, COLOR_BODY
            lvr $t0, PIXEL_DOWN($t5)
            # detecta morte por contato proprio
            beq $t0, $t1, *set_dead
            inc $t3          # player y += 1
            sb $t3, PLAYER_Y # salve player y
            return
    else_down:

    # if player.direction == left { }
    li $t1, DIR_LEFT
    bne $t0, $t1, *else_left
        # if player.x > 1 { player.x -= 1 }
        li $t1, FIRST_INSIDE
        beq $t2, $t1, *set_dead
            li $t1, COLOR_BODY
            lvr $t0, PIXEL_LEFT($t5)
            # detecta morte por contato proprio
            beq $t0, $t1, *set_dead

            # player x -= 1
            dec $t2
            sb $t2, PLAYER_X
            return
    else_left:

    set_dead:
        sb $zero, PLAYER_ALIVE
        return

_spawn_fuit:
    try_to_spawn:
        rand $t0 # random number
        move $t1, $t0 # t1 = t0

        # t0 = t0 % 3600
        divui $t1, $t1, SCREEN_PIXELS
        multi $t1, $t1, SCREEN_PIXELS
        sub $t0, $t0, $t1

        lvr $t1, 0($t0)
        bne $t1, $zero, *try_to_spawn
    spawn:
        sw $t0, FRUIT_POSITION
        return

_new_segment:
    lb $t0, PLAYER_LENGTH  # t0 = lenght
    li $t1, SEGMENTS       # first segment on ram
    li $t2, SEGMENT_SIZE   # size of each segment
    mult $t2, $t2, $t0     # size * lenght
    add $t1, $t1, $t2      # new segment position
    sb $zero, SEGMENT_X($t1) # new_segment.x = 0
    sb $zero, SEGMENT_Y($t1) # new_segment.y = 0
    return

# Desenha um numero sem sinal, alinhado a direita
# t0 = numero
# t1 = posição na VRAM do ultimo digito
# Altera t0, t1, t2 e t4
_draw_number:
    loop_draw_number:
        divui $t2, $t0, 10      # t2 = numero / 10
        multui $t4, $t2, 10
        subu $t4, $t0, $t4      # t4 = numero % 10

        addi $t4, $t4, CHAR_ZERO
        sll $t4, $t4, 8         # [digito][00000000]
        ori $t4, $t4, COLOR_TEXT
        svr $t4, 0($t1)

        dec $t1
        move $t0, $t2           # numero /= 10
        bne $t0, $zero, *loop_draw_number
    return

_draw_highscores:
    # push ra
    addi $sp, $sp, -2
    sw $ra, 0($sp)

    la $t5, *highscores_text
    li $t1, HIGHSCORES_TITLE_POSITION
    loop_draw_highscores_title:
        lrb $t0, 0($t5)
        beq $t0, $zero, *end_draw_highscores_title
            sll $t0, $t0, 8  # [xxxxxxxx][00000000]
            ori $t0, $t0, COLOR_MENU_TEXT
            svr $t0, 0($t1)
            inc $t5
            inc $t1
            j *loop_draw_highscores_title
    end_draw_highscores_title:

    li $t5, 0 # t5 = i
    loop_draw_highscores:
        # t3 = inicio da linha i na VRAM
        addi $t3, $t5, HIGHSCORES_FIRST_LINE
        multi $t3, $t3, SCREEN_WIDTH

        # posição, i + 1
        addi $t0, $t5, 1
        addi $t1, $t3, RANK_COLUMN
        jal *_draw_number

        li $t0, CHAR_DOT
        sll $t0, $t0, 8
        ori $t0, $t0, COLOR_TEXT
        svr $t0, DOT_COLUMN($t3)

        lb $t0, HIGHSCORES($t5)
        addi $t1, $t3, SCORE_COLUMN
        jal *_draw_number

        inc $t5
        li $t0, HIGHSCORE_COUNT
        bne $t5, $t0, *loop_draw_highscores

    # pop ra
    lw $ra, 0($sp)
    addi $sp, $sp, 2
    return

# Le os highscores do save para a RAM
_load_highscores:
    request_save *save_name # rt = save_id, 0 sem memory card
    move $t5, $rt
    sb $t5, SAVE_ID

    li $t1, 0 # t1 = i
    li $t2, HIGHSCORE_COUNT
    loop_load_highscores:
        addi $t4, $t1, SAVE_HIGHSCORES
        read_save $t5, $t4
        sb $rt, HIGHSCORES($t1)
        inc $t1
        bne $t1, $t2, *loop_load_highscores
    return

# Insere o score na lista, mantendo a ordem do maior para o menor, e grava
# a lista no save. Um score que não é maior que o menor da lista é descartado.
_save_score:
    lb $t0, PLAYER_LENGTH # t0 = score

    # t1 = primeira posição com score menor que o novo
    li $t1, 0
    li $t2, HIGHSCORE_COUNT
    loop_find_score_position:
        beq $t1, $t2, *end_save_score
        lb $t3, HIGHSCORES($t1)
        bgt $t0, $t3, *found_score_position
        inc $t1
        j *loop_find_score_position
    found_score_position:

    # Desloca os scores menores uma posição para baixo, descartando o ultimo
    li $t2, LAST_HIGHSCORE
    loop_shift_scores:
        beq $t2, $t1, *end_shift_scores
        subi $t4, $t2, 1
        lb $t3, HIGHSCORES($t4)
        sb $t3, HIGHSCORES($t2)
        dec $t2
        j *loop_shift_scores
    end_shift_scores:
    sb $t0, HIGHSCORES($t1)

    lb $t5, SAVE_ID
    li $t1, 0 # t1 = i
    li $t2, HIGHSCORE_COUNT
    loop_write_highscores:
        addi $t4, $t1, SAVE_HIGHSCORES
        lb $t3, HIGHSCORES($t1)
        write_save $t5, $t4, $t3
        inc $t1
        bne $t1, $t2, *loop_write_highscores

    end_save_score:
    return

_main:
    jal *_load_highscores
    jal *_draw_border
    _start_over:
        li $t0, START_X
        sb $t0, PLAYER_X

        li $t0, START_Y
        sb $t0, PLAYER_Y

        li $t0, DIR_RIGHT
        sb $t0, PLAYER_DIRECTION

        li $t0, 1 # player alive
        sb $t0, PLAYER_ALIVE

        li $t0, 0 # player length
        sb $t0, PLAYER_LENGTH

        li $t0, 0 # fruit don't exist
        sb $t0, FRUIT_EXISTS

        li $t0, 0
        sw $t0, FRUIT_POSITION

        li $t0, MOVE_DELAY
        sb $t0, WAIT_FRAMES

        li $t0, MAX_LENGTH

        li $t0, STATE_START_MENU
        sb $t0, GAME_STATE

        li $t0, PAUSE_COOLDOWN
        sb $t0, PAUSE_DELAY

        sb $zero, SCORE_SAVED

        game_loop:
            _start_logic:
                lb $t0, GAME_STATE
                li $t1, STATE_RUNNING
                beq $t0, $t1, *_game_running_logic

                li $t1, STATE_START_MENU
                beq $t0, $t1, *_start_menu_logic

                li $t1, STATE_PAUSE_MENU
                beq $t0, $t1, *_pause_menu_logic
                    _game_running_logic:
                        _start_moviment_logic:
                            # if player_alive { move }
                            lb $t0, PLAYER_ALIVE
                            beq $t0, $zero, *game_over
                                lb $t1, PLAYER_DIRECTION

                                lb $t0, KEY_UP
                                beq $t0, $zero, *else_set_up
                                    # if direction != down { set up }
                                    li $t2, DIR_DOWN
                                    beq $t1, $t2, *else_set_up
                                        li $t0, DIR_UP
                                        sb $t0, PLAYER_DIRECTION
                                        j *go_move
                                else_set_up:

                                lb $t0, KEY_DOWN
                                beq $t0, $zero, *else_set_down
                                    # if direction != up { set down }
                                    li $t2, DIR_UP
                                    beq $t1, $t2, *else_set_down
                                        li $t0, DIR_DOWN
                                        sb $t0, PLAYER_DIRECTION
                                        j *go_move
                                else_set_down:

                                lb $t0, KEY_LEFT
                                beq $t0, $zero, *else_set_left
                                    # if direction != right { set left }
                                    li $t2, DIR_RIGHT
                                    beq $t1, $t2, *else_set_left
                                        li $t0, DIR_LEFT
                                        sb $t0, PLAYER_DIRECTION
                                        j *go_move
                                else_set_left:

                                lb $t0, KEY_RIGHT
                                beq $t0, $zero, *else_set_right
                                    # if direction != left { set right }
                                    li $t2, DIR_LEFT
                                    beq $t1, $t2, *else_set_right
                                        li $t0, DIR_RIGHT
                                        sb $t0, PLAYER_DIRECTION
                                        j *go_move
                                else_set_right:

                                go_move:
                                    lb $t0, WAIT_FRAMES
                                    bne $t0, $zero, *_waiting_frames
                                        li $t0, MOVE_DELAY
                                        sb $t0, WAIT_FRAMES
                                        jal *_move
                                        j *_end_logic
                                    _waiting_frames:
                                        lb $t0, WAIT_FRAMES
                                        dec $t0
                                        sb $t0, WAIT_FRAMES
                        _end_moviment_logic:

                        _check_pause_game:
                            lb $t0, PAUSE_DELAY
                            beq $t0, $zero, *no_decrement_on_game_running
                                dec $t0
                            no_decrement_on_game_running:
                            sb $t0, PAUSE_DELAY

                            lb $t0, KEY_ENTER
                            beq $t0, $zero, *no_pause
                                lb $t0, PAUSE_DELAY
                                bne $t0, $zero, *_end_pause_menu_logic
                                    li $t0, STATE_PAUSE_MENU
                                    sb $t0, GAME_STATE
                                    li $t0, PAUSE_COOLDOWN
                                    sb $t0, PAUSE_DELAY
                                    j *_end_logic
                            no_pause:

                        _start_fruit_colision:
                            lb $t0, PLAYER_Y                # t0 = player.y
                            multi $t0, $t0, SCREEN_WIDTH    # line in the vram
                            lb $t1, PLAYER_X                # t1 = player.x
                            add $t0, $t0, $t1               # t0 = player position in the vram
                            lw $t1, FRUIT_POSITION          # t1 = fruit position
                            # if fruit.position == player.position { }
                            bne $t0, $t1, *_end_fruit_colision
                                lb $t0, PLAYER_LENGTH # t0 = player.lenght
                                inc $t0               # t0 += 1
                                li $t1, MAX_LENGTH
                                # if player.lenght == 255 { win }
                                bne $t0, $t1, *not_win
                                    sb $t0, PLAYER_LENGTH  # score maximo
                                    sb $zero, PLAYER_ALIVE # player win
                                    j *game_over
                                not_win:
                                    sb $t0, PLAYER_LENGTH
                                    sb $zero, FRUIT_EXISTS # despawn fruit
                                    jal *_new_segment
                        _end_fruit_colision:

                        _start_spawn_fuit_logic:
                            # if fuit_exit == 0 { spawn_fuit }
                            lb $t0, FRUIT_EXISTS
                            bne $t0, $zero, *_end_spawn_fuit_logic
                                jal *_spawn_fuit
                        _end_spawn_fuit_logic:

                        j *_end_logic

                        game_over:
                            # salva o score só no primeiro frame do game over
                            lb $t0, SCORE_SAVED
                            bne $t0, $zero, *score_already_saved
                                jal *_save_score
                                li $t0, 1
                                sb $t0, SCORE_SAVED
                            score_already_saved:

                            lb $t0, KEY_ENTER
                            beq $t0, $zero, *no_restart
                                j *_start_over
                            no_restart:
                                j *_end_logic

                    _end_game_running_logic:

                    _start_menu_logic:
                        lb $t0, KEY_ENTER
                        beq $t0, $zero, *no_start
                            li $t0, STATE_RUNNING
                            sb $t0, GAME_STATE
                        no_start:
                            j *_end_logic
                    _end_start_menu_logic:

                    _pause_menu_logic:
                        lb $t0, PAUSE_DELAY
                        beq $t0, $zero, *no_decrement_pause_menu
                            dec $t0
                        no_decrement_pause_menu:
                        sb $t0, PAUSE_DELAY

                        lb $t0, KEY_ENTER
                        beq $t0, $zero, *_end_pause_menu_logic
                            lb $t0, PAUSE_DELAY
                            bne $t0, $zero, *_end_pause_menu_logic
                                li $t0, STATE_RUNNING
                                sb $t0, GAME_STATE
                                li $t0, PAUSE_COOLDOWN
                                sb $t0, PAUSE_DELAY
                                j *_end_logic

                    _end_pause_menu_logic:

            _end_logic:

            _start_drawing:
                jal *_clear_scream
                jal *_draw_fruit
                jal *_draw_player
                jal *_draw_border

                lb $t0, GAME_STATE

                li $t1, STATE_RUNNING
                beq $t0, $t1, *_game_runing

                li $t1, STATE_START_MENU
                beq $t0, $t1, *_start_menu

                li $t1, STATE_PAUSE_MENU
                beq $t0, $t1, *_game_pause

                _game_runing:
                    lb $t0, PLAYER_ALIVE
                    bne $t0, $zero, *not_game_over
                        jal *_draw_game_over
                        jal *_draw_highscores
                    not_game_over:
                    j *_render

                _start_menu:
                    jal *_draw_start_menu
                    j *_render

                _game_pause:
                    jal *_draw_pause_menu
                    j *_render

                _render:
                    li $sc, RENDER
                    syscall
            _end_drawing:

            j *game_loop



.data
    press_start: .string "[PRESS START!]"
    resume_game: .string "[RESUME GAME!]"
    save_name: .string "snake\0"
    highscores_text: .string "HIGHSCORES\0"
    # Strings não ASCII devem ficar no fim, labels depois delas ficam com o
    # endereço errado
    game_over_text_1: .string "■■■■■ ■■■■■ ■■■■■ ■■■■■■     ■   ■ ■ ■ ■ ■    ■   ■ ■■■■■ ■ ■ ■ ■■■  ■   ■ ■   ■ ■   ■ ■    ■■■■■ ■   ■ ■   ■ ■■■■■                       ■■■■■ ■   ■ ■■■■■ ■■■■ ■   ■ ■   ■ ■     ■   ■■   ■ ■   ■ ■■■   ■■■■ ■   ■  ■ ■  ■     ■   ■■■■■■   ■   ■■■■■ ■   ■"
