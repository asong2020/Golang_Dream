package third

import (
	"sync"
)

type singleton struct {

}

var instance *singleton
var lock sync.Mutex

func GetInstance() *singleton {
	if instance == nil{
		lock.Lock()
		instance = new(singleton)
		lock.Unlock()
	}
	return instance
}