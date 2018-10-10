package main

import (
	"fmt"
	"math"
)

type Animal interface {
	Speak() string
}

type Dog struct {
}

func (d Dog) Speak() string {
	return "Woof!"
}

type Cat struct {
}

func (c *Cat) Speak() string {
	return "Meow!"
}

func (c Cat) Fly() string {
	return "Flying!"
}

type Llama struct {
}

func (l Llama) Speak() string {
	return "?????"
}

type JavaProgrammer struct {
}

func (j JavaProgrammer) Speak() string {
	return "Design patterns!"
}

type Vertex struct {
	X, Y float64
}

func (v Vertex) Abs() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (v *Vertex) Scale(f float64) {
	v.X = v.X * f
	v.Y = v.Y * f
}

func Abs(v Vertex) float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func Scale(v *Vertex, f float64) {
	v.X = v.X * f
	v.Y = v.Y * f
}

func main() {
	v := Vertex{3, 4}
	Scale(&v, 10)
	fmt.Println(Abs(v))
	v.Scale(3)

	c := Cat{}
	fmt.Println(
		c.Fly(),
		(&c).Fly(),

		c.Speak(),
		(&c).Speak(),
	)

	//animals := []Animal{Dog{}, Cat{}, Llama{}, JavaProgrammer{}}
	//for _, animal := range animals {
	//	fmt.Println(animal.Speak())
	//}

	var p = 12
	var q = 12
	fmt.Printf("%p %p\n", &p, &q)

	a := Animal(Dog{})
	fmt.Println(a.Speak())
}
