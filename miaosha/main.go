package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"log"
	"net/http"
	"silentcxl/go-shop/conf"
	_ "silentcxl/go-shop/conf"
	"silentcxl/go-shop/util/connect_mysql"
	"silentcxl/go-shop/util/redis_cmd"
	"strconv"
	"time"
)

var (
	db       *gorm.DB
	redisCmd redis.Cmdable
)

func main() {
	db, _ = connect_mysql.ConnectMysql(conf.ConfigData.Ms.MysqlDsn)
	redisCmd = redis_cmd.NewRedisCmd(conf.ConfigData.Ms.Redis)
	// 模拟预热
	loadGoodsInCache()

	engine := gin.Default()

	engine.Use(FlowRestriction).POST("/goods", ms())

	engine.Run(":23568")
}

type Goods struct {
	GoodsId  uint64 `json:"goods_id"`
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
		redisGoodsKey := fmt.Sprintf("ms_goods:%d", goods.GoodsId)
		goodsNumCache := redisCmd.HGet(ctx, redisGoodsKey, "goods_num").Val()
		goodsNum, err := strconv.ParseUint(goodsNumCache, 0, 64)
		if err != nil {
			log.Println("库存取失败了：" + err.Error())
			err = db.Table("goods").Select("goods_num").Where("goods_id = ?", goods.GoodsId).Scan(&goodsNum).Error
			if err != nil {
				log.Fatal("库存取失败了：" + err.Error())
			}
			//redisCmd.HSet(ctx, redisGoodsKey, "goods_num", goodsNum)
		}
		if goodsNum == 0 {
			log.Fatal("商品没有库存了", goods.GoodsId)
		}

		//db.Table("")
	}
}

type GoodsModel struct {
	GoodsId  uint64 `json:"goods_id"`
	GoodsNum uint64 `json:"goods_num"`
}

// 产品数据预先加载进换成
func loadGoodsInCache() {
	var goodsList []*GoodsModel
	err := db.Table("goods").Select("goods_id,goods_num").Find(&goodsList).Error
	if err != nil {
		log.Fatal("产品加载失败")
	}

	for _, goods := range goodsList {
		goodsCache := []interface{}{
			"goods_num", goods.GoodsNum,
		}
		cached := false
		for i := 0; i < 3; i++ {
			val := redisCmd.HSet(context.Background(), fmt.Sprintf("ms_goods:%d", goods.GoodsId), goodsCache...).Val()
			if val > 0 {
				cached = true
				break
			}
		}
		if !cached {
			log.Fatal(fmt.Sprintf("产品写缓存失败：%d", goods.GoodsId))
		}
	}
}

// FlowRestriction 流量限制
func FlowRestriction(ctx *gin.Context) {
	// 使用redis实现一个tcp的滑动窗口限流，要求每分钟最多可以处理300个请求
	flowRestrictionKey := "flow_restriction:"
	go clearUp(ctx, flowRestrictionKey)

	if redisCmd.Incr(ctx, flowRestrictionKey+ctx.ClientIP()).Val() > 300 {
		ctx.JSON(http.StatusTooManyRequests, gin.H{"msg": "流量限制"})
		ctx.Abort()
		return
	}

	ctx.Next()
}

func clearUp(ctx *gin.Context, flowRestrictionKey string) {
	ticker := time.NewTicker(time.Second * 15)
	for range ticker.C {
		keys := redisCmd.Keys(ctx, flowRestrictionKey+"*").Val()
		for _, key := range keys {
			val := redisCmd.TTL(ctx, key).Val()
			if val.Seconds() < 0 {
				redisCmd.Del(ctx, key)
			}
		}
	}
}
