package git

import (
	"github.com/ldez/go-git-cmd-wrapper/v2/checkout"
	"github.com/ldez/go-git-cmd-wrapper/v2/commit"
	"github.com/ldez/go-git-cmd-wrapper/v2/git"
	"github.com/ldez/go-git-cmd-wrapper/v2/merge"
	"github.com/ldez/go-git-cmd-wrapper/v2/push"
	"github.com/ldez/go-git-cmd-wrapper/v2/rebase"
	"github.com/ldez/go-git-cmd-wrapper/v2/revparse"
	"os/exec"
	"task-runner/config"
)

func Release(release *config.Branch, cnf config.Yaml) error {
	if len(cnf.Git.Intermediate) > 0 {
		var last *config.Branch
		for _, intermediate := range cnf.Git.Intermediate {
			branchB := release
			if last != nil {
				branchB = last
			}
			err := doRelease(intermediate.Name, branchB.Name, func() error {
				return exec.Command(intermediate.Command.Main, intermediate.Command.Args...).Run()
			}, intermediate.Amend)
			if err != nil {
				return err
			}
			last = &intermediate
		}
	}
	if release.Name == "master" {
		pushTo(release.Name)
	}
	return nil
}

func pushTo(branch string) {
	git.Push(push.Remote("origin"), push.Repo(branch))
}

func doRelease(branchA string, branchB string, command func() error, amend bool) error {
	git.Checkout(checkout.Branch(branchA))
	git.Pull()
	git.Checkout(checkout.Branch(branchB))
	if err := command(); err != nil {
		return err
	}
	if amend == true {
		git.Commit(commit.Amend, commit.NoEdit)
	}
	git.Rebase(rebase.Branch(branchA))
	git.Checkout(checkout.Branch(branchA))
	git.Merge(merge.Commits(branchB))
	pushTo(branchA)
	return nil
}

func CurrentBranch() (string, error) {
	return git.RevParse(revparse.AbbrevRef(""), revparse.Args("HEAD"))
}
