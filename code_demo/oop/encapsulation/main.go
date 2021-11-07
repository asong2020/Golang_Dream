package main

type IdCard string

func NewIdCard(card string) IdCard {
	return IdCard(card)
}

func (i IdCard) GetPlaceOfBirth() string {
	return string(i[:6])
}

func (i IdCard) GetBirthDay() string {
	return string(i[6:14])
}


func main()  {

}
