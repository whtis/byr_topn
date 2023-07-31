package service

import (
	"byr_topn/consts"
	"byr_topn/dal/kv"
	"bytes"
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var client = &http.Client{}

func Login(ctx context.Context) bool {
	log.Println("start login...")
	// Create a new multipart writer
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	// Add any additional form fields
	extraParams := map[string]string{
		"username": consts.ByrUserName,
		"password": consts.ByrPasswd,
	}
	for key, val := range extraParams {
		_ = bodyWriter.WriteField(key, val)
	}
	// Close the multipart writer
	err := bodyWriter.Close()
	if err != nil {
		return false
	}
	url := "https://bbs.byr.cn/n/b/auth/login.json"
	req, err := http.NewRequest("POST", url, bodyBuf)
	if err != nil {
		return false
	}
	req.Header.Set("x-requested-with", "XMLHttpRequest")
	req.Header.Set("Content-Type", bodyWriter.FormDataContentType())
	req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 16_1_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.1 Mobile/15E148 Safari/604.1")
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		cookies := resp.Header["Set-Cookie"]
		setLoginMap(ctx, cookies)
	}
	return true
}

func setLoginMap(ctx context.Context, cookies []string) {
	var keyValuePairs []string
	expiresRegex := regexp.MustCompile(`expires=([^;]+)`)
	keyValueRegex := regexp.MustCompile(`([^=\s]+)=([^;\s]+)`)
	for _, s := range cookies {
		match := keyValueRegex.FindStringSubmatch(s)
		if len(match) == 3 {
			keyValuePairs = append(keyValuePairs, fmt.Sprintf("%s=%s", match[1], match[2]))
		}
		expiresMatch := expiresRegex.FindStringSubmatch(s)
		if len(expiresMatch) == 2 {
			expires := expiresMatch[1]
			layout := "Mon, 02-Jan-2006 15:04:05 MST"
			timestamp, err := time.Parse(layout, expires)
			if err == nil {
				kv.Set(ctx, consts.ExpireKey, strconv.FormatInt(timestamp.Unix(), 10), 0)
			}
		}
	}
	kv.Set(ctx, consts.CookieKey, strings.Join(keyValuePairs, ";"), 0)
}
