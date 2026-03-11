package template

import (
	"bytes"
	"errors"
	"io"
	"os"
	"text/template"
)

type (
	VarName string
	Value   string
)

type Template struct {
	Vars  map[VarName]Value
	Stdin io.Reader

	t *template.Template
}

func (t *Template) Load(path string) (err error) {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() {
		err = errors.Join(err, f.Close())
	}()
	var b bytes.Buffer
	_, err = io.Copy(&b, f)
	if err != nil {
		return err
	}
	tmpl, err := template.New("").Parse(b.String())
	if err != nil {
		return err
	}
	t.t = tmpl
	return nil
}

func (t *Template) Execute(out io.Writer) (err error) {
	var b bytes.Buffer
	_, err = io.Copy(&b, t.Stdin)
	vars := make(map[string]string)
	for name, value := range t.Vars {
		vars[string(name)] = string(value)
	}
	return t.t.Option("missingkey=zero").Execute(out, struct {
		Vars  map[string]string
		Stdin string
	}{
		Vars:  vars,
		Stdin: b.String(),
	})
}
