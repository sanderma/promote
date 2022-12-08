package main

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
)

var environments = []string{"non-prod", "prod", "cde"}

func main() {

	out, err := exec.Command("git", "diff", "--name-only", "master").Output()

	if err != nil {
		log.Fatal(err)
	}

	pattern := regexp.MustCompile(`(kubernetes/|terraform/)staging(/.*)`)

	candidates := pattern.FindAll(out, -1)

	for _, candidate := range candidates {

		for _, env := range environments {
			fmt.Println("cp " + string(candidate) + " " + pattern.ReplaceAllString(string(candidate), "${1}"+env+"${2}"))
		}
	}
}
