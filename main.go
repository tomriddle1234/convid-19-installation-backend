package main

import (
	"github.com/gin-gonic/gin"
	limit "github.com/yangxikun/gin-limit-by-key"
	"github.com/gin-contrib/cors"
	"golang.org/x/time/rate"
	"io"
	"log"
	"os"
	"time"
)

var router *gin.Engine

func main() {
	t := time.Now()
	f, _ := os.Create("convid-19-installation-backend-" + t.Format("2006-01-02 15:04:05") + ".log")
	gin.DefaultWriter = io.MultiWriter(f)
	log.SetOutput(gin.DefaultWriter)

	GetSimpleProdEngine().Run(":8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}


func GetSimpleProdEngine() *gin.Engine {
	router = gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Content-Length", "Accept"}
	config.AllowMethods = []string{"GET", "POST", "OPTIONS"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true
	router.Use(cors.New(config))

	router.Use(gin.Logger())


	// limit request per IP
	router.Use(limit.NewRateLimiter(func(c *gin.Context) string {
		return c.ClientIP() // limit rate by client ip
	}, func(c *gin.Context) (*rate.Limiter, time.Duration) {
		return rate.NewLimiter(rate.Every(100*time.Millisecond), 200), time.Hour // limit permit bursts of at most 200 burst per token in 0.1s, 2000 qps per IP, and the limiter liveness time duration is 1 hour
	}, func(c *gin.Context) {
		c.AbortWithStatus(429) // handle exceed rate limit request
	}))

	initializeRoutes()
	initData()
	go listenToChannel()

	return router
}

// now this is just for tests
func GetMainEngine() *gin.Engine {
	router = gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://localhost", "http://localhost", "http://localhost:8000", "http://localhost:8100", "https://127.0.0.1", "http://127.0.0.1"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Content-Length", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://localhost"
		},
		MaxAge: 12 * time.Hour,
	}))

	router.Use(gin.Logger())

	initializeRoutes()
	initData()
	go listenToChannel()

	return router
}