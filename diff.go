package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

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

func filterCommits(a *cli.Context) []*object.Commit {
	sourceBranch := a.Args().Get(0)
	targetBranch := a.Args().Get(1)
	commits := []*object.Commit{}
	for _, c := range diffCommits(sourceBranch, targetBranch) {
		if !strings.Contains(c.Author.String(), a.String("author")) {
			continue
		}
		issues := strings.Split(a.String("issue"), ",")
		found := false
		for _, issue := range issues {
			if strings.Contains(c.Message, issue) {
				found = true
				break
			}
		}
		if !found {
			continue
		}
		commits = append(commits, c)
	}
	return commits
}

func diff(a *cli.Context) {
	if a.NArg() < 2 {
		log.Fatal("argument not enough")
	}
	commits := filterCommits(a)
	loc, _ := time.LoadLocation("Asia/Shanghai")
	csvFile := a.String("csv")
	if len(csvFile) > 0 {
		f, err := os.Create(csvFile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		f.WriteString("\xEF\xBB\xBF")
		w := csv.NewWriter(f)
		for _, c := range commits {
			message := strings.Replace(c.Message, "\n", " ", -1)
			when := c.Author.When.In(loc)
			w.Write([]string{c.Author.Name, fmt.Sprint(c.Hash), message, fmt.Sprint(when)})
		}
		w.Flush()
	} else {
		for _, c := range commits {
			message := strings.Replace(c.Message, "\n", " ", -1)
			when := c.Author.When.In(loc)
			fmt.Printf("%s\t%s\t%s\t%s\n", c.Hash, when, c.Author.Name, message)
		}
	}
}
