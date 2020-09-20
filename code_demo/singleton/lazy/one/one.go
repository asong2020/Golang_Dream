package one

type singleton struct {

}

var  instance *singleton
func GetInstance() *singleton {
	if instance == nil{
		instance = new(singleton)
	}
	return instance
}
