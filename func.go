package rewrite_config

import (
	"context"
	"strconv"
	"strings"
)

func rewrite(ctx context.Context, data map[string]interface{}, testParams map[string]interface{}) (newData map[string]interface{}) {

	for rewriteKey,rewriteVal := range testParams {
		//分割字段
		fields := strings.Split(rewriteKey,".")
		lastIdx := len(fields) - 1

		var preFieldValue interface{} = data

		fieldAddr := make([]interface{},0,len(fields) + 1)
		fieldAddr = append(fieldAddr, preFieldValue)

		for currFieldIdx,fieldName := range fields {
			var ok bool
			//如果是数字，代表是数组的key，还要判断是否为数组类型，因为有的人会把数字当做对象的key
			idx,err := strconv.Atoi(fieldName)
			if err == nil {
				_,ok = preFieldValue.([]interface{})
			}
			if ok {
				//超出数组长度,后续的字段不会再处理
				if idx >= len(preFieldValue.([]interface{})) {
					break
				}
				//如果是最后一个key替换
				if currFieldIdx == lastIdx {
					preFieldValue.([]interface{})[idx] = rewriteVal
				}else{
					preFieldValue = preFieldValue.([]interface{})[idx]
				}
			}else{
				//优先命令处理
				if fieldName == "$append" && currFieldIdx == lastIdx && currFieldIdx > 0 {
					//当前 currFieldValue 一定是数组，追加一个新元素(格式：likes.$append，likes.2.2.$append)
					if _,ok = preFieldValue.([]interface{});ok {
						//append后，preFieldValue地址会指向新的切片地址
						newValue := append(preFieldValue.([]interface{}), rewriteVal)
						//根据上个字段，判断上上个是数组还是map
						preFieldStr := fields[currFieldIdx-1]
						if preFieldInt,err := strconv.Atoi(preFieldStr);err != nil {
							fieldAddr[len(fieldAddr) - 2].(map[string]interface{})[preFieldStr] = newValue
						}else{
							fieldAddr[len(fieldAddr) - 2].([]interface{})[preFieldInt] = newValue
						}
					}
				} else if fieldName == "$delete" && currFieldIdx == lastIdx && currFieldIdx > 0 {
					//根据上个字段，判断上上个是数组还是map
					preFieldStr := fields[currFieldIdx-1]
					preFieldInt,err := strconv.Atoi(preFieldStr)
					if err == nil {
						_,ok = fieldAddr[len(fieldAddr) - 2].([]interface{})
					}
					if ok {
						//数组类型
						a := fieldAddr[len(fieldAddr) - 2].([]interface{})[:preFieldInt]
						b := fieldAddr[len(fieldAddr) - 2].([]interface{})[preFieldInt+1:]
						a = append(a,b...)
						//根据上上个字段，判断上上上个是数组还是map
						preFieldStr = fields[currFieldIdx-2]
						if preFieldInt,err = strconv.Atoi(preFieldStr); err != nil {
							//map类型
							fieldAddr[len(fieldAddr) - 3].(map[string]interface{})[preFieldStr] = a
						}else{
							fieldAddr[len(fieldAddr) - 3].([]interface{})[preFieldInt] = a
						}
					}else{
						//map 类型
						delete(fieldAddr[len(fieldAddr) - 2].(map[string]interface{}), fields[currFieldIdx-1])
					}
				} else if fieldName == "$insert" && currFieldIdx == lastIdx && currFieldIdx > 0 {
					//数组元素插入
					//根据上个字段，判断上上个是不是数组
					preFieldStr := fields[currFieldIdx-1]
					preFieldInt,err := strconv.Atoi(preFieldStr)
					if err == nil {
						arr,ok := fieldAddr[len(fieldAddr) - 2].([]interface{})
						if ok {
							insertArr,ok := rewriteVal.([]interface{})
							if ok {
								newArr := append(arr[:preFieldInt], append(insertArr, arr[preFieldInt:]...)...)
								//根据上上个字段，判断上上上个是数组还是map
								preFieldStr = fields[currFieldIdx-2]
								if preFieldInt,err = strconv.Atoi(preFieldStr); err != nil {
									//map类型
									fieldAddr[len(fieldAddr) - 3].(map[string]interface{})[preFieldStr] = newArr
								}else{
									fieldAddr[len(fieldAddr) - 3].([]interface{})[preFieldInt] = newArr
								}
							}
						}
					}
				} else {
					//代表是map
					if _,ok = preFieldValue.(map[string]interface{}); !ok{
						break
					}
					//如果是最后一个key替换
					if currFieldIdx == lastIdx {
						preFieldValue.(map[string]interface{})[fieldName] = rewriteVal
					}else{
						preFieldValue = preFieldValue.(map[string]interface{})[fieldName]
					}
				}
			}
			//保存到map中，以做反查
			fieldAddr = append(fieldAddr, preFieldValue)
		}
	}

	newData = data
	//fmt.Println(newData)
	return
}