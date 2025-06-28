package test

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"testing"
	"time"
)

func TestReq(t *testing.T) {
	client := resty.New()
	location, _ := time.LoadLocation("Asia/Shanghai")
	s, _ := time.ParseInLocation(time.DateTime, "2025-04-05 00:00:00", location)
	startTime := s.UnixMilli()

	endTime := s.AddDate(0, 0, 62).UnixMilli()
	resp, _ := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Access-Token", "d7fac5a01542f6fb5d1ff5f3b7975502").
		SetBody(fmt.Sprintf(
			`{"advertiser_id":65467682,"start_time":%d,"end_time":%d,"page_num":1,"page_size":20,"unit_ids":[10013029918]}`,
			startTime,
			endTime,
		)).
		Post("https://ad.e.kuaishou.com/rest/openapi/gw/esp/promotion/mapi/unit/detail")

	fmt.Println(time.Unix(startTime/1000, 0).Format(time.DateTime), time.Unix(endTime/1000, 0).Format(time.DateTime))
	fmt.Println(string(resp.Body()))
}
