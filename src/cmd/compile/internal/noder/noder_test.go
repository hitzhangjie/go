package noder

import (
	"cmd/compile/internal/syntax"
	"os"
	"path/filepath"
	"testing"
)

func Test_syntax_parse(t *testing.T) {
	n := noder{
		err: make(chan syntax.Error),
	}

	wd, _ := os.Getwd()
	filename := filepath.Join(wd, "testdata/main.go")
	fbase := syntax.NewFileBase(filename)

	fin, err := os.Open(filename)
	if err != nil {
		n.error(syntax.Error{Msg: err.Error()})
		return
	}
	defer fin.Close()

	// syntax.Parse do lexical analysis and syntax analysis, after this the AST is built and ready to use.
	// And, there's no any symbol table created yet.
	n.file, err = syntax.Parse(fbase, fin, n.error, n.pragma, syntax.CheckBranches)
	if err != nil {
		n.error(syntax.Error{Msg: err.Error()})
	}

	select {
	case err := <-n.err:
		t.Fatalf("parse fail: %v", err)
	default:
	}

	ast := n.file
	_ = ast
}
