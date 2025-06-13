package main

import (
	"os"

	"github.com/Faith-Kiv/Ticketing-Backend/middlewares"
	"github.com/Faith-Kiv/Ticketing-Backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// SetupRouter specifiy routes
func SetupRouter() *gin.Engine {
	//Register custom validation with gin
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("notblank", utils.NotBlank)
	}

	// gin initialization and middleware configuration
	r := gin.Default()
	r.Use(middlewares.ValidateToken())
	// Set a lower memory limit for multipart forms (default is 32 MiB)
	r.MaxMultipartMemory = 8 << 20 // 8 MiB
	r.Use(CORSMiddleware())
	r.Use(gin.Recovery())

	for path, handlers := range Routes {
		for method, handler := range handlers {
			switch method {
			case "GET":
				r.GET(path, handler)

			case "POST":
				r.POST(path, handler)

			case "PUT":
				r.PUT(path, handler)

			case "PATCH":
				r.PATCH(path, handler)

			case "DELETE":
				r.DELETE(path, handler)
			}
		}
	}
	return r
}

func main() {
	r := SetupRouter()
	r.Run("0.0.0.0:8080")

}

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetOutput(os.Stderr)
	logrus.SetReportCaller(true)
	logrus.SetLevel(logrus.DebugLevel)

}
