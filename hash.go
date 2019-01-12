package main

import (
	"crypto/sha256"
	"fmt"

	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

func hashCommit(c *object.Commit) string {
	author := c.Author
	seri := []byte(author.Name + author.Email + author.When.String())
	return string(fmt.Sprintf("%x", sha256.Sum256(seri)))
}
