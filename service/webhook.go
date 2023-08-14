package service

import (
	"byr_topn/consts"
	"byr_topn/dal/kv"
	"byr_topn/utils"
	"context"
	"fmt"
	"github.com/go-lark/lark"
	"github.com/go-lark/lark/card"
	"time"
)

var bot *lark.Bot

func init() {
	bot = lark.NewNotificationBot(consts.WebHookUrl)
}

func SendTopN(ctx context.Context, byrResp *ByrResp) {
	b := lark.NewCardBuilder()
	var elements []card.Element
	for i, data := range byrResp.Data {
		div := b.Div(
			b.Field(b.Text(fmt.Sprintf("[%d. %s](%s)", i+1, data.Title, genLink(data.BoardName, data.ID))).LarkMd()),
			b.Field(b.Text(utils.PureText(data.Content))),
		)
		elements = append(elements, div)
	}
	elements = append(elements, baseDiv(b))
	postCard := b.Card(elements...).
		Wathet().
		Title(fmt.Sprintf("今日十大 %s", utils.GetYMD()))
	msg := lark.NewMsgBuffer(lark.MsgInteractive)
	resp, err := bot.PostNotificationV2(msg.Card(postCard.String()).Build())
	if err == nil && resp.StatusCode == 0 && resp.Code == 0 {
		// 如果确实发过了，则存储一下当日的发送情况，定时任务直接跳过当天
		kv.Set(ctx, fmt.Sprintf(consts.SendKey, utils.GetYMD()), "true", time.Hour*12)
	}
}

func baseDiv(b lark.CardBuilder) card.Element {
	// 贴一下开源地址
	return b.Div(
		b.Field(b.Text(fmt.Sprintf("([本项目已开源，欢迎贡献代码~](%s))", "https://github.com/whtis/byr_topn")).LarkMd()),
	)
}

func genLink(boardName string, id int64) string {
	return fmt.Sprintf("https://bbs.byr.cn/#!article/%s/%d", boardName, id)
}
