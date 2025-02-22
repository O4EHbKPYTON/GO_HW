package main

import "fmt"

func main() {
	users := []string{"Tom", "Alice", "Kate"}
	users = append(users, "Bob")
	users = append(users[:3], users[4:]...)
	for _, value := range users {
		fmt.Println(value, 1)
	}
}
