package render

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/CloudyKit/jet/v6"
)

type Render struct {
	Renderer   string
	RootPath   string
	Secure     bool
	Port       string
	ServerName string
	JetViews   *jet.Set
}

type TemplateData struct {
	IsAuthenticated bool
	IntMap          map[string]int
	StringMap       map[string]string
	FloatMap        map[string]float32
	Data            map[string]interface{}
	CSRFToken       string
	Port            string
	Secure          bool
}

func (c *Render) Page(rw http.ResponseWriter, r *http.Request, view string, variables, data interface{}) error {
	switch strings.ToLower(c.Renderer) {
	case "go":
		return c.GoPage(rw, r, view, data)
	case "jet":
		return c.JetPage(rw, r, view, variables, data)
	}

	return nil
}

// GoPage renders a standard Go template
func (c *Render) GoPage(rw http.ResponseWriter, r *http.Request, view string, data interface{}) error {
	tmpl, err := template.ParseFiles(fmt.Sprintf("%s/views/%s.page.tmpl", c.RootPath, view))

	if err != nil {
		return err
	}

	td := &TemplateData{}

	if data != nil {
		td = data.(*TemplateData)
	}

	err = tmpl.Execute(rw, &td)

	if err != nil {
		return err
	}

	return nil
}

// JetPage renders a template using Jet templating engine
func (c *Render) JetPage(rw http.ResponseWriter, r *http.Request, templateName string, variables, data interface{}) error {
	var vars jet.VarMap

	if variables == nil {
		vars = make(jet.VarMap)
	} else {
		vars = variables.(jet.VarMap)
	}

	td := &TemplateData{}

	if data != nil {
		td = data.(*TemplateData)
	}

	t, err := c.JetViews.GetTemplate(fmt.Sprintf("%s.jet", templateName))

	if err != nil {
		log.Println(err)

		return err
	}

	if err = t.Execute(rw, vars, td); err != nil {
		log.Println(err)

		return err
	}

	return nil
}
