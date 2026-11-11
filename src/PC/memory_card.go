package pc

import (
	"bytes"
	"fmt"
	"os"
)

// Memory Card
// O total de memoria disponivel é 512 KiB, dividido em 32 saves de 16 KiB.
// O save 0 é reservado para o sistema de arquivos, sobrando os saves 1..31
// para os programas. Cada programa usa no maximo 1 save, identificado por
// uma string de até 16 caracteres dada pelo desenvolvedor.
//
// Tabela do sistema de arquivos (inicio do save 0), 1 entrada por save:
// [save_id 1B][em uso 1B][id do programa 17B]
//
// Um memory card zerado é valido, os ids da tabela são preenchidos ao abrir.

const MEMORY_CARD_SIZE = 524288 // 512 KiB
const SAVE_SIZE = 16384         // 16 KiB
const SAVES = 32
const PROGRAM_ID_SIZE = 17 // 16 caracteres + nulo
const TABLE_ENTRY_SIZE = 2 + PROGRAM_ID_SIZE

type MemoryCard struct {
	file *os.File
	data []byte
}

// Cria um memory card vazio e já formatado
func CreateMemoryCard(file_path string) error {
	data := make([]byte, MEMORY_CARD_SIZE)
	for save_id := 1; save_id < SAVES; save_id++ {
		data[tableEntry(save_id)] = uint8(save_id)
	}
	return os.WriteFile(file_path, data, 0644)
}

func OpenMemoryCard(file_path string) (*MemoryCard, error) {
	data, err := os.ReadFile(file_path)
	if err != nil {
		return nil, err
	}

	if len(data) != MEMORY_CARD_SIZE {
		return nil, fmt.Errorf("Memory Card \"%s\" deve ter %d bytes, mas tem %d", file_path, MEMORY_CARD_SIZE, len(data))
	}

	file, err := os.OpenFile(file_path, os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	mc := MemoryCard{file, data}
	err = mc.checkTable()
	if err != nil {
		file.Close()
		return nil, err
	}
	return &mc, nil
}

func (mc *MemoryCard) Close() error {
	return mc.file.Close()
}

func tableEntry(save_id int) int {
	return (save_id - 1) * TABLE_ENTRY_SIZE
}

func (mc *MemoryCard) checkTable() error {
	for save_id := 1; save_id < SAVES; save_id++ {
		entry := tableEntry(save_id)
		id := mc.data[entry]
		used := mc.data[entry+1]

		// Entrada nunca formatada
		if id == 0 && used == 0 {
			err := mc.writeByte(entry, uint8(save_id))
			if err != nil {
				return err
			}
			continue
		}

		if int(id) != save_id || used > 1 {
			return fmt.Errorf("Memory Card corrompido, entrada %d da tabela invalida", save_id)
		}
	}
	return nil
}

// Escreve na memoria e no arquivo, para que nada se perca caso o programa
// termine sem aviso
func (mc *MemoryCard) writeByte(position int, value uint8) error {
	mc.data[position] = value
	_, err := mc.file.WriteAt([]byte{value}, int64(position))
	return err
}

func (mc *MemoryCard) programId(save_id int) []byte {
	entry := tableEntry(save_id) + 2
	id := mc.data[entry : entry+PROGRAM_ID_SIZE]
	end := bytes.IndexByte(id, 0)
	if end == -1 {
		return id
	}
	return id[:end]
}

// Retorna o save do programa, reservando um save livre caso o programa ainda
// não tenha um. Retorna 0 caso todos os saves estejam em uso.
func (mc *MemoryCard) Request(program_id []byte) (int, error) {
	if len(program_id) == 0 || len(program_id) >= PROGRAM_ID_SIZE {
		return 0, fmt.Errorf("id do programa deve ter entre 1 e %d caracteres", PROGRAM_ID_SIZE-1)
	}

	free := 0
	for save_id := 1; save_id < SAVES; save_id++ {
		used := mc.data[tableEntry(save_id)+1] == 1
		if used && bytes.Equal(mc.programId(save_id), program_id) {
			return save_id, nil
		}
		if !used && free == 0 {
			free = save_id
		}
	}

	if free == 0 {
		return 0, nil
	}

	entry := tableEntry(free)
	for i := range PROGRAM_ID_SIZE {
		var char uint8 = 0
		if i < len(program_id) {
			char = program_id[i]
		}
		err := mc.writeByte(entry+2+i, char)
		if err != nil {
			return 0, err
		}
	}

	err := mc.writeByte(entry+1, 1)
	if err != nil {
		return 0, err
	}
	return free, nil
}

func (mc *MemoryCard) savePosition(save_id int, address uint16) (int, error) {
	if save_id < 1 || save_id >= SAVES {
		return 0, fmt.Errorf("save %d não existe, saves vão de 1 a %d", save_id, SAVES-1)
	}

	if mc.data[tableEntry(save_id)+1] != 1 {
		return 0, fmt.Errorf("save %d não foi requisitado", save_id)
	}

	if int(address) >= SAVE_SIZE {
		return 0, fmt.Errorf("endereço %d fora do save, saves vão de 0 a %d", address, SAVE_SIZE-1)
	}

	return save_id*SAVE_SIZE + int(address), nil
}

func (mc *MemoryCard) Read(save_id int, address uint16) (uint8, error) {
	position, err := mc.savePosition(save_id, address)
	if err != nil {
		return 0, err
	}
	return mc.data[position], nil
}

func (mc *MemoryCard) Write(save_id int, address uint16, value uint8) error {
	position, err := mc.savePosition(save_id, address)
	if err != nil {
		return err
	}
	return mc.writeByte(position, value)
}
