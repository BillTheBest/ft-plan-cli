package main

import "github.com/flowthings/ft-plan-cli/flowthings"

var version string

func main() {
	// This is where the CLI enters.
	// We're using Cobra for our CLI. It's very nice and makes nested commands very easy to do.
	flowthings.Execute(version)
}
