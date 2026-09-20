package main

import (
	"fmt"
	"os"
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

func main() {

	if len(os.Args) < 3 {
		fmt.Println("Use: zprime")
		fmt.Println("    Compile: zprime as file.asm")
		fmt.Println("    Run: zprime run file.krom memorycard.zsave [optional]")
		os.Exit(0)
	}

	option := os.Args[1]
	if option != "as" && option != "run" {
		fmt.Println("Use: zprime")
		fmt.Println("    Compile: zprime as file.asm")
		fmt.Println("    Run: zprime run file.krom zprime.zsave [optional]")
		os.Exit(0)
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
		load_save := true
		if len(os.Args) < 4 {
			println("Save file not found; continue without Memory Card? [Y\\N]")
			var escolha string
			fmt.Scan(&escolha)
			if escolha != "Y" {
				return
			}
			load_save = false
		}
		cpu, err := pc.NewCPU(strings.Split(data, "\n"))
		if err != nil {
			println(err.Error())
			return
		}

		if load_save {
			saves_raw, err := os.ReadFile(os.Args[3])
			if err != nil {
				if strings.Split(err.Error(), ":")[1] == " no such file or directory" {
					println("Memory Card: ", os.Args[3], " not found.")
					println("Want to create Memory Card file? [Y\\N]")
					var escolha string
					fmt.Scan(&escolha)
					if escolha != "Y" {
						return
					}

					file := make([]byte, 1048576)
					err := os.WriteFile(os.Args[3], file, 0644)
					if err != nil {
						panic(err)
					}
					saves_raw, err = os.ReadFile(os.Args[3])
					if err != nil {
						panic(err)
					}
				} else {
					panic(err)
				}
			}
			cpu.LoadSave(saves_raw)
		}

		// cpu.ShowRom()
		cpu.Run()
	}
}
