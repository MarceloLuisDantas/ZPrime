package assembler

import (
	"fmt"
	"strings"
)

const DEFINE = "#define"

type Preprocessor struct {
	Data    string
	defines map[string]string
}

func newPreprocessor(data string) *Preprocessor {
	preprocessor := Preprocessor{data, make(map[string]string)}
	return &preprocessor
}

// Retorna o indice logo após a string que começa em line[start]
func stringEnd(line string, start int) int {
	for i := start + 1; i < len(line); i++ {
		if line[i] == '\\' {
			i++
		} else if line[i] == '"' {
			return i + 1
		}
	}
	return len(line)
}

// Remove o comentário da linha, ignorando '#' dentro de strings
func removeComment(line string) string {
	for i := 0; i < len(line); i++ {
		if line[i] == '"' {
			i = stringEnd(line, i) - 1
		} else if line[i] == '#' {
			return line[:i]
		}
	}
	return line
}

func isDefine(line string) bool {
	line = strings.TrimLeft(line, " \t")
	if !strings.HasPrefix(line, DEFINE) {
		return false
	}

	// "#defineX" é só um comentário
	rest := line[len(DEFINE):]
	return rest == "" || rest[0] == ' ' || rest[0] == '\t'
}

func isValidDefineName(name string) bool {
	if name[0] >= '0' && name[0] <= '9' {
		return false
	}

	for i := 0; i < len(name); i++ {
		if !isValidCharacterToIdentifier(name[i]) {
			return false
		}
	}
	return true
}

// #define NOME valor
func (preprocessor *Preprocessor) handleDefine(line string, line_number int) error {
	line = strings.TrimLeft(line, " \t")
	body := strings.TrimSpace(removeComment(line[len(DEFINE):]))
	if body == "" {
		return fmt.Errorf("Nome esperado após #define. Linha %d", line_number)
	}

	end_name := strings.IndexAny(body, " \t")
	if end_name == -1 {
		return fmt.Errorf("Valor esperado para \"%s\". Linha %d", body, line_number)
	}

	name := body[:end_name]
	if !isValidDefineName(name) {
		return fmt.Errorf("Nome \"%s\" invalido para #define. Linha %d", name, line_number)
	}

	if KEY_WORDS[name] {
		return fmt.Errorf("Nome invalido para #define, '%s' é palavra reservada. Linha %d", name, line_number)
	}

	_, exists := preprocessor.defines[name]
	if exists {
		return fmt.Errorf("\"%s\" já foi definido. Linha %d", name, line_number)
	}

	// Expandir aqui permite usar defines anteriores no valor e impede ciclos
	value := strings.TrimSpace(body[end_name:])
	preprocessor.defines[name] = preprocessor.expand(value)
	return nil
}

// Palavras coladas em '$', '*' ou '.', ou seguidas de ':', são registradores,
// referencias a labels, seções, tipos ou definições de label.
func canReplace(line string, start, end int) bool {
	if start > 0 {
		last := line[start-1]
		if last == '$' || last == '*' || last == '.' {
			return false
		}
	}
	return end == len(line) || line[end] != ':'
}

// Substitui os defines da linha, mantendo strings e comentários intactos
func (preprocessor *Preprocessor) expand(line string) string {
	var result strings.Builder

	position := 0
	for position < len(line) {
		current := line[position]
		if current == '#' {
			result.WriteString(line[position:])
			break
		}

		if current == '"' {
			end := stringEnd(line, position)
			result.WriteString(line[position:end])
			position = end
			continue
		}

		if !isValidCharacterToIdentifier(current) {
			result.WriteByte(current)
			position += 1
			continue
		}

		start := position
		for position < len(line) && isValidCharacterToIdentifier(line[position]) {
			position += 1
		}

		word := line[start:position]
		value, exists := preprocessor.defines[word]
		if exists && canReplace(line, start, position) {
			result.WriteString(value)
		} else {
			result.WriteString(word)
		}
	}

	return result.String()
}

// A linha do #define vira uma linha em branco, para manter a numeração das
// linhas nas mensagens de erro do tokenizer e do parser.
func (preprocessor *Preprocessor) Process() (string, error) {
	lines := strings.Split(preprocessor.Data, "\n")
	for i, line := range lines {
		if isDefine(line) {
			err := preprocessor.handleDefine(line, i+1)
			if err != nil {
				return "", err
			}
			lines[i] = ""
		} else {
			lines[i] = preprocessor.expand(line)
		}
	}
	return strings.Join(lines, "\n"), nil
}
