package rewrite_config

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
)

func TestRewriteConfigByAbtest(t *testing.T) {
	content,err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}
	data := map[string]interface{}{}
	err = json.Unmarshal(content, &data)
	if err != nil {
		panic(err)
	}

	//请求的中间件里把ab设置到ctx中
	abMap := map[string]string{
		"mytest_1":"1",
	}
	ctx := context.WithValue(context.Background(), "abtests", abMap)

	//定义一个返回指定ab key的value的函数
	getABValFunc := func(ctx context.Context, key string) (val string) {
		return ctx.Value("abtests").(map[string]string)[key]
	}

	//用ab对应的参数覆盖原始参数
	data = RewriteConfigByAbtest(ctx, data, getABValFunc)
	fmt.Println(data)
	return
}
