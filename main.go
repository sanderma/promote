package main

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
)

var environments = []string{"non-prod", "prod", "cde"}

func main() {
	content := `
	# comment line
	kubernetes/staging/policies/opa/templates/k8srequiredannotations-template.yaml
	terraform/staging/files.tf

	# another comment line
	onzinfiles/staging/lalalala
`

	out, err := exec.Command("git", "diff", "--name-only", "master").Output()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(out))

	// Regex pattern captures "key: value" pair from the content.
	//pattern := regexp.MustCompile(`(^kubernetes\/|^terraform\/)(staging)`)
	pattern := regexp.MustCompile(`(kubernetes/|terraform/)staging`)

	for _, env := range environments {
		fmt.Println(pattern.ReplaceAllString(content, "${1}"+env))
	}
}
