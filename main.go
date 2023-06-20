package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"appota/web-builder/config"
	"appota/web-builder/content"
	"appota/web-builder/contentsvc"
	"appota/web-builder/db/mongohelper"
	"appota/web-builder/media"
	"appota/web-builder/mediasvc"
	"appota/web-builder/res"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	err := config.GetConfig()
	if err != nil {
		log.Warn(".env file not found.")
	}

	debug, _ := strconv.ParseBool(os.Getenv("DEBUG"))
	if debug {
		log.SetReportCaller(true)
	}
}

func main() {
	helper := mongohelper.NewHelper()
	err := helper.Connect(config.GetMongoURI())
	if err != nil {
		log.Panic(err)
	}

	err = os.MkdirAll(config.GetStorePath(), os.ModePerm)
	if err != nil {
		log.Panic(err)
	}

	router := httprouter.New()

	mediaRepo := media.NewRepository(helper)
	mediaService := mediasvc.New(mediaRepo)
	mediasvc.HandleHTTP(mediaService, router)

	contentRepo := content.NewRepository(helper)
	builderService := contentsvc.New(contentRepo)
	contentsvc.HandleHTTP(builderService, router)

	router.ServeFiles("/www/*filepath", http.FS(res.APIDoc))

	handler := cors.AllowAll().Handler(router)
	port := config.GetPort()

	fmt.Printf("\n🚀 WEB BUILDER STARTED\n")
	fmt.Printf("\nLocal: http://localhost:%s/ \nDocumentation: http://localhost:%s/www/apidoc/\n\n", port, port)
	http.ListenAndServe(":"+port, handler)
}
