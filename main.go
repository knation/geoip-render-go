package main

import (
  "context"
  "errors"
  "log"
  "net"
	"net/http"
  "os"
  "os/signal"
  "syscall"
  "time"
  "github.com/gin-gonic/gin"
  "github.com/oschwald/geoip2-golang"
)

var serviceMode string = os.Getenv("MODE")
var port string = os.Getenv("PORT")

var geoDb *geoip2.Reader

func main() {

  if serviceMode == "" {
    serviceMode = "release"
  }

  if port == "" {
    port = "3000"
  }

  log.Printf("Starting `geoip` service in '%s' mode...\n", serviceMode)

  // Set the run mode of gin (release/debug)
  gin.SetMode(serviceMode)

  router := gin.New()

  // Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.Recovery())

  router.GET("/healthz", func (c *gin.Context) {
    c.String(200, "OK")
  })

	router.GET("/:ip", handler)

  srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

  // Open Maxmind database
  var geoErr error
  geoDb, geoErr = geoip2.Open(os.Getenv("GEO_FILE"))
	if geoErr != nil {
		log.Fatal(geoErr)
	}
	defer geoDb.Close()

  // Start webserver in background to allow for graceful shutdown code below
  go func() {
    log.Printf("Listening on port %v...\n", port)
    if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
      log.Panicln(err.Error());
		}
  }()

  // Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
  log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
    log.Panicf("Server forced to shutdown: %s\n", err.Error())
	}

  log.Println("Server exiting")
}

func handler(c *gin.Context) {

	ip := net.ParseIP(c.Param("ip"))
  if ip == nil {
    c.AbortWithStatus(400)
    return
  }

  // Update below this line to get the data you want

	record, err := geoDb.City(ip)
	if err != nil {
		log.Fatal(err)
    c.AbortWithStatus(500)
    return
	}

  c.JSON(200, gin.H{
		"zip": record.Postal.Code,
	})

}
