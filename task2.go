package main

import (
	"fmt"
)
type str string 
type A struct{
	Name string
	Str str
}
func main() {
	s := "hey"
 d := A{
	Name: "vivek",
	// Str: s, won't work; variables, go requires explicit type conversion
	Str : str(s),
	// Str: "a", // will work; direct literal initialization - implicit conversion of string literal to named string types allowed
 }
 fmt.Printf("%+v", d)
}