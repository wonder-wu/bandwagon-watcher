package main

import "encoding/json"

type Json struct {
	m interface{}
}

//创建新json
func NewJson(jsonstr []byte) *Json {
	_json := &Json{}
	if err := json.Unmarshal(jsonstr, &_json.m); err != nil {
		panic(err)
	}
	return _json
}

func (j *Json) getMap() map[string]interface{} {
	return (j.m).(map[string]interface{})
}

//读取
func (j *Json) Get(key string) *Json {
	m := j.getMap()
	return &Json{m[key]}
}

//获取整数
func (j *Json) Number() float64 {

	val, ok := j.m.(float64)
	if ok != true {
		panic("type assertion to float64 failed")
	}
	return val
}

//获取字符串
func (j *Json) String() string {
	if val, ok := j.m.(string); ok {
		return val
	} else {
		panic("type assertion to string failed")
		return ""
	}
}

//获取全部array
func (j *Json) AllArray() []interface{} {
	val, ok := j.m.([]interface{})
	if ok != true {
		panic("type assertion to array failed")
	}
	return val
}

//获取能继续链式调用的array
func (j *Json) Item(i int) *Json {
	val, ok := j.m.([]interface{})
	if !ok {
		panic("target can't assertion to array")
	}

	return &Json{val[i]}

	//
}
