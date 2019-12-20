package logpeck

import (
	"encoding/json"
	"errors"

	log "github.com/Sirupsen/logrus"
	lua "github.com/yuin/gopher-lua"
	luajson "layeh.com/gopher-json"
)

// LuaExtractorConfig .
type LuaExtractorConfig struct {
	LuaString string
	Fields    []PeckField
}

// LuaExtractor .
type LuaExtractor struct {
	state  *lua.LState
	fields map[string]bool
}

var luaExtractorFuncName = "extract"

// NewLuaExtractorConfig .
func NewLuaExtractorConfig(configStr []byte) (LuaExtractorConfig, error) {
	c := LuaExtractorConfig{}
	err := json.Unmarshal(configStr, &c)
	if err != nil {
		return c, err
	}
	return c, nil
}

// NewLuaExtractor .
func NewLuaExtractor(config interface{}) (LuaExtractor, error) {
	c, ok := config.(LuaExtractorConfig)
	if !ok {
		return LuaExtractor{}, errors.New("LuaExtractor config error")
	}

	l := LuaExtractor{
		state:  lua.NewState(),
		fields: make(map[string]bool),
	}
	c.LuaString = "json = require(\"luajson.json\") " + c.LuaString
	l.state.PreloadModule("luajson.json", luajson.Loader)
	if err := l.state.DoString(c.LuaString); err != nil {
		return l, err
	}
	for _, f := range c.Fields {
		l.fields[f.Name] = true
	}
	log.Infof("[LuoExtractor] Init extractor finished %#v", l)
	return l, nil
}

// Extract .
func (le LuaExtractor) Extract(content string) (map[string]interface{}, error) {
	param := lua.P{
		Fn:      le.state.GetGlobal(luaExtractorFuncName),
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
	enable := true
	key := ""
	lT.ForEach(func(k, v lua.LValue) {
		if _, ok := le.fields[k.String()]; !ok {
			enable = false
			key = k.String()
		}
		ret[k.String()] = v.String()
	})
	if !enable {
		return map[string]interface{}{}, errors.New(key + " is not in Fields")
	}
	return ret, nil
}

// Close .
func (le LuaExtractor) Close() {
	le.state.Close()
}
