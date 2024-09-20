package main

import (
	"fmt"
)

/*
type animal interface {
	sleep()
	eat()
}

type dog struct {
	name string
	food string
}

func (mydog dog) sleep() {
	fmt.Printf("%s is sleeping\n", mydog.name)
}

func (mydog dog) eat() {
	fmt.Printf("%s is eating %s\n", mydog.name, mydog.food)
}

type cat struct {
	name string
	food string
}

func (mycat cat) sleep() {
	fmt.Printf("%s is sleeping \n", mycat.name)
}

*/

/*
//错误示例
func main() {
	var ani animal
	kitty := cat{
		name: "kitty",
		food: "fish",
	}

	ani = kitty

	ani.sleep()
}
*/

/*
//正确示例

	func main() {
		xiaobai := dog{
		name: "xiaobai",
		food: "bone",
		}
		xiaobai.sleep()
		xiaobai.eat()
		}
*/
type empty_interface interface {
}

func example(empty empty_interface) {
	fmt.Printf("example", empty)
}

func main() {
	var data []interface{}
	data = append(data, 42)
	data = append(data, "hello")

	fmt.Println("data...........", data)
}
