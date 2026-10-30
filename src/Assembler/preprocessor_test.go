package assembler

import (
	"strings"
	"testing"
)

func TestProcess(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			"exemplo do hp",
			"#define HP 2($zero) # endereço do hp\n" +
				"lb $t0, HP # t0 = hp\n" +
				"addi $t0, $t0, 10 # hp += 10\n" +
				"sb $t0, HP\n",
			"\n" +
				"lb $t0, 2($zero) # t0 = hp\n" +
				"addi $t0, $t0, 10 # hp += 10\n" +
				"sb $t0, 2($zero)\n",
		},
		{
			"define com tab e indentação",
			"    #define\tMAX\t3600\t# pixels\nli $t0, MAX\n",
			"\nli $t0, 3600\n",
		},
		{
			"define usando define anterior",
			"#define A 10\n#define B A\n#define C B\nli $t0, C\n",
			"\n\n\nli $t0, 10\n",
		},
		{
			"valor é fixado na definição",
			"#define B A\n#define A 10\nli $t0, B\n",
			"\n\nli $t0, A\n",
		},
		{
			"não substitui em comentário",
			"#define HP 2\n# HP\nli $t0, 1 # HP\n",
			"\n# HP\nli $t0, 1 # HP\n",
		},
		{
			"não substitui em string",
			"#define HP 2\ns: .string \"HP \\\" HP\" # HP\n",
			"\ns: .string \"HP \\\" HP\" # HP\n",
		},
		{
			"não substitui parte de outro nome",
			"#define HP 2\nli $t0, HP_MAX\nli $t0, MAX_HP\nli $t0, HP2\n",
			"\nli $t0, HP_MAX\nli $t0, MAX_HP\nli $t0, HP2\n",
		},
		{
			"não substitui label, registrador e seção",
			"#define HP 2\n#define t0 3\n#define text 4\n" +
				".text\nHP:\n    j *HP\n    li $t0, HP\n",
			"\n\n\n.text\nHP:\n    j *HP\n    li $t0, 2\n",
		},
		{
			"valor com string contendo #",
			"#define MSG \"a # b\" # comentário\ns: .string MSG\n",
			"\ns: .string \"a # b\"\n",
		},
		{
			"não é define",
			"#defineX 1\n# define Y 2\nli $t0, X\n",
			"#defineX 1\n# define Y 2\nli $t0, X\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := newPreprocessor(test.input).Process()
			if err != nil {
				t.Fatal(err)
			}
			if result != test.expected {
				t.Errorf("\nesperado: %q\nobtido:   %q", test.expected, result)
			}
		})
	}
}

func TestProcessErrors(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"sem nome", "\n#define\n", "Linha 2"},
		{"só comentário", "#define # HP 2\n", "Linha 1"},
		{"sem valor", "#define HP\n", "Valor esperado"},
		{"sem valor com comentário", "#define HP # hp\n", "Valor esperado"},
		{"redefinição", "#define HP 2\n#define HP 2\n", "já foi definido"},
		{"nome com digito", "#define 2HP 2\n", "invalido"},
		{"nome com caractere invalido", "#define H-P 2\n", "invalido"},
		{"palavra reservada", "#define add 2\n", "reservada"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := newPreprocessor(test.input).Process()
			if err == nil {
				t.Fatal("erro esperado")
			}
			if !strings.Contains(err.Error(), test.expected) {
				t.Errorf("erro %q não contém %q", err.Error(), test.expected)
			}
		})
	}
}

func TestProcessKeepsLineNumbers(t *testing.T) {
	input := "#define A 1\n\n#define B 2\nli $t0, A\n"
	result, err := newPreprocessor(input).Process()
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(result, "\n") != strings.Count(input, "\n") {
		t.Errorf("numero de linhas alterado: %q", result)
	}
}

func parseSource(t *testing.T, source string) ([]Instruction, int) {
	t.Helper()

	data, err := newPreprocessor(source).Process()
	if err != nil {
		t.Fatal(err)
	}

	tokenizer := newTokenizer(data)
	err = tokenizer.Tokenize()
	if err != nil {
		t.Fatal(err)
	}
	return newParser(tokenizer.tokens).Parse()
}

func TestDefineBeforeText(t *testing.T) {
	instructions, _ := parseSource(t, "# comentário\n#define HP 2($zero)\n\n.text\n    lb $t0, HP\n")
	if instructions == nil {
		t.Fatal("parser falhou")
	}
	if len(instructions) != 1 {
		t.Fatalf("esperado 1 instrução, obtido %d", len(instructions))
	}

	expected := Instruction{"lb", "$t0", "2", "$zero"}
	if strings.Join(instructions[0], " ") != strings.Join(expected, " ") {
		t.Errorf("esperado %v, obtido %v", expected, instructions[0])
	}
}

func TestFileMustStartWithText(t *testing.T) {
	instructions, _ := parseSource(t, ".data\n    a: .int8 1\n")
	if instructions != nil {
		t.Error("arquivo sem .text no começo deveria falhar")
	}
}
