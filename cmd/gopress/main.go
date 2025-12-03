package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func main() {
	cobra.MousetrapHelpText = ""
	Execute()
	fmt.Println("\nNaciśnij Enter, aby zamknąć okno...")
	fmt.Scanln()
}
