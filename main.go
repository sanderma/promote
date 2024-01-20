package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type TagCommit map[string]string

type Apps map[string]TagCommit

func (a Apps) List() {
	for k := range a {
		fmt.Println(k)
	}
}

func main() {

	deployments := make(Apps)

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	dir, err := os.ReadDir("./deployments")
	if err != nil {
		fmt.Println(err)
	}
	for _, d := range dir {
		if !strings.HasPrefix(d.Name(), ".") {
			deployments[d.Name()] = TagCommit{}
		}
	}

	deployments.List()

	repo, err := git.PlainOpen(path)
	if err != nil {
		log.Println(err)
	}
	tagrefs, err := repo.Tags()
	if err != nil {
		log.Println(err)
	}

	regexEnv := regexp.MustCompile(`\w*$`)

	//repo.Tag()
	err = tagrefs.ForEach(func(t *plumbing.Reference) error {
		n, found := strings.CutPrefix(t.Name().String(), "refs/tags/")

		if found {
			fmt.Println(regexEnv.FindString(n))
		}
		return nil
	})
	if err != nil {
		log.Println(err)
	}
}
