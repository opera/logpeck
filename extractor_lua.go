package logpeck

import (
	"encoding/json"
	"errors"
	log "github.com/Sirupsen/logrus"
	lua "github.com/yuin/gopher-lua"
)

type LuaExtractorConfig struct {
	LuaString string
}

type LuaExtractor struct {
	state *lua.LState
}

var LuaExtractorFuncName string = "extract"

func NewLuaExtractorConfig(configStr []byte) (LuaExtractorConfig, error) {
	c := LuaExtractorConfig{}
	err := json.Unmarshal(configStr, &c)
	if err != nil {
		return c, err
	}
	return c, nil
}

func NewLuaExtractor(config interface{}) (*LuaExtractor, error) {
	c, ok := config.(LuaExtractorConfig)
	if !ok {
		return nil, errors.New("LuaExtractor config error")
	}
	return newLuaExtractor(c.LuaString)
}

func newLuaExtractor(luaStr string) (*LuaExtractor, error) {
	l := &LuaExtractor{
		state: lua.NewState(),
	}
	if err := l.state.DoString(luaStr); err != nil {
		return nil, err
	}
	return l, nil
}

func (le *LuaExtractor) Extract(content string) (map[string]interface{}, error) {
	param := lua.P{
		Fn:      le.state.GetGlobal(LuaExtractorFuncName),
		NRet:    1,
		Protect: true,
	}
	if err := le.state.CallByParam(param, lua.LString(content)); err != nil {
		return nil, err
	}
	lRet := le.state.Get(-1)
	lT, ok := lRet.(*lua.LTable)
	if !ok {
		return nil, errors.New("lua return type error " + lRet.String())
	}
	le.state.Pop(1)
	log.Debugf("[LuaExtractor] %s %#v", content, lT)
	ret := make(map[string]interface{})
	lT.ForEach(func(k, v lua.LValue) {
		ret[k.String()] = v.String()
	})
	return ret, nil
}

func (le *LuaExtractor) Close() {
	le.state.Close()
}
