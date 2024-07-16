package objects

import (
	"github.com/NikoMalik/Tod-go-compiler/src/print2"
	"github.com/llir/llvm/ir"
)

type PackageObject struct {
	Objects
	Exists        bool
	Original      *PackageObject
	Name          string
	Functions     []FunctionObject
	Module        *ir.Module
	ErrorLocation print2.TextSpan
}

func (PackageObject) ObjectType() ObjectType {
	return Package
}

func (p PackageObject) ObjectName() string {
	return p.Name
}

func (p PackageObject) Print(indent string) {
	print2.PrintC(print2.Green, indent+"└ PackageSymbol ["+p.Name+"]")
}

func (p PackageObject) FingerPrint() string {
	id := "P_" + p.Name + "_"
	return id

}

func CreatePackageObject(name string, functions []FunctionObject, module *ir.Module, errorLocation print2.TextSpan) PackageObject {
	return PackageObject{
		Exists:        true,
		Name:          name,
		Functions:     functions,
		Module:        module,
		ErrorLocation: errorLocation,
	}
}
