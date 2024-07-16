package semantic

import "github.com/NikoMalik/Tod-go-compiler/src/objects"

type Scope struct {
	Parent  *Scope
	Objects map[string]objects.Objects
}

func (s *Scope) TryDeclareObject(sym objects.Objects) bool {
	lookup := s.TryLookupObject(sym.ObjectName())

	if lookup != nil {
		return false // symbol already exists
	} else {
		s.Objects[sym.ObjectName()] = sym
		return true
	}
}

func (s *Scope) TryLookupObject(name string) objects.Objects {
	sym, found := s.Objects[name]

	if found {
		return sym
	}

	if s.Parent != nil {
		return s.Parent.TryLookupObject(name)
	}

	return nil
}

func (s *Scope) InsertFunctionObject(obj []objects.FunctionObject) {
	for _, sym := range obj {
		s.TryDeclareObject(sym)
	}
}

func (s *Scope) InsertVariableObject(obj []objects.VariableObjects) {
	for _, sym := range obj {
		s.TryDeclareObject(sym)
	}
}

func (s *Scope) GetAllFunctions() []objects.FunctionObject {
	functions := make([]objects.FunctionObject, 0)

	for _, sym := range s.Objects {
		if sym.ObjectType() == objects.Function {
			functions = append(functions, sym.(objects.FunctionObject))
		}
	}

	moreFunctions := make([]objects.FunctionObject, 0)
	if s.Parent != nil {
		moreFunctions = s.Parent.GetAllFunctions()
	}

	functions = append(functions, moreFunctions...)

	return functions

}

func (s *Scope) GetAllVariables() []objects.VariableObjects {
	variables := make([]objects.VariableObjects, 0)

	for _, sym := range s.Objects {
		if sym.ObjectType() == objects.LocalVariable ||
			sym.ObjectType() == objects.GlobalVariable ||

			sym.ObjectType() == objects.Parameter {

			{
				variables = append(variables, sym.(objects.VariableObjects))
			}
		}

		moreVariables := make([]objects.VariableObjects, 0)
		if s.Parent != nil {
			moreVariables = s.Parent.GetAllVariables()
		}
		variables = append(variables, moreVariables...)
	}

	return variables

}

func (s *Scope) GetAllStructs() []objects.StructObject {
	structs := make([]objects.StructObject, 0)

	for _, sym := range s.Objects {
		if sym.ObjectType() == objects.Struct {
			structs = append(structs, sym.(objects.StructObject))
		}
	}

	moreStructs := make([]objects.StructObject, 0)

	if s.Parent != nil {
		moreStructs = s.Parent.GetAllStructs()
	}

	structs = append(structs, moreStructs...)

	return structs
}

func (s *Scope) GetAllPackage() []objects.PackageObject {
	packages := make([]objects.PackageObject, 0)

	for _, sym := range s.Objects {
		if sym.ObjectType() == objects.Package {
			packages = append(packages, sym.(objects.PackageObject))
		}
	}

	morePackages := make([]objects.PackageObject, 0)

	if s.Parent != nil {
		morePackages = s.Parent.GetAllPackage()
	}

	packages = append(packages, morePackages...)

	return packages
}

func CreateScope(parent *Scope) Scope {
	return Scope{
		Parent:  parent,
		Objects: make(map[string]objects.Objects),
	}
}
