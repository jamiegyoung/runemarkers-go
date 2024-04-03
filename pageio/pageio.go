package pageio

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/jamiegyoung/runemarkers-go/logger"
	"github.com/jamiegyoung/runemarkers-go/templating"
)

var log = logger.New("pageio")

type Page struct {
	ShowInfoButton bool
	CardHidden     bool
}

type renderable interface {
	Data() map[string]interface{}
}

func RenderPage[T renderable](
	page_name string,
	page_string string,
	out_file *os.File,
	page_data T) {
	templ, err := templating.TemplateWithComponents(page_name, page_string)
	if err != nil {
		panic(err)
	}

	// create an interface containing both page_data and additional_data
	log("Rendering " + page_name + " to " + out_file.Name())
	err = templ.ExecuteTemplate(out_file, page_name, page_data.Data())
	if err != nil {
		panic(err)
	}
}

func ReadPageString(page_path string) (string, error) {
	log("Reading " + page_path)
	page_bytes, err := os.ReadFile(page_path)
	if err != nil {
		return "", err
	}

	return string(page_bytes), nil
}

func CreateOutFile(output_path string, page_path string) *os.File {
	if _, err := os.Stat(output_path); os.IsNotExist(err) {
		err := os.Mkdir(output_path, 0755)
		if err != nil {
			panic(err)
		}
	}

	out_file, err := os.Create(
		output_path + "/" + replaceTmplWithHtml(filepath.Base(page_path)),
	)
	if err != nil {
		panic(err)
	}
	return out_file
}

func replaceTmplWithHtml(tmp string) string {
	return strings.Replace(tmp, ".tmpl", ".html", -1)
}
