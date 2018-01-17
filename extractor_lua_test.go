package logpeck

import (
	"fmt"
	lua "github.com/yuin/gopher-lua"
	"testing"
)

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
