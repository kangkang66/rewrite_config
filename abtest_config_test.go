package rewrite_config

import (
	"context"
	"fmt"
	"testing"
)

//00. 指定map字段删除
func Test_00(t *testing.T) {
	//请求的中间件里把ab设置到ctx中
	abMap := map[string]string{
		"mytest_1":"1",
	}
	ctx := context.WithValue(context.Background(), "abtests", abMap)
	//定义一个返回指定ab key的value的函数
	getABValFunc := func(ctx context.Context, key string) (val string) {
		return ctx.Value("abtests").(map[string]string)[key]
	}

	input := `
{
	"watch_video": {
		"show": 1,
		"reward": 1,
		"daily_max": 1,
		"weight": 1
	},
	"abtests":[
        {
            "enable": true,
            "ab_key": "mytest_1",
            "tests":[
                {
                    "ab_val": "1",
                    "params": {
                        "watch_video.show.$delete":""
                    }
                }
            ]
        }
    ]
}`
	output := `{"watch_video":{"daily_max":1,"reward":1,"weight":1}}`
	testOutput := RewriteConfigByAbtestByte(ctx, []byte(input), getABValFunc)

	if string(testOutput) != output {
		fmt.Println("Test_00 not pass")
	}else{
		fmt.Println("Test_00 pass")
	}
}

//0. 字段新增
func Test_0(t *testing.T) {
	//请求的中间件里把ab设置到ctx中
	abMap := map[string]string{
		"mytest_1":"1",
	}
	ctx := context.WithValue(context.Background(), "abtests", abMap)
	//定义一个返回指定ab key的value的函数
	getABValFunc := func(ctx context.Context, key string) (val string) {
		return ctx.Value("abtests").(map[string]string)[key]
	}

	input := `
{
	"watch_video": {
		"show": 1,
		"reward": 1,
		"daily_max": 1,
		"weight": 1
	},
	"abtests":[
        {
            "enable": true,
            "ab_key": "mytest_1",
            "tests":[
                {
                    "ab_val": "1",
                    "params": {
                        "play_game":{
                            "show": 1,
                            "reward": 2,
                            "daily_max": 1,
                            "weight": 1
                        }
                    }
                }
            ]
        }
    ]
}`
	output := `{"play_game":{"daily_max":1,"reward":2,"show":1,"weight":1},"watch_video":{"daily_max":1,"reward":1,"show":1,"weight":1}}`
	testOutput := RewriteConfigByAbtestByte(ctx, []byte(input), getABValFunc)

	if string(testOutput) != output {
		fmt.Println("Test_0 not pass")
	}else{
		fmt.Println("Test_0 pass")
	}
}

//1. 整个字段直接覆盖
func Test_1(t *testing.T) {
	//请求的中间件里把ab设置到ctx中
	abMap := map[string]string{
		"mytest_1":"1",
	}
	ctx := context.WithValue(context.Background(), "abtests", abMap)
	//定义一个返回指定ab key的value的函数
	getABValFunc := func(ctx context.Context, key string) (val string) {
		return ctx.Value("abtests").(map[string]string)[key]
	}

	input := `
{
	"watch_video": {
		"show": 1,
		"reward": 1,
		"daily_max": 1,
		"weight": 1
	},
	"abtests":[
        {
            "enable": true,
            "ab_key": "mytest_1",
            "tests":[
                {
                    "ab_val": "1",
                    "params": {
                        "watch_video":{
                            "show": 2,
                            "reward": 2,
                            "daily_max": 2,
                            "weight": 2
                        }
                    }
                }
            ]
        }
    ]
}`
	output := `{"watch_video":{"daily_max":2,"reward":2,"show":2,"weight":2}}`
	testOutput := RewriteConfigByAbtestByte(ctx, []byte(input), getABValFunc)

	if string(testOutput) != output {
		fmt.Println("Test_1 not pass")
	}else{
		fmt.Println("Test_1 pass")
	}
}

//2. 指定map字段覆盖
func Test_2(t *testing.T) {
	//请求的中间件里把ab设置到ctx中
	abMap := map[string]string{
		"mytest_1":"1",
	}
	ctx := context.WithValue(context.Background(), "abtests", abMap)
	//定义一个返回指定ab key的value的函数
	getABValFunc := func(ctx context.Context, key string) (val string) {
		return ctx.Value("abtests").(map[string]string)[key]
	}

	input := `
{
	"watch_video": {
		"show": 1,
		"reward": 1,
		"daily_max": 1,
		"weight": 1
	},
	"abtests":[
        {
            "enable": true,
            "ab_key": "mytest_1",
            "tests":[
                {
                    "ab_val": "1",
                    "params": {
                        "watch_video.show":2
                    }
                }
            ]
        }
    ]
}`
	output := `{"watch_video":{"daily_max":1,"reward":1,"show":2,"weight":1}}`
	testOutput := RewriteConfigByAbtestByte(ctx, []byte(input), getABValFunc)

	if string(testOutput) != output {
		fmt.Println("Test_2 not pass")
	}else{
		fmt.Println("Test_2 pass")
	}
}

//3. 数组字段覆盖
func Test_3(t *testing.T) {
	//请求的中间件里把ab设置到ctx中
	abMap := map[string]string{
		"mytest_1":"1",
	}
	ctx := context.WithValue(context.Background(), "abtests", abMap)
	//定义一个返回指定ab key的value的函数
	getABValFunc := func(ctx context.Context, key string) (val string) {
		return ctx.Value("abtests").(map[string]string)[key]
	}

	input := `
{
	"watch_video": {
		"show": 1,
		"reward": 1,
		"daily_max": 1,
		"weight": 1
	},
	"like":["apple","banana"],
	"abtests":[
        {
            "enable": true,
            "ab_key": "mytest_1",
            "tests":[
                {
                    "ab_val": "1",
                    "params": {
                        "like.0":"bigApple"
                    }
                }
            ]
        }
    ]
}`
	output := `{"like":["bigApple","banana"],"watch_video":{"daily_max":1,"reward":1,"show":1,"weight":1}}`
	testOutput := RewriteConfigByAbtestByte(ctx, []byte(input), getABValFunc)

	if string(testOutput) != output {
		fmt.Println("Test_3 not pass")
	}else{
		fmt.Println("Test_3 pass")
	}
}

//4. 数组字段追加
func Test_4(t *testing.T) {
	//请求的中间件里把ab设置到ctx中
	abMap := map[string]string{
		"mytest_1":"1",
	}
	ctx := context.WithValue(context.Background(), "abtests", abMap)
	//定义一个返回指定ab key的value的函数
	getABValFunc := func(ctx context.Context, key string) (val string) {
		return ctx.Value("abtests").(map[string]string)[key]
	}

	input := `
{
	"watch_video": {
		"show": 1,
		"reward": 1,
		"daily_max": 1,
		"weight": 1
	},
	"like":["apple","banana"],
	"abtests":[
        {
            "enable": true,
            "ab_key": "mytest_1",
            "tests":[
                {
                    "ab_val": "1",
                    "params": {
                        "like.$append":"orange"
                    }
                }
            ]
        }
    ]
}`
	output := `{"like":["apple","banana","orange"],"watch_video":{"daily_max":1,"reward":1,"show":1,"weight":1}}`
	testOutput := RewriteConfigByAbtestByte(ctx, []byte(input), getABValFunc)

	if string(testOutput) != output {
		fmt.Println("Test_4 not pass")
	}else{
		fmt.Println("Test_4 pass")
	}
}

//5. 数组字段指定元素删除
func Test_5(t *testing.T) {
	//请求的中间件里把ab设置到ctx中
	abMap := map[string]string{
		"mytest_1":"1",
	}
	ctx := context.WithValue(context.Background(), "abtests", abMap)
	//定义一个返回指定ab key的value的函数
	getABValFunc := func(ctx context.Context, key string) (val string) {
		return ctx.Value("abtests").(map[string]string)[key]
	}

	input := `
{
	"watch_video": {
		"show": 1,
		"reward": 1,
		"daily_max": 1,
		"weight": 1
	},
	"like":["apple","banana"],
	"abtests":[
        {
            "enable": true,
            "ab_key": "mytest_1",
            "tests":[
                {
                    "ab_val": "1",
                    "params": {
                        "like.0.$delete":""
                    }
                }
            ]
        }
    ]
}`
	output := `{"like":["banana"],"watch_video":{"daily_max":1,"reward":1,"show":1,"weight":1}}`
	testOutput := RewriteConfigByAbtestByte(ctx, []byte(input), getABValFunc)

	if string(testOutput) != output {
		fmt.Println("Test_5 not pass")
	}else{
		fmt.Println("Test_5 pass")
	}
}

func Test_6(t *testing.T) {
	//请求的中间件里把ab设置到ctx中
	abMap := map[string]string{
		"mytest_1":"1",
	}
	ctx := context.WithValue(context.Background(), "abtests", abMap)
	//定义一个返回指定ab key的value的函数
	getABValFunc := func(ctx context.Context, key string) (val string) {
		return ctx.Value("abtests").(map[string]string)[key]
	}

	input := `
{
	"watch_video": {
		"show": 1,
		"reward": 1,
		"daily_max": 1,
		"weight": 1
	},
	"like":["apple","banana"],
	"abtests":[
        {
            "enable": true,
            "ab_key": "mytest_1",
            "tests":[
                {
                    "ab_val": "1",
                    "params": {
                        "like.0":"orange"
                    }
                }
            ]
        }
    ]
}`
	output := `{"like":["orange","banana"],"watch_video":{"daily_max":1,"reward":1,"show":1,"weight":1}}`
	testOutput := RewriteConfigByAbtestByte(ctx, []byte(input), getABValFunc)

	if string(testOutput) != output {
		fmt.Println("Test_6 not pass")
	}else{
		fmt.Println("Test_6 pass")
	}
}

func Test_7(t *testing.T) {
	//请求的中间件里把ab设置到ctx中
	abMap := map[string]string{
		"mytest_1":"1",
	}
	ctx := context.WithValue(context.Background(), "abtests", abMap)
	//定义一个返回指定ab key的value的函数
	getABValFunc := func(ctx context.Context, key string) (val string) {
		return ctx.Value("abtests").(map[string]string)[key]
	}

	input := `
{
	"watch_video": {
		"show": 1,
		"reward": 1,
		"daily_max": 1,
		"weight": 1
	},
	"like":["apple","banana"],
	"abtests":[
        {
            "enable": true,
            "ab_key": "mytest_1",
            "tests":[
                {
                    "ab_val": "1",
                    "params": {
                        "like.1.$insert":["orange"]
                    }
                }
            ]
        }
    ]
}`
	output := `{"like":["apple","orange","banana"],"watch_video":{"daily_max":1,"reward":1,"show":1,"weight":1}}`
	testOutput := RewriteConfigByAbtestByte(ctx, []byte(input), getABValFunc)
	fmt.Println(string(testOutput))

	if string(testOutput) != output {
		fmt.Println("Test_7 not pass")
	}else{
		fmt.Println("Test_7 pass")
	}
}









