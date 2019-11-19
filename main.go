package main

import (
	"fmt"
	"os"
	"os/user"
	"strings"

	"github.com/vita-dounai/Firework/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	greeting := fmt.Sprintf("Hello %s, This is the REPL for Firework programming language.", user.Name)
	greetingWidth := len(greeting)

	logoWidth := 42
	startPosition := (greetingWidth - logoWidth) / 2
	whitePrefix := strings.Repeat(" ", startPosition)

	fmt.Println(whitePrefix + "  ___ _                             _     ")
	fmt.Println(whitePrefix + "  / __(_)_ __ _____      _____  _ __| | __")
	fmt.Println(whitePrefix + " / _\\ | | '__/ _ \\ \\ /\\ / / _ \\| '__| |/ /")
	fmt.Println(whitePrefix + "/ /   | | | |  __/\\ V  V / (_) | |  |   < ")
	fmt.Println(whitePrefix + "\\/    |_|_|  \\___| \\_/\\_/ \\___/|_|  |_|\\_\\")

	fmt.Println(greeting)
	repl.Start(os.Stdin, os.Stdout)
}
