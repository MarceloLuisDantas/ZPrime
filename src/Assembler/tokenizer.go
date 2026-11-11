package assembler

import (
	"fmt"
	"strconv"
)

var KEY_WORDS = map[string]bool{
	"add": true, "addi": true, "sub": true, "subi": true,
	"mult": true, "multi": true, "div": true, "divi": true,
	"addu": true, "addui": true, "subu": true, "subui": true,
	"multu": true, "multui": true, "divu": true, "divui": true,
	"or": true, "ori": true, "and": true, "andi": true,
	"slt": true, "slti": true, "sltu": true, "sltui": true,
	"sll": true, "srl": true,
	"li": true, "la": true, "move": true,
	"j": true, "jr": true, "jal": true,
	"beq": true, "bne": true, "bgt": true,
	"bge": true, "blt": true, "ble": true,
	"lw": true, "lb": true,
	"sw": true, "sb": true,
	"lvr": true, "svr": true,
	"lrw": true, "lrb": true,
	"inc": true, "dec": true,
	"syscall": true, "return": true, "rand": true,
	"nop": true, "request_save": true, "read_save": true, "write_save": true,
}

var SECTIONS = map[string]bool{
	".text": true, ".data": true,
}

var TYPES = map[string]bool{
	".int8": true, ".uint8": true,
	".int16": true, ".uint16": true,
	".string": true,
}

var REGISTRERS = map[string]bool{
	"$zero": true, "$t0": true, "$t1": true, "$t2": true, "$t3": true, "$t4": true, "$t5": true,
	"$rt": true, "$sp": true, "$fp": true, "$sc": true, "$ra": true, "$pc": true, "$ir": true, "gp": true,
}

type TokenType string

const (
	IDENTIFIER   = "IDENTIFIER"
	INSTRUCTION  = "INSTRUCTION"
	REGISTER     = "REGISTER"
	LABEL_DEF    = "LABEL_DEF"
	LABEL_REF    = "LABEL_REF"
	STR_LABEL    = "STR_LABEL"
	INT8_LABEL   = "INT8_LABEL"
	UINT8_LABEL  = "UINT8_LABEL"
	INT16_LABEL  = "INT16_LABEL"
	UINT16_LABEL = "UINT16_LABEL"
	SECTION      = "SECTION"
	TYPE         = "TYPE"
	STRING       = "STRING"
	NUMBER       = "NUMBER"
	VIRGULA      = "VIRGULA"
	NEW_LINE     = "NEW_LINE"
	OPEN_PAREN   = "OPEN_PAREN"
	CLOSE_PAREN  = "CLOSE_PAREN"
)

type Token struct {
	TokenType TokenType
	Value     string
	Line      int
	Column    int
}

func newToken(tt TokenType, value string, line int, column int) *Token {
	t := Token{tt, value, line, column}
	return &t
}

type Tokenizer struct {
	Data     string
	position int
	column   int
	line     int
	len      int
	tokens   []Token
}

func newTokenizer(data string) *Tokenizer {
	tokenizer := Tokenizer{data, 0, 0, 0, 0, []Token{}}
	return &tokenizer
}

func (tokenizer *Tokenizer) addToken(token_type TokenType, value string) {
	token := newToken(token_type, value, tokenizer.line+1, tokenizer.column+1)
	tokenizer.tokens = append(tokenizer.tokens, *token)
	tokenizer.len += 1
}

func (tokenizer *Tokenizer) advance() {
	tokenizer.position += 1
	tokenizer.column += 1
}

// func (Tokenizer *Tokenizer) advanceX(value int) {
// 	Tokenizer.position += value
// 	Tokenizer.column += value
// }

func (tokenizer *Tokenizer) nextLine() {
	tokenizer.position += 1
	tokenizer.column = 0
	tokenizer.line += 1
}

func (tokenizer *Tokenizer) getCurrentChar() byte {
	return tokenizer.Data[tokenizer.position]
}

func (tokenizer *Tokenizer) setLastTokenType(tk TokenType) {
	tokenizer.tokens[tokenizer.len-1].TokenType = tk
}

func (tokenizer *Tokenizer) handleSpace() {
	tokenizer.advance()
}

func (tokenizer *Tokenizer) handleNewLine() {
	tokenizer.addToken(NEW_LINE, "\n")
	tokenizer.nextLine()
}

func (tokenizer *Tokenizer) handleVirgula() error {
	last_token := tokenizer.tokens[tokenizer.len-1]

	if last_token.TokenType == REGISTER || last_token.TokenType == NUMBER || last_token.TokenType == LABEL_REF {
		tokenizer.addToken(VIRGULA, ",")
	} else {
		return fmt.Errorf("Virgula servem para separar argumentos em funções. Linha: %d Coluna %d", tokenizer.line+1, tokenizer.column+1)
	}

	tokenizer.advance()
	return nil
}

func (tokenizer *Tokenizer) handleComment() {
	for tokenizer.getCurrentChar() != '\n' {
		tokenizer.advance()
	}
	// tokenizer.nextLine()
}

func isValidCharacterToIdentifier(s byte) bool {
	return (s >= 'a' && s <= 'z') || (s >= 'A' && s <= 'Z') || s == '_' || (s >= '0' && s <= '9')
}

func (tokenizer *Tokenizer) handleIdentifier() error {
	start := tokenizer.position

	current := tokenizer.getCurrentChar()
	for current != ' ' && current != ',' && current != ':' && current != '\n' {
		if !isValidCharacterToIdentifier(current) {
			return fmt.Errorf("Character '%c' invalido para identificador. Linha %d Coluna %d", current, tokenizer.line+1, tokenizer.column+1)
		}
		tokenizer.advance()
		current = tokenizer.getCurrentChar()
	}

	key_world := tokenizer.Data[start:tokenizer.position]
	if tokenizer.getCurrentChar() == ':' {
		if KEY_WORDS[key_world] {
			return fmt.Errorf("Nome invalido para label, '%s' é palavra reservada.", key_world)
		}
		tokenizer.addToken(LABEL_DEF, key_world)
		tokenizer.advance()
	} else if KEY_WORDS[key_world] {
		tokenizer.addToken(INSTRUCTION, key_world)
	} else {
		return fmt.Errorf("'%s' não é um comando valido. Linha %d Coluna %d", key_world, tokenizer.line+1, tokenizer.column+1)
	}
	return nil
}

func (tokenizer *Tokenizer) handleRegister() error {
	start := tokenizer.position
	tokenizer.advance() // Pula o $

	current := tokenizer.getCurrentChar()
	for current != ' ' && current != ',' && current != ')' && current != '(' && current != '\n' {
		if !isValidCharacterToIdentifier(current) {
			return fmt.Errorf("Character '%c' invalido ao referenciar registrador. Linha %d Coluna %d", current, tokenizer.line+1, tokenizer.column+1)
		}
		tokenizer.advance()
		current = tokenizer.getCurrentChar()
	}

	register := tokenizer.Data[start:tokenizer.position]
	if !REGISTRERS[register] {
		return fmt.Errorf("Registrador \"%s\" não é um registrador valido. Linha: %d, Coluna: %d", register, tokenizer.line+1, tokenizer.column+1)
	}

	tokenizer.addToken(REGISTER, register)
	return nil
}

func (tokenizer *Tokenizer) handleLabelRef() error {
	tokenizer.advance() // Pula o *
	start := tokenizer.position

	for tokenizer.getCurrentChar() != ' ' && tokenizer.getCurrentChar() != '\n' && tokenizer.getCurrentChar() != ')' {
		if !isValidCharacterToIdentifier(tokenizer.getCurrentChar()) {
			fmt.Printf("%c\n", tokenizer.getCurrentChar())
			return fmt.Errorf("Character invalido ao referenciar label. Linha %d Coluna %d", tokenizer.line+1, tokenizer.column+1)
		}
		tokenizer.advance()
	}

	label := tokenizer.Data[start:tokenizer.position]
	if label == "*" {
		return fmt.Errorf("Definição de referencia a label sem nome. Linha %d Coluna %d", tokenizer.line+1, tokenizer.column+1)
	}

	tokenizer.addToken(LABEL_REF, label)
	// println(label)
	return nil
}

func (tokenizer *Tokenizer) handleSectionOrType() error {
	start := tokenizer.position
	tokenizer.advance() // Pula o '.'

	var current byte
	for {
		current = tokenizer.getCurrentChar()
		if current == ' ' || current == '\n' {
			break
		}

		if !isValidCharacterToIdentifier(current) {
			return fmt.Errorf("Character '%c' invalido ao definir seção ou tipo. Linha %d Coluna %d", current, tokenizer.line+1, tokenizer.column+1)
		}
		tokenizer.advance()
	}

	token := tokenizer.Data[start:tokenizer.position]
	if SECTIONS[token] {
		tokenizer.addToken(SECTION, token)

		var current byte
		for {
			current = tokenizer.getCurrentChar()
			if current == '\n' {
				break
			}

			if current != ' ' {
				return fmt.Errorf("Declaração de seção não recebe argumentos. Linha %d Coluna %d", tokenizer.line+1, tokenizer.column+1)
			}
			tokenizer.advance()
		}

	} else if TYPES[token] {
		switch token {
		case ".string":
			tokenizer.setLastTokenType(STR_LABEL)
		case ".int8":
			tokenizer.setLastTokenType(INT8_LABEL)
		case ".int16":
			tokenizer.setLastTokenType(INT16_LABEL)
		case ".uint8":
			tokenizer.setLastTokenType(UINT8_LABEL)
		case ".uint16":
			tokenizer.setLastTokenType(UINT16_LABEL)
		}
	} else {
		return fmt.Errorf("Nome \"%s\" não é valido para definição de seção ou tipo. Linha %d, Coluna %d", token, tokenizer.line+1, tokenizer.column+1)
	}

	return nil
}

func (tokenizer *Tokenizer) handleNumber() error {
	start := tokenizer.position
	current := tokenizer.getCurrentChar()
	if current == '-' {
		tokenizer.advance()
		current = tokenizer.getCurrentChar()
	}

	if current < '0' || current > '9' {
		if tokenizer.Data[start] == '-' {
			return fmt.Errorf("Dígito esperado após '-'. Linha %d, Coluna %d", tokenizer.line, tokenizer.column)
		}
		return fmt.Errorf("Caractere '%c' não é um dígito. Linha %d, Coluna %d", current, tokenizer.line, tokenizer.column)
	}

	for tokenizer.position < len(tokenizer.Data) {
		current = tokenizer.getCurrentChar()
		if current < '0' || current > '9' {
			break
		}
		tokenizer.advance()
	}
	num := tokenizer.Data[start:tokenizer.position]

	value, err := strconv.Atoi(num)
	if err != nil {
		return fmt.Errorf("Valor invalido ao definir numero. Linha %d, Coluna %d", tokenizer.line+1, tokenizer.column+1)
	}

	tokenizer.addToken(NUMBER, strconv.Itoa(int(int16(value))))
	return nil
}

func (tokenizer *Tokenizer) handleStrings() error {
	tokenizer.advance() // Pula primeira aspas
	start := tokenizer.position

	var last_char byte
	current := tokenizer.getCurrentChar()
	for {
		if current == '\n' {
			return fmt.Errorf("Quebra de linha ao declara string. Linha %d Coluna %d", tokenizer.line+1, tokenizer.column+1)
		}

		tokenizer.advance()
		current = tokenizer.getCurrentChar()
		if current == '"' && last_char != '\\' {
			break
		}

		// println(current, last_char)
		last_char = current

	}

	str := tokenizer.Data[start:tokenizer.position]
	tokenizer.addToken(STRING, str)
	tokenizer.advance()
	return nil
}

func (tokenizer *Tokenizer) handleOpenParenteses() error {
	last_token := tokenizer.tokens[tokenizer.len-1]
	if last_token.TokenType != NUMBER && last_token.TokenType != REGISTER {
		return fmt.Errorf("Mal uso de abertura de parenteses. É esperado offset ou registrador como valor anterior. Linha %d Coluna %d", tokenizer.line+1, tokenizer.column+1)
	}

	tokenizer.advance()
	tokenizer.addToken(OPEN_PAREN, "(")
	return nil
}

func (tokenizer *Tokenizer) handleCloseParenteses() error {
	last_token := tokenizer.tokens[tokenizer.len-1]
	// println(last_token.Value)
	if last_token.TokenType != LABEL_REF && last_token.TokenType != REGISTER {
		return fmt.Errorf("Mal uso de fechamento de parenteses. É esperado label ou registrador como valor anterior. Linha %d Coluna %d", tokenizer.line+1, tokenizer.column+1)
	}

	second_last := tokenizer.tokens[tokenizer.len-2]
	if second_last.TokenType != OPEN_PAREN {
		return fmt.Errorf("Fechamento de parenteses sem abertura. Linha %d Coluna %d", tokenizer.line+1, tokenizer.column+1)
	}

	tokenizer.advance()
	tokenizer.addToken(CLOSE_PAREN, ")")
	return nil
}

func (tokenizer *Tokenizer) Tokenize() error {
	erros := 0
	var err error = nil
	for tokenizer.position != len(tokenizer.Data) {
		current := tokenizer.getCurrentChar()
		if current == ' ' || current == '\t' {
			tokenizer.handleSpace()
		} else if current == '\n' {
			tokenizer.handleNewLine()
		} else if current == '#' {
			tokenizer.handleComment()
		} else if current == ',' {
			err = tokenizer.handleVirgula()
		} else if (current >= 'a' && current <= 'z') || (current >= 'A' && current <= 'Z') || current == '_' {
			err = tokenizer.handleIdentifier()
		} else if current == '$' {
			err = tokenizer.handleRegister()
		} else if current == '*' {
			err = tokenizer.handleLabelRef()
		} else if current == '.' {
			err = tokenizer.handleSectionOrType()
		} else if current == '-' || (current >= '0' && current <= '9') {
			err = tokenizer.handleNumber()
		} else if current == '"' {
			err = tokenizer.handleStrings()
		} else if current == '(' {
			err = tokenizer.handleOpenParenteses()
		} else if current == ')' {
			err = tokenizer.handleCloseParenteses()
		} else {
			return fmt.Errorf("Character invalido na linha %d coluna %d", tokenizer.line+1, tokenizer.column+1)
		}

		if err != nil {
			println(err.Error())
			return fmt.Errorf("%d tokenizer erro(s)", erros)
		}
	}

	return nil
}
