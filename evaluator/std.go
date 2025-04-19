package evaluator

import (
	"bytes"
	"fmt"
	"interpreter/object"
)

var stdFunc = map[string]*object.StdFunction{
	"print": {
		Fun: func(params ...object.Object) object.Object {
			var buf bytes.Buffer
			for _, param := range params {
				buf.WriteString(param.Inspect())
				buf.WriteString(" ")
			}
			fmt.Println(buf.String())
			return nil
		},
	},
}
