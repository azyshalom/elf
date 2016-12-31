package main

import (
        "debug/elf"
        "fmt"
        "os"
)

func main() {
        if len(os.Args) < 2 {
                fmt.Println("Usage: elf elf_file")
                os.Exit(1)
        }
        f, err := os.Open(os.Args[1])
        if err != nil {
                fmt.Printf("Failed to open file, %s\n", err)
                os.Exit(1)
        }
        defer f.Close()
        _elf, err := elf.NewFile(f)
        if err != nil {
                fmt.Printf("Failed to open elf, %s\n", err)
                os.Exit(1)
        }
        // Read and decode ELF identifier
        var ident [16]uint8
        _, err = f.ReadAt(ident[0:], 0)
        if err != nil {
                fmt.Printf("Failed to read file, %s\n", err)
                os.Exit(1)
        }
        if ident[0] != '\x7f' || ident[1] != 'E' || ident[2] != 'L' || ident[3] != 'F' {
                fmt.Printf("Bad magic number at %d\n", ident[0:4])
                os.Exit(1)
        }
        var arch string
        switch _elf.Class.String() {
        case "ELFCLASS64":
                arch = "64 bits"
        case "ELFCLASS32":
                arch = "32 bits"
        }
        var mach string
        switch _elf.Machine.String() {
        case "EM_AARCH64":
                mach = "ARM64"
        case "EM_386":
                mach = "x86"
        case "EM_X86_64":
                mach = "x86_64"
        }

        fmt.Printf("\n")
        fmt.Printf("File Header: ")
        fmt.Println(_elf.FileHeader)
        fmt.Printf("ELF Class           : %s\n", arch)
        fmt.Printf("Machine             : %s\n", mach)
        fmt.Printf("ELF Type            : %s\n", _elf.Type)
        fmt.Printf("ELF Data            : %s\n", _elf.Data)
        fmt.Printf("Entry Point         : %d\n", _elf.Entry)
        fmt.Printf("Section Addresses   : %d\n", _elf.Sections)

        symbols, err := _elf.Symbols()
        fmt.Printf("\nSymbols:\n")
        for _, sym := range symbols {
                fmt.Printf("\t%s\n", sym.Name)
        }

        dynamicSymbols, err := _elf.DynamicSymbols()
        fmt.Printf("\nDynamic Symbols:\n")
        for _, sym := range dynamicSymbols {
                fmt.Printf("\t%s\n", sym.Name)
        }

        libraries, err := _elf.ImportedLibraries()
        fmt.Printf("\nImported Libraries:\n")
        for _, lib := range libraries {
                fmt.Printf("\t%s\n", lib)
        }

        importedSymbols, err := _elf.ImportedSymbols()
        fmt.Printf("\nImported Symbols:\n")
        for _, sym := range importedSymbols {
                fmt.Printf("\t%s.%s\n", sym.Library, sym.Name)
        }
}
