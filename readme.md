# 安装
> go get github.com/kangkang66/rewrite_config


# 介绍

```
当需要做ABtest时，只需要按如下方法修改配置即可：

在正常的配置内容中，增加"abtests"项，如：

 {
      // 原有配置项目（主配置）...

      ...

      // 加入ab相关配置
      "abtests":[
          //这里每一项对应一个实验
          {
       
          "enable": 1 ,   // 可选，是否开启ABtest，默认是1

          // 可选，默认值为当前配置项的key。当多个配置要共用一组实验时，可以将此配置设置成同一个值。
          //  当一个配置中进行多组实验时，每一组对应的cfg_test_tag需要不同 
          "ab_key": "mytest",

          // 各组实验参数，如果某一组实验使用默认的全局配置，可以不对其进行配置
          "tests":[
              {
                  "ab_key":2,    // 对应在实验配置中的 cfg_test_tag指向的值

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
      //按版本号替换
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
        ]

 }
```

# 使用

> 参考test.go

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

   //也可以换成filters
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
            "watch_video_2.show":22,
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
