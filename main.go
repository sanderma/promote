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
	var cmd *exec.Cmd

	out, err := exec.Command("git", "diff", "--name-only", "master").Output()
	if err != nil {
		log.Fatal(err)
	}

	pattern := regexp.MustCompile(`(kubernetes/|terraform/)staging(/.*)`)

	candidates := pattern.FindAll(out, -1)

	// currentBranch, err := exec.Command("git", "branch", "--show-current").Output()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	for _, env := range environments {

		log.Println("Switch to " + env)
		cmd = exec.Command("git", "switch", "-C", env)
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}

		for _, candidate := range candidates {

			source := string(candidate)
			dest := pattern.ReplaceAllString(string(candidate), "${1}"+env+"${2}")

			fmt.Println("cp " + source + " " + dest)
			copyFile(source, dest)
		}

		log.Println("add kubernetes/" + env)
		cmd = exec.Command("git", "add", "kubernetes/"+env)
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}

		log.Println("add terraform/" + env)
		cmd = exec.Command("git", "add", "terraform/"+env)
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}

		log.Println("commit promote to " + env)
		cmd = exec.Command("git", "commit", "-m", "promote to "+env)
		if err := cmd.Run(); err != nil {
			log.Println("Nothing to commit...")
		}

		cmd = exec.Command("git", "checkout", "dev") //TODO dynamic
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}

		log.Println("Done " + env)
	}
}

func copyFile(s string, d string) {
	var new *os.File
	// Open original file
	original, err := os.Open(s)
	if err != nil {
		log.Fatal(err)
	}
	defer original.Close()

	// Create new file
	new, err = os.Create(d)
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
