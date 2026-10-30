package assembler

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func replaceLabelsRef(ins []Instruction, parser *Parser) error {
	for i, value := range ins {
		var label = ""
		op := value[0]
		if op == "j" || op == "jal" {
			label = value[1]
		} else if op == "la" {
			label = value[2]
		} else if op[0] == 'b' {
			label = value[3]
		} else if op == "lrw" || op == "lrb" {
			if value[3][0] != '$' {
				label = value[3]
			}
		}

		if label != "" {
			if op == "j" || op == "jal" || op[0] == 'b' {
				line, exists := parser.JumpLabels[label]
				if !exists {
					return fmt.Errorf("Label \"%s\" não encontrando", label)
				}

				if op == "j" || op == "jal" {
					ins[i][1] = strconv.Itoa(line)
				} else {
					ins[i][3] = strconv.Itoa(line)
				}
			} else {
				line, exists := parser.RomLabels[label]
				if !exists {
					return fmt.Errorf("Label \"%s\" não encontrando", label)
				}

				if op == "la" {
					ins[i][2] = strconv.Itoa(line)
				} else {
					ins[i][3] = strconv.Itoa(line)
				}
			}
			label = ""
		}
	}
	return nil
}

func exportFile(instructions []Instruction, file_name string, gp int) error {
	file, err := os.Create(file_name)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Fprintf(file, "%s", fmt.Sprintf("%d\n", gp))

	for _, ins := range instructions {

		switch ins[0] {
		case "STR":
			value := ins[1]
			// println(value[len(value)-1])
			for _, char := range value {
				fmt.Fprintf(file, "%s", fmt.Sprintf("%d\n", char))
			}
		case "INT8", "UINT8":
			values := ins[1:]
			for _, v := range values {
				value, _ := strconv.ParseInt(v, 10, 16)
				fmt.Fprintf(file, "%s", fmt.Sprintf("%d\n", uint8(value)))
			}

		case "INT16", "UINT16":
			values := ins[1:]
			for _, v := range values {
				value, _ := strconv.ParseInt(v, 10, 32)
				numui16 := uint16(value)
				lsb := uint8(numui16)
				msb := uint8(numui16 >> 8)
				fmt.Fprintf(file, "%s", fmt.Sprintf("%d\n", msb))
				fmt.Fprintf(file, "%s", fmt.Sprintf("%d\n", lsb))
			}
		case "syscall", "return":
			fmt.Fprintln(file, ins[0][0:3])
			fmt.Fprintln(file, ins[0][3:])

		default:
			fmt.Fprintf(file, "%s", fmt.Sprintf("%s\n", ins[0]))

			values := ins[1:]
			for _, v := range values {
				fmt.Fprintf(file, "%s", fmt.Sprintf("%s ", v))
			}
			fmt.Fprintf(file, "\n")
		}
	}

	return nil
}

func Assembler(data, file_name string) {

	data, err := newPreprocessor(data).Process()
	if err != nil {
		println("Erro ao preprocessar: ", err.Error())
		return
	}
	println("Preprocessor: OK")

	tk := newTokenizer(data)
	err = tk.Tokenize()
	if err != nil {
		println("Erro ao tokenizar: ", err.Error())
		return
	}
	println("Tokenizer: OK")

	parser := newParser(tk.tokens)
	instructions, gp := parser.Parse()
	if instructions == nil {
		println("Erro ao realisar parser")
		return
	}
	println("Parser: OK")

	err = replaceLabelsRef(instructions, parser)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Assembler: OK")

	final_name := strings.Split(file_name, ".")[0]
	final_name += ".krom"
	err = exportFile(instructions, final_name, gp)
	if err != nil {
		println(err)
	} else {
		fmt.Printf("ROM exportada em \"%s\"\n", final_name)
	}
}
