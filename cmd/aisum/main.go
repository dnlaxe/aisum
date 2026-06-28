package main

import (
	"os"

	"github.com/dnlaxe/aisum/internal/app/aisum"
)

var Version = "dev"

func main() {
	os.Exit(aisum.Run(os.Stdin, os.Stdout, os.Stderr))
}
