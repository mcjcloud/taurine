package util

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func CompileIR(file *os.File, outpath string) error {
	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("error compiling file '%s': %s", file.Name(), err)
	}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("clang", info.Name(), "-o", outpath)
	default:
		return fmt.Errorf("unknown OS '%s'", runtime.GOOS)
	}

	out, err := cmd.Output()
	fmt.Println(out)
	return err
}
