package rewrite_config

import (
	"bytes"
	"context"
	"encoding/json"
)

type ABFunc func(ctx context.Context, key string) (val string)

func RewriteConfigByAbtestByte(ctx context.Context, configData []byte, getABValueFunc ABFunc) (newJsonByte []byte) {
	newJsonByte = configData
	//判断是否包含"abtests"，没有直接返回
	if !bytes.Contains(configData, []byte("abtests")) {
		return
	}

	cfgData := map[string]interface{}{}
	err := json.Unmarshal(configData, &cfgData)
	if err != nil {
		return
	}
	cfgData = RewriteConfigByAbtest(ctx, cfgData, getABValueFunc)
	configData,err = json.Marshal(cfgData)
	if err != nil {
		return
	}
	newJsonByte = configData
	return
}

func RewriteConfigByAbtest(ctx context.Context, configData map[string]interface{}, getABValueFunc ABFunc) (newData map[string]interface{}) {
	newData = configData
	//不存在直接返回
	_,ok := configData["abtests"]
	if !ok {
		return
	}
	for _,abtestVal := range configData["abtests"].([]interface{}) {
		abtest,ok := abtestVal.(map[string]interface{})
		if !ok {
			continue
		}
		//获取命中的实验params
		testParams := getABTestParams(ctx,abtest,getABValueFunc)
		if len(testParams) == 0 {
			continue
		}
		//把params覆盖当前的data
		newData = rewrite(ctx, newData, testParams)
	}

	delete(newData,"abtests")
	return
}

func getABTestParams(ctx context.Context, abtest map[string]interface{}, getABValueFunc func(ctx context.Context, key string) (val string)) (testParams map[string]interface{}) {
	//是否启用这个实验
	enableField,ok := abtest["enable"]
	if !ok {
		return
	}
	enable,ok := enableField.(bool)
	if !ok || !enable {
		return
	}

	//获取abKey
	abKeyField,ok := abtest["ab_key"]
	if !ok || abKeyField.(string) == "" {
		return
	}
	abVal := getABValueFunc(ctx, abKeyField.(string))

	//获取tests
	testsField,ok := abtest["tests"]
	if !ok {
		return
	}
	tests,ok := testsField.([]interface{})
	if !ok {
		return
	}
	for _,testVal := range tests {
		test,ok := testVal.(map[string]interface{})
		if !ok {
			continue
		}
		abValField,ok := test["ab_val"]
		if !ok || abValField.(string) != abVal{
			continue
		}
		//配置的abval == 当前请求的abval
		paramsField,ok := test["params"]
		if !ok {
			continue
		}
		//把当前的params赋值，返回出去
		testParams,ok = paramsField.(map[string]interface{})
		if ok {
			return
		}
	}

	return
}

