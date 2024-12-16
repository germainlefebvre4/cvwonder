package serve

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	reloader "github.com/MPMcIntyre/go-again"
	"github.com/gin-gonic/gin"
)

//	func StartServer(app *gin.Engine, outputDirectory string) {
//		http.Handle("/", http.FileServer(http.Dir(outputDirectory)))
//		err := http.ListenAndServe(listenPort, nil)
//		utils.CheckError(err)
//	}

func StartServer(outputDirectory string) {
	fmt.Println("Serving CV as website")

	var app = gin.Default()

	port := 3000
	listenPort := fmt.Sprintf(":%d", port)

	fmt.Printf("Listening on :%d", port)

	// Initialize reloader
	// Watch template directories
	rel := StartLiveReloader(app)

	// Register LiveReload function
	app.SetFuncMap(template.FuncMap{
		"LiveReload": rel.TemplateFunc()["LiveReload"],
	})

	// Load templates
	// app.LoadHTMLGlob("generated/index.html")
	app.LoadHTMLGlob("../../generated/index.html")

	// Routes
	count := 0
	app.GET("/", func(g *gin.Context) {
		count += 1
		g.HTML(http.StatusOK, "index.html", gin.H{
			// "title": "Go Again",
			// "count": count,
		})
	})

	app.Run(listenPort)

}

func StartLiveReloader(app *gin.Engine) *reloader.Reloader {
	rel, err := reloader.New(
		// func() { app.LoadHTMLGlob("../../generated/index.html") },
		func() { app.LoadHTMLGlob("../../generated/index.html") },
		9000,
		reloader.WithLogs(true),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer rel.Close()

	// rel.Add("generated/index.html")
	rel.Add("../../generated/index.html")
	return rel
}

func ListeningUrl(port int) string {
	return fmt.Sprintf("http://localhost:%d", port)
}
