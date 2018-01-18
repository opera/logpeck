package logpeck

import (
	"fmt"
	lua "github.com/yuin/gopher-lua"
	"testing"
)

func TestExtractor(*testing.T) {
	confStr := `{ 
		"Name":"lua",
		"Config":{
			"LuaString":"function conv(s) local ret = {} ret['haha'] = string.sub(s, 2, -2) return ret end" }
		}
	}`
	config, err := NewExtractorConfig(confStr)
	if err != nil {
		panic(err)
	}
	fmt.Printf("[NewExtractorConfig]%#v\n", config)

	extractor, err := NewExtractor(config, nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", extractor)
}

func TestLuaExtractorConfig(*testing.T) {
	confStr := `{ "LuaString":"function conv(s) local ret = {} ret['haha'] = string.sub(s, 2, -2) return ret end" }`
	config, err := NewLuaExtractorConfig([]byte(confStr))
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", config)

	extractor, err := NewLuaExtractor(config)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", extractor)
}

func TestLuaExtractor(*testing.T) {
	luaStr := `
		function extract(s)
		  local ret = {}
			ret["haha"] = string.sub(s, 2, -2)
			return ret
    end`

	le, err := newLuaExtractor(luaStr)
	if err != nil {
		panic(err)
	}
	defer le.Close()

	ret, err := le.Extract("12345678")
	if err != nil || ret["haha"] != "234567" {
		panic(err)
	}
	fmt.Printf("%#v\n", ret)

	ret, err = le.Extract("87654321")
	if err != nil || ret["haha"] != "765432" {
		panic(err)
	}
	fmt.Printf("%#v\n", ret)
}

func TestLua(*testing.T) {
	lua_str := `
		function conv(s)
		  local ret = {}
			ret["haha"] = string.sub(s, 2, -2)
			return ret
    end`

	test_str := `12345678`

	L := lua.NewState()
	defer L.Close()
	if err := L.DoString(lua_str); err != nil {
		panic(err)
	}
	if err := L.CallByParam(lua.P{
		Fn:      L.GetGlobal("conv"),
		NRet:    1,
		Protect: true,
	}, lua.LString(test_str)); err != nil {
		panic(err)
	}
	ret := L.Get(-1).(*lua.LTable)
	L.Pop(1)
	fmt.Println(ret.Type())
	fmt.Println(ret.RawGetString("haha"))
}

func TestTextExtractor(*testing.T) {
	confStr := `{ "Delimiters":" " }`
	config, err := NewTextExtractorConfig([]byte(confStr))
	if err != nil {
		panic(err)
	}
	fmt.Printf("[NewTextExtractorConfig] %#v\n", config)

	fields := []PeckField{
		{Name: "col2", Value: "$2"},
		{Name: "col3", Value: "$3"},
		{Name: "col4", Value: "$4"},
	}

	extractor, err := NewTextExtractor(config, fields)
	if err != nil {
		panic(err)
	}
	fmt.Printf("[NewTextExtractor] %#v\n", extractor)

	content := "this is an text extractor"
	m, err := extractor.Extract(content)
	if err != nil {
		panic(err)
	}
	if m["col2"] != "is" || m["col3"] != "an" || m["col4"] != "text" {
		panic(m)
	}
	fmt.Printf("[Extract] %#v\n", m)
}

func TestJsonExtractor(*testing.T) {
	confStr := `{}`
	config, err := NewJsonExtractorConfig([]byte(confStr))
	if err != nil {
		panic(err)
	}
	fmt.Printf("[NewJsonExtractorConfig] %#v\n", config)

	fields := []PeckField{
		{Name: "k1"},
		{Name: "k2.1"},
		{Name: "k3.2.3"},
	}

	extractor, err := NewJsonExtractor(config, fields)
	if err != nil {
		panic(err)
	}
	fmt.Printf("[NewJsonExtractor] %#v\n", extractor)

	content := `{
		"k1":"v1",
		"k2":{
			"1":"v2"
		},
		"k3":{
			"2":{
				"3":"v3"
			}
		},
		"k4":"v4"
	}`
	m, err := extractor.Extract(content)
	if err != nil {
		panic(err)
	}
	if m["k1"] != "v1" || m["k2.1"] != "v2" || m["k3.2.3"] != "v3" {
		panic(m)
	}
	fmt.Printf("[Extract] %#v\n", m)
}
