package objects

import (
	"github.com/NikoMalik/Tod-go-compiler/src/ast"
	"github.com/NikoMalik/Tod-go-compiler/src/print2"
	"github.com/llir/llvm/ir"
)

type FunctionObject struct {
	Objects
	Exists     bool
	BuiltIn    bool
	External   bool
	Public     bool
	IRFunction *ir.Func

	Name string

	Parameters  []ParameterObject
	TypeObject  TypeObject
	Declaration ast.FunctionDeclarationMember
}

func (FunctionObject) ObjectType() ObjectType {
	return Function
}

func (f FunctionObject) ObjectName() string {
	return f.Name
}

func (f FunctionObject) Print(indent string) {
	if f.BuiltIn {
		print2.PrintC(print2.Cyan, indent+"└ FunctionObject ["+f.Name+"]")
	} else {
		print2.PrintC(print2.Magenta, indent+"└ FunctionObject ["+f.Name+"]")
	}
}

func (f FunctionObject) FingerPrint() string {
	id := "F_" + f.Name + "_"

	for _, param := range f.Parameters {
		id += param.Type.FingerPrint()
	}

	id += f.TypeObject.Name

	return id

}

func CreateFunctionObject(name string, params []ParameterObject, typeObject TypeObject, declaration ast.FunctionDeclarationMember, public bool) FunctionObject {
	return FunctionObject{
		Exists:      true,
		Name:        name,
		Parameters:  params,
		TypeObject:  typeObject,
		Declaration: declaration,
		Public:      public,
	}
}

func CreateExternalFunctionObject(name string, params []ParameterObject, typeObject TypeObject, declaration ast.FunctionDeclarationMember) FunctionObject {
	return FunctionObject{
		Exists: true,
		Name:   name,

		Parameters:  params,
		TypeObject:  typeObject,
		Declaration: declaration,
		External:    true,
		Public:      true,
	}
}

func CreateBuiltInFunctionObject(name string, params []ParameterObject, typeObject TypeObject, declaration ast.FunctionDeclarationMember) FunctionObject {
	return FunctionObject{
		Exists:      true,
		BuiltIn:     true,
		Name:        name,
		Parameters:  params,
		TypeObject:  typeObject,
		Declaration: declaration,
		Public:      true,
	}
}
