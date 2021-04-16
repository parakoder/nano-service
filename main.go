package main

import (
	"fmt"
	"nano-service/config"
	"nano-service/controllers"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main(){
	if err := godotenv.Load(".env.dev"); err != nil {
		panic(".env not exists")
	} else {

		conn, errConn := config.ConnectSQL()
		if errConn != nil {
			fmt.Println(errConn)
			os.Exit(-1)
		}

		r := gin.Default()
		r.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"messages": "Wellcome to Documents management service",
			})
		})

		pController := controllers.NewNanoHandler(conn)
		p := r.Group("v1/api/nano")
		{
			p.GET("/pelayanan", pController.GetPelayanan)
			p.GET("/getAntrian", pController.GetAntrian)
			p.GET("/getPDF", pController.DownloadPdf)
			p.POST("/createAntrian", pController.CreateAntrian)
			p.GET("/cekAntrian", pController.CekAntrian)
		}
		r.Use(cors.New(cors.Config{
			AllowOrigins:     []string{"*"},
			AllowMethods:     []string{"POST, GET, OPTIONS, PUT, DELETE"},
			AllowHeaders:     []string{"*"},
			ExposeHeaders:    []string{"*"},
			AllowCredentials: true,
			AllowOriginFunc: func(origin string) bool {
				return origin == "*"
			},
			MaxAge: 12 * time.Hour,
		}))
		r.Run(":" + os.Getenv("PORT"))
	}
}
