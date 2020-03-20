package log

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"
)

type Logger struct {
	Level       string `json:"level"`
	Identifier  string `json:"identifier"`
	TimeFormat  string `json:"time_format"`
	FileCording bool   `json:"file_cording"`
	FileName    string `json:"file_name"`
}

func newConfig(path string) (*Logger, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	//start decode json config
	var decode map[string]interface{}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &decode)
	if err != nil {
		return nil, err
	}
	var set Logger
	startDecode(decode, reflect.TypeOf(&set), reflect.ValueOf(&set))
	return &set, nil
}

func startDecode(obj interface{}, rSetType reflect.Type, rSetValue reflect.Value) {
	rObj := reflect.ValueOf(obj)
	OneKey := rObj.MapKeys()

	for _, k := range OneKey {
		TowKey := rObj.MapIndex(k).Interface().(map[string]interface{})

		rvalue1 := reflect.ValueOf(TowKey)
		//take two key
		OkTowKey := rvalue1.MapKeys()

		for _, k1 := range OkTowKey {
			for i := 0; i < rSetValue.Elem().NumField(); i++ {
				if rSetType.Elem().Field(i).Tag.Get("json") == k1.String() {

					switch rSetType.Elem().Field(i).Type.Kind() {
					case reflect.String:
						rSetValue.Elem().Field(i).SetString(rvalue1.MapIndex(k1).Interface().(string))
					case reflect.Bool:
						rSetValue.Elem().Field(i).SetBool(rvalue1.MapIndex(k1).Interface().(bool))
					case reflect.Int64:
						rSetValue.Elem().Field(i).SetInt(int64(rvalue1.MapIndex(k1).Interface().(float64)))
					}

				}
			}

		}
	}

}
