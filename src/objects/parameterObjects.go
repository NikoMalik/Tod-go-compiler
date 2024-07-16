package objects

import (
	"fmt"

	"github.com/NikoMalik/Tod-go-compiler/src/print2"
)

type ParameterObject struct {
	VariableObjects

	Name string

	Ordinal  int
	Type     TypeObject
	UniqueID int
}

func (ParameterObject) ObjectType() ObjectType {
	return Parameter
}

func (p ParameterObject) ObjectName() string {
	return p.Name
}

func (p ParameterObject) Print(indent string) {
	print2.PrintC(print2.Green, indent+"└ ParameterSymbol ["+p.Name+"]")
}

func (ParameterObject) IsGlobal() bool {
	return false
}

func (p ParameterObject) IsReadOnly() bool {
	return true
}

func (p ParameterObject) VarType() TypeObject {
	return p.Type
}

func (p ParameterObject) FingerPrint() string {
	return fmt.Sprintf("%s_%d", p.Name, p.UniqueID)
}

func CreateParameterObject(name string, ordinal int, typeObject TypeObject, uniqueID int) ParameterObject {
	variableCounter++
	return ParameterObject{
		Name:     name,
		Ordinal:  ordinal,
		Type:     typeObject,
		UniqueID: variableCounter,
	}
}
