package rewrite_config

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
)

func TestRewriteConfigByFilter(t *testing.T) {
	//定义一个返回指定ab key的value的函数
	getVersionFunc := func(ctx context.Context) (version int64) {
		return 120
	}

	//获取配置
	content,err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}

	//方法1
	content = RewriteConfigByFilterByte(context.Background(), content, getVersionFunc)
	fmt.Println(string(content))

	//方法2
	data := map[string]interface{}{}
	err = json.Unmarshal(content, &data)
	if err != nil {
		panic(err)
	}
	data = RewriteConfigByFilter(context.Background(), data, getVersionFunc)
	fmt.Println(data)
}
