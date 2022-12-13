package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
)

var environments = []string{"non-prod", "prod", "cde"}
var currentBranch string

func main() {
	var cmd *exec.Cmd

	out, err := exec.Command("git", "diff", "--name-only", "master").Output()
	if err != nil {
		log.Fatal(err)
	}

	pattern := regexp.MustCompile(`(kubernetes/|terraform/)staging(/.*)`)

	candidates := pattern.FindAll(out, -1)

	b, err := exec.Command("git", "branch", "--show-current").Output()
	if err != nil {
		log.Fatal(err)
	}

	currentBranch = string(bytes.TrimSuffix(b, []byte("\n")))

	log.Println("Currentbranch : " + currentBranch)

	for _, env := range environments {

		promoteBranchName := currentBranch + "-" + env

		log.Println("Switch to " + promoteBranchName)
		cmd = exec.Command("git", "switch", "-C", promoteBranchName)
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}

		for _, candidate := range candidates {

			source := string(candidate)
			dest := pattern.ReplaceAllString(string(candidate), "${1}"+env+"${2}")

			fmt.Println("cp " + source + " " + dest)
			copyFile(source, dest)
		}

		gitPromote(env, currentBranch)

		log.Println("Done " + env)
	}
}

func gitPromote(env string, sourceBranch string) {
	var cmd *exec.Cmd

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

	cmd = exec.Command("git", "checkout", currentBranch)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
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

	//ensure the folder structure exists
	ensureDir(d)

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

func ensureDir(fileName string) {
	dirName := filepath.Dir(fileName)
	if _, serr := os.Stat(dirName); serr != nil {
		merr := os.MkdirAll(dirName, os.ModePerm)
		if merr != nil {
			panic(merr)
		}
	}
}
