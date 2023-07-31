package service

import (
	"byr_topn/consts"
	"byr_topn/dal/kv"
	"context"
	"io"
	"log"
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

type Data struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	BoardName string `json:"board_name"`
}
type ByrResp struct {
	Success bool   `json:"success"`
	Data    []Data `json:"data"`
}

func GetTopN(ctx context.Context) *ByrResp {
	log.Println("start get topN...")
	tUrl := "https://bbs.byr.cn/n/b/home/topten.json"
	req, _ := http.NewRequest("GET", tUrl, nil)
	req.Header.Set("x-requested-with", "XMLHttpRequest")
	req.Header.Set("user-agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 16_1_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.1 Mobile/15E148 Safari/604.1")
	cookie := kv.Get(ctx, consts.CookieKey)
	if cookie != nil {
		req.Header.Set("Cookie", *cookie)
		tClient := &http.Client{}
		resp, _ := tClient.Do(req)
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			log.Println("success get topN")
			body, _ := io.ReadAll(resp.Body)
			retString := string(body)
			var byrResp ByrResp
			_ = jsoniter.UnmarshalFromString(retString, &byrResp)
			return &byrResp
		}
	}
	return nil
}
