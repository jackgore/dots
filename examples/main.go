package main

import (
	"fmt"

	"github.com/JonathonGore/dots"
)

func main() {
	y, err = dots.NewYaml("example.yml")
	if err != nil {
		fmt.Fatalf(err)
	}

	fmt.Println(y.GetString("a.b.c")
	fmt.Println(y.GetInt("a.b.d")
}
