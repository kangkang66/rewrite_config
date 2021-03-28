package rewrite_config

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
)

func TestRewriteConfigByFilter(t *testing.T) {
	content,err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}
	data := map[string]interface{}{}
	err = json.Unmarshal(content, &data)
	if err != nil {
		panic(err)
	}

	//定义一个返回指定ab key的value的函数
	getVersionFunc := func(ctx context.Context) (version int64) {
		return 120
	}

	data = RewriteConfigByFilter(context.Background(), data, getVersionFunc)
	fmt.Println(data)
}
