package main

import "fmt"

var mp = make(map[string]int)

func init() {
	mp["John"] = 25
	mp["Jane"] = 26
	mp["Jim"] = 27
	mp["Jill"] = 28
	mp["Jack"] = 29
}

func main() {
	fmt.Println(GetAge("John"))
	AddPerson("Jill", 40)
	DeletePerson("John")
	PrintAll()
}

func GetAge(name string) (int, bool) {
	age, ok := mp[name]
	return age, ok
}

func AddPerson(name string, age int) {
	mp[name] = age
}

func DeletePerson(name string) {
	delete(mp, name)
}

func PrintAll() {
	for name, age := range mp {
		fmt.Println(name, age)
	}
}
