package objects

import (
	"fmt"

	"github.com/NikoMalik/Tod-go-compiler/src/print2"
)

type LocalVariableObject struct {
	VariableObjects

	Name     string
	ReadOnly bool

	Type TypeObject

	UniqueID int
}

func (LocalVariableObject) ObjectType() ObjectType {
	return LocalVariable
}

func (l LocalVariableObject) ObjectName() string {
	return l.Name
}

func (l LocalVariableObject) Print(indent string) {
	print2.PrintC(print2.Yellow, indent+"└ LocalVariableObject ["+l.Name+"]")
}

func (LocalVariableObject) IsGlobal() bool {
	return false
}

func (s LocalVariableObject) IsReadOnly() bool {
	return s.ReadOnly
}

func (l LocalVariableObject) VarType() TypeObject {
	return l.Type
}

func (l LocalVariableObject) FingerPrint() string {
	return fmt.Sprintf("%s_%d", l.Name, l.UniqueID)
}

func CreateLocalVariableObject(name string, readonly bool, typeObj TypeObject) LocalVariableObject {
	variableCounter++
	return LocalVariableObject{
		Name:     name,
		ReadOnly: readonly,
		Type:     typeObj,
		UniqueID: variableCounter,
	}
}
