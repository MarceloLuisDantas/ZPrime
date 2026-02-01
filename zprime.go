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
		fmt.Println("uso: zprime [as/run] file.[asm/krom]")
		os.Exit(0)
	}

	option := os.Args[1]
	if option != "as" && option != "run" {
		fmt.Println("uso: zprime [as/run] file.[asm/krom]")
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
		cpu, err := pc.NewCPU(strings.Split(data, "\n"))
		if err != nil {
			println(err.Error())
			return
		}

		// cpu.ShowRom()
		cpu.Run()
	}
}
