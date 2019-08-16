package main

import (
	"fmt"
	"log"
	"strings"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	cli "gopkg.in/urfave/cli.v1"
)

func reverse(commits []*object.Commit) []*object.Commit {
	left := 0
	right := len(commits) - 1
	for left < right {
		commits[left], commits[right] = commits[right], commits[left]
		left++
		right--
	}
	return commits
}

func diffCommits(sourceBranch string, targetBranch string) []*object.Commit {
	r, err := git.PlainOpenWithOptions(".", &git.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil {
		log.Fatal(err)
	}
	ref, err := r.Reference(plumbing.NewBranchReferenceName(targetBranch), true)
	if err != nil {
		log.Fatal(err)
	}
	targetIter, err := r.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		log.Fatal(err)
	}
	picked := make(map[string]bool)
	targetIter.ForEach(func(c *object.Commit) error {
		picked[hashCommit(c)] = true
		return nil
	})
	ref, err = r.Reference(plumbing.NewBranchReferenceName(sourceBranch), true)
	if err != nil {
		log.Fatal(err)
	}
	sourceIter, err := r.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		log.Fatal(err)
	}
	commits := []*object.Commit{}
	sourceIter.ForEach(func(c *object.Commit) error {
		if !picked[hashCommit(c)] {
			commits = append(commits, c)
		}
		return nil
	})
	return reverse(commits)
}

func diff(a *cli.Context) {
	if a.NArg() < 2 {
		log.Fatal("argument not enough")
	}
	sourceBranch := a.Args().Get(0)
	targetBranch := a.Args().Get(1)
	for _, c := range diffCommits(sourceBranch, targetBranch) {
		if !strings.Contains(c.Author.String(), a.String("author")) {
			continue
		}
		issues := strings.Split(a.String("issue"), ",")
		find := false
		for _, issue := range issues {
			if strings.Contains(c.Message, issue) {
				find = true
				break
			}
		}
		if !find {
			continue
		}
		message := strings.Replace(c.Message, "\n", " ", -1)
		fmt.Printf("%s\t%s\t%s\n", c.Hash, c.Author.Name, message)
	}
}
