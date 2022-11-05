package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mcjcloud/taurine/pkg/evaluator"
	"github.com/mcjcloud/taurine/pkg/parser"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "taurine <file.tc>",
	Short: "taurine is a simple language, fueled by caffiene",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("missing source file")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		absPath, err := filepath.Abs(args[0])
		if err != nil {
			fmt.Printf("Could not get absolute path to source file: %s\n", err.Error())
			os.Exit(1)
		}

		// create parse context
		ctx, err := parser.NewParseContext(absPath)
		if err != nil {
			fmt.Printf("Could not create parse context: %s\n", err.Error())
			os.Exit(1)
		}

		// parse using context
		tree := parser.Parse(ctx)
		ctx.PopImportWithTree(tree)

		// check for import cycles
		if cycles := ctx.ImportGraph.FindCycles(); len(cycles) > 0 {
			fmt.Println("import cycle found.")
			for _, n := range cycles {
				fmt.Println(n)
			}
			os.Exit(1)
		}

		// print any errors during parsing
		if ctx.HasErrors() {
			ctx.PrintErrors()
			os.Exit(1)
		}

		// evaluate
		err = evaluator.Evaluate(tree, ctx.ImportGraph)
		if err != nil {
			fmt.Printf("eval error: %s", err)
		}
	},
}

func Execute() {
	rootCmd.AddCommand(buildAstCommand())
	rootCmd.AddCommand(buildTokenCommand())
	rootCmd.AddCommand(buildCompileCommand())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
