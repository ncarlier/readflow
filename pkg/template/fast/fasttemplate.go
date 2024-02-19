package gotemplate

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/valyala/fasttemplate"

	"github.com/ncarlier/readflow/pkg/html"
	templateEngine "github.com/ncarlier/readflow/pkg/template"
)

// fastTemplateProvider is the structure definition of the Fast Template provider
type fastTemplateProvider struct {
	tpl *fasttemplate.Template
}

func newFastTemplateProvider(text string) (templateEngine.Provider, error) {
	return &fastTemplateProvider{
		tpl: fasttemplate.New(text, "{{", "}}"),
	}, nil
}

// Execute template engine on article
func (t *fastTemplateProvider) Execute(w io.Writer, data map[string]interface{}) error {
	content := t.tpl.ExecuteFuncString(func(w io.Writer, tag string) (int, error) {
		value, err := evalTemplateStatement(data, tag)
		if err != nil {
			return 0, err
		}
		return w.Write([]byte(value))
	})

	_, err := w.Write([]byte(content))
	return err
}

func evalTemplateStatement(data map[string]interface{}, statement string) (string, error) {
	keywords := strings.Split(statement, "|")
	attribute := strings.TrimSpace(keywords[0])
	v := data[attribute]
	if v == nil {
		return "", fmt.Errorf("unknown attribute: %s", attribute)
	}
	filters := keywords[1:]

	switch value := v.(type) {
	case string:
		return evalTemplateFilters(value, filters)
	default:
		return "", fmt.Errorf("unexpected value type: %#v", value)
	}
}

func evalTemplateFilters(value string, filters []string) (string, error) {
	var err error
	for _, filter := range filters {
		filter := strings.TrimSpace(filter)
		value, err = evalTemplateFilter(value, filter)
		if err != nil {
			return value, err
		}
	}
	return value, err
}

func evalTemplateFilter(value, filter string) (string, error) {
	switch filter {
	case "urlquery":
		return url.QueryEscape(value), nil
	case "base64":
		return base64.StdEncoding.EncodeToString([]byte(value)), nil
	case "html2text":
		return html.HTML2Text(value)
	case "json":
		buf, err := json.Marshal(value)
		if err != nil {
			return "", err
		}
		return string(buf), nil
	default:
		return "", fmt.Errorf("undefined filter function: %s", filter)
	}
}

func init() {
	templateEngine.Register("fast", newFastTemplateProvider)
}
