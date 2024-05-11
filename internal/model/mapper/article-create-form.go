package mapper

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/ncarlier/readflow/internal/model"
	"github.com/ncarlier/readflow/pkg/reflection"
)

var fieldNameRe = regexp.MustCompile("^(_|url|title|text|html|tags)$")
var fieldPathRe = regexp.MustCompile(`^\w+(\[\d\]|\.\w+)*$`)

type builderFunc func(string) *model.ArticleCreateFormBuilder

func setFieldValue(fn builderFunc, field *reflection.Field) {
	if val, ok := field.String(); ok {
		fn(val)
	}
}

// ArticleCreateFormMapper is a mapper used to create "article create form" object
type ArticleCreateFormMapper struct {
	config map[string]string
}

// NewArticleCreateFormMapper create a new mapper with mapping rules
func NewArticleCreateFormMapper(config string) (*ArticleCreateFormMapper, error) {
	lines := strings.Split(config, "\n")
	conf := make(map[string]string)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}
		parts := strings.Split(line, "=")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid mapping configutration: %s", line)
		}
		fieldName := parts[0]
		if !fieldNameRe.MatchString(fieldName) {
			return nil, fmt.Errorf("invalid field name: %s", fieldName)
		}
		fieldPath := parts[1]
		if !fieldPathRe.MatchString(fieldPath) {
			return nil, fmt.Errorf("invalid field path: %s", fieldPath)
		}
		conf[fieldName] = fieldPath
	}
	return &ArticleCreateFormMapper{
		config: conf,
	}, nil
}

func (m *ArticleCreateFormMapper) build(from interface{}) *model.ArticleCreateForm {
	builder := model.NewArticleCreateFormBuilder()
	for fieldName, path := range m.config {
		field := reflection.GetField(from, path)
		if !field.Exists() {
			continue
		}
		switch fieldName {
		case "title":
			setFieldValue(builder.Title, field)
		case "url":
			setFieldValue(builder.URL, field)
		case "text":
			setFieldValue(builder.Text, field)
		case "html":
			setFieldValue(builder.HTML, field)
		case "tags":
			setFieldValue(builder.Tags, field)
		}
	}
	return builder.Build()
}

// Build an "article create form" from another object using the mapping rules
func (m *ArticleCreateFormMapper) Build(from interface{}) ([]*model.ArticleCreateForm, error) {
	result := []*model.ArticleCreateForm{}
	if rootPath, ok := m.config["_"]; ok {
		// seek root object
		root := reflection.GetField(from, rootPath)
		if !root.Exists() || !root.IsValidObject() {
			return nil, fmt.Errorf("invalid root attribute: %s", rootPath)
		}
		if arr, ok := root.Slice(); ok {
			for _, obj := range arr {
				result = append(result, m.build(obj))
			}
		} else if obj, ok := root.Map(); ok {
			result = append(result, m.build(obj))
		} else {
			return nil, fmt.Errorf("invalid root attribute: %s", rootPath)
		}
	} else {
		result = append(result, m.build(from))
	}
	return result, nil
}
