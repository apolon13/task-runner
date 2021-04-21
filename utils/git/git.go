package git

import (
	"fmt"
	"github.com/ldez/go-git-cmd-wrapper/v2/checkout"
	"github.com/ldez/go-git-cmd-wrapper/v2/commit"
	"github.com/ldez/go-git-cmd-wrapper/v2/git"
	"github.com/ldez/go-git-cmd-wrapper/v2/merge"
	"github.com/ldez/go-git-cmd-wrapper/v2/push"
	"github.com/ldez/go-git-cmd-wrapper/v2/rebase"
	"github.com/ldez/go-git-cmd-wrapper/v2/revparse"
	"strings"
	"task-runner/cmd"
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
			if intermediate.Name != branchB.Name {
				err := doRelease(strings.Trim(intermediate.Name, "\n"), branchB.Name, func() error {
					if intermediate.Command.Main == "" {
						return nil
					}
					return cmd.Handle(intermediate.Command)
				}, intermediate.Amend)
				if err != nil {
					return err
				}
				last = &intermediate
			}
		}
		if last != nil {
			return doRelease("master", last.Name, func() error {
				return nil
			}, false)
		}
	}
	if release.Name == "master" {
		return pushTo(release.Name)
	}
	return doRelease("master", release.Name, func() error {
		return nil
	}, false)
}

func pushTo(branch string) error {
	out, err := git.Push(push.Remote("origin"), push.Remote(branch), git.Debug)
	fmt.Println(out)
	if err != nil {
		return err
	}
	return nil
}

func doRelease(branchA string, branchB string, command func() error, amend bool) error {
	checkoutAOut, checkoutAErr := git.Checkout(checkout.Branch(branchA), git.Debug)

	fmt.Println(checkoutAOut)
	if checkoutAErr != nil {
		return checkoutAErr
	}
	pullOut, pullErr := git.Pull()
	fmt.Println(pullOut)
	if pullErr != nil {
		return pullErr
	}
	checkoutBOut, checkoutBErr := git.Checkout(checkout.Branch(branchB), git.Debug)
	fmt.Println(checkoutBOut)
	if checkoutBErr != nil {
		return checkoutBErr
	}
	if err := command(); err != nil {
		return err
	}
	if amend == true {
		commitOut, commitErr := git.Commit(commit.Amend, commit.NoEdit, git.Debug)
		fmt.Println(commitOut)
		if commitErr != nil {
			return commitErr
		}
	}
	rebaseOut, rebaseErr := git.Rebase(rebase.Branch(branchA), git.Debug)
	fmt.Println(rebaseOut)
	if rebaseErr != nil {
		return rebaseErr
	}
	checkoutAOut, checkoutAErr = git.Checkout(checkout.Branch(branchA), git.Debug)
	fmt.Println(checkoutAOut)
	if checkoutAErr != nil {
		return checkoutAErr
	}
	mergeOut, mergeErr := git.Merge(merge.Commits(branchB), git.Debug)
	fmt.Println(mergeOut)
	if mergeErr != nil {
		return mergeErr
	}
	return pushTo(branchA)
}

func CurrentBranch() (string, error) {
	return git.RevParse(revparse.AbbrevRef(""), revparse.Args("HEAD"))
}
