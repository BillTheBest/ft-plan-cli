package utils

import (
	"os"

	"github.com/flowthings/ft-plan-cli/Godeps/_workspace/src/github.com/fatih/color"
)

// Checkerror checks for an error. It exits with an error and prints it out if there was one.
// Otherwise, it does nothing.
// Error handling needs to be improved before we release to the public
func Checkerror(err error) {
	if err != nil {
		color.Red(err.Error())
		os.Exit(1)
	}
}
