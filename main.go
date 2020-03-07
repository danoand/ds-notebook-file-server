package main

import (
	"fmt"
	"log"
	"os"

	"github.com/danoand/utils"
	"github.com/gin-gonic/gin"
)

// forceNoCache is a middleware handler that directs the browser to not use cached objects (but fetch from the web server)
func forceNoCache() gin.HandlerFunc {

	return func(c *gin.Context) {
		// TODO NOTE: Possible cause of long response times (?). May need to revisit at some point (Dan)

		// Set response headers
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate") // HTTP 1.1
		c.Header("Pragma", "no-cache")                                   // HTTP 1.0 (probably rare)
		c.Header("Expires", "0")                                         // proxy servers
		c.Header("Access-Control-Allow-Origin", "*")                     // Allow CORS across site

		// Serve the next handler in the middleware chain
		c.Next()
	}
}

var err error

func main() {
	// Get the file directory root location
	fileRoot := os.Getenv("DANODS_FILE_ROOT")
	if len(fileRoot) == 0 {
		// missing environment variable (use a default)
		fileRoot = "/app/prod/dsfiles"
		log.Printf("WARN: %v - missing env var: %v - using default value: %v\n",
			utils.FileLine(),
			"DANODS_FILE_ROOT",
			fileRoot)
	}
	log.Printf("INFO: %v - serving files from local directory: %v\n",
		utils.FileLine(),
		fileRoot)

	port := os.Getenv("DANODS_PORT")
	if len(port) == 0 {
		// missing environment variable (use a default)
		port = "localhost:9000"
		log.Printf("WARN: %v - missing env var: %v - using default value: %v\n",
			utils.FileLine(),
			"DANODS_PORT",
			port)
	}

	// Check to see if the specified files directory exists
	_, err = os.Stat(fileRoot)
	if err != nil {
		log.Fatalf("ERROR: %v - file directory not found. Check this out!\n", utils.FileLine())
	}

	// Define Gin routes
	rtr := gin.Default()
	rtr.Use(forceNoCache()) // NOTE: directing the browser to not use cached objects; always fetch from web server

	rtr.Static("/files", fileRoot)

	// Start the web server
	log.Printf("INFO - %v - Starting the webserver on: %v", utils.FileLine(), port)
	serveErr := rtr.Run(fmt.Sprintf("%v", port))
	if serveErr != nil {
		log.Printf("INFO: %v - Program terminating with error (if any see: %v)", utils.FileLine(), serveErr)
	}

	log.Printf("INFO: Stop processing\n")
}
