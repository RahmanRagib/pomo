package main

import (
	"os"

	"github.com/RahmanRagib/pomo/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
