package importService

import (
	"bytes"
	"time"

	"github.com/Monkey-Mouse/mo2/server/model"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
)

// NoTitle 未命名文章题目修改为空，方便前端处理
const NoTitle = ""

// Transform to parse md file for model.Blog
func Transform(file []byte) (blog model.Blog) {
	markdown := goldmark.New(
		goldmark.WithExtensions(
			meta.Meta,
		),
	)
	var buf bytes.Buffer
	context := parser.NewContext()
	if err := markdown.Convert(file, &buf, parser.WithContext(context)); err != nil {
		panic(err)
	}
	metaData := meta.Get(context)
	title := metaData["title"]
	date := metaData["date"]
	//categories:=metaData["categories"]
	entityInfo := model.Entity{}
	entityInfo.Set(getBlogDate(date))
	blog = model.Blog{
		Title:      getBlogTitle(title),
		EntityInfo: entityInfo,
		Content:    buf.String(),
	}
	return
}

func getBlogTitle(title interface{}) (titleStr string) {
	titleStr, ok := title.(string)
	if !ok {
		titleStr = NoTitle
	}
	return
}

func getBlogDate(date interface{}) (res time.Time) {
	dateString, ok := date.(string)
	if !ok {
		res = time.Now()
		return
	}
	var layout string
	var err error
	if len(dateString) < 23 {
		dateString = dateString[:10]
		layout = "2006-01-02"
	} else {
		dateString = dateString[:23]
		layout = "2006-01-02T15:04:05.000"
	}
	if res, err = time.Parse(layout, dateString); err != nil {
		res = time.Now()
	}
	return
}
