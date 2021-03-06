// team-suse shows all the pull requests where either the user from the used
// Github token or Team SUSE was requestet to review in saltstack/salt.

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func init() {
	// The syscall restriciton is only available for Linux right now via
	// seccomp.
	applySyscallRestrictions()
}

func main() {
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
