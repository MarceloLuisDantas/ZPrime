package assembler

import (
	"strings"
	"testing"
)

func TestSaveInstructions(t *testing.T) {
	source := ".text\n" +
		"    request_save *program_name\n" +
		"    request_save $t0\n" +
		"    read_save $t1, $t0\n" +
		"    write_save $t1, $t0, $rt\n" +
		".data\n" +
		"    other: .int8 1 2\n" +
		"    program_name: .string \"programa1\\0\"\n"

	tokenizer := newTokenizer(source)
	err := tokenizer.Tokenize()
	if err != nil {
		t.Fatal(err)
	}

	parser := newParser(tokenizer.tokens)
	instructions, _ := parser.Parse()
	if instructions == nil {
		t.Fatal("parser falhou")
	}

	err = replaceLabelsRef(instructions, parser)
	if err != nil {
		t.Fatal(err)
	}

	expected := []string{
		"request_save 2",
		"request_save $t0",
		"read_save $t1 $t0",
		"write_save $t1 $t0 $rt",
	}
	for i, ins := range expected {
		if strings.Join(instructions[i], " ") != ins {
			t.Errorf("esperado %q, obtido %q", ins, strings.Join(instructions[i], " "))
		}
	}
}

func TestSaveInstructionsErrors(t *testing.T) {
	tests := []string{
		".text\n    request_save 10\n",
		".text\n    request_save *a, $t0\n",
		".text\n    read_save $t1\n",
		".text\n    read_save $t1, *a\n",
		".text\n    write_save $t1, $t0\n",
		".text\n    write_save $t1, $t0, 1\n",
	}

	for _, source := range tests {
		tokenizer := newTokenizer(source)
		if tokenizer.Tokenize() != nil {
			continue
		}

		instructions, _ := newParser(tokenizer.tokens).Parse()
		if instructions != nil {
			t.Errorf("deveria falhar: %q", source)
		}
	}
}
