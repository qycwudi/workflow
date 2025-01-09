package rulego

import (
	"errors"
	"github.com/rulego/rulego/api/types"
	"sync"
)

type RolePool struct {
	Pool sync.Map
}

func Init() *RolePool {
	return new(RolePool)
}

func (p *RolePool) Get(key string) (types.RuleEngine, error) {
	if v, ok := p.Pool.Load(key); ok {
		return v.(types.RuleEngine), nil
	}
	return nil, errors.New("not exist")
}

func (p *RolePool) Put(key string, rule *types.RuleEngine) {
	p.Pool.Store(key, rule)
}
