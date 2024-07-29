package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	dsn = "root:pass@tcp(127.0.0.1:3306)/go_shop?charset=utf8mb4&parseTime=True&loc=Local"
)

func main() {
	engine := gin.Default()

	engine.POST("/goods", ms())

	engine.Run(":9999")
}

type Goods struct {
	GoodsId  string `json:"goods_id"`
	GoodsNum uint64 `json:"goods_num"`
	UserId   uint64 `json:"user_id"`
}

func ms() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var goods Goods
		if err := ctx.ShouldBindJSON(&goods); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误: " + err.Error()})
			ctx.Abort()
			return
		}
		//db, _ := connect_mysql.ConnectMysql(dsn)

		//db.Table("")
	}
}
