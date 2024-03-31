package contract

import (
	"fmt"
	"sync"
	"time"

	"github.com/tonindexer/anton/abi"
	"github.com/tonindexer/anton/internal/core"
)

var cacheInvalidation = 10 * time.Second

type cache struct {
	interfaces  []*core.ContractInterface
	getMethods  map[abi.ContractName]map[string]abi.GetMethodDesc
	operations  []*core.ContractOperation
	lastCleared time.Time
	sync.Mutex
}

func newCache() *cache {
	return &cache{
		getMethods:  map[abi.ContractName]map[string]abi.GetMethodDesc{},
		lastCleared: time.Now(),
	}
}

func (c *cache) clearCaches() {
	if time.Since(c.lastCleared) < cacheInvalidation {
		return
	}
	c.interfaces = nil
	c.operations = nil
	c.getMethods = map[abi.ContractName]map[string]abi.GetMethodDesc{}
	c.lastCleared = time.Now()
}

func (c *cache) setInterfaces(interfaces []*core.ContractInterface) {
	c.Lock()
	defer c.Unlock()

	fmt.Println("cache interfaces")

	for _, i := range interfaces {
		if _, ok := c.getMethods[i.Name]; !ok {
			c.getMethods[i.Name] = map[string]abi.GetMethodDesc{}
		}
		for it := range i.GetMethodsDesc {
			c.getMethods[i.Name][i.GetMethodsDesc[it].Name] = i.GetMethodsDesc[it]
		}
	}

	c.interfaces = interfaces
}

func (c *cache) getInterfaces() []*core.ContractInterface {
	c.Lock()
	defer c.Unlock()

	c.clearCaches()

	return c.interfaces
}

func (c *cache) getMethodDesc(name abi.ContractName, method string) (abi.GetMethodDesc, bool) {
	c.Lock()
	defer c.Unlock()

	if _, ok := c.getMethods[name]; !ok {
		return abi.GetMethodDesc{}, false
	}

	d, ok := c.getMethods[name][method]
	return d, ok
}

func (c *cache) setOperations(operations []*core.ContractOperation) {
	c.Lock()
	defer c.Unlock()

	c.operations = operations
}

func (c *cache) getOperations() []*core.ContractOperation {
	c.Lock()
	defer c.Unlock()

	c.clearCaches()

	return c.operations
}

func (c *cache) getOperationByID(types []abi.ContractName, outgoing bool, id uint32) *core.ContractOperation {
	c.Lock()
	defer c.Unlock()

	c.clearCaches()

	for _, op := range c.operations {
		if op.Outgoing != outgoing {
			continue
		}
		if op.OperationID != id {
			continue
		}
		for _, t := range types {
			if op.ContractName == t {
				return op
			}
		}
	}

	return nil
}
