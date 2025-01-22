package docker

import (
	"fmt"
	"os"
)

var INFO = "\033[32m✔\033[0m"
var ERROR = "\033[31m✘\033[0m"

func Print(types string, format string, a ...any) {
	fmt.Printf("[%s] %s", types, fmt.Sprintf(format, a...))
}

func Error(types string, format string, a ...any) {
	fmt.Printf("[%s] %s", types, fmt.Sprintf(format, a...))
	os.Exit(1)
}
