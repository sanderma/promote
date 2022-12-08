package main

import (
	"fmt"
	"io"
	"log"
	"os"
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

	for _, env := range environments {

		for _, candidate := range candidates {

			source := string(candidate)
			dest := pattern.ReplaceAllString(string(candidate), "${1}"+env+"${2}")
			fmt.Println("cp " + source + " " + dest)
			copyFile(source, dest)
		}
		cmd := exec.Command("git", "add", "--all")
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}

		cmd = exec.Command("git", "commit", "-m", "promote to "+env)
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}
	}
}

func copyFile(s string, d string) {
	// Open original file
	original, err := os.Open(s)
	if err != nil {
		log.Fatal(err)
	}
	defer original.Close()

	// Create new file
	new, err := os.Create(d)
	if err != nil {
		log.Fatal(err)
	}
	defer new.Close()

	//This will copy
	_, err = io.Copy(new, original)
	if err != nil {
		log.Fatal(err)
	}
}
