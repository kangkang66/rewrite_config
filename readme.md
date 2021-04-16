# 安装
> go get github.com/kangkang66/rewrite_config


# 介绍

当需要做ABtest、不同版本不同配置、或者满足某个自定义条件时，只需要按如下方法修改配置即可
```
在正常的配置内容中，增加"abtests"项，如：

 {
      // 原有配置项目（主配置）
      ...
      //1. 加入ab相关配置
      "abtests":[
          //这里每一项对应一个实验
          {
          "enable": 1 ,   // 可选，是否开启ABtest，默认是1
          // 可选，默认值为当前配置项的key。当多个配置要共用一组实验时，可以将此配置设置成同一个值。
          //  当一个配置中进行多组实验时，每一组对应的ab_key需要不同 
          "ab_key": "mytest",
          // 各组实验参数，如果某一组实验使用默认的全局配置，可以不对其进行配置
          "tests":[
              {
                  "ab_val":2,    // 对应在实验配置中的 ab_key指向的值
                  //指定要覆盖主配置的参数，如果命中此组试验，参数会被合并到主配置中
                  "params":{
                      "test1": "value",
                      "test2.test3.test4": "value"  //  可以用 . 指定要覆盖的子节点 
                  }
              },
              {
                  "id":3,
                  "params":{...}
              }
          ]
      }
      // 可以配置多组
      ...
      ],
      //2. 按版本号替换
      "filters": [
            {
              "enable": true,
              "version": 120,
              //支持 in，=，!=，>，>=，<，<=
              "operator": "=",
              "params":{
                "watch_video.weight": 120,
                "likes.1":"banana120"
              }
            }
            // 可以配置多组
            ...
      ],
    //3. 高级扩展用法，按照自定义条件来覆盖配置
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
        // 可以配置多组
        ...
    ]

 }
```

# 使用

## 覆盖配置的用法

```
{
  "watch_video_1": {
    "show": 1
  },
   "watch_video_2": {
    "show": 1,
    "award":100
  },
  "likes":[1,2,3],
  "likes_2":[6,7,8],

   //也可以换成filters，expands
  "abtests" : [
    {
      "enable": true,
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
  ]
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
    //定义一个返回指定ab key的value的函数
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
    //定义一个返回指定ab key的value的函数
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
    //定义一个返回指定ab key的value的函数
    getVersionFunc := func(ctx context.Context) (version int64) {
        return 120
    }

    content = RewriteConfigByFilterByte(context.Background(), content, getVersionFunc)
	fmt.Println(string(content))
```

### RewriteConfigByFilter
处理version配置，传入的配置数据类型为map[string]interface{}
```
    //定义一个返回指定ab key的value的函数
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