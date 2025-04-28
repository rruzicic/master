package evaluator

import (
	"bufio"
	"bytes"
	"fmt"
	"interpreter/internal/object"
	"os"
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
	"len": {
		Fun: func(params ...object.Object) object.Object {
			switch params[0].Type() {
			case object.ARRAY_OBJ:
				arr := params[0].(*object.Array)
				return &object.Integer{Value: int64(len(arr.Elements))}
			case object.STRING_OBJ:
				str := params[0].(*object.String)
				return &object.Integer{Value: int64(len(str.Value))}
			default:
				return &object.Integer{Value: 0}
			}
		},
	},
	"panic": {
		Fun: func(params ...object.Object) object.Object {
			var buf bytes.Buffer
			for _, param := range params {
				buf.WriteString(param.Inspect())
				buf.WriteString(" ")
			}
			fmt.Println(buf.String())
			os.Exit(1)
			return nil
		},
	},
	"read": {
		Fun: func(params ...object.Object) object.Object {
			if len(params) != 1 {
				return &object.Error{Error: "read function only accepts one parameter"}
			}
			filename, ok := params[0].(*object.String)
			if !ok {
				return &object.Error{Error: "filename must be a string"}
			}
			f, err := os.ReadFile(filename.Value)
			if err != nil {
				return &object.Error{Error: "could not open file: " + err.Error()}
			}
			return &object.String{Value: string(f)}
		},
	},
	"write": {
		Fun: func(params ...object.Object) object.Object {
			if len(params) != 2 {
				return &object.Error{Error: "read function only accepts two parameters"}
			}
			filename, ok := params[0].(*object.String)
			if !ok {
				return &object.Error{Error: "filename must be a string"}
			}
			data := params[1].(*object.String)
			if !ok {
				return &object.Error{Error: "data must be string"}
			}
			f, err := os.Create(filename.Value)
			if err != nil {
				return &object.Error{Error: "could not open file: " + err.Error()}
			}
			defer f.Close()
			w := bufio.NewWriter(f)
			_, err = w.WriteString(data.Value)
			if err != nil {
				return &object.Error{Error: "could not write to file: " + err.Error()}
			}
			w.Flush()
			return nil
		},
	},
}
