package objects

import (
	"github.com/NikoMalik/Tod-go-compiler/src/print2"
)

type ObjectType string

const (
	Function       ObjectType = "FunctionSymbol"
	Class          ObjectType = "ClassSymbol"
	Struct         ObjectType = "StructSymbol"
	Enum           ObjectType = "EnumSymbol"
	GlobalVariable ObjectType = "GlobalVariableSymbol"
	LocalVariable  ObjectType = "LocalVariableSymbol"
	Parameter      ObjectType = "ParameterSymbol"
	Type           ObjectType = "TypeSymbol"
	Package        ObjectType = "PackageSymbol"
)

type Objects interface {
	ObjectType() ObjectType
	ObjectName() string
	Print(indent string)
	FingerPrint() string
}

var (
	variableCounter = 0
)

type VariableObjects interface {
	Objects
	IsReadOnly() bool
	IsGlobal() bool
	VarType() TypeObject

	Declaration() print2.TextSpan
}
