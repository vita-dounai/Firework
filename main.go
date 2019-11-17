package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/vita-dounai/Firework/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf(`
	    ___ _                             _    
	    / __(_)_ __ _____      _____  _ __| | __
	   / _\ | | '__/ _ \ \ /\ / / _ \| '__| |/ /
	  / /   | | | |  __/\ V  V / (_) | |  |   < 
	  \/    |_|_|  \___| \_/\_/ \___/|_|  |_|\_\				
  `)
	fmt.Printf("\nHello %s, This is the REPL for Firework programming language.\n", user.Name)
	repl.Start(os.Stdin, os.Stdout)
}
