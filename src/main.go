package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/gorp.v1"
	_ "github.com/lib/pq"
	"database/sql"
	"seal_online_go_server/src/innotrio"
	"seal_online_go_server/src/internal"
	"seal_online_go_server/src/config"
	"github.com/robfig/cron"
	"seal_online_go_server/src/innotrio/router"
	"seal_online_go_server/src/cases"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		//c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Origin, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		//c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		//c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			fmt.Println("OPTIONS")
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}

func GetMainEngine(dbmap *gorp.DbMap) (*gin.Engine, *cron.Cron) {
	gin.SetMode(config.GIN_RELEASE_MODE)

	ginEngine := gin.New()
	//ginEngine := gin.Default()
	ginEngine.Use(CORSMiddleware())

	internalModel := &internal.Model{&innotrio.Model{dbmap}}
	internalCtrl := &internal.Ctrl{internalModel, router.NewApiRouter(ginEngine, "REQUESTS")}
	internalCtrl.Init()

	casesCalculationModel := &cases.CalculationModel{&innotrio.Model{dbmap}, make(map[string]bool)}
	casesCtrl := &cases.Ctrl{casesCalculationModel, internalModel, router.NewApiRouter(ginEngine, "CASES")}
	casesCtrl.Init()

	cron := cron.New()

	cron.AddFunc("0 * * * * *", func() {
		casesCtrl.Refresh5mCases()
	})

	return ginEngine, cron
}

func main() {
	dbinfo := fmt.Sprintf("user=%s %s host=%s port=%s dbname=%s sslmode=disable",
		config.DB_USER, config.DB_PASS, config.DB_HOST, config.DB_PORT, config.DB_NAME)
	db, _ := sql.Open("postgres", dbinfo)

	dbmap := &gorp.DbMap{Db: db, Dialect:gorp.PostgresDialect{}}
	defer dbmap.Db.Close()

	router, cron := GetMainEngine(dbmap)
	if config.CRON_ENABLED {
		cron.Start()
	}
	router.Run(":8080")
}
