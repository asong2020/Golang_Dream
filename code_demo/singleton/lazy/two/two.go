package two

import (
	"sync"
)

type singleton struct {

}

var instance *singleton
var lock sync.Mutex

func GetInstance() *singleton {
	lock.Lock()
	defer lock.Unlock()
	if instance == nil{
		instance = new(singleton)
	}
	return instance
}
