package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/AadumKhor/bitespeed-backend-task/src/app"
)

func main() {
	banner()
	bite.Run()
}

func banner() {
	title := `	
    _______  ___  _______  _______  _______  _______  _______  _______  ______  
   |  _    ||   ||       ||       ||       ||       ||       ||       ||      | 
   | |_|   ||   ||_     _||    ___||  _____||    _  ||    ___||    ___||  _    |
   |       ||   |  |   |  |   |___ | |_____ |   |_| ||   |___ |   |___ | | |   |
   |  _   | |   |  |   |  |    ___||_____  ||    ___||    ___||    ___|| |_|   |
   | |_|   ||   |  |   |  |   |___  _____| ||   |    |   |___ |   |___ |       |
   |_______||___|  |___|  |_______||_______||___|    |_______||_______||______| 
   
   `
	fmt.Println(title)
	fmt.Printf("GoVersion: %s\n", runtime.Version())
	fmt.Printf("GOOS: %s\n", runtime.GOOS)
	fmt.Printf("GOARCH: %s\n", runtime.GOARCH)
	fmt.Printf("NumCPU: %d\n", runtime.NumCPU())
	fmt.Printf("GOROOT: %s\n", runtime.GOROOT())
	fmt.Printf("Compiler: %s\n", runtime.Compiler)
	fmt.Printf("Compiler: %s\n", time.Now().Format("Monday, 2 Jan 2006"))
	fmt.Println("-----------------")
}
