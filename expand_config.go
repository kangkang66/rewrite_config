package rewrite_config

import (
	"bytes"
	"context"
	"encoding/json"
)

type ExpandsFunc func(ctx context.Context, field, operator, value string) (result bool)

func RewriteConfigByExplandsByte(ctx context.Context, configData []byte, getFunc ExpandsFunc) (newJsonByte []byte) {
	newJsonByte = configData
	//判断是否包含"expands"，没有直接返回
	if !bytes.Contains(configData, []byte("expands")) {
		return
	}
	cfgData := map[string]interface{}{}
	err := json.Unmarshal(configData, &cfgData)
	if err != nil {
		return
	}
	cfgData = RewriteConfigByExplands(ctx, cfgData, getFunc)
	configData,err = json.Marshal(cfgData)
	if err != nil {
		return
	}
	newJsonByte = configData
	return
}

func RewriteConfigByExplands(ctx context.Context, configData map[string]interface{}, getFunc ExpandsFunc) (newData map[string]interface{}) {
	newData = configData
	//不存在直接返回
	_,ok := configData["expands"]
	if !ok {
		return
	}
	for _,filterVal := range configData["expands"].([]interface{}) {
		filter,ok := filterVal.(map[string]interface{})
		if !ok {
			continue
		}
		params := getExpandParams(ctx, filter, getFunc)
		if len(params) == 0 {
			continue
		}
		//把params覆盖当前的data
		newData = rewrite(ctx, newData, params)
	}
	delete(newData,"expands")
	return
}

func getExpandParams(ctx context.Context, expand map[string]interface{}, getFunc ExpandsFunc) (filterParams map[string]interface{})  {
	enableField,ok := expand["enable"]
	if !ok {
		//默认为false
		return
	}
	enable,ok := enableField.(bool)
	if !ok || !enable {
		return
	}

	fieldInter,ok := expand["field"]
	if !ok {
		return
	}
	field := fieldInter.(string)

	operatorInter,ok := expand["operator"]
	if !ok {
		return
	}
	operator,ok := operatorInter.(string)
	if !ok {
		return
	}

	valueInter,ok := expand["value"]
	if !ok {
		return
	}
	value,ok := valueInter.(string)
	if !ok {
		return
	}

	if getFunc(ctx, field, operator, value) {
		filterParams,_ = expand["params"].(map[string]interface{})
	}
	return
}
