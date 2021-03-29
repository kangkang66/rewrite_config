package rewrite_config

import (
	"context"
	"strconv"
	"strings"
)

func rewrite(ctx context.Context, data map[string]interface{}, testParams map[string]interface{}) (newData map[string]interface{}) {
	newData = data

	for rewriteKey,rewriteVal := range testParams {
		var sonV interface{}
		sonV = newData

		//分割字段
		fields := strings.Split(rewriteKey,".")
		lastKey := len(fields) - 1

		for currFieldkey,fieldName := range fields {
			//如果是数字，代表是数组的key
			if idx,err := strconv.Atoi(fieldName);err == nil {
				_,ok := sonV.([]interface{})
				if ok {
					//别超出数组长度
					if idx < len(sonV.([]interface{})) {
						//如果是最后一个key替换
						if currFieldkey == lastKey {
							sonV.([]interface{})[idx] = rewriteVal
						}else{
							sonV = sonV.([]interface{})[idx]
						}
					}
				}
			}else{
				//代表是map
				//如果是最后一个key替换
				_,ok := sonV.(map[string]interface{})
				if ok{
					if currFieldkey == lastKey {
						sonV.(map[string]interface{})[fieldName] = rewriteVal
					}else{
						sonV = sonV.(map[string]interface{})[fieldName]
					}
				}
			}
		}
	}

	//fmt.Println(newData)
	return
}