package irtools

import (
	"fmt"
	"regexp"

	"os"
	"strings"

	"github.com/NikoMalik/Tod-go-compiler/src/print2"
	"github.com/llir/llvm/asm"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
)

func ReadModule(path string) *ir.Module {
	// read out contents of the the file
	moduleBytes, _ := os.ReadFile(path)
	module := string(moduleBytes)

	// do a regex replace to replace all function content with 'ret void'
	module = regexp.MustCompile(`{\n.*?\n}`).ReplaceAllString(module, " {\n  ret void\n}\n")

	// do a regex replace to remove invalid function declarations
	module = regexp.MustCompile(`(?m)(declare|define).*?align [0-9]*`).ReplaceAllString(module, "")

	// do a regex replace to remove invalid function declarations

	os.WriteFile("./mod.ll", []byte(module), os.ModePerm)

	// do a regex replace to remove new fangled sret

	// parse the module using llir/llvm
	irModule, err := asm.ParseString(path, module)
	if err != nil {
		print2.PrintC(print2.Red, "Couldnt load module '"+path+"'")
		fmt.Println(module)
		panic(err)
	}

	return irModule
}

func FindFunctionsWithPrefix(module *ir.Module, prefix string) []*ir.Func {
	funcs := make([]*ir.Func, 0)

	for _, fnc := range module.Funcs {
		if strings.HasPrefix(fnc.Name(), prefix) {
			funcs = append(funcs, fnc)
		}
	}

	return funcs
}

func FindFunction(module *ir.Module, name string) *ir.Func {
	for _, fnc := range module.Funcs {
		if fnc.Name() == name {
			return fnc
		}
	}

	print2.PrintC(print2.Red, "Couldnt find function '"+name+"'")
	return nil
}

func TryFindFunction(module *ir.Module, name string) *ir.Func {
	for _, fnc := range module.Funcs {
		if fnc.Name() == name {
			return fnc
		}
	}

	return nil
}

func FindGlobal(module *ir.Module, name string) *ir.Global {
	for _, glb := range module.Globals {
		if glb.Name() == name {
			return glb
		}
	}

	print2.PrintC(print2.Red, "Couldnt find global '"+name+"'")
	return nil
}

func FindGlobalSuffix(module *ir.Module, name string) *ir.Global {
	for _, glb := range module.Globals {
		if strings.HasSuffix(glb.Name(), name) {
			return glb
		}
	}

	print2.PrintC(print2.Red, "Couldnt find global '"+name+"'")
	return nil
}

func TryFindGlobal(module *ir.Module, name string) *ir.Global {
	for _, glb := range module.Globals {
		if glb.Name() == name {
			return glb
		}
	}

	return nil
}

func ReadConstStringArray(module *ir.Module, glb *ir.Global) ([]string, bool) {
	arr, ok := glb.Init.(*constant.Array)
	if !ok {
		return make([]string, 0), false
	}

	names := make([]string, 0)

	for _, elem := range arr.Elems {
		gep, ok := elem.(*constant.ExprGetElementPtr)
		if !ok {
			return make([]string, 0), false
		}

		cFld := TryFindGlobal(module, strings.TrimPrefix(gep.Src.Ident(), "@"))
		if cFld == nil {
			return make([]string, 0), false
		}

		chr, ok := cFld.Init.(*constant.CharArray)
		if !ok {
			return make([]string, 0), false
		}

		names = append(names, string(chr.X))
	}

	return names, true
}
