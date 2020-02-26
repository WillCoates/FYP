package framework

import (
	"errors"
	"html/template"
	"io"
	"path"

	"github.com/pelletier/go-toml"
)

var ErrNoSuchTemplate = errors.New("No such template")

type TemplateManager struct {
	templates map[string]*template.Template
}

type templatePayload struct {
	Data interface{}
	Page string
}

func NewTemplateManager(dir string) (*TemplateManager, error) {
	manager := new(TemplateManager)
	manager.templates = make(map[string]*template.Template)

	cfg, err := toml.LoadFile(path.Join(dir, "templates.toml"))
	if err != nil {
		return nil, err
	}

	for _, key := range cfg.Keys() {
		filesRaw, ok := cfg.Get(key).([]interface{})
		if ok {
			files := make([]string, len(filesRaw))

			for i, file := range filesRaw {
				files[i] = path.Join(dir, file.(string))
			}

			manager.templates[key], err = template.ParseFiles(files...)
			if err != nil {
				return nil, err
			}
		}
	}

	return manager, nil
}

func (manager *TemplateManager) Execute(name string, wr io.Writer, data interface{}) error {
	template, ok := manager.templates[name]
	if !ok {
		return ErrNoSuchTemplate
	}

	var payload templatePayload
	payload.Data = data
	payload.Page = name

	return template.Execute(wr, &payload)
}
