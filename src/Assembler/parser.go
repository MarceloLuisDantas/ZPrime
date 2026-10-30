package assembler

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	NOP     string = "nop"
	ADD     string = "add"
	ADDI    string = "addi"
	ADDU    string = "addu"
	ADDUI   string = "addui"
	SUB     string = "sub"
	SUBI    string = "subi"
	SUBU    string = "subu"
	SUBUI   string = "subui"
	MULT    string = "mult"
	MULTI   string = "multi"
	MULTU   string = "multu"
	MULTUI  string = "multui"
	DIV     string = "div"
	DIVI    string = "divi"
	DIVU    string = "divu"
	DIVUI   string = "divui"
	OR      string = "or"
	ORI     string = "ori"
	AND     string = "and"
	ANDI    string = "andi"
	SLL     string = "sll"
	SRL     string = "srl"
	SLT     string = "slt"
	SLTI    string = "slti"
	SLTU    string = "sltu"
	SLTUI   string = "sltui"
	LI      string = "li"
	LA      string = "la"
	MOVE    string = "move"
	J       string = "j"
	JR      string = "jr"
	JAL     string = "jal"
	BEQ     string = "beq"
	BNE     string = "bne"
	BGT     string = "bgt"
	BGE     string = "bge"
	BLT     string = "blt"
	BLE     string = "ble"
	LW      string = "lw"
	LB      string = "lb"
	SW      string = "sw"
	SB      string = "sb"
	LVR     string = "lvr"
	SVR     string = "svr"
	LRW     string = "lrw"
	LRB     string = "lrb"
	INC     string = "inc"
	DEC     string = "dec"
	SYSCALL string = "syscall"
	RETURN  string = "return"
	RAND    string = "rand"
)

var RegistersInstructions = map[string]bool{
	ADD: true, SUB: true, MULT: true, DIV: true,
	ADDU: true, SUBU: true, MULTU: true, DIVU: true,
	OR: true, AND: true, SLT: true, SLTU: true,
}

var ImmediateInstructions = map[string]bool{
	ADDI: true, SUBI: true, MULTI: true, DIVI: true,
	ADDUI: true, SUBUI: true, MULTUI: true, DIVUI: true,
	ORI: true, ANDI: true, SLTI: true, SLTUI: true,
	SLL: true, SRL: true,
}

var BranchInstructions = map[string]bool{
	BEQ: true, BNE: true, BGT: true, BGE: true, BLT: true, BLE: true,
}

var MemoriInstructions = map[string]bool{
	LW: true, LB: true, SW: true, SB: true, LVR: true, SVR: true, LRW: true, LRB: true,
}

var JumpInstructions = map[string]bool{
	J: true, JR: true, JAL: true,
}

type Type string

const (
	STR    Type = "string"
	INT8   Type = "int8"
	INT16  Type = "int16"
	UINT8  Type = "uint8"
	UINT16 Type = "uint16"
)

type Section string

const (
	TEXT Section = ".text"
	DATA Section = ".data"
)

type Instruction []string

type Parser struct {
	Tokens       []Token
	Instructions []Instruction
	Position     int
	Len          int
	Gp           int
	JumpLabels   map[string]int
	RomLabels    map[string]int
}

func newParser(tokens []Token) *Parser {
	p := Parser{tokens, []Instruction{}, 0, 0, -1, make(map[string]int), make(map[string]int)}
	return &p
}

func (parser *Parser) getCurrentToken() (Token, error) {
	if parser.Position >= len(parser.Tokens) {
		return Token{}, fmt.Errorf("end of file")
	}
	return parser.Tokens[parser.Position], nil
}

func (parser *Parser) getNextXToken(x int) Token {
	return parser.Tokens[parser.Position+x]
}

func invalidRegister(token Token) error {
	fmt.Printf("Registrador esperado em - linha: %d coluna: %d, mas foi encontrado %s \"%s\".\n",
		token.Line, token.Column, token.TokenType, token.Value)
	return fmt.Errorf("invalid reg")
}

func invalidImmediate(token Token) error {
	fmt.Printf("Valor imediato esperado em - linha: %d coluna: %d, mas foi encontrado %s.\n",
		token.Line, token.Column, token.TokenType)
	return fmt.Errorf("invalid imme")
}

func invalidLabel(token Token) error {
	fmt.Printf("Label esperado em - linha %d coluna: %d, mas foi encontrado %s.\n",
		token.Line, token.Column, token.TokenType)
	return fmt.Errorf("invalid label")
}

func invalidRegOrLabel(token Token) error {
	fmt.Printf("Label ou Registrador esperando em linha: %d coluna: %d, mas foi encontrado %s.\n",
		token.Line, token.Column, token.TokenType)
	return fmt.Errorf("invalid register or label")
}

func toMuchArgs(token Token, args int) error {
	fmt.Printf("Instrução \"%s\" recebe apena %d argumentos.\n", token.Value, args)
	return fmt.Errorf("to much args")
}

func virgulaMissing(line, column int) error {
	fmt.Printf("Virgula esperada em linha: %d coluna: %d.\n", line, column)
	return fmt.Errorf("missing separator")
}

func openParenMissing(line, column int) error {
	fmt.Printf("Abertura de parenteses esperada ao acessar memoria linha: %d coluna: %d.\n", line, column)
	return fmt.Errorf("missing separator")
}

func closeParenMissing(line, column int) error {
	fmt.Printf("Fechamento de parenteses esperada ao acessar memoria linha: %d coluna: %d.\n", line, column)
	return fmt.Errorf("missing separator")
}

// add, sub, mult, div, or, and, slt, slt
// add $r1, $r2, $r2 \n
func (parser *Parser) parseRegisterInstruciton() error {
	current_token, _ := parser.getCurrentToken()
	ins := Instruction{current_token.Value}

	arg1 := parser.getNextXToken(1)
	if arg1.TokenType != REGISTER {
		return invalidRegister(arg1)
	}

	virgula1 := parser.getNextXToken(2)
	if virgula1.TokenType != VIRGULA {
		return virgulaMissing(virgula1.Line, virgula1.Column)
	}

	arg2 := parser.getNextXToken(3)
	if arg2.TokenType != REGISTER {
		return invalidRegister(arg2)
	}

	virgula2 := parser.getNextXToken(4)
	if virgula2.TokenType != VIRGULA {
		return virgulaMissing(virgula2.Line, virgula2.Column)
	}

	arg3 := parser.getNextXToken(5)
	if arg3.TokenType != REGISTER {
		return invalidRegister(arg3)
	}

	nl := parser.getNextXToken(6)
	if nl.TokenType != NEW_LINE {
		return toMuchArgs(current_token, 3)
	}

	ins = append(ins, arg1.Value)
	ins = append(ins, arg2.Value)
	ins = append(ins, arg3.Value)
	parser.Instructions = append(parser.Instructions, ins)
	parser.Len += 2
	parser.Position += 7
	return nil
}

// addi, subi, multi, divi, ori, andi, slti, sltui, sll, srl
// addi $r1, $r2, value \n
func (parser *Parser) parseImediateInstruciton() error {
	current_token, _ := parser.getCurrentToken()
	ins := Instruction{current_token.Value}

	arg1 := parser.getNextXToken(1)
	if arg1.TokenType != REGISTER {
		return invalidRegister(arg1)
	}

	virgula1 := parser.getNextXToken(2)
	if virgula1.TokenType != VIRGULA {
		return virgulaMissing(virgula1.Line, virgula1.Column)
	}

	arg2 := parser.getNextXToken(3)
	if arg2.TokenType != REGISTER {
		return invalidRegister(arg2)
	}

	virgula2 := parser.getNextXToken(4)
	if virgula2.TokenType != VIRGULA {
		return virgulaMissing(virgula2.Line, virgula2.Column)
	}

	arg3 := parser.getNextXToken(5)
	if arg3.TokenType != NUMBER {
		return invalidImmediate(arg3)
	}

	nl := parser.getNextXToken(6)
	if nl.TokenType != NEW_LINE {
		return toMuchArgs(current_token, 3)
	}

	ins = append(ins, arg1.Value)
	ins = append(ins, arg2.Value)
	ins = append(ins, arg3.Value)
	parser.Instructions = append(parser.Instructions, ins)
	parser.Len += 2
	parser.Position += 7
	return nil
}

// j, jr, jal
func (parser *Parser) parseJump() error {
	current_token, _ := parser.getCurrentToken()
	ins := Instruction{current_token.Value}

	arg := parser.getNextXToken(1)
	switch current_token.Value {
	case "j", "jal":
		if arg.TokenType != LABEL_REF {
			return invalidLabel(arg)
		}
	case "jr":
		if arg.TokenType != REGISTER {
			return invalidRegister(arg)
		}
	}

	nl := parser.getNextXToken(2)
	if nl.TokenType != NEW_LINE {
		return toMuchArgs(current_token, 1)
	}

	ins = append(ins, arg.Value)
	parser.Instructions = append(parser.Instructions, ins)
	parser.Len += 2
	parser.Position += 3
	return nil
}

// beq, bne, bgt, bge, blt, ble
func (parser *Parser) parseBranch() error {
	current_token, _ := parser.getCurrentToken()
	ins := Instruction{current_token.Value}

	arg1 := parser.getNextXToken(1)
	if arg1.TokenType != REGISTER {
		return invalidRegister(arg1)
	}

	virgula1 := parser.getNextXToken(2)
	if virgula1.TokenType != VIRGULA {
		return virgulaMissing(virgula1.Line, virgula1.Column)
	}

	arg2 := parser.getNextXToken(3)
	if arg2.TokenType != REGISTER {
		return invalidRegister(arg2)
	}

	virgula2 := parser.getNextXToken(4)
	if virgula2.TokenType != VIRGULA {
		return virgulaMissing(virgula2.Line, virgula2.Column)
	}

	arg3 := parser.getNextXToken(5)
	if arg3.TokenType != LABEL_REF {
		return invalidLabel(arg3)
	}

	nl := parser.getNextXToken(6)
	if nl.TokenType != NEW_LINE {
		return toMuchArgs(current_token, 3)
	}

	ins = append(ins, arg1.Value)
	ins = append(ins, arg2.Value)
	ins = append(ins, arg3.Value)
	parser.Instructions = append(parser.Instructions, ins)
	parser.Len += 2
	parser.Position += 7
	return nil
}

// lw, lb, sw, sb, lv, sv, lrw, lrb
// lw $r1, 0(*label)
// lw $r1, 0($r2)
// lw $r1, $r2($r3)
func (parser *Parser) parseMemorie() error {
	current_token, _ := parser.getCurrentToken()
	ins := Instruction{current_token.Value}

	arg1 := parser.getNextXToken(1)
	if arg1.TokenType != REGISTER {
		return invalidRegister(arg1)
	}

	virgula1 := parser.getNextXToken(2)
	if virgula1.TokenType != VIRGULA {
		return virgulaMissing(virgula1.Line, virgula1.Column)
	}

	offset := parser.getNextXToken(3)
	if offset.TokenType != NUMBER && offset.TokenType != REGISTER {
		return invalidRegOrLabel(offset)
	}

	open_paren := parser.getNextXToken(4)
	if open_paren.TokenType != OPEN_PAREN {
		return openParenMissing(open_paren.Line, open_paren.Column)
	}

	arg2 := parser.getNextXToken(5)
	if arg2.TokenType != REGISTER && arg2.TokenType != LABEL_REF {
		return invalidRegister(arg2)
	}

	close_paren := parser.getNextXToken(6)
	if close_paren.TokenType != CLOSE_PAREN {
		return closeParenMissing(close_paren.Line, close_paren.Column)
	}

	nl := parser.getNextXToken(7)
	if nl.TokenType != NEW_LINE {
		return toMuchArgs(current_token, 3)
	}

	ins = append(ins, arg1.Value)
	ins = append(ins, offset.Value)
	ins = append(ins, arg2.Value)
	parser.Instructions = append(parser.Instructions, ins)
	parser.Len += 2
	parser.Position += 8
	return nil
}

// inc, dec
func (parser *Parser) parseIncDec() error {
	current_token, _ := parser.getCurrentToken()
	ins := Instruction{current_token.Value}

	arg := parser.getNextXToken(1)
	if arg.TokenType != REGISTER {
		return invalidRegister(arg)
	}

	nl := parser.getNextXToken(2)
	if nl.TokenType != NEW_LINE {
		return toMuchArgs(current_token, 1)
	}

	ins = append(ins, arg.Value)
	parser.Instructions = append(parser.Instructions, ins)
	parser.Len += 2
	parser.Position += 3
	return nil
}

// Syscall, Return
func (parser *Parser) parseSyscallReturn() error {
	current_token, _ := parser.getCurrentToken()
	ins := Instruction{current_token.Value}

	nl := parser.getNextXToken(1)
	if nl.TokenType != NEW_LINE {
		return toMuchArgs(current_token, 1)
	}

	parser.Instructions = append(parser.Instructions, ins)
	parser.Len += 2
	parser.Position += 2
	return nil
}

func (parser *Parser) parseMove() error {
	current_token, _ := parser.getCurrentToken()
	ins := Instruction{current_token.Value}

	arg1 := parser.getNextXToken(1)
	if arg1.TokenType != REGISTER {
		return invalidRegister(arg1)
	}

	virgula1 := parser.getNextXToken(2)
	if virgula1.TokenType != VIRGULA {
		return virgulaMissing(virgula1.Line, virgula1.Column)
	}

	arg2 := parser.getNextXToken(3)
	if arg2.TokenType != REGISTER {
		return invalidRegister(arg2)
	}

	nl := parser.getNextXToken(4)
	if nl.TokenType != NEW_LINE {
		return toMuchArgs(current_token, 1)

	}

	ins = append(ins, arg1.Value)
	ins = append(ins, arg2.Value)
	parser.Instructions = append(parser.Instructions, ins)
	parser.Len += 2
	parser.Position += 5
	return nil
}

func (parser *Parser) parseLi() error {
	current_token, _ := parser.getCurrentToken()
	ins := Instruction{current_token.Value}

	arg1 := parser.getNextXToken(1)
	if arg1.TokenType != REGISTER {
		return invalidRegister(arg1)
	}

	virgula1 := parser.getNextXToken(2)
	if virgula1.TokenType != VIRGULA {
		return virgulaMissing(virgula1.Line, virgula1.Column)
	}

	arg2 := parser.getNextXToken(3)
	if arg2.TokenType != NUMBER {
		return invalidImmediate(arg2)
	}

	nl := parser.getNextXToken(4)
	if nl.TokenType != NEW_LINE {
		return toMuchArgs(current_token, 1)
	}

	ins = append(ins, arg1.Value)
	ins = append(ins, arg2.Value)
	parser.Instructions = append(parser.Instructions, ins)
	parser.Len += 2
	parser.Position += 5
	return nil
}

func (parser *Parser) parseLa() error {
	current_token, _ := parser.getCurrentToken()
	ins := Instruction{current_token.Value}

	arg1 := parser.getNextXToken(1)
	if arg1.TokenType != REGISTER {
		return invalidRegister(arg1)
	}

	virgula1 := parser.getNextXToken(2)
	if virgula1.TokenType != VIRGULA {
		return virgulaMissing(virgula1.Line, virgula1.Column)
	}

	arg2 := parser.getNextXToken(3)
	if arg2.TokenType != LABEL_REF {
		return invalidLabel(arg2)
	}

	nl := parser.getNextXToken(4)
	if nl.TokenType != NEW_LINE {
		return toMuchArgs(current_token, 1)
	}

	ins = append(ins, arg1.Value)
	ins = append(ins, arg2.Value)
	parser.Instructions = append(parser.Instructions, ins)
	parser.Len += 2
	parser.Position += 5
	return nil
}

func (parser *Parser) parseRand() error {
	current_token, _ := parser.getCurrentToken()
	ins := Instruction{current_token.Value}

	arg1 := parser.getNextXToken(1)
	if arg1.TokenType != REGISTER {
		return invalidRegister(arg1)
	}

	nl := parser.getNextXToken(2)
	if nl.TokenType != NEW_LINE {
		return toMuchArgs(current_token, 1)
	}

	ins = append(ins, arg1.Value)
	parser.Instructions = append(parser.Instructions, ins)
	parser.Len += 2
	parser.Position += 3
	return nil
}

func (parser *Parser) parseJumpLabelDef() error {
	current_token, _ := parser.getCurrentToken()
	parser.JumpLabels[current_token.Value] = parser.Len
	parser.Position += 1
	return nil
}

func (parser *Parser) parseStrings() error {
	current, err := parser.getCurrentToken()
	label := current.Value
	parser.RomLabels[label] = parser.Len - parser.Gp
	// println(parser.Len - parser.Gp)
	parser.Position += 1
	current, _ = parser.getCurrentToken()

	// fmt.Println(current)
	ins := Instruction{"STR"}
	str := ""
	for current.TokenType != NEW_LINE {
		escape := false
		for _, char := range current.Value {
			// print(string(char))
			if char == '\\' {
				if escape == false {
					escape = true
				} else {
					str += "\\"
				}
			} else {
				if escape == false {
					str += string(char)
				} else {
					switch char {
					case 'n':
						str += "\n"
					case '"':
						str += "\""
					case '\\':
						str += "\\"
					case '0':
						str += "\x00"
					default:
						return fmt.Errorf("Character \"%s\" de escape não reconhecido. Linha %d Coluna %d",
							string(char), current.Line+1, current.Column)
					}
					escape = false
				}
			}
		}

		// fmt.Printf("-%s-", str)
		parser.Position += 1
		current, err = parser.getCurrentToken()
		if err != nil {
			panic(err)
		}
	}
	ins = append(ins, str)
	parser.Len += len(str)
	parser.Instructions = append(parser.Instructions, ins)
	return nil
}

func ConvertWithOverflowAny(value string, t string) (string, error) {
	num, err := strconv.ParseInt(value, 0, 64)
	if err != nil {
		return "", fmt.Errorf("Erro ao ler inteiro")
	}

	switch t {
	case "int8":
		return strconv.Itoa(int(int8(num))), nil
	case "uint8":
		// println(int(uint8(num)))
		return strconv.Itoa(int(uint8(num))), nil
	case "int16":
		return strconv.Itoa(int(int16(num))), nil
	case "uint16":
		return strconv.Itoa(int(uint16(num))), nil
	default:
		return "", fmt.Errorf("tipo \"%s\" não é um tipo de inteiro valido", t)
	}
}

func (parser *Parser) parseInt(t string) error {
	current, err := parser.getCurrentToken()
	parser.RomLabels[current.Value] = parser.Len - parser.Gp
	parser.Position += 1

	total := 0
	current, err = parser.getCurrentToken()
	if err != nil {
		panic(err)
	}

	ins := Instruction{strings.ToUpper(t)}
	for current.TokenType != NEW_LINE {
		num, err := ConvertWithOverflowAny(current.Value, t)
		if err != nil {
			return err
		}

		ins = append(ins, num)
		parser.Position += 1
		current, err = parser.getCurrentToken()
		if err != nil {
			panic(err)
		}
		total += 1
	}

	switch t {
	case "int8", "uint8":
		parser.Len += total
	case "int16", "uint16":
		parser.Len += total * 2
	}

	parser.Instructions = append(parser.Instructions, ins)
	return nil
}

func (parser *Parser) Parse() ([]Instruction, int) {
	currentToken, err := parser.getCurrentToken()
	if err != nil {
		println(err.Error())
		return nil, -1
	}

	// Linhas de #define e comentários antes do .text viram linhas em branco
	for currentToken.TokenType == NEW_LINE {
		parser.Position += 1
		currentToken, err = parser.getCurrentToken()
		if err != nil {
			println(err.Error())
			return nil, -1
		}
	}

	if (currentToken.TokenType != SECTION) || (currentToken.Value != ".text") {
		println("Seção de text deve ser devinida no começo do arquivo")
		return nil, -1
	}

	validos := map[TokenType]bool{
		SECTION:      true,
		INSTRUCTION:  true,
		LABEL_DEF:    true,
		NEW_LINE:     true,
		STR_LABEL:    true,
		INT8_LABEL:   true,
		INT16_LABEL:  true,
		UINT8_LABEL:  true,
		UINT16_LABEL: true,
	}

	parser.Position += 2
	for parser.Position < len(parser.Tokens) {
		currentToken, err = parser.getCurrentToken()
		// fmt.Println(currentToken)

		if !validos[currentToken.TokenType] {
			fmt.Printf("Valor inesperado na linha: %d coluna: %d, %s.\n",
				currentToken.Line, currentToken.Column, currentToken.TokenType)
			return nil, -1
		}

		var err error = nil
		if RegistersInstructions[currentToken.Value] {
			err = parser.parseRegisterInstruciton()
		} else if ImmediateInstructions[currentToken.Value] {
			err = parser.parseImediateInstruciton()
		} else if JumpInstructions[currentToken.Value] {
			err = parser.parseJump()
		} else if BranchInstructions[currentToken.Value] {
			err = parser.parseBranch()
		} else if MemoriInstructions[currentToken.Value] {
			err = parser.parseMemorie()
		} else if currentToken.Value == "inc" || currentToken.Value == "dec" {
			err = parser.parseIncDec()
		} else if currentToken.Value == "return" || currentToken.Value == "syscall" {
			err = parser.parseSyscallReturn()
		} else if currentToken.Value == "move" {
			err = parser.parseMove()
		} else if currentToken.Value == "li" {
			err = parser.parseLi()
		} else if currentToken.Value == "la" {
			err = parser.parseLa()
		} else if currentToken.Value == "rand" {
			err = parser.parseRand()
		} else if currentToken.TokenType == LABEL_DEF {
			err = parser.parseJumpLabelDef()
		} else if currentToken.TokenType == NEW_LINE {
			parser.Position += 1
		} else if currentToken.TokenType == SECTION {
			parser.Gp = parser.Len
			parser.Position += 1
		} else if currentToken.TokenType == STR_LABEL {
			err = parser.parseStrings()
		} else if currentToken.TokenType == INT8_LABEL {
			err = parser.parseInt("int8")
		} else if currentToken.TokenType == UINT8_LABEL {
			err = parser.parseInt("uint8")
		} else if currentToken.TokenType == INT16_LABEL {
			err = parser.parseInt("int16")
		} else if currentToken.TokenType == UINT16_LABEL {
			err = parser.parseInt("uint16")
		} else {
			err = fmt.Errorf("Instrução ou Seção invalida em linha %d coluna %d.",
				currentToken.Line, currentToken.Column)
		}

		if err != nil {
			fmt.Println(err)
			return nil, -1
		}
	}

	return parser.Instructions, parser.Gp
}
