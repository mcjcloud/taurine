package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	compile "github.com/mcjcloud/taurine/pkg/llvm"
	"github.com/mcjcloud/taurine/pkg/parser"
	"github.com/mcjcloud/taurine/pkg/util"
	"github.com/spf13/cobra"
)

var compileCmd = &cobra.Command{
	Use: "c <file.tc>",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("missing source file")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		outpath, err := cmd.Flags().GetString("o")
		if err != nil {
			fmt.Printf("Could not parse flag 'o': %s\n", err.Error())
			os.Exit(1)
		}

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

		// compile
		m, err := compile.Compile(tree, ctx.ImportGraph)
		if err != nil {
			fmt.Printf("compile error: %s\n", err)
			os.Exit(1)
		}
		fmt.Println(m)

		outFileName := fmt.Sprintf("%s.ll", uuid.New())
		llOut, err := os.Create(outFileName)
		if err != nil {
			fmt.Printf("write error: %s\n", err)
			os.Exit(1)
		}

		if _, err := m.WriteTo(llOut); err != nil {
			fmt.Printf("write error: %s\n", err)
			os.Exit(1)
		}

		// shell out to local gcc/clang/mingw/etc
		if err := util.CompileIR(llOut, outpath); err != nil {
			fmt.Printf("compile error: %s\n", err)
			os.Exit(1)
		}
	},
}

func buildCompileCommand() *cobra.Command {
	compileCmd.PersistentFlags().String("o", "a.out", "output path")
	return compileCmd
}
