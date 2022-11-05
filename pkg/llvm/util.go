package llvm

import (
	"fmt"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/value"
)

func wrapInCharArray(val value.Value) (*constant.CharArray, error) {
	var nullTerminated []byte
	switch t := val.(type) {
	case *constant.CharArray:
		return t, nil
	case *constant.Float:
		if t.NaN {
			nullTerminated = append([]byte("NaN"), 0x00)
		} else {
			nullTerminated = append([]byte(t.X.String()), 0x00)
		}
	case *constant.Int:
		nullTerminated = append([]byte(t.X.String()), 0x00)
	default:
		return nil, fmt.Errorf("unknown type %s", t)
	}

	return constant.NewCharArrayFromString(string(nullTerminated)), nil
}
