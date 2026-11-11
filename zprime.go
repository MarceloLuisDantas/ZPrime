package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	assembler "zprime/src/Assembler"
	pc "zprime/src/PC"
)

func load_file(file_name string) string {
	data, err := os.ReadFile(file_name)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func usage() {
	fmt.Println("Use: zprime")
	fmt.Println("    Compile: zprime as file.asm")
	fmt.Println("    Run: zprime run file.krom memorycard.zsave [optional]")
	fmt.Println("         Without memorycard.zsave, save.zsave next to zprime is used")
	os.Exit(0)
}

// Retorna o memory card dado na linha de comando, ou o save.zsave que fica
// na mesma pasta do executavel
func memoryCardPath() (string, bool) {
	if len(os.Args) >= 4 {
		return os.Args[3], false
	}

	exe, err := os.Executable()
	if err != nil {
		return "save.zsave", true
	}
	return filepath.Join(filepath.Dir(exe), "save.zsave"), true
}

func main() {

	if len(os.Args) < 3 {
		usage()
	}

	option := os.Args[1]
	if option != "as" && option != "run" {
		usage()
	}

	file_name := os.Args[2]
	data := load_file(file_name)
	switch option {
	case "as":
		if data[len(data)-1] != '\n' {
			data += "\n"
		}
		assembler.Assembler(data, file_name)
	case "run":
		cpu, err := pc.NewCPU(strings.Split(data, "\n"))
		if err != nil {
			println(err.Error())
			return
		}

		card_path, default_card := memoryCardPath()
		mc, err := pc.OpenMemoryCard(card_path)
		if errors.Is(err, os.ErrNotExist) {
			println("Memory Card: ", card_path, " not found.")
			println("Want to create Memory Card file? [Y\\N]")
			var escolha string
			fmt.Scan(&escolha)
			if escolha == "Y" {
				err = pc.CreateMemoryCard(card_path)
				if err != nil {
					panic(err)
				}
				mc, err = pc.OpenMemoryCard(card_path)
			} else if default_card {
				println("Running without Memory Card.")
				err = nil
			} else {
				return
			}
		}

		if err != nil {
			println(err.Error())
			return
		}

		if mc != nil {
			defer mc.Close()
			cpu.SetMemoryCard(mc)
		}

		// cpu.ShowRom()
		cpu.Run()
	}
}
