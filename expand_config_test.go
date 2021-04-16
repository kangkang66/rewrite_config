package rewrite_config

import (
	"context"
	"fmt"
	"io/ioutil"
	"testing"
)

func TestRewriteConfigByExplandsByte(t *testing.T) {
	//获取配置
	content,err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}

	//回调函数
	callBackFunc := func(ctx context.Context, field, operator, value string) (result bool) {
		uid := "120"
		switch field {
		case "uid":
			switch operator {
			case "=":
				result = uid == value
			case "!=":
				result = uid != value
			}
		}
		return
	}

	//方法1
	content = RewriteConfigByExplandsByte(context.Background(), content, callBackFunc)
	fmt.Println(string(content))

}
