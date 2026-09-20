package pc

import "os"

// Memory Card
// O total de memoria disponivel é 1MiB.
// Os 1MiB é dividido em 32 Saves de 64 KiB cada
// A primeira pagina é reservada para o sistema de arquivos

type SaveMetaData struct {
	Name string
	Used bool
}

type Save struct {
	SMD  SaveMetaData
	Data [65536]byte // 64 KiB
}

type MemoryCard struct {
	Saves [16]Save // 1 MiB dividido em 32 arquivos
}

func NewMemoryCard() *MemoryCard {
	var mc MemoryCard
	return &mc
}

func (mc *MemoryCard) Save(file_path string) error {
	file := make([]byte, 1048576)
	save := 0
	cell := 0
	for i := 0; i < 1048576; i++ {
		file[i] = mc.Saves[save].Data[cell]
		cell++
		if cell == 65536 {
			save++
			cell = 0
		}
	}
	return os.WriteFile(file_path, file, 0644)
}

func (mc *MemoryCard) Load(file []byte) {
	save := 0
	cell := 0
	println(save, 0)
	for i := range 1048576 {
		mc.Saves[save].Data[cell] = file[i]
		cell++
		if cell == 65536 {
			save++
			println(save, i)
			cell = 0
		}
	}
}
