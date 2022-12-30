package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	command := flag.String("command", "", "The git command")
	ignoreErros := flag.Bool(
		"ignore-errors",
		false,
		"Keep running after error if true")
	flag.Parse()

	// Get managed repos from environment variables
	root := os.Getenv("MG_ROOT")
	if root[len(root)-1] != '/' {
		root += "/"
	}

	repo_names := strings.Split(os.Getenv("MG_REPOS"), ",")
	var repos []string
	// Verify all repos exist and are actually git repos (have .git sub-dir)
	for _, r := range repo_names {
		path := root + r
		_, err := os.Stat(path + "/.git")
		if err != nil {
			log.Fatal(err)
		}
		repos = append(repos, path)
	}

	// Break the git command into components (needed to execute)
	var git_components []string
	for _, component := range strings.Split(*command, " ") {
		git_components = append(git_components, component)
	}
	command_string := "git " + *command

	for _, r := range repos {
		// Go to the repo's directory
		os.Chdir(r)

		// Print the command
		fmt.Printf("[%s] %s\n", r, command_string)

		// Execute the command
		out, err := exec.Command("git", git_components...).CombinedOutput()

		// Print the result
		fmt.Println(string(out))

		// Bail out if there was an error and NOT ignoring errors
		if err != nil && !*ignoreErros {
			os.Exit(1)
		}
	}

	fmt.Println("Done.")
}
