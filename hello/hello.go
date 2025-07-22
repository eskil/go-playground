package main

import (
	"fmt"
	"log"

	"example.com/greetings"
)

func main() {
	log.SetPrefix("greetings: ")
	log.SetFlags(0)

	names := []string{"Gladys", "Samantha", "Darrin"}
	msgs, err := greetings.Hellos(names)
	if err != nil {
		log.Fatal(err)
	}

	for _, msg := range msgs {
		fmt.Println(msg)
	}
}
