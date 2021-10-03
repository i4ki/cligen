package cligen

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"html/template"
	"strings"
)

type (
	Arg struct {
		Name string
		Type string
		Desc string
	}

	Flag struct {
		Name  string
		Short string
		Desc  string
	}

	Cli struct {
		Name  string
		Desc  string
		Flags []Flag
		Args  []Arg
		Cmds  []Cli
	}
)

func Parse(fname string, src string) ([]Cli, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fname, src, parser.ParseComments)
	if err != nil {
		return []Cli{}, fmt.Errorf("parsing go file \"%s\": %w", fname, err)
	}

	var (
		clis     []Cli
		parseErr error
	)

	ast.Inspect(f, func(n ast.Node) bool {
		if parseErr != nil {
			return false
		}

		switch x := n.(type) {
		case *ast.FuncDecl:
			if x.Recv != nil {
				return false
			}

			cli, err := parseFuncDecl(src, x)
			if err != nil {
				parseErr = fmt.Errorf("parsing func %s: %w", x.Name.Name, err)
				return false
			}

			clis = append(clis, cli)
		}

		return true
	})

	if parseErr != nil {
		return nil, parseErr
	}

	return clis, nil
}

func parseFuncDecl(src string, fn *ast.FuncDecl) (Cli, error) {
	cli := Cli{
		Name: fn.Name.Name,
		Desc: strings.TrimSpace(fn.Doc.Text()),
	}

	for _, argDecl := range fn.Type.Params.List {
		for _, name := range argDecl.Names {
			atype := string(src[argDecl.Type.Pos()-1 : argDecl.End()-1])

			if atype != "bool" {
				cli.Args = append(cli.Args, Arg{
					Name: name.Name,
					Type: atype,
					Desc: argDecl.Comment.Text(),
				})
			} else {
				cli.Flags = append(cli.Flags, Flag{
					Name: name.Name,
					Desc: argDecl.Comment.Text(),
				})
			}
		}
	}

	return cli, nil
}

// Help returns the help message that Generate() will possibly create.
// It's used while inspecting the file's cli possibilities.
func (cli *Cli) Help() (string, error) {
	tmpl, err := template.New("cli").Parse(`{{.Name}}: {{.Desc}}

{{.Name}} [flags] {{range .Args}}{{.Name}} {{end}}

{{with $len := len .Flags}}{{if gt $len 0}}Options:{{end}}{{end}}

{{range .Flags}}  {{if ne .Short ""}}-{{.Short}} {{.Name}}{{else}}-{{.Name}}{{end}}	{{.Desc}}
{{end}}
{{with $len := len .Cmds}}{{if gt $len 0}}Commands:{{end}}{{end}}{{range .Cmds}}  {{.Name}}{{end}}`)
	if err != nil {
		return "", fmt.Errorf("parsing help template: %w", err)
	}

	var r bytes.Buffer

	err = tmpl.Execute(&r, cli)
	if err != nil {
		return "", fmt.Errorf("executing help template: %w", err)
	}

	return r.String(), nil
}
