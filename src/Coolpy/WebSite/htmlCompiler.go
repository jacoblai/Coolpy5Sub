package WebSite

import (
	"io/ioutil"
	"path/filepath"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/html"
	"html/template"
)

func compileTemplates(filenames ...string) (*template.Template, error) {
	m := minify.New()
	m.AddFunc("text/html", html.Minify)
	//m.AddFunc("text/css", css.Minify)
	//m.AddFunc("text/javascript", js.Minify)

	var tmpl *template.Template
	for _, filename := range filenames {
		name := filepath.Base(filename)
		if tmpl == nil {
			tmpl = template.New(name)
		} else {
			tmpl = tmpl.New(name)
		}

		b, err := ioutil.ReadFile(filename)
		if err != nil {
			return nil, err
		}

		mb, err := m.Bytes("text/html", b)
		if err != nil {
			return nil, err
		}
		tmpl.Parse(string(mb))
	}
	return tmpl, nil
}
