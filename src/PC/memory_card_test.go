package pc

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func newTestCard(t *testing.T) (*MemoryCard, string) {
	t.Helper()

	file_path := filepath.Join(t.TempDir(), "card.zsave")
	err := CreateMemoryCard(file_path)
	if err != nil {
		t.Fatal(err)
	}

	mc, err := OpenMemoryCard(file_path)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { mc.Close() })
	return mc, file_path
}

func TestCreateMemoryCard(t *testing.T) {
	_, file_path := newTestCard(t)

	data, err := os.ReadFile(file_path)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) != MEMORY_CARD_SIZE {
		t.Fatalf("tamanho %d, esperado %d", len(data), MEMORY_CARD_SIZE)
	}

	for save_id := 1; save_id < SAVES; save_id++ {
		entry := tableEntry(save_id)
		if data[entry] != uint8(save_id) || data[entry+1] != 0 {
			t.Errorf("entrada %d: id %d, em uso %d", save_id, data[entry], data[entry+1])
		}
	}
}

func TestRequest(t *testing.T) {
	mc, _ := newTestCard(t)

	first, _ := mc.Request([]byte("programa1"))
	second, _ := mc.Request([]byte("programa2"))
	again, _ := mc.Request([]byte("programa1"))

	if first != 1 || second != 2 || again != 1 {
		t.Errorf("saves: %d %d %d, esperado 1 2 1", first, second, again)
	}

	// Prefixo de um id existente é outro programa
	prefix, _ := mc.Request([]byte("programa"))
	if prefix != 3 {
		t.Errorf("save do prefixo: %d, esperado 3", prefix)
	}
}

func TestRequestFullCard(t *testing.T) {
	mc, _ := newTestCard(t)

	for save_id := 1; save_id < SAVES; save_id++ {
		got, err := mc.Request([]byte{'p', uint8('A' + save_id)})
		if err != nil || got != save_id {
			t.Fatalf("save %d, esperado %d (%v)", got, save_id, err)
		}
	}

	got, err := mc.Request([]byte("sem_espaco"))
	if err != nil || got != 0 {
		t.Errorf("card cheio retornou %d (%v), esperado 0", got, err)
	}

	// Quem já tem save continua conseguindo
	got, _ = mc.Request([]byte{'p', 'B'})
	if got != 1 {
		t.Errorf("save existente com card cheio: %d, esperado 1", got)
	}
}

func TestRequestProgramIdSize(t *testing.T) {
	mc, _ := newTestCard(t)

	_, err := mc.Request([]byte("1234567890123456"))
	if err != nil {
		t.Errorf("id com 16 caracteres deveria ser aceito: %v", err)
	}

	_, err = mc.Request([]byte("12345678901234567"))
	if err == nil {
		t.Error("id com 17 caracteres deveria falhar")
	}

	_, err = mc.Request([]byte{})
	if err == nil {
		t.Error("id vazio deveria falhar")
	}
}

func TestReadWrite(t *testing.T) {
	mc, _ := newTestCard(t)
	save_id, _ := mc.Request([]byte("programa1"))

	err := mc.Write(save_id, 0, 42)
	if err != nil {
		t.Fatal(err)
	}
	err = mc.Write(save_id, SAVE_SIZE-1, 255)
	if err != nil {
		t.Fatal(err)
	}

	first, _ := mc.Read(save_id, 0)
	last, _ := mc.Read(save_id, SAVE_SIZE-1)
	if first != 42 || last != 255 {
		t.Errorf("lido %d e %d, esperado 42 e 255", first, last)
	}

	// O save seguinte não pode ser afetado
	other, _ := mc.Request([]byte("programa2"))
	value, _ := mc.Read(other, 0)
	if value != 0 {
		t.Errorf("save vizinho alterado: %d", value)
	}
}

func TestReadWriteErrors(t *testing.T) {
	mc, _ := newTestCard(t)
	save_id, _ := mc.Request([]byte("programa1"))

	tests := []struct {
		name     string
		save_id  int
		address  uint16
		expected string
	}{
		{"save do sistema", 0, 0, "não existe"},
		{"save inexistente", SAVES, 0, "não existe"},
		{"save não requisitado", save_id + 1, 0, "não foi requisitado"},
		{"endereço fora do save", save_id, SAVE_SIZE, "fora do save"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := mc.Read(test.save_id, test.address)
			if err == nil || !strings.Contains(err.Error(), test.expected) {
				t.Errorf("leitura: erro %v, esperado %q", err, test.expected)
			}

			err = mc.Write(test.save_id, test.address, 1)
			if err == nil || !strings.Contains(err.Error(), test.expected) {
				t.Errorf("escrita: erro %v, esperado %q", err, test.expected)
			}
		})
	}
}

func TestPersistence(t *testing.T) {
	mc, file_path := newTestCard(t)
	save_id, _ := mc.Request([]byte("programa1"))
	mc.Write(save_id, 100, 7)
	mc.Close()

	reopened, err := OpenMemoryCard(file_path)
	if err != nil {
		t.Fatal(err)
	}
	defer reopened.Close()

	again, _ := reopened.Request([]byte("programa1"))
	value, _ := reopened.Read(again, 100)
	if again != save_id || value != 7 {
		t.Errorf("save %d valor %d, esperado %d e 7", again, value, save_id)
	}
}

func TestOpenZeroedCard(t *testing.T) {
	file_path := filepath.Join(t.TempDir(), "zerado.zsave")
	err := os.WriteFile(file_path, make([]byte, MEMORY_CARD_SIZE), 0644)
	if err != nil {
		t.Fatal(err)
	}

	mc, err := OpenMemoryCard(file_path)
	if err != nil {
		t.Fatal(err)
	}
	mc.Close()

	data, _ := os.ReadFile(file_path)
	if data[tableEntry(31)] != 31 {
		t.Error("tabela não foi formatada no arquivo")
	}
}

func TestOpenInvalidCard(t *testing.T) {
	dir := t.TempDir()

	wrong_size := filepath.Join(dir, "antigo.zsave")
	os.WriteFile(wrong_size, make([]byte, 1048576), 0644)
	_, err := OpenMemoryCard(wrong_size)
	if err == nil || !strings.Contains(err.Error(), "deve ter") {
		t.Errorf("tamanho errado: %v", err)
	}

	corrupted := filepath.Join(dir, "corrompido.zsave")
	data := make([]byte, MEMORY_CARD_SIZE)
	data[tableEntry(2)] = 9
	os.WriteFile(corrupted, data, 0644)
	_, err = OpenMemoryCard(corrupted)
	if err == nil || !strings.Contains(err.Error(), "corrompido") {
		t.Errorf("tabela corrompida: %v", err)
	}
}
