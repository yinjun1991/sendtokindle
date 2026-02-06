package web

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"path/filepath"
	"strconv"
	"time"
)

//go:embed templates/*
var templatesFS embed.FS

type Renderer struct {
	indexTmpl *template.Template
	adminTmpl *template.Template
}

func NewRenderer() (*Renderer, error) {
	funcs := template.FuncMap{
		"formatTime": func(t time.Time) string { return t.Format("2006-01-02 15:04") },
		"formatBytes": func(n int64) string {
			const unit = 1024
			if n < unit {
				return strconv.FormatInt(n, 10) + " B"
			}
			div, exp := int64(unit), 0
			for v := n / unit; v >= unit; v /= unit {
				div *= unit
				exp++
			}
			value := float64(n) / float64(div)
			suffix := "KMGTPE"[exp : exp+1]
			return fmt.Sprintf("%.1f %siB", value, suffix)
		},
		"base": filepath.Base,
	}

	indexTmpl, err := template.New("index.html").Funcs(funcs).ParseFS(templatesFS, "templates/index.html")
	if err != nil {
		return nil, fmt.Errorf("parse index template: %w", err)
	}

	adminTmpl, err := template.New("admin.html").Funcs(funcs).ParseFS(templatesFS, "templates/admin.html")
	if err != nil {
		return nil, fmt.Errorf("parse admin template: %w", err)
	}

	return &Renderer{
		indexTmpl: indexTmpl,
		adminTmpl: adminTmpl,
	}, nil
}

func (r *Renderer) RenderIndex(w io.Writer, data any) error {
	return r.indexTmpl.Execute(w, data)
}

func (r *Renderer) RenderAdmin(w io.Writer, data any) error {
	return r.adminTmpl.Execute(w, data)
}

