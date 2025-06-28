package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/go-resty/resty/v2"
	"gorm.io/gorm"
	"log"
	_ "silentcxl/go-shop/conf"
	"silentcxl/go-shop/util/connect_mysql"
	"time"
)

var (
	db       *gorm.DB
	redisCmd redis.Cmdable
)

type MsgParam struct {
	Cursor   int64 `form:"cursor"`
	PageSize int   `form:"page_size"`
}
type MsgBody struct {
	ID         int64  `json:"id"`
	MainUserId int64  `json:"main_user_id"`
	GroupName  string `json:"group_name"`
	Message    string `json:"message"`
}

type CallbackItem struct {
	Callback  string `json:"callback"`
	ClickTime int64  `json:"click_time"`
}

type Req struct {
	EventType string     `json:"event_type"`
	Context   ContextReq `json:"context"`
	Timestamp int64      `json:"timestamp"`
}

type ContextReq struct {
	Ad AdReq `json:"ad"`
}

type AdReq struct {
	Callback string `json:"callback"`
}

func main() {
	var err error
	db, err = connect_mysql.ConnectMysql("cl_readOnly:Chuangliang@2023@mysql+tcp(mysql23e742b714ab.rds.ivolces.com:3306)/chuangliang_cid_cpa?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal("数据库连接失败", err)
		return
	}

	//engine := gin.Default()
	//
	//engine.Use(func(c *gin.Context) {
	//	method := c.Request.Method
	//	c.Header("Access-Control-Allow-Origin", "*")
	//	c.Header("Access-Control-Allow-Headers", "Access-Control-Allow-Headers,Authorization,User-Agent, Keep-Alive, Content-Type, X-Requested-With,X-CSRF-Token,AccessToken,Token")
	//	c.Header("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT, PATCH, OPTIONS")
	//	c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
	//	c.Header("Access-Control-Allow-Credentials", "true")
	//
	//	// 放行所有OPTIONS方法
	//	if method == "OPTIONS" {
	//		c.AbortWithStatus(http.StatusAccepted)
	//	}
	//	c.Next()
	//})
	//
	//engine.GET("/api/messages", func(context *gin.Context) {
	//	//var params MsgParam
	//	//if err = context.ShouldBindQuery(&params); err != nil {
	//	//	context.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误: " + err.Error(), "data": nil})
	//	//	context.Abort()
	//	//	return
	//	//}
	//	//var messages []MsgBody
	//	//pageSie := params.PageSize
	//	//if pageSie > 100 {
	//	//	pageSie = 100
	//	//}
	//	//err = db.Table("notify_wechat_msg").Where("id > ?", params.Cursor).Limit(pageSie).Find(&messages).Error
	//	//if err != nil {
	//	//	context.JSON(http.StatusBadRequest, gin.H{"msg": "查询失败: " + err.Error(), "data": nil})
	//	//	context.Abort()
	//	//	return
	//	//}
	//
	//	//context.JSON(http.StatusOK, gin.H{"data": 1, "msg": "OK"})
	//	log.Println("跳转")
	//
	//	context.Redirect(http.StatusFound, "https://www.baidu.com")
	//})
	//engine.POST("/api/stop", func(context *gin.Context) {
	//
	//})
	//
	//engine.Run(":8989")

	adIds := map[int64]int{
		//7520716887003217962: 3,
		//7520696151568023591: 5,
		//7520696101071814697: 4,
		//7520696269380190262: 6,
		//7520696130302165046: 3,
		//7520696194727198756: 3,
		//7520696240897065003: 5,
		//7520696265328508991: 3,
		//7520696231321108523: 5,
		//7520696227860693035: 5,
		//7520696266042621993: 6,
		//7520696099888939071: 5,
		//7520696245040873511: 4,
		//7520696087918297142: 5,
		//7520696237206880310: 5,
		//7520696151412260918: 5,
		//7520696049801396287: 3,
		//7520696273663311926: 4,
		//7520696103932461110: 5,
		//7520696228716855335: 6,
		//7520696254327078931: 6,
		//7520696266477535286: 6,
		//7520696107715411987: 2,
		//7520696218235224107: 5,
		//7520696047079145526: 5,
		//7520696083690569771: 3,
		//7520696155543978020: 6,
		//7520696192222101547: 6,
		//7520696170790191140: 3,
		//7520696245527543827: 6,
		//7520696067740729363: 2,
		//7520696197036867620: 6,
		//7520696223822299155: 5,
		//7520696060861218852: 6,
		//7520696196177821739: 5,
		//7520696172600705060: 6,
		//7520696263374880831: 6,
		//7520696245926101035: 5,
		//7520696117005320211: 5,
		//7520696263633911844: 4,
		//7520696102283165759: 6,
	}
	//adIdsList := make([]int64, 0)
	//for adId, _ := range adIds {
	//	adIdsList = append(adIdsList, adId)
	//}
	d := time.Now()
	date := time.Date(d.Year(), d.Month(), d.Day(), 0, 10, 0, 0, time.Local)
	//db.Table("monitor_links").
	//	Where("ad_id in (?)", adIdsList).
	//	Where("click_time >= ?", date.Unix()).
	//	Limit(num).
	//	Find(&clicks)

	for adId, num := range adIds {
		var clicks []*CallbackItem
		db.Debug().Table("cpa_ad_click").
			Where("ad_id = ?", adId).
			Where("click_time >= ?", date.Unix()).
			Limit(num).
			Find(&clicks)

		if len(clicks) < num {
			fmt.Println(adId, "点击不足: ", num, "，只查到: ", len(clicks))
		}

		for _, click := range clicks {
			// in_app_uv
			reqData := &Req{
				EventType: "in_app_uv",
				Context:   ContextReq{Ad: AdReq{Callback: click.Callback}},
				Timestamp: click.ClickTime * 1000,
			}
			marshal, err := json.Marshal(reqData)
			if err != nil {
				fmt.Println("JSON 失败:", err)
				continue
			}
			fmt.Println("请求参数: ", string(marshal))
			post, err := resty.New().R().SetBody(marshal).Post("https://analytics.oceanengine.com/api/v2/conversion")
			if err != nil {
				fmt.Println("请求失败:", err)
				continue
			}
			fmt.Println(post.StatusCode(), post.String())
		}
		fmt.Println(clicks)
	}
}
