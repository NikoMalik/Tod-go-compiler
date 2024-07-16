package objects

import (
	"github.com/NikoMalik/Tod-go-compiler/src/ast"
	"github.com/NikoMalik/Tod-go-compiler/src/print2"
	"github.com/llir/llvm/ir/types"
)

type StructObject struct {
	Objects
	Exists bool

	Type TypeObject

	IRType      types.Type
	Name        string
	Declaration ast.StructDeclarationMember
	Fields      []VariableObjects
}

func (s StructObject) ObjectType() ObjectType {
	return Struct
}

func (s StructObject) ObjectName() string {
	return s.Name
}

func (s StructObject) Print(indent string) {
	print2.PrintC(print2.Green, indent+"└ StructObject ["+s.Name+"]")
}

func (s StructObject) FingerPrint() string {
	id := "S_" + s.Name + "_"
	return id
}

func CreateStructObject(name string, declaration ast.StructDeclarationMember, fields []VariableObjects) StructObject {
	sym := StructObject{
		Exists:      true,
		Name:        name,
		Declaration: declaration,
		Fields:      fields,
	}
	sym.Type = CreateTypeObject(name, make([]TypeObject, 0), false, true, PackageObject{}, sym)
	return sym
}
