package main

import (
	"fmt"
	"os"

	"github.com/4sp1/snippi/internal/template"
)

func main() {
	if err := execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func execute() error {
	vars, file, err := parseFlags()
	if err != nil {
		return err
	}
	tmpl := template.Template{
		Stdin: os.Stdin,
		Vars:  vars,
	}
	if err := tmpl.Load(file); err != nil {
		return err
	}
	if err := tmpl.Execute(os.Stdout); err != nil {
		return err
	}
	return nil
}

func parseFlags() (
	vars map[template.VarName]template.Value,
	templateFile string,
	err error,
) {
	if len(os.Args) == 1 {
		err = errNotEnoughArgs
		return
	}
	if len(os.Args)%2 == 0 {
		err = errNonMatchingFlags
		return
	}
	vars = make(map[template.VarName]template.Value)
	for i := 1; i+1 < len(os.Args); i += 2 {
		flag := os.Args[i]
		value := os.Args[i+1]
		if len(flag) == 1 {
			err = errZeroFlagName
			return
		}
		if flag[0] != '-' {
			err = errIncorrectflagPrefix
			return
		}
		if flag == "-template" {
			templateFile = value
			continue
		}
		vars[template.VarName(flag[1:])] = template.Value(value)
	}
	if templateFile == "" {
		err = errMissingTemplateFlag
		return
	}
	return vars, templateFile, nil
}

//go:generate stringer -type=errParseFlag
type errParseFlag int

func (err errParseFlag) Error() string {
	return err.String()
}

const (
	errNotEnoughArgs errParseFlag = iota
	errNonMatchingFlags
	errZeroFlagName
	errIncorrectflagPrefix
	errMissingTemplateFlag
)

func usage() {
	fmt.Println(`
Usage:
	snippi -template some.gotmpl [-varname value]... 
	`)
}
