package rewrite_config

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
)

func TestRewriteConfigByAbtest(t *testing.T) {


	//请求的中间件里把ab设置到ctx中
	abMap := map[string]string{
		"mytest_1":"1",
	}
	ctx := context.WithValue(context.Background(), "abtests", abMap)

	//定义一个返回指定ab key的value的函数
	getABValFunc := func(ctx context.Context, key string) (val string) {
		return ctx.Value("abtests").(map[string]string)[key]
	}

	//获取配置文件内容直接
	configContent,err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}

	//第一种方法： 直接通过[]byte方法处理配置
	configContent = RewriteConfigByAbtestByte(ctx, configContent, getABValFunc)
	fmt.Println(string(configContent))

	//第二种方法： 通过map类型处理配置
	data := map[string]interface{}{}
	err = json.Unmarshal(configContent, &data)
	if err != nil {
		panic(err)
	}
	data = RewriteConfigByAbtest(ctx, data, getABValFunc)
	fmt.Println(data)
	return
}
