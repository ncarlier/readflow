package gotemplate

import (
	"fmt"
	"io"
	"text/template"

	"github.com/ncarlier/readflow/pkg/helper"
	templateEngine "github.com/ncarlier/readflow/pkg/template"
)

// goTemplateProvider is the structure definition of the Go template provider
type goTemplateProvider struct {
	tpl *template.Template
}

func newGoTemplateProvider(text string) (templateEngine.Provider, error) {
	tplName := fmt.Sprintf("gotemplate-%s", helper.Hash(text))
	tpl, err := template.New(tplName).Parse(text)
	if err != nil {
		return nil, err
	}
	return &goTemplateProvider{
		tpl: tpl,
	}, nil
}

// Execute template engine on article
func (t *goTemplateProvider) Execute(w io.Writer, data map[string]interface{}) error {
	return t.tpl.Execute(w, data)
}

func init() {
	templateEngine.Register("gotemplate", newGoTemplateProvider)
}
