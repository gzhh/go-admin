package controllers

import (
	"go-admin/app/admin/handlers"
	"go-admin/internal/lib/logger"

	"github.com/gin-gonic/gin"
)

func HealthCheck(ctx *gin.Context) {
	/*
		// test db
		db, err := mysql.DBInstance("go-admin", true)
		if err != nil {
			panic(err)
		}
		var rows []map[string]interface{}
		err = db.Table("test").Exec("select * from test").Find(&rows).Error
		if err != nil {
			panic(err)
		}
		fmt.Printf("----rows %+v\n", rows)
	*/
	logger.Infof(ctx.Request.Context(), "health check")

	handlers.Success(ctx, nil)
}
