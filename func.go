package rewrite_config

import (
	"context"
	"strconv"
	"strings"
)

func rewrite(ctx context.Context, data map[string]interface{}, testParams map[string]interface{}) (newData map[string]interface{}) {
	newData = data

	for rewriteKey,rewriteVal := range testParams {
		var preSonV, sonV interface{}
		preIdx := -1

		sonV = newData
		//分割字段
		fields := strings.Split(rewriteKey,".")
		lastKey := len(fields) - 1

		//data.$append  data一定是数组，追加一个新数组
		//data.2.$delete

		for currFieldkey,fieldName := range fields {
			var ok bool
			//如果是数字，代表是数组的key
			if idx,err := strconv.Atoi(fieldName);err == nil {
				if _,ok = sonV.([]interface{});ok {
					//别超出数组长度
					if idx > len(sonV.([]interface{})) {
						continue
					}
					preSonV = sonV
					preIdx = idx
					//如果是最后一个key替换
					if currFieldkey == lastKey {
						sonV.([]interface{})[idx] = rewriteVal
					}else{
						sonV = sonV.([]interface{})[idx]
					}
				}
			}else{
				//优先命令处理
				if fieldName == "$append" && currFieldkey == lastKey && currFieldkey != 0 {
					//一定是数组，追加一个新元素(格式：likes.$append，likes.2.$append)
					if _,ok = sonV.([]interface{});ok {
						if preIdx == -1 {
							preSonV.(map[string]interface{})[fields[currFieldkey-1]] = append(sonV.([]interface{}), rewriteVal)
						}else{
							preSonV.([]interface{})[preIdx] = append(sonV.([]interface{}), rewriteVal)
						}
					}
				}else {
					//代表是map
					//如果是最后一个key替换
					if _,ok = sonV.(map[string]interface{});ok{
						preSonV = sonV
						if currFieldkey == lastKey {
							sonV.(map[string]interface{})[fieldName] = rewriteVal
						}else{
							sonV = sonV.(map[string]interface{})[fieldName]
						}
					}
				}
			}
		}
	}

	//fmt.Println(newData)
	return
}