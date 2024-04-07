package main

import "fmt"

func main() {
	fmt.Println(Hello("sehyeong", ""))
}

func Hello(name string, language string) string {
	if name == "" {
		name = "world"
	}

	if language == "Spanish" {
		return "Hola, " + name
	}
	return "Hello, " + name
}
