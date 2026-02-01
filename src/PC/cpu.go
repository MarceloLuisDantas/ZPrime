package pc

import (
	"fmt"
	"image/color"
	"math/rand"
	"os"
	"strconv"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const SCREAM_CHARS = 60
const FONT_SIZE = 6
const PADDING_BORDER = 4
const PADDING_BETWEEN_LINES = 1
const PADDING_BETWEEN_CHARS = 1
const SCALING = 2
const SCREAM_H = (SCREAM_CHARS * FONT_SIZE) + (PADDING_BETWEEN_LINES * SCREAM_CHARS) + PADDING_BORDER*2
const SCREAM_W = (SCREAM_CHARS * FONT_SIZE) + (PADDING_BETWEEN_CHARS * SCREAM_CHARS) + PADDING_BORDER*2

var cores = map[COLOR]color.RGBA{
	BRANCO:   rl.White,
	PRETO:    rl.Black,
	AZUL:     rl.Blue,
	VERMELHO: rl.Red,
	VERDE:    rl.Green,
	CINZA:    rl.Gray,
	MARROM:   rl.Brown,
	AMARELO:  rl.Yellow,
}

type CPU struct {
	zero int16
	t0   int16
	t1   int16
	t2   int16
	t3   int16
	t4   int16
	t5   int16
	rt   int16
	sc   int16
	sp   uint16
	fp   uint16
	ra   uint16
	pc   uint16
	gp   uint16
	ir   string
	rom  *ROM
	ram  *RAM
	vram *VRAM
}

func NewCPU(file []string) (*CPU, error) {
	var cpu CPU

	gp, err := strconv.ParseInt(file[0], 10, 16)
	if err != nil {
		return nil, err
	}
	cpu.gp = uint16(gp)

	rom, err := NewRom(file[1:])
	if err != nil {
		return nil, err
	}
	cpu.rom = rom

	cpu.ram = NewRam()
	cpu.vram = NewVram()

	cpu.sp = 65535
	cpu.fp = 65535
	return &cpu, nil
}

func (cpu *CPU) SetRegister(dest string, value int16) {
	switch dest {
	case "$t0":
		cpu.t0 = value
	case "$t1":
		cpu.t1 = value
	case "$t2":
		cpu.t2 = value
	case "$t3":
		cpu.t3 = value
	case "$t4":
		cpu.t4 = value
	case "$t5":
		cpu.t5 = value
	case "$pc":
		cpu.pc = uint16(value)
	case "$zero":
		cpu.zero = 0
	case "$sp":
		cpu.sp = uint16(value)
	case "$fp":
		cpu.fp = uint16(value)
	case "$sc":
		cpu.sc = value
	case "$ra":
		cpu.ra = uint16(value)
	case "$gp":
		cpu.gp = uint16(value)
	case "$rt":
		cpu.rt = value
	}
}

func (cpu *CPU) GetRegister(dest string) int16 {
	switch dest {
	case "$t0":
		return cpu.t0
	case "$t1":
		return cpu.t1
	case "$t2":
		return cpu.t2
	case "$t3":
		return cpu.t3
	case "$t4":
		return cpu.t4
	case "$t5":
		return cpu.t5
	case "$pc":
		return int16(cpu.pc)
	case "$zero":
		return 0
	case "$sp":
		return int16(cpu.sp)
	case "$fp":
		return int16(cpu.fp)
	case "$sc":
		return cpu.sc
	case "$ra":
		return int16(cpu.ra)
	case "$gp":
		return int16(cpu.gp)
	case "$rt":
		return cpu.rt
	}

	println("registrador inválido: ", dest)
	panic("")
}

func (cpu *CPU) Add(dest, src1, src2 string) {
	v1 := cpu.GetRegister(src1)
	v2 := cpu.GetRegister(src2)
	cpu.SetRegister(dest, v1+v2)
}

func (cpu *CPU) Addi(dest, src string, value int16) {
	v := cpu.GetRegister(src)
	cpu.SetRegister(dest, v+value)
}

func (cpu *CPU) Addu(dest, src1, src2 string) {
	v1 := uint16(cpu.GetRegister(src1))
	v2 := uint16(cpu.GetRegister(src2))
	cpu.SetRegister(dest, int16(v1+v2))
}

func (cpu *CPU) Addui(dest, src string, value uint16) {
	v := uint16(cpu.GetRegister(src))
	result := v + uint16(value)
	cpu.SetRegister(dest, int16(result))
}

func (cpu *CPU) Sub(dest, src1, src2 string) {
	v1 := cpu.GetRegister(src1)
	v2 := cpu.GetRegister(src2)
	cpu.SetRegister(dest, v1-v2)
}

func (cpu *CPU) Subi(dest, src string, value int16) {
	v := cpu.GetRegister(src)
	cpu.SetRegister(dest, v-value)
}

func (cpu *CPU) Subu(dest, src1, src2 string) {
	v1 := uint16(cpu.GetRegister(src1))
	v2 := uint16(cpu.GetRegister(src2))
	cpu.SetRegister(dest, int16(v1-v2))
}

func (cpu *CPU) Subui(dest, src string, value uint16) {
	v := uint16(cpu.GetRegister(src))
	result := v - uint16(value)
	cpu.SetRegister(dest, int16(result))
}

func (cpu *CPU) Mult(dest, src1, src2 string) {
	v1 := cpu.GetRegister(src1)
	v2 := cpu.GetRegister(src2)
	cpu.SetRegister(dest, v1*v2)
}

func (cpu *CPU) Multi(dest, src string, value int16) {
	v := cpu.GetRegister(src)
	cpu.SetRegister(dest, v*value)
}

func (cpu *CPU) Multu(dest, src1, src2 string) {
	v1 := uint16(cpu.GetRegister(src1))
	v2 := uint16(cpu.GetRegister(src2))
	cpu.SetRegister(dest, int16(v1*v2))
}

func (cpu *CPU) Multui(dest, src string, value uint16) {
	v := uint16(cpu.GetRegister(src))
	cpu.SetRegister(dest, int16(v*value))
}

func (cpu *CPU) Div(dest, src1, src2 string) {
	v1 := cpu.GetRegister(src1)
	v2 := cpu.GetRegister(src2)
	cpu.SetRegister(dest, v1/v2)
}

func (cpu *CPU) Divi(dest, src string, value int16) {
	v := cpu.GetRegister(src)
	cpu.SetRegister(dest, v/value)
}

func (cpu *CPU) Divu(dest, src1, src2 string) {
	v1 := uint16(cpu.GetRegister(src1))
	v2 := uint16(cpu.GetRegister(src2))
	cpu.SetRegister(dest, int16(v1/v2))
}

func (cpu *CPU) Divui(dest, src string, value uint16) {
	v := uint16(cpu.GetRegister(src))
	cpu.SetRegister(dest, int16(v/value))
}

func (cpu *CPU) And(dest, src1, src2 string) {
	v1 := cpu.GetRegister(src1)
	v2 := cpu.GetRegister(src2)
	cpu.SetRegister(dest, v1&v2)
}

func (cpu *CPU) Andi(dest, src string, value int16) {
	v := cpu.GetRegister(src)
	cpu.SetRegister(dest, v&value)
}

func (cpu *CPU) Or(dest, src1, src2 string) {
	v1 := cpu.GetRegister(src1)
	v2 := cpu.GetRegister(src2)
	cpu.SetRegister(dest, v1|v2)
}

func (cpu *CPU) Ori(dest, src string, value int16) {
	v := cpu.GetRegister(src)
	cpu.SetRegister(dest, v|value)
}

func (cpu *CPU) Slt(dest, src1, src2 string) {
	if cpu.GetRegister(src1) < cpu.GetRegister(src2) {
		cpu.SetRegister(dest, 1)
	} else {
		cpu.SetRegister(dest, 0)
	}
}

func (cpu *CPU) Slti(dest, src string, value int16) {
	if cpu.GetRegister(src) < value {
		cpu.SetRegister(dest, 1)
	} else {
		cpu.SetRegister(dest, 0)
	}
}

func (cpu *CPU) Sltu(dest, src1, src2 string) {
	if uint16(cpu.GetRegister(src1)) < uint16(cpu.GetRegister(src2)) {
		cpu.SetRegister(dest, 1)
	} else {
		cpu.SetRegister(dest, 0)
	}
}

func (cpu *CPU) Sltui(dest, src string, value uint16) {
	if uint16(cpu.GetRegister(src)) < value {
		cpu.SetRegister(dest, 1)
	} else {
		cpu.SetRegister(dest, 0)
	}
}

func (cpu *CPU) Move(dest, src string) {
	cpu.SetRegister(dest, cpu.GetRegister(src))
}

func (cpu *CPU) Li(dest string, value int16) {
	cpu.SetRegister(dest, value)
}

func (cpu *CPU) La(dest string, value uint16) {
	cpu.SetRegister(dest, int16(value))
}

func (cpu *CPU) Jump(point uint16) {
	cpu.pc = point
}

func (cpu *CPU) Jr(src string) {
	cpu.pc = uint16(cpu.GetRegister(src))
}

func (cpu *CPU) Jal(point uint16) {
	cpu.ra = cpu.pc
	cpu.pc = point
}

func (cpu *CPU) Ret() {
	cpu.pc = cpu.ra
}

func (cpu *CPU) Beq(src1, src2 string, point uint16) {
	if cpu.GetRegister(src1) == cpu.GetRegister(src2) {
		cpu.pc = point
	}
}

func (cpu *CPU) Bne(src1, src2 string, point uint16) {
	// println(cpu.GetRegister(src1), " == ", cpu.GetRegister(src2))
	if cpu.GetRegister(src1) != cpu.GetRegister(src2) {
		cpu.pc = point
	}
}

func (cpu *CPU) Bgt(src1, src2 string, point uint16) {
	if cpu.GetRegister(src1) > cpu.GetRegister(src2) {
		cpu.pc = point
	}
}

func (cpu *CPU) Bge(src1, src2 string, point uint16) {
	if cpu.GetRegister(src1) >= cpu.GetRegister(src2) {
		cpu.pc = point
	}
}

func (cpu *CPU) Blt(src1, src2 string, point uint16) {
	if cpu.GetRegister(src1) < cpu.GetRegister(src2) {
		cpu.pc = point
	}
}

func (cpu *CPU) Ble(src1, src2 string, point uint16) {
	if cpu.GetRegister(src1) <= cpu.GetRegister(src2) {
		cpu.pc = point
	}
}

func (cpu *CPU) Inc(dest string) {
	cpu.SetRegister(dest, cpu.GetRegister(dest)+1)
}

func (cpu *CPU) Dec(dest string) {
	cpu.SetRegister(dest, cpu.GetRegister(dest)-1)
}

func (cpu *CPU) Rand(dest string) {
	cpu.SetRegister(dest, int16(rand.Intn(65536)-32768))
}

func (cpu *CPU) Lrb(dest string, offset, point uint16) {
	addr := cpu.gp + point + offset
	v, err := strconv.Atoi(cpu.rom.rom[addr])
	if err != nil {
		panic(err)
	}

	cpu.SetRegister(dest, int16(v))
}

func (cpu *CPU) Lrw(dest string, offset, point uint16) {
	addr := cpu.gp + point + offset

	h1, err := strconv.Atoi(cpu.rom.rom[addr])
	if err != nil {
		panic(err)
	}

	h2, err := strconv.Atoi(cpu.rom.rom[addr+1])
	if err != nil {
		panic(err)
	}

	value := (int16(h1) << 8) | int16(h2)
	// print(value)
	cpu.SetRegister(dest, value)
}

func (cpu *CPU) Sll(dest, src string, value uint16) {
	v := int16(cpu.GetRegister(src))
	cpu.SetRegister(dest, int16(v<<value))
}

func (cpu *CPU) Srl(dest, src string, value uint16) {
	v := int16(cpu.GetRegister(src))
	cpu.SetRegister(dest, int16(v>>value))
}

func (cpu *CPU) Sb(src string, offset int16, reg_index string) {
	value := uint8(cpu.GetRegister(src))
	index := int32(cpu.GetRegister(reg_index)) + int32(offset)

	cpu.ram.save_byte(value, uint16(index))
}

func (cpu *CPU) Sw(src string, offset int16, reg_index string) {
	value := uint16(cpu.GetRegister(src))
	index := int32(cpu.GetRegister(reg_index)) + int32(offset)
	cpu.ram.save_world(value, uint16(index))
}

func (cpu *CPU) Lb(dest string, offset int16, reg_index string) {
	index := int32(cpu.GetRegister(reg_index)) + int32(offset)
	value := cpu.ram.load_byte(uint16(index))
	cpu.SetRegister(dest, int16(value))
}

func (cpu *CPU) Lw(dest string, offset int16, reg_index string) {
	index := int32(cpu.GetRegister(reg_index)) + int32(offset)
	value := cpu.ram.load_world(uint16(index))
	cpu.SetRegister(dest, int16(value))
}

func (cpu *CPU) Svr(src string, offset int16, reg_index string) {
	value := uint16(cpu.GetRegister(src))
	index := int32(cpu.GetRegister(reg_index)) + int32(offset)
	y := index / 60
	x := index % 60
	cpu.vram.vram[y][x] = Character(value)
}

func (cpu *CPU) Lvr(dest string, offset int16, reg_index string) {
	index := int32(cpu.GetRegister(reg_index)) + int32(offset)
	y := index / 60
	x := index % 60
	cpu.SetRegister(dest, int16(cpu.vram.vram[y][x]))
}

func (cpu *CPU) RenderFrame() {
	rl.ClearBackground(rl.Black)
	rl.BeginDrawing()
	var pos_y int32 = PADDING_BORDER * SCALING
	for y := range 60 {
		var pos_x int32 = PADDING_BORDER * SCALING
		for x := range 60 {
			_, char := cpu.vram.GetChar(x, y)
			_, color := cpu.vram.GetCharColor(x, y)
			_, bcolor := cpu.vram.GetCharBackColor(x, y)

			rl.DrawRectangle(pos_x, pos_y, FONT_SIZE*SCALING, FONT_SIZE*SCALING, cores[bcolor])
			rl.DrawText(string(char), pos_x, pos_y, FONT_SIZE*SCALING, cores[color])

			pos_x += (FONT_SIZE + PADDING_BETWEEN_CHARS) * SCALING
		}
		pos_y += (FONT_SIZE + PADDING_BETWEEN_LINES) * SCALING
	}
	rl.EndDrawing()
}

// 0 - exit
// 1 - print t0 como int8
// 2 - print t0 como utin8
// 3 - print t0 como int16
// 4 - print t0 como uint16
// 5 - print t0 como char utf8
// 100 - render frame
func (cpu *CPU) Syscall() {
	value := cpu.t0
	switch cpu.sc {
	case 0:
		os.Exit(0)
	case 1:
		print(int8(value))
	case 2:
		print(uint8(value))
	case 3:
		print(value)
	case 4:
		print(uint16(value))
	case 5:
		fmt.Printf("%c", int8(value))

	case 100:
		cpu.RenderFrame()

	case 1001:
		println(int8(value))
	case 1002:
		println(uint8(value))
	case 1003:
		println(value)
	case 1004:
		println(uint16(value))
	case 1005:
		fmt.Printf("%c\n", int8(value))

	case -1:
		cpu.ram.show_ram()
	}
}

func (cpu *CPU) ExecCurrentInstruction(tokens []string) int {
	switch tokens[0] {

	case "add":
		cpu.Add(tokens[1], tokens[2], tokens[3])

	case "sub":
		cpu.Sub(tokens[1], tokens[2], tokens[3])

	case "mult":
		cpu.Mult(tokens[1], tokens[2], tokens[3])

	case "div":
		cpu.Div(tokens[1], tokens[2], tokens[3])

	case "addu":
		cpu.Addu(tokens[1], tokens[2], tokens[3])

	case "subu":
		cpu.Subu(tokens[1], tokens[2], tokens[3])

	case "multu":
		cpu.Multu(tokens[1], tokens[2], tokens[3])

	case "divu":
		cpu.Divu(tokens[1], tokens[2], tokens[3])

	case "addi":
		v, _ := strconv.Atoi(tokens[3])
		cpu.Addi(tokens[1], tokens[2], int16(v))

	case "subi":
		v, _ := strconv.Atoi(tokens[3])
		cpu.Subi(tokens[1], tokens[2], int16(v))

	case "multi":
		v, _ := strconv.Atoi(tokens[3])
		cpu.Multi(tokens[1], tokens[2], int16(v))

	case "divi":
		v, _ := strconv.Atoi(tokens[3])
		cpu.Divi(tokens[1], tokens[2], int16(v))

	case "addui":
		v, _ := strconv.Atoi(tokens[3])
		cpu.Addui(tokens[1], tokens[2], uint16(v))

	case "subui":
		v, _ := strconv.Atoi(tokens[3])
		cpu.Subui(tokens[1], tokens[2], uint16(v))

	case "multui":
		v, _ := strconv.Atoi(tokens[3])
		cpu.Multui(tokens[1], tokens[2], uint16(v))

	case "divui":
		v, _ := strconv.Atoi(tokens[3])
		cpu.Divui(tokens[1], tokens[2], uint16(v))

	case "move":
		cpu.Move(tokens[1], tokens[2])

	case "li":
		v, _ := strconv.Atoi(tokens[2])
		cpu.Li(tokens[1], int16(v))

	case "la":
		v, _ := strconv.Atoi(tokens[2])
		cpu.La(tokens[1], uint16(v))

	case "j":
		v, _ := strconv.Atoi(tokens[1])
		cpu.Jump(uint16(v))

	case "jal":
		v, _ := strconv.Atoi(tokens[1])
		cpu.Jal(uint16(v))

	case "jr":
		cpu.Jr(tokens[1])

	case "return":
		cpu.Ret()

	case "beq":
		v, _ := strconv.Atoi(tokens[3])
		cpu.Beq(tokens[1], tokens[2], uint16(v))

	case "bne":
		v, _ := strconv.Atoi(tokens[3])
		cpu.Bne(tokens[1], tokens[2], uint16(v))

	case "bgt":
		v, _ := strconv.Atoi(tokens[3])
		cpu.Bgt(tokens[1], tokens[2], uint16(v))

	case "bge":
		v, _ := strconv.Atoi(tokens[3])
		cpu.Bge(tokens[1], tokens[2], uint16(v))

	case "blt":
		v, _ := strconv.Atoi(tokens[3])
		cpu.Blt(tokens[1], tokens[2], uint16(v))

	case "ble":
		v, _ := strconv.Atoi(tokens[3])
		cpu.Ble(tokens[1], tokens[2], uint16(v))

	case "slt":
		cpu.Slt(tokens[1], tokens[2], tokens[3])

	case "slti":
		v, _ := strconv.Atoi(tokens[3])
		cpu.Slti(tokens[1], tokens[2], int16(v))

	case "sltu":
		cpu.Sltu(tokens[1], tokens[2], tokens[3])

	case "sltui":
		v, _ := strconv.Atoi(tokens[3])
		cpu.Sltui(tokens[1], tokens[2], uint16(v))

	case "and":
		cpu.And(tokens[1], tokens[2], tokens[3])

	case "andi":
		v, _ := strconv.Atoi(tokens[3])
		cpu.Andi(tokens[1], tokens[2], int16(v))

	case "or":
		cpu.Or(tokens[1], tokens[2], tokens[3])

	case "ori":
		v, _ := strconv.Atoi(tokens[3])
		cpu.Ori(tokens[1], tokens[2], int16(v))

	case "sll":
		v, _ := strconv.Atoi(tokens[3])
		cpu.Sll(tokens[1], tokens[2], uint16(v))

	case "srl":
		v, _ := strconv.Atoi(tokens[3])
		cpu.Srl(tokens[1], tokens[2], uint16(v))

	case "inc":
		cpu.Inc(tokens[1])

	case "dec":
		cpu.Dec(tokens[1])

	case "rand":
		cpu.Rand(tokens[1])

	case "lrb":
		var offset uint16
		if tokens[2][0] == '$' {
			offset = uint16(cpu.GetRegister(tokens[2]))
		} else {
			v, _ := strconv.Atoi(tokens[2])
			offset = uint16(v)
		}

		if tokens[3][0] == '$' {
			cpu.Lrb(tokens[1], uint16(offset), uint16(cpu.GetRegister(tokens[3])))
		} else {
			point, _ := strconv.Atoi(tokens[3])
			cpu.Lrb(tokens[1], uint16(offset), uint16(point))
		}

	case "lrw":
		offset, _ := strconv.Atoi(tokens[2])

		if tokens[3][0] == '$' {
			cpu.Lrw(tokens[1], uint16(offset), uint16(cpu.GetRegister(tokens[3])))
		} else {
			point, _ := strconv.Atoi(tokens[3])
			cpu.Lrw(tokens[1], uint16(offset), uint16(point))
		}

	case "lvr":
		offset, _ := strconv.Atoi(tokens[2])
		cpu.Lvr(tokens[1], int16(offset), tokens[3])

	case "svr":
		offset, _ := strconv.Atoi(tokens[2])
		cpu.Svr(tokens[1], int16(offset), tokens[3])

	case "sb":
		offset, _ := strconv.Atoi(tokens[2])
		cpu.Sb(tokens[1], int16(offset), tokens[3])

	case "sw":
		offset, _ := strconv.Atoi(tokens[2])
		cpu.Sw(tokens[1], int16(offset), tokens[3])

	case "lb":
		offset, _ := strconv.Atoi(tokens[2])
		cpu.Lb(tokens[1], int16(offset), tokens[3])

	case "lw":
		offset, _ := strconv.Atoi(tokens[2])
		cpu.Lw(tokens[1], int16(offset), tokens[3])

	case "syscall":
		cpu.Syscall()

	default:
		fmt.Printf("Instrução invalida \"%s\".\n", tokens[0])
		return -2
	}

	return 0
}

func (cpu *CPU) CarryNextInstruction() {
	h1 := cpu.rom.rom[cpu.pc]
	h2 := cpu.rom.rom[cpu.pc+1]
	if h1 == "sys" || h1 == "ret" {
		cpu.ir = h1 + h2
	} else {
		cpu.ir = h1 + " " + h2
	}
	cpu.pc += 2
}

// REGISTRADORES EM MEMORIA
// 0x0000 [ KEY UP ]
// 0x0001 [ KEY DOWN ]
// 0x0002 [ KEY LEFT ]
// 0x0003 [ KEY RIGHT ]
// 0x0004 [ KEY SPACE ]
// 0x0005 [ KEY ENTER ]
// 0x0006 [ KEY BACKSPACE ]
// 0x0007 [ KEY W ]
// 0x0008 [ KEY A ]
// 0x0009 [ KEY S ]
// 0x000A [ KEY D ]
// 0x000B [ KEY Q ]
// 0x000C [ KEY E ]
// 0x000D [ KEY I ]
// 0x000E [ KEY O ]
// 0x000F [ KEY P ]
func (cpu *CPU) ReadInput() {
	if rl.IsKeyDown(rl.KeyUp) {
		cpu.ram.ram[0] = 1
	} else if rl.IsKeyDown(rl.KeyDown) {
		cpu.ram.ram[1] = 1
	} else if rl.IsKeyDown(rl.KeyLeft) {
		cpu.ram.ram[2] = 1
	} else if rl.IsKeyDown(rl.KeyRight) {
		cpu.ram.ram[3] = 1
	} else if rl.IsKeyDown(rl.KeySpace) {
		cpu.ram.ram[4] = 1
	} else if rl.IsKeyDown(rl.KeyEnter) {
		cpu.ram.ram[5] = 1
	} else if rl.IsKeyDown(rl.KeyBackspace) {
		cpu.ram.ram[6] = 1
	} else if rl.IsKeyDown(rl.KeyW) {
		cpu.ram.ram[7] = 1
	} else if rl.IsKeyDown(rl.KeyA) {
		cpu.ram.ram[8] = 1
	} else if rl.IsKeyDown(rl.KeyS) {
		cpu.ram.ram[9] = 1
	} else if rl.IsKeyDown(rl.KeyS) {
		cpu.ram.ram[10] = 1
	} else if rl.IsKeyDown(rl.KeyD) {
		cpu.ram.ram[11] = 1
	} else if rl.IsKeyDown(rl.KeyQ) {
		cpu.ram.ram[12] = 1
	} else if rl.IsKeyDown(rl.KeyE) {
		cpu.ram.ram[13] = 1
	} else if rl.IsKeyDown(rl.KeyI) {
		cpu.ram.ram[14] = 1
	} else if rl.IsKeyDown(rl.KeyO) {
		cpu.ram.ram[15] = 1
	} else if rl.IsKeyDown(rl.KeyP) {
		cpu.ram.ram[16] = 1
	}

	if rl.IsKeyReleased(rl.KeyUp) {
		cpu.ram.ram[0] = 0
	} else if rl.IsKeyReleased(rl.KeyDown) {
		cpu.ram.ram[1] = 0
	} else if rl.IsKeyReleased(rl.KeyLeft) {
		cpu.ram.ram[2] = 0
	} else if rl.IsKeyReleased(rl.KeyRight) {
		cpu.ram.ram[3] = 0
	} else if rl.IsKeyReleased(rl.KeySpace) {
		cpu.ram.ram[4] = 0
	} else if rl.IsKeyReleased(rl.KeyEnter) {
		cpu.ram.ram[5] = 0
	} else if rl.IsKeyReleased(rl.KeyBackspace) {
		cpu.ram.ram[6] = 0
	} else if rl.IsKeyReleased(rl.KeyW) {
		cpu.ram.ram[7] = 0
	} else if rl.IsKeyReleased(rl.KeyA) {
		cpu.ram.ram[8] = 0
	} else if rl.IsKeyReleased(rl.KeyS) {
		cpu.ram.ram[9] = 0
	} else if rl.IsKeyReleased(rl.KeyS) {
		cpu.ram.ram[10] = 0
	} else if rl.IsKeyReleased(rl.KeyD) {
		cpu.ram.ram[11] = 0
	} else if rl.IsKeyReleased(rl.KeyQ) {
		cpu.ram.ram[12] = 0
	} else if rl.IsKeyReleased(rl.KeyE) {
		cpu.ram.ram[13] = 0
	} else if rl.IsKeyReleased(rl.KeyI) {
		cpu.ram.ram[14] = 0
	} else if rl.IsKeyReleased(rl.KeyO) {
		cpu.ram.ram[15] = 0
	} else if rl.IsKeyReleased(rl.KeyP) {
		cpu.ram.ram[16] = 0
	}
}

func (cpu *CPU) Run() {
	rl.InitWindow(SCREAM_W*SCALING, SCREAM_H*SCALING, "Teste raylib")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		cpu.ReadInput()

		cpu.CarryNextInstruction()
		ok := cpu.ExecCurrentInstruction(strings.Split(cpu.ir, " "))
		if ok != 0 {
			return
		}

		if cpu.pc >= cpu.gp {
			return
		}
	}
}

func (cpu *CPU) ShowStack() {

}
