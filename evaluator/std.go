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
}
