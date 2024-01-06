package main

import (
	"os"

	"github.com/kkoch986/ai-skeletons-output-generation/command"
)

func main() {
	app := command.App()
	_ = app.Run(os.Args)
}
