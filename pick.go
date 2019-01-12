package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	cli "gopkg.in/urfave/cli.v1"
)

const (
	red   = 31
	green = 32
)

func color(code int, msg string) string {
	return fmt.Sprintf("\x1b[%dm\x1b[49m%s\x1b[39m\x1b[49m", code, msg)
}

func pick(a *cli.Context) {
	msg := "Cherry-Pick %s %s\n"
	if a.NArg() < 2 {
		log.Fatal("argument not enough")
	}
	sourceBranch := a.Args().Get(0)
	targetBranch := a.Args().Get(1)
	out, err := exec.Command("git", "checkout", targetBranch).CombinedOutput()
	if err != nil {
		log.Println(string(out))
		return
	}
	out, err = exec.Command("git", "symbolic-ref", "--short", "HEAD").CombinedOutput()
	if err != nil || !strings.Contains(string(out), targetBranch) {
		fmt.Println(string(out))
		fmt.Println(targetBranch)
		log.Println("Git checkout targetBranch failed")
		return
	}
	for _, c := range diffCommits(sourceBranch, targetBranch) {
		if !strings.Contains(c.Author.String(), a.String("author")) {
			continue
		}
		if !strings.Contains(c.Message, a.String("issue")) {
			continue
		}
		log.Printf("Cherry-Pick %s\n", c.Hash.String())
		out, err := exec.Command("git", "cherry-pick", "-x", c.Hash.String()).CombinedOutput()
		if err != nil {
			log.Printf(msg, c.Hash.String(), color(red, "ERR"))
			log.Println(string(out))
			break
		}
		out, err = exec.Command("git", "log", "--grep", c.Hash.String()).CombinedOutput()
		if err != nil || !strings.Contains(string(out), c.Hash.String()) {
			log.Printf(msg, c.Hash.String(), color(red, "ERR"))
			log.Println(string(out))
			break
		}
		log.Printf(msg, c.Hash.String(), color(green, "OK"))
		if a.Bool("step") {
			break
		}
	}
}
