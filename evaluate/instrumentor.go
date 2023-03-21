// detectors
package evaluate

import (
	"bytes"
	"github.com/staheri/goatlib/instrument"
	"go/ast"
	"go/printer"
	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/loader"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func builtinDL_inst(src, dest string) []*instrument.ConcurrencyUsage {
	var conf loader.Config
	var astfiles []*ast.File

	// obtain MAXPROCS
	MAXPROCS := os.Getenv("GOATMAXPROCS")
	if MAXPROCS == "" {
		panic("GOATMAXPROCS is not set!")
	} else {
		mp, err := strconv.Atoi(MAXPROCS)
		check(err)
		if mp < 1 || mp > 64 {
			panic("GOATMAXPROCS out of range!")
		}
	}

	// copy all files and inject to its AST
	paths, err := filepath.Glob(src + "/*.go")
	check(err)
	if _, err := conf.FromArgs(paths, false); err != nil {
		panic(err)
	}
	prog, err := conf.Load()
	check(err)
	for _, crt := range prog.Created {
		for _, ast := range crt.Files {
			astfiles = append(astfiles, ast)
		}
	}
	for _, astF := range astfiles {
		var entryFunc *ast.FuncDecl
		astutil.Apply(astF, func(cur *astutil.Cursor) bool {
			if node, ok := cur.Node().(*ast.FuncDecl); ok {
				if strings.HasPrefix(node.Name.Name, "Test") {
					entryFunc = node
					return false
				}
			}
			return true
		}, nil)

		entryFunc.Body.List = append([]ast.Stmt{&ast.ExprStmt{
			X: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "runtime"},
					Sel: &ast.Ident{Name: "GOMAXPROCS"},
				},
				Args: []ast.Expr{
					&ast.BasicLit{Value: MAXPROCS},
				},
			}}}, entryFunc.Body.List...)

		var buf bytes.Buffer
		astutil.AddImport(prog.Fset, astF, "runtime")
		printer.Fprint(&buf, prog.Fset, astF)
		filename := filepath.Join(dest, strings.Split(filepath.Base(prog.Fset.Position(astF.Pos()).Filename), ".")[0]+".go")
		err = ioutil.WriteFile(filename, buf.Bytes(), 0644)
		check(err)
	}
	return nil
}

func lockDL_inst(src, dest string) []*instrument.ConcurrencyUsage {

	// obtain MAXPROCS
	MAXPROCS := os.Getenv("GOATMAXPROCS")
	if MAXPROCS == "" {
		panic("GOATMAXPROCS is not set!")
	} else {
		mp, err := strconv.Atoi(MAXPROCS)
		check(err)
		if mp < 1 || mp > 64 {
			panic("GOATMAXPROCS out of range!")
		}
	}

	// copy all files to dest && rewrite (sed)
	// copy from source to dest
	files, err := filepath.Glob(src + "/*.go")
	check(err)
	for _, file := range files {
		cmd := exec.Command("cp", file, dest)
		if err1 := cmd.Run(); err1 != nil {
			panic("lockDL_inst cp failed")
		}
	}

	var conf loader.Config
	var astfiles []*ast.File
	// copy all files and inject to its AST
	paths, err := filepath.Glob(src + "/*.go")
	check(err)
	if _, err := conf.FromArgs(paths, false); err != nil {
		panic(err)
	}
	prog, err := conf.Load()
	check(err)
	for _, crt := range prog.Created {
		for _, ast := range crt.Files {
			astfiles = append(astfiles, ast)
		}
	}
	for _, astF := range astfiles {
		var entryFunc *ast.FuncDecl
		astutil.Apply(astF, func(cur *astutil.Cursor) bool {
			if node, ok := cur.Node().(*ast.FuncDecl); ok {
				if strings.HasPrefix(node.Name.Name, "Test") {
					entryFunc = node
					return false
				}
			}
			return true
		}, nil)

		entryFunc.Body.List = append([]ast.Stmt{&ast.ExprStmt{
			X: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "runtime"},
					Sel: &ast.Ident{Name: "GOMAXPROCS"},
				},
				Args: []ast.Expr{
					&ast.BasicLit{Value: MAXPROCS},
				},
			}}}, entryFunc.Body.List...)

		var buf bytes.Buffer
		astutil.AddImport(prog.Fset, astF, "runtime")
		printer.Fprint(&buf, prog.Fset, astF)
		filename := filepath.Join(dest, strings.Split(filepath.Base(prog.Fset.Position(astF.Pos()).Filename), ".")[0]+".go")
		err = ioutil.WriteFile(filename, buf.Bytes(), 0644)
		check(err)
	}

	files, err = filepath.Glob(dest + "/*.go")
	check(err)
	for _, file := range files {
		commands := [][]string{
			[]string{"sed", "-i ", "s/sync.RWMutex/deadlock.RWMutex/", file},
			[]string{"sed", "-i ", "s/sync.Mutex/deadlock.Mutex/", file},
		}
		for _, args := range commands {
			if out, err := exec.Command(args[0], args[1:]...).CombinedOutput(); err != nil {
				log.Println(string(out))
				panic(err)
			}
		}

		abspath, _ := filepath.Abs(file)
		cmd := exec.Command("goimports", "-w", abspath)
		if out, err := cmd.CombinedOutput(); err != nil {
			log.Println(string(out))
			panic(err)
		}
	}
	return nil
}

func goleak_inst(src, dest string) []*instrument.ConcurrencyUsage {
	var conf loader.Config
	var astfiles []*ast.File

	// obtain MAXPROCS
	MAXPROCS := os.Getenv("GOATMAXPROCS")
	if MAXPROCS == "" {
		panic("GOATMAXPROCS is not set!")
	} else {
		mp, err := strconv.Atoi(MAXPROCS)
		check(err)
		if mp < 1 || mp > 64 {
			panic("GOATMAXPROCS out of range!")
		}
	}

	// copy all files and inject to its AST
	paths, err := filepath.Glob(src + "/*.go")
	check(err)
	if _, err := conf.FromArgs(paths, false); err != nil {
		panic(err)
	}
	prog, err := conf.Load()
	check(err)
	for _, crt := range prog.Created {
		for _, ast := range crt.Files {
			astfiles = append(astfiles, ast)
		}
	}
	for _, astF := range astfiles {
		var entryFunc *ast.FuncDecl
		astutil.Apply(astF, func(cur *astutil.Cursor) bool {
			if node, ok := cur.Node().(*ast.FuncDecl); ok {
				if strings.HasPrefix(node.Name.Name, "Test") {
					entryFunc = node
					return false
				}
			}
			return true
		}, nil)

		entryFunc.Body.List = append([]ast.Stmt{&ast.DeferStmt{
			Call: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "goleak"},
					Sel: &ast.Ident{Name: "VerifyNone"},
				},
				Args: []ast.Expr{
					&ast.BasicLit{Value: "t"},
				},
			}}}, entryFunc.Body.List...)
		entryFunc.Body.List = append([]ast.Stmt{&ast.ExprStmt{
			X: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "runtime"},
					Sel: &ast.Ident{Name: "GOMAXPROCS"},
				},
				Args: []ast.Expr{
					&ast.BasicLit{Value: MAXPROCS},
				},
			}}}, entryFunc.Body.List...)

		var buf bytes.Buffer
		astutil.AddImport(prog.Fset, astF, "go.uber.org/goleak")
		astutil.AddImport(prog.Fset, astF, "runtime")
		printer.Fprint(&buf, prog.Fset, astF)
		filename := filepath.Join(dest, strings.Split(filepath.Base(prog.Fset.Position(astF.Pos()).Filename), ".")[0]+".go")
		err = ioutil.WriteFile(filename, buf.Bytes(), 0644)
		check(err)
	}
	return nil
}

func goat_trace_inst(src, dest string) []*instrument.ConcurrencyUsage {
	return instrument.InstrumentTraceOnly(src, dest)
}

func goat_critic_inst(src, dest string) []*instrument.ConcurrencyUsage {
	return instrument.InstrumentCriticOnly(src, dest) // critical
}

func goat_delay_inst(src, dest string) []*instrument.ConcurrencyUsage {
	return instrument.InstrumentCriticalPoints(src, dest) // critical
}
