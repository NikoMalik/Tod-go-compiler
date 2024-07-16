package objects

import (
	"fmt"

	"github.com/NikoMalik/Tod-go-compiler/src/print2"
)

type GlobalVariableObject struct {
	VariableObjects
	Name     string
	ReadOnly bool

	Type     TypeObject
	UniqueID int
}

func (GlobalVariableObject) ObjectType() ObjectType {
	return GlobalVariable
}

func (g GlobalVariableObject) ObjectName() string {
	return g.Name
}

func (g GlobalVariableObject) Print(indent string) {
	print2.PrintC(print2.Green, indent+"└ GlobalVariableObject ["+g.Name+"]")
}

func (g GlobalVariableObject) IsGlobal() bool {
	return true
}

func (g GlobalVariableObject) IsReadOnly() bool {
	return g.ReadOnly
}

func (g GlobalVariableObject) VarType() TypeObject {
	return g.Type
}

func (g GlobalVariableObject) FingerPrint() string {
	return fmt.Sprintf("%s_%d", g.Name, g.UniqueID)
}

func CreateGlobalVariableObject(name string, readonly bool, typeObj TypeObject) GlobalVariableObject {
	variableCounter++
	return GlobalVariableObject{
		Name:     name,
		ReadOnly: readonly,
		Type:     typeObj,
		UniqueID: variableCounter,
	}
}
