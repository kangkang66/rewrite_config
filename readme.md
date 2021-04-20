# 安装
> go get github.com/kangkang66/rewrite_config


# 介绍
当需要做AB、不同版本不同配置、或者满足某个自定义条件时，配置里的某些内容需要不一样，只需要按如下方法修改配置即可。

## AB的配置规则
不同的实验组用户看到的配置不一样，例如命中实验组1的用户奖励200，实验组2的用户奖励300。  
在正常的配置内容中，增加"abtests"项
```
{
    "watch_video": {
        "show": 1,
        "reward": 100,
        "daily_max": 50,
        "weight": 10
    },
    "abtests" : [
        {
          "enable": true,
          //获取实验的key
          "ab_key": "mytest_1",
          "tests":[
            {
              //命中的实验组值
              "ab_val": "1",
              "params": {
                "watch_video.reward": 200,
              }
            },
            {
              //命中的实验组值
              "ab_val": "2",
              "params": {
                "watch_video.reward": 300,
              }
            }
          ]
        }
      ],
}
```


## 版本号的配置规则
不同版本号的用户看到的配置不一样，例如120版本的用户看不到任务
```
{
    "watch_video": {
        "show": 1,
        "reward": 100,
        "daily_max": 50,
        "weight": 10
    },
    "filters": [
        {
          "enable": true,
          "version": 120,
          "operator": "=",
          "params":{
            "watch_video.show": 0,
          }
        }
    ]
}
```

## 自定义条件的配置规则
除了上面常用的ab、版本号的条件外，还有些例如像一些其他的需求，比如不同渠道、不同设备展示的用户配置要不一样。这就可以通过自定义条件来满足。
例如要求uid=120这个用户的奖励是3000。
```
{
    "watch_video": {
        "show": 1,
        "reward": 100,
        "daily_max": 50,
        "weight": 10
    },
    "expands": [
        {
          "enable": true,
          "field": "uid",
          "operator": "=",
          "value": "120",
          "params":{
            "watch_video.reward": 3000,
          }
        }
    ]
}
```


# 使用

## 原始配置
```
{
  "watch_video": {
    "show": 1,
    "reward": 100,
    "daily_max": 50,
    "weight": 10
  },
  "likes2": {
    "shuma": {
      "apple": [
        [
          "apple11"
        ],
        [
          "apple12"
        ]
      ]
    }
  },
  "likes": [
    "b","b","bigApple"
  ],
  "goods": {
    "1": {
      "id": 1,
      "amount": 1
    },
    "2": {
      "id": 2,
      "amount": 2
    }
  },
  "abtests" : [
    {
      "enable": false,
      "ab_key": "mytest_1",
      "tests":[
        {
          "ab_val": "1",
          "params": {
            //直接完整覆盖
            "watch": {
                 "show": 2
            },
            //指定字段覆盖
            "watch_video.show":22,
            //指定数组元素替换
            "likes.0": 11,
            //命令$append，针对数组元素追加
            "likes.$append": 4,
            //命令$delete，针对数组元素删除
            "likes_2.$delete": "",
            //命令$delete，针对对象元素删除
            "watch_video_2.award.$delete":""
          }
        }
      ]
    }
  ],
  "filters": [
    {
      "enable": true,
      "version": 120,
      "operator": "=",
      "params":{
        "watch_video.weight": 120,
        "likes.0":["banana120"]
      }
    }
  ],
  "expands": [
    {
      "enable": true,
      "field": "uid",
      "operator": "=",
      "value": "120",
      "params":{
        "watch_video.weight": 333,
        "likes.0":["3333"]
      }
    }
  ]
}
```

## params覆盖配置的用法

- 指定对象字段删除
```
"params": {
    "watch_video.show.$delete":""
}
```

- 指定数组元素删除
```
"params": {
    "likes.0.$delete":""
}
```

- 指定对象字段的修改
```
"params": {
    "watch_video.show":0
}
```

- 指定数组元素的修改
```
"params": {
    "likes.0":"bbb"
}
```

- 添加新的对象字段
```
"params": {
    "play_game":{
        "show": 1,
        "reward": 100,
        "daily_max": 50,
        "weight": 10
    }
}
```

- 添加新的数组元素
```
"params": {
    "likes.$append":"yyy"
}
```

- 替换指定下标的数组元素
```
"params": {
    "likes.0.$replace":"yyy"
}
```



## 方法

- 前期准备数据
```
    //请求的中间件里把ab设置到ctx中
	abMap := map[string]string{
		"mytest_1":"1",
	}
	ctx := context.WithValue(context.Background(), "abtests", abMap)


	//获取配置文件内容直接
	configContent,err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}
```

### RewriteConfigByAbtestByte
处理ab配置，传入的配置数据类型为[]byte

```
    //定义一个返回指定ab key的value的回调函数
    getABValFunc := func(ctx context.Context, key string) (val string) {
        return ctx.Value("abtests").(map[string]string)[key]
    }

	//直接通过[]byte方法处理配置
	configContent = RewriteConfigByAbtestByte(ctx, configContent, getABValFunc)
	fmt.Println(string(configContent))
```

### RewriteConfigByAbtest
处理ab配置，传入的配置数据类型为map[string]interface{}
```
    //定义一个返回指定ab key的value的回调函数
    getABValFunc := func(ctx context.Context, key string) (val string) {
        return ctx.Value("abtests").(map[string]string)[key]
    }

	data := map[string]interface{}{}
    err = json.Unmarshal(configContent, &data)
    if err != nil {
        panic(err)
    }
    data = RewriteConfigByAbtest(ctx, data, getABValFunc)
    fmt.Println(data)
```

### RewriteConfigByFilterByte
处理version配置，传入的配置数据类型为[]byte
```
    //定义一个返回版本号的回调函数
    getVersionFunc := func(ctx context.Context) (version int64) {
        return 120
    }

    content = RewriteConfigByFilterByte(context.Background(), content, getVersionFunc)
	fmt.Println(string(content))
```

### RewriteConfigByFilter
处理version配置，传入的配置数据类型为map[string]interface{}
```
    //定义一个返回版本号的回调函数
    getVersionFunc := func(ctx context.Context) (version int64) {
        return 120
    }

    data := map[string]interface{}{}
	err = json.Unmarshal(content, &data)
	if err != nil {
		panic(err)
	}
	data = RewriteConfigByFilter(context.Background(), data, getVersionFunc)
	fmt.Println(data)
```

### RewriteConfigByExplandsByte
处理自定义条件覆盖
```
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

```

### RewriteConfigByExplands
处理自定义条件覆盖
```
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
    
    data := map[string]interface{}{}
    err = json.Unmarshal(content, &data)
    if err != nil {
        panic(err)
    }
    data = RewriteConfigByExplands(context.Background(), data, callBackFunc)
    fmt.Println(data)

```