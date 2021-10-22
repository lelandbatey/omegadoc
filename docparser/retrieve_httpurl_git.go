package docparser

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	git "github.com/go-git/go-git/v5"
	log "github.com/sirupsen/logrus"
)

// retrieve_httpurl_git holds the code for guessing the URL to view a file in
// the browser. It only works when a provided file is inside a Git repo, and
// only when that Git repo has a remote pointing to github or gitlab. In such a
// case, it will try to find the "current commit" of that repo and generate a
// hard-link to that file in that repo at that commit. This is based on the
// implementation of the 'GitLink' command provided by the gitlink-vim
// extension for the Vim text editor:
//     https://github.com/iautom8things/gitlink-vim

type gitURLFinder struct {
	// holds all the paths we've checked before to see if they contain .git
	// folders
	checkedpaths map[string]bool
}

func newGitURLFinder() gitURLFinder {
	return gitURLFinder{
		checkedpaths: map[string]bool{},
	}
}

func (guf *gitURLFinder) GetURL(filepath string, lineno int) (string, error) {
	var pth string = filepath
	var repourl string = ""
	var hash string = ""
	var gitfilepath string = ""
	for {
		if pth == "/" || pth == "" {
			break
		}

		log.Infof("checking path to see if it's a git repo: %s", pth)

		pth, _ = path.Split(pth)
		pth = strings.TrimSuffix(pth, "/")
		lookp := path.Join(pth, ".git")
		if _, ok := guf.checkedpaths[pth]; !ok {
			_, err := os.Stat(lookp)
			if err != nil && !errors.Is(err, os.ErrNotExist) {
				return "", fmt.Errorf("cannot inspect folder %q, err: %w", lookp, err)
			}
			if err == nil {
				guf.checkedpaths[pth] = true
			} else {
				guf.checkedpaths[pth] = false
			}
		}

		isgit := guf.checkedpaths[pth]
		log.Infof("is %q a git repo?: %v", pth, isgit)
		if !isgit {
			continue
		}

		gitfilepath = strings.TrimPrefix(filepath, pth)

		r, err := git.PlainOpen(lookp)
		if err != nil {
			return "", fmt.Errorf("cannot PlainOpen git repo on disk at path %q, error: %w", lookp, err)
		}
		ref, err := r.Head()
		if err != nil {
			return "", fmt.Errorf("cannot access HEAD of repository at path %q, error: %w", lookp, err)
		}
		hash = ref.Hash().String()

		remotes, err := r.Remotes()
		if err != nil {
			return "", fmt.Errorf("cannot get Remotes of repository at path %q, error: %w", lookp, err)
		}
		for _, rem := range remotes {
			if rem.Config().Name != "origin" {
				continue
			}
			remurl := rem.Config().URLs[0]
			remurl = strings.TrimSuffix(remurl, ".git")
			if strings.HasPrefix(remurl, "https://") {
				repourl = remurl
			} else if strings.HasPrefix(remurl, "git@") {
				repourl = strings.Replace(remurl, ":", "/", 1)
				repourl = strings.Replace(repourl, "git@", "https://", 1)
			} else if strings.HasPrefix(remurl, "ssh://") {
				repourl = strings.Replace(remurl, "ssh://", "https://", 1)
			} else if strings.HasPrefix(remurl, "git:") {
				repourl = strings.Replace(remurl, "git:", "https://", 1)
			}
		}
		break
	}
	if repourl != "" && hash != "" && gitfilepath != "" {
		// have to add 1 to the line numbers because when displaying code you
		// start from line 1, not line 0
		return fmt.Sprintf("%s/tree/%s%s#L%d", repourl, hash, gitfilepath, lineno+1), nil
	}
	return "", nil
}
