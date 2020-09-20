package one

type singleton struct {

}

var instance *singleton

func init()  {
	instance = new(singleton)
}

func GetInstance()  *singleton{
	return instance
}
