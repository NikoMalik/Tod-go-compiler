package objects

import "github.com/NikoMalik/Tod-go-compiler/src/print2"

type TypeObject struct {
	Objects
	Name          string
	SubTypes      []TypeObject
	IsObject      bool
	IsUserDefined bool
	Package       PackageObject
	SourceObject  Objects
}

func (TypeObject) ObjectType() ObjectType {
	return Type
}

func (t TypeObject) ObjectName() string {
	return t.Name
}

func (t TypeObject) Print(indent string) {
	print2.PrintC(print2.Green, indent+"└ TypeSymbol ["+t.FingerPrint()+"]")
}

func (t TypeObject) FingerPrint() string {
	id := "T"
	if t.IsObject {
		id += "O"
	}
	id += "_" + t.Name + "_["
	for _, subtype := range t.SubTypes {
		if subtype.Name == "array" {
			id += subtype.FingerPrint() + ";"
		} else {
			id += subtype.Name + ";"
		}
	}
	id += "]"
	return id
}

func CreateTypeObject(name string, subtypes []TypeObject, isObject bool, isUserDefined bool, pck PackageObject, src Objects) TypeObject {
	return TypeObject{
		Name:          name,
		SubTypes:      subtypes,
		IsObject:      isObject,
		IsUserDefined: isUserDefined,
		Package:       pck,
		SourceObject:  src,
	}
}
