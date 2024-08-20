package evaluator

import (
	"fmt"

	"monkey/object"
)

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: lenFn,
	},
	"first": {
		Fn: arrayFirstFn,
	},
	"last": {
		Fn: arrayLasttFn,
	},
	"push": {
		Fn: arrayPushFn,
	},
	"rest": {
		Fn: arrayRestFn,
	},
	"put": {
		Fn: putFn,
	},
}

func putFn(args ...object.Object) object.Object {
	for _, args := range args {
		fmt.Println(args.Inspect())
	}
	return NULL
}

func arrayRestFn(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}
	if args[0].Type() != object.ARRAY_OBJ {
		return newError("argument to `rest` must be ARRAY, got %s", args[0].Type())
	}

	arr := args[0].(*object.Array)
	length := len(arr.Elements)

	if length > 0 {
		newElements := make([]object.Object, length-1)
		copy(newElements, arr.Elements[1:length])
		return &object.Array{Elements: newElements}
	}

	return NULL
}

func arrayPushFn(args ...object.Object) object.Object {
	if len(args) != 2 {
		return newError("wrong number of arguments. got=%d, want=2",
			len(args))
	}
	if args[0].Type() != object.ARRAY_OBJ {
		return newError("argument to `push` must be ARRAY, got %s", args[0].Type())
	}

	arr := args[0].(*object.Array)
	length := len(arr.Elements)

	newElements := make([]object.Object, length+1)
	copy(newElements, arr.Elements)
	newElements[length] = args[1]

	return &object.Array{Elements: newElements}
}

func arrayFirstFn(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}

	array, ok := args[0].(*object.Array)
	if !ok {
		return newError("argument to `first` must be an ARRAY, got %s instead", args[0].Type())
	}

	if len(array.Elements) > 0 {
		return array.Elements[0]
	}
	return NULL
}

func arrayLasttFn(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}

	array, ok := args[0].(*object.Array)
	if !ok {
		return newError("argument to `first` must be an ARRAY, got %s instead", args[0].Type())
	}

	if len(array.Elements) > 0 {
		return array.Elements[len(array.Elements)-1]
	}
	return NULL
}

func lenFn(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1", len(args))
	}

	switch arg := args[0].(type) {
	case *object.String:
		return &object.Integer{Value: int64(len(arg.Value))}
	case *object.Array:
		return &object.Integer{Value: int64(len(arg.Elements))}
	default:
		return newError("argument to `len` not supported, got %s", args[0].Type())
	}
}
