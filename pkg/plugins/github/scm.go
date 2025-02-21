package github

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	git "github.com/updatecli/updatecli/pkg/plugins/git/generic"
)

// Init set default Github parameters if not set.
func (g *Github) Init(source string, pipelineID string) error {
	g.VersionFilter.Pattern = source
	g.remoteBranch = git.SanitizeBranchName(fmt.Sprintf("updatecli_%v", pipelineID))
	g.setDirectory()

	errs := g.Check()
	if len(errs) > 0 {
		for _, e := range errs {
			logrus.Errorf("%s\n", e)
		}
		return fmt.Errorf("wrong github configuration")
	}
	return nil
}

// GetDirectory returns the local git repository path.
func (g *Github) GetDirectory() (directory string) {
	return g.Directory
}

// Clean deletes github working directory.
func (g *Github) Clean() error {
	err := os.RemoveAll(g.Directory)
	if err != nil {
		return err
	}
	return nil
}

// Clone run `git clone`.
func (g *Github) Clone() (string, error) {

	URL := fmt.Sprintf("https://github.com/%v/%v.git",
		g.Owner,
		g.Repository)

	g.setDirectory()

	err := git.Clone(g.Username, g.Token, URL, g.GetDirectory())

	if err != nil {
		return "", err
	}

	if len(g.remoteBranch) > 0 && len(g.GetDirectory()) > 0 {
		err = git.Checkout(g.Branch, g.remoteBranch, g.GetDirectory())
	}

	if err != nil {
		return "", err
	}

	return g.Directory, nil
}

// Commit run `git commit`.
func (g *Github) Commit(message string) error {

	// Generate the conventional commit message
	commitMessage, err := g.CommitMessage.Generate(message)
	if err != nil {
		return err
	}

	err = git.Commit(g.User, g.Email, commitMessage, g.GetDirectory())
	if err != nil {
		return err
	}
	return nil
}

// Checkout create and then uses a temporary git branch.
func (g *Github) Checkout() error {
	err := git.Checkout(g.Branch, g.remoteBranch, g.Directory)
	if err != nil {
		return err
	}
	return nil
}

// Add run `git add`.
func (g *Github) Add(files []string) error {

	err := git.Add(files, g.Directory)
	if err != nil {
		return err
	}
	return nil
}

// Push run `git push` then open a pull request on Github if not already created.
func (g *Github) Push() error {

	err := git.Push(g.Username, g.Token, g.GetDirectory(), g.Force)
	if err != nil {
		return err
	}

	return nil
}

// PushTag push tags
func (g *Github) PushTag(tag string) error {

	err := git.PushTag(tag, g.Username, g.Token, g.GetDirectory(), g.Force)
	if err != nil {
		return err
	}

	return nil
}
