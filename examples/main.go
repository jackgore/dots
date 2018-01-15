package main

import (
	"fmt"
	"log"

	"github.com/JonathonGore/dots/yaml"
)

func main() {
	y, err := yaml.New("example.yml")
	if err != nil {
		log.Fatal(err)
	}

	s := y.GetString("a.b.c")
	i := y.GetInt("a.b.d")

	fmt.Printf("%v\n", s)
	fmt.Printf("%v\n", i)
}
