package views

import (
	"embed"
	"fmt"
	"html/template"
	"io"

	"github.com/TerrenceHo/monorepo/utils-go/stackerrors"
	"github.com/labstack/echo/v4"
)

type View struct {
	*template.Template // embed Template to store the templates for this view
}

func NewView(contents embed.FS, files ...string) (*View, error) {
	t, err := template.New("").ParseFS(contents, files...)
	if err != nil {
		return nil, stackerrors.Wrap(err, "failed to parse templates in NewView")
	}
	view := &View{
		Template: t,
	}
	return view, nil
}

type Renderer struct {
	RootView *View
	NewView  *View
}

func NewRenderer(RootView *View, NewView *View) *Renderer {
	return &Renderer{
		RootView: RootView,
		NewView:  NewView,
	}
}

// Render implements the Renderer interface
// See https://github.com/labstack/echo/blob/d48197db7af19becf2363496493ed0e2a8d1caea/echo.go#L134
// Docs: https://echo.labstack.com/guide/templates/
func (r *Renderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	td, ok := data.(TemplateData)
	if !ok {
		return stackerrors.New("failed to cast template data to TemplateData")
	}
	var err error
	switch td.Type {
	case ROOT_VIEW:
		err = r.RootView.ExecuteTemplate(w, name, td)
	case NEW_VIEW:
		err = r.NewView.ExecuteTemplate(w, name, td)
	default:
		err = stackerrors.New("unknown template type")
	}
	if err != nil {
		fmt.Println(err.Error())
	}
	return err
}
