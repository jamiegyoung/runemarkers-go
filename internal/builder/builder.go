package builder

import (
	"fmt"
	"os"

	"github.com/jamiegyoung/runemarkers-go/internal/api"
	"github.com/jamiegyoung/runemarkers-go/internal/entities"
	"github.com/jamiegyoung/runemarkers-go/internal/entitypages"
	"github.com/jamiegyoung/runemarkers-go/internal/logger"
	"github.com/jamiegyoung/runemarkers-go/internal/pages"
	"github.com/jamiegyoung/runemarkers-go/internal/templating"
	"github.com/jamiegyoung/runemarkers-go/internal/thumbnails"
)

const destination = "public"

var log = logger.New("build")

func logErr(ctx string, err error) {
	if err != nil {
		log(fmt.Sprintf("An error occured when %v!", ctx))
		log(err.Error())
	}
}

func Build(skipThumbs bool) {
	templating.ClearCache()

	if _, err := os.Stat(destination); os.IsNotExist(err) {
		err := os.Mkdir(destination, 0755)
		if err != nil {
			logErr("making destination file", err)
			return
		}
	}
	log("Reading entities")

	ents, err := entities.ReadAllEntities()
	if err != nil {
		logErr("reading entities", err)
		return
	}

	log("Generating entities")
	err = api.Generate(ents)
	if err != nil {
		logErr("generating entities", err)
		return
	}

	log("Collecting thumbnails")
	thumbnails.Collect(ents, destination, skipThumbs)

	log("Generating pages")
	err = pages.GeneratePages(destination, ents)
	if err != nil {
		logErr("generating pages", err)
		return
	}

	err = entitypages.GeneratePages(destination, ents)
	if err != nil {
		logErr("generating entity pages", err)
		return
	}
}
