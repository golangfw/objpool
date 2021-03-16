package objpool

import (
	"golang.org/x/sync/singleflight"
	"sync"
)

type pool struct {
	sf       singleflight.Group
	valueMap sync.Map
	New      func(string) (interface{}, error)
}

func NewPool(newFunc func(string) (interface{}, error)) *pool {
	p := &pool{
		New: newFunc,
	}
	return p
}

func (p *pool) Get(name string) (v interface{}, err error) {
	if val, ok := p.valueMap.Load(name); ok {
		return val, err
	}
	val, err, _ := p.sf.Do(name, func() (interface{}, error) {
		if val, ok := p.valueMap.Load(name); ok {
			return val, err
		}
		v, err = p.New(name)
		if err != nil {
			return nil, err
		}
		p.valueMap.Store(name, v)
		return v, err
	})
	return val, err
}
