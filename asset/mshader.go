package asset

import (
	"korok.io/korok/gfx/bk"

	"log"
	"errors"
)

type ShaderManager struct {
	repo map[string]idCount
}

func NewShaderManager() *ShaderManager {
	return &ShaderManager{make(map[string]idCount)}
}

func (sm *ShaderManager) LoadDefaultShader() {
	sm.LoadShader("dft", vertex, color)
	sm.LoadShader("batch", bVertex, bColor)
}

// 引用计数 +1
func (sm *ShaderManager) LoadShader(name string, vertex, color string) {
	var rid, cnt uint16
	if v, ok := sm.repo[name]; ok {
		cnt = v.cnt
	} else {
		id, err := sm.load(vertex, color)
		if err != nil {
			log.Println(err)
		}
		rid = id
	}
	sm.repo[name] = idCount{rid, cnt + 1}
}

// 引用计数 -1
func (sm *ShaderManager) Unload(name string) {
	if v, ok := sm.repo[name]; ok {
		if v.cnt > 1 {
			sm.repo[name] = idCount{v.rid, v.cnt - 1}
		} else {
			delete(sm.repo, name)
			bk.R.Free(v.rid)
		}
	}
}

func (sm *ShaderManager) GetShader(key string) (uint16, *bk.Shader) {
	if ref, ok := sm.repo[key]; ok {
		if ok, sh := bk.R.Shader(ref.rid); ok {
			return ref.rid, sh
		}
	}
	return bk.InvalidId, nil
}

func (sm *ShaderManager) GetShaderStr(key string) (string, string) {
	switch key {
	case "dft", "mesh":
		return vertex, color
	case "batch":
		return bVertex, bColor
	}
	return "", ""
}

func (sm *ShaderManager) load(vertex, fragment string) (uint16, error){
	if id, _ :=  bk.R.AllocShader(vertex, fragment); id != bk.InvalidId {
		return id, nil
	}
	return bk.InvalidId, errors.New("fail to load shader")
}
