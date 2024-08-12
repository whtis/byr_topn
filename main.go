package main

import (
	"byr_topn/consts"
	"byr_topn/dal/kv"
	"byr_topn/service"
	"byr_topn/utils"
	"context"
	"fmt"
	"github.com/robfig/cron"
	"log"
	"strconv"
	"time"
)

func main() {
	// Start a goroutine to execute the scheduled task
	go func() {
		c := cron.New()
		c.AddFunc("0 0 */1 * * ?", func() { // 每小时一次
			topN()
		})
		c.Start()
	}()
	// Keep the main goroutine running indefinitely
	select {}
}

func topN() {
	// 每天十点发送，如果当前时间小于10点直接return
	if time.Now().Hour() < 10 {
		return
	}
	ctx := context.Background()
	// 判断当天是否发过了，如果发过了，直接return
	sendFlag := kv.Get(ctx, fmt.Sprintf(consts.SendKey, utils.GetYMD()))
	if sendFlag != nil {
		log.Println("today already send topN to feishu")
		return
	}
	// 1.判断内存中是否存在cookie,如果存在cookie，同时cookie没有失效，直接发起十大的请求
	cookie := kv.Get(ctx, consts.CookieKey)
	expireStr := kv.Get(ctx, consts.ExpireKey)
	if cookie != nil && expireStr != nil {
		expire, _ := strconv.ParseInt(*expireStr, 10, 64)
		if expire > time.Now().Unix() {
			byrResp := service.GetTopN(ctx)
			if byrResp != nil {
				service.SendTopN(ctx, byrResp)
			}
		}
	}
	if service.Login(ctx) {
		byrResp := service.GetTopN(ctx)
		if byrResp != nil {
			service.SendTopN(ctx, byrResp)
		}
	}
}
