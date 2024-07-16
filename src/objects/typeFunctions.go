package objects

import (
	"github.com/NikoMalik/Tod-go-compiler/src/ast"
	"github.com/NikoMalik/Tod-go-compiler/src/print2"
)

type TypeFunctionObject struct {
	Objects

	Exist bool

	BuiltIn bool

	Name string

	Parameters []ParameterObject

	Type TypeObject

	Declaration ast.FunctionDeclarationMember

	OriginType TypeObject
}

func (TypeFunctionObject) ObjectType() ObjectType {
	return Function
}

func (t TypeFunctionObject) ObjectName() string {
	return t.Name
}

func (t TypeFunctionObject) Print(indent string) {
	if t.BuiltIn {
		print2.PrintC(print2.Cyan, indent+"└ TypeFunctionObject ["+t.Name+"]")
	} else {
		print2.PrintC(print2.Magenta, indent+"└ TypeFunctionObject ["+t.Name+"]")
	}
}

func (t TypeFunctionObject) FingerPrint() string {
	id := "TF_" + t.OriginType.FingerPrint() + "_" + t.Name + "_"

	for _, param := range t.Parameters {
		id += "[" + param.Type.FingerPrint() + "]"
	}

	id += t.Type.Name

	return id

}

func CreateTypeFunctionObject(name string, parameters []ParameterObject, typeObject TypeObject, declaration ast.FunctionDeclarationMember) TypeFunctionObject {
	return TypeFunctionObject{
		Exist:       true,
		Name:        name,
		Parameters:  parameters,
		Type:        typeObject,
		Declaration: declaration,
	}
}

func CreateBuiltInTypeFunctionObject(name string, parameters []ParameterObject, typeObject TypeObject, declaration ast.FunctionDeclarationMember, origin TypeObject) TypeFunctionObject {
	return TypeFunctionObject{
		Exist:       true,
		BuiltIn:     true,
		Name:        name,
		Parameters:  parameters,
		Type:        typeObject,
		Declaration: declaration,
		OriginType:  origin,
	}
}
