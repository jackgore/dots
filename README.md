# dots

## Introduction
The dots package creates a configuration object that allows you to parse yaml files and access them using dot syntax in Go. 
It is dependent on the Yaml parsing ability of the [go-yaml](https://github.com/go-yaml/yaml) library.

## Example
Given the following `example.yml`:


    ---
    a:
      b:
        c: "Hello, World"
        d: 600



In the `main.go` file we access fields like so:

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

        s, _ := y.GetString("a.b.c")
        i, _ := y.GetInt("a.b.d")

        fmt.Printf("%v\n", s)
        fmt.Printf("%v\n", i)
    }
