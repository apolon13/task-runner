package git

import (
	"fmt"
	"github.com/ldez/go-git-cmd-wrapper/v2/checkout"
	"github.com/ldez/go-git-cmd-wrapper/v2/commit"
	"github.com/ldez/go-git-cmd-wrapper/v2/git"
	"github.com/ldez/go-git-cmd-wrapper/v2/merge"
	"github.com/ldez/go-git-cmd-wrapper/v2/pull"
	"github.com/ldez/go-git-cmd-wrapper/v2/push"
	"github.com/ldez/go-git-cmd-wrapper/v2/rebase"
	"github.com/ldez/go-git-cmd-wrapper/v2/reset"
	"github.com/ldez/go-git-cmd-wrapper/v2/revparse"
	"strings"
	"task-runner/cmd"
	"task-runner/config"
)

func handleGitCommand(out string, err error) {
	fmt.Println(fmt.Sprintf("Git: %s", out))
	if err != nil {
		panic(err)
	}
}

func Deploy(deploy *config.Branch, testStand string) {
	handleGitCommand(git.Pull(pull.Repository("origin"), pull.Repository("master")))
	handleGitCommand(git.Checkout(checkout.Branch(testStand), git.Debug))
	handleGitCommand(git.Reset(reset.Hard, reset.Path("origin", "master")))
	handleGitCommand(git.Checkout(checkout.Branch(deploy.Name), git.Debug))
	handleGitCommand(git.Rebase(rebase.Branch(testStand), git.Debug))
	handleGitCommand(git.Checkout(checkout.Branch(testStand), git.Debug))
	handleGitCommand(git.Merge(merge.Commits(deploy.Name), git.Debug))
	handleGitCommand(git.Push(push.Force, push.Remote("origin"), push.Remote(testStand), git.Debug))
	handleGitCommand(git.Checkout(checkout.Branch(deploy.Name), git.Debug))
}

func Release(release *config.Branch, cnf config.Yaml) {
	if len(cnf.Git.Release.Intermediate) > 0 {
		var last *config.Branch
		for _, intermediate := range cnf.Git.Release.Intermediate {
			branchB := release
			if last != nil {
				branchB = last
			}
			if intermediate.Name != branchB.Name {
				doRelease(strings.Trim(intermediate.Name, "\n"), branchB.Name, func() error {
					if intermediate.Command.Main == "" {
						return nil
					}
					return cmd.Handle(intermediate.Command)
				}, intermediate.Amend)
				last = &intermediate
			}
		}
		if last != nil {
			doRelease("master", last.Name, func() error {
				return nil
			}, false)
			return
		}
	}
	if release.Name == "master" {
		handleGitCommand(pushTo(release.Name))
		return
	}
	doRelease("master", release.Name, func() error {
		return nil
	}, false)
}

func pushTo(branch string) (string, error) {
	return git.Push(push.Remote("origin"), push.Remote(branch), git.Debug)
}

func doRelease(branchA string, branchB string, command func() error, amend bool) {
	handleGitCommand(git.Checkout(checkout.Branch(branchA), git.Debug))
	handleGitCommand(git.Pull(pull.Repository("origin"), pull.Repository(branchA)))
	handleGitCommand(git.Checkout(checkout.Branch(branchB), git.Debug))
	err := command()
	if err != nil {
		panic(err)
	}
	if amend == true {
		handleGitCommand(git.Commit(commit.Amend, commit.NoEdit, git.Debug))
	}
	handleGitCommand(git.Rebase(rebase.Branch(branchA), git.Debug))
	handleGitCommand(git.Checkout(checkout.Branch(branchA), git.Debug))
	handleGitCommand(git.Merge(merge.Commits(branchB), git.Debug))
	handleGitCommand(pushTo(branchA))
}

func CurrentBranch() (string, error) {
	return git.RevParse(revparse.AbbrevRef(""), revparse.Args("HEAD"))
}
