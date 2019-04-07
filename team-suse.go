// get_my_prs shows all the pull requests where either you
// or Team SUSE was requestet to review in saltstack/salt.

package main

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"syscall"

	"github.com/google/go-github/github"
	libseccomp "github.com/seccomp/libseccomp-golang"
	"golang.org/x/oauth2"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	if runtime.GOOS == "Linux" {
		var syscalls = []string{"futex", "epoll_pwait", "nanosleep", "read",
			"write", "openat", "epoll_ctl", "close", "rt_sigaction", "mmap",
			"sched_yield", "lstat", "fstat", "mprotect", "rt_sigprocmask",
			"connect", "munmap", "sigaltstack", "set_robust_list", "clone",
			"setsockopt", "socket", "getsockname", "gettid", "getpeername",
			"fcntl", "readlinkat", "getrandom", "newfstatat", "getsockopt",
			"epoll_create1", "brk", "access", "execve", "arch_prctl",
			"sched_getaffinity", "getdents64", "set_tid_address", "prlimit64",
			"exit_group"}
		whiteList(syscalls)
	}

	const (
		owner      = "saltstack"
		repo       = "salt"
		repoID     = 1390248
		teamSuseID = 2582043

		OKBLUE    = "\033[94m"
		OKGREEN   = "\033[92m"
		WARNING   = "\033[93m"
		FAIL      = "\033[91m"
		ENDC      = "\033[0m"
		BOLD      = "\033[1m"
		UNDERLINE = "\033[4m"

		userIcon = "👤"
		teamIcon = "👥"
	)

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		fmt.Fprintln(os.Stderr, "You have to set the env. variable GITHUB_TOKEN.")
		os.Exit(1)
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// Get the authenticated user.
	currentUser, _, err := client.Users.Get(ctx, "")
	checkErr(err)

	// Getting all the PRs for saltstack/salt.
	page := 1 // counter for the 'page' we are fetching from GitHub
	lastPage := -1
	for {
		// Seems like we have to take pagination into account. That's why
		// we loop here and try to get a second page.
		opts := &github.PullRequestListOptions{
			State: "open",
			ListOptions: github.ListOptions{
				// Unfortunately 100 is the max. Pagination is needed.
				PerPage: 100,
				Page:    page,
			},
		}
		prs, response, err := client.PullRequests.List(ctx, owner, repo, opts)
		if err != nil {
			checkErr(err)
		}
		if lastPage == -1 {
			lastPage = response.LastPage
		}

		for _, pr := range prs {
			reviewers, _, err := client.PullRequests.ListReviewers(ctx, owner, repo, pr.GetNumber(), nil)
			checkErr(err)

			needsUserReview, needsTeamReview := false, false

			for _, user := range reviewers.Users {
				if currentUser.GetLogin() == user.GetLogin() {
					needsUserReview = true
				}
			}
			for _, team := range reviewers.Teams {
				if team.GetID() == teamSuseID {
					needsTeamReview = true
				}
			}

			// Print PR if either user or team was requested.
			if needsTeamReview || needsUserReview {
				if needsUserReview {
					fmt.Printf("%s ", userIcon)
				} else if needsTeamReview {
					fmt.Printf("%s ", teamIcon)
				} else {
					fmt.Print(" ")
				}
				fmt.Printf(" %s%s%s\n", BOLD, pr.GetTitle(), ENDC)
				fmt.Printf("🔗  %s%s%s\n\n", OKBLUE, pr.GetHTMLURL(), ENDC)
			}
		}
		if page == lastPage {
			break
		}
		page = page + 1 // incrementing the page counter
	}
}

// Load the seccomp whitelist.
func whiteList(syscalls []string) {

	filter, err := libseccomp.NewFilter(
		libseccomp.ActErrno.SetReturnCode(int16(syscall.EPERM)))
	if err != nil {
		fmt.Printf("Error creating filter: %s\n", err)
	}
	for _, element := range syscalls {
		// fmt.Printf("[+] Whitelisting: %s\n", element)
		syscallID, err := libseccomp.GetSyscallFromName(element)
		if err != nil {
			panic(err)
		}
		filter.AddRule(syscallID, libseccomp.ActAllow)
	}
	filter.Load()
}
