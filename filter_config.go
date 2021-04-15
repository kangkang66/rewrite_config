package rewrite_config

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/kangkang66/xcompare"
)

type VersionFunc func(ctx context.Context) (version int64)

func RewriteConfigByFilterByte(ctx context.Context, configData []byte, getVersionFunc VersionFunc) (newJsonByte []byte) {
	newJsonByte = configData
	//判断是否包含"filters"，没有直接返回
	if !bytes.Contains(configData, []byte("filters")) {
		return
	}
	cfgData := map[string]interface{}{}
	err := json.Unmarshal(configData, &cfgData)
	if err != nil {
		return
	}
	cfgData = RewriteConfigByFilter(ctx, cfgData, getVersionFunc)
	configData,err = json.Marshal(cfgData)
	if err != nil {
		return
	}
	newJsonByte = configData
	return
}

func RewriteConfigByFilter(ctx context.Context, configData map[string]interface{}, getVersionFunc VersionFunc) (newData map[string]interface{}) {
	newData = configData
	//不存在直接返回
	_,ok := configData["filters"]
	if !ok {
		return
	}
	for _,filterVal := range configData["filters"].([]interface{}) {
		filter,ok := filterVal.(map[string]interface{})
		if !ok {
			continue
		}
		params := getFilterParams(ctx, filter, getVersionFunc)
		if len(params) == 0 {
			continue
		}
		//把params覆盖当前的data
		newData = rewrite(ctx, newData, params)
	}
	delete(newData,"filters")
	return
}

func getFilterParams(ctx context.Context, filter map[string]interface{}, getVersionFunc func(ctx context.Context) (version int64)) (filterParams map[string]interface{})  {
	enableField,ok := filter["enable"]
	if !ok {
		//默认为false
		return
	}
	enable,ok := enableField.(bool)
	if !ok || !enable {
		return
	}

	filterVersionInter,ok := filter["version"]
	if !ok {
		return
	}

	operatorInter,ok := filter["operator"]
	if !ok {
		return
	}
	operator,ok := operatorInter.(string)
	if !ok {
		return
	}

	if versionOperator(float64(getVersionFunc(ctx)), filterVersionInter, operator) {
		filterParams,_ = filter["params"].(map[string]interface{})
	}
	return
}

//因为go把json中的数字全部都按float64类型接收，所以为了方便也把userVersion换成float64
func versionOperator(userVersion float64, filterVersion interface{}, operator string) (hit bool) {
	switch operator {
	case "in":
		if version,ok := filterVersion.([]float64);ok {
			return xcompare.IN.Float64(userVersion, version)
		}
	case "=","equal":
		if version,ok := filterVersion.(float64);ok {
			return xcompare.Equal.Float64(userVersion, version)
		}
	case "!=","not_equal":
		if version,ok := filterVersion.(float64);ok {
			return xcompare.NotEqual.Float64(userVersion, version)
		}
	case ">","great":
		if version,ok := filterVersion.(float64);ok {
			return xcompare.Great.Float64(userVersion, version)
		}
	case ">=","great_equal":
		if version,ok := filterVersion.(float64);ok {
			return xcompare.GreatEqual.Float64(userVersion, version)
		}
	case "<","little":
		if version,ok := filterVersion.(float64);ok {
			return xcompare.Litter.Float64(userVersion, version)
		}
	case "<=","little_equal":
		if version,ok := filterVersion.(float64);ok {
			return xcompare.LitterEqual.Float64(userVersion, version)
		}
	}
	return
}