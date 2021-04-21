# 安装
> go get github.com/kangkang66/rewrite_config


# 介绍
本库用来通过配置来修改json，常用的用处
1. 一套配置需要按不同的ab实验组有不同的值
2. 不同的版本用户拉到的配置有不同的值
3. 不同的用户群体有不同的配置值


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


## params配置介绍

- 原始配置
```
{
  "watch_video": {
    "show": 1,
    "reward": 100,
    "daily_max": 50,
    "weight": 10
  },
  "likes": [
    "a","b"
  ]
}
```


### 指定对象字段删除
```
"params": {
    "watch_video.show.$delete":""
}
```
- 输出
```
{
  "watch_video": {
    "reward": 100,
    "daily_max": 50,
    "weight": 10
  },
  "likes": [
    "a","b"
  ]
}
```

### 指定数组元素删除
```
"params": {
    "likes.0.$delete":""
}
```

- 输出
```
{
  "watch_video": {
    "show": 1,
    "reward": 100,
    "daily_max": 50,
    "weight": 10
  },
  "likes": [
    "b"
  ]
}
```

### 指定对象字段的修改
```
"params": {
    "watch_video.show":0
}
```

- 输出
```
{
  "watch_video": {
    "show": 0,
    "reward": 100,
    "daily_max": 50,
    "weight": 10
  },
  "likes": [
    "a","b"
  ]
}
```

### 指定对象全覆盖
```
"params": {
    "watch_video":{
      "show": 0,
      "reward": 0,
      "daily_max": 0,
      "weight": 0
    }
}
```

- 输出
```
{
  "watch_video":{
    "show": 0,
    "reward": 0,
    "daily_max": 0,
    "weight": 0
  },
  "likes": [
    "a","b"
  ]
}
```

### 指定数组元素的修改
```
"params": {
    "likes.0":"bbb"
}
```

- 输出
```
{
  "watch_video": {
    "show": 1,
    "reward": 100,
    "daily_max": 50,
    "weight": 10
  },
  "likes": [
    "bbb","b"
  ]
}
```

### 添加新的对象字段
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

- 输出
```
{
  "watch_video": {
    "show": 1,
    "reward": 100,
    "daily_max": 50,
    "weight": 10
  },
  "play_game":{
        "show": 1,
        "reward": 100,
        "daily_max": 50,
        "weight": 10
  },
  "likes": [
    "a","b"
  ]
}
```

### 添加新的数组元素
```
"params": {
    "likes.$append":"yyy"
}
```

- 输出
```
{
  "watch_video": {
    "show": 1,
    "reward": 100,
    "daily_max": 50,
    "weight": 10
  },
  "likes": [
    "a","b","yyy"
  ]
}
```

### 替换指定下标的数组元素
```
"params": {
    "likes.0":"yyy"
}
```

- 输出
```
{
  "watch_video": {
    "show": 1,
    "reward": 100,
    "daily_max": 50,
    "weight": 10
  },
  "likes": [
    "yyy","b"
  ]
}
```

### 数组元素，指定位置插入新元素
```
"params": {
    "likes.1.$insert":["orange"]
}
```

- 输出
```
{
  "watch_video": {
    "show": 1,
    "reward": 100,
    "daily_max": 50,
    "weight": 10
  },
  "likes": [
    "a","orange","b"
  ]
}
```


## 支持的方法

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