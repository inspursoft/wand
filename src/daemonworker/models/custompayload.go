package models

import (
	"encoding/json"
	"fmt"
	"strings"
)

type CustomRepo struct {
	RepoName string `json:"repo_name"`
	HTMLURL  string `json:"html_url"`
	CloneURL string `json:"clone_url"`
	Branch   string `json:"branch"`
}

type CustomWebhookPayload struct {
	GroupName    string     `json:"group_name"`
	FullName     string     `json:"full_name"`
	RepoName     string     `json:"repo_name"`
	Username     string     `json:"username"`
	GogsURL      string     `json:"gogs_url"`
	MasterURL    string     `json:"master_url"`
	NodeIP       string     `json:"node_ip"`
	RegistryURL  string     `json:"registry_url"`
	NativeAction string     `json:"native_action"`
	Action       string     `json:"action"`
	BaseRepo     CustomRepo `json:"base_repo"`
	HeadRepo     CustomRepo `json:"head_repo"`
	HTMLURL      string     `json:"html_url"`
	APIURL       string     `json:"api_url"`
	CommentURL   string     `json:"comment_url"`
	Affinity     string     `json:"affinity"`
	LastCoverage string     `json:"last_coverage"`
	Collaborator string     `json:"collaborator"`
	AccessToken  string     `json:"access_token"`
}

func (c *CustomWebhookPayload) AdaptToCustom(event string, payload []byte) (trigger bool, err error) {
	switch event {
	case "push":
		err = c.fromPushPayload(payload)
	case "pull_request":
		err = c.fromPullRequestPayload(payload)
	}
	trigger = c.filterAction()
	return
}

func (c *CustomWebhookPayload) filterAction() bool {
	for _, e := range []string{"push", "opened", "synchronized", "reopened"} {
		if e == c.NativeAction {
			return true
		}
	}
	return false
}

func (c *CustomWebhookPayload) fromPushPayload(payload []byte) (err error) {
	var push PushPayload
	err = json.Unmarshal(payload, &push)
	c.Action = "push"
	c.NativeAction = "push"
	c.FullName = push.Repo.FullName
	c.RepoName = push.Repo.Name
	c.Username = push.Repo.Owner.UserName
	c.BaseRepo.RepoName = push.Repo.Name
	c.BaseRepo.HTMLURL = push.Repo.HTMLURL
	c.BaseRepo.CloneURL = push.Repo.CloneURL
	c.BaseRepo.Branch = strings.Split(push.Ref, "/")[2]
	c.HeadRepo.RepoName = "-"
	c.HeadRepo.HTMLURL = "-"
	c.HeadRepo.CloneURL = "-"
	c.HeadRepo.Branch = "-"
	c.HTMLURL = push.Repo.HTMLURL
	c.CommentURL = "-"
	c.Collaborator = push.Repo.Owner.UserName
	return
}

func (c *CustomWebhookPayload) fromPullRequestPayload(payload []byte) (err error) {
	var pullRequest PullRequestPayload
	err = json.Unmarshal(payload, &pullRequest)
	c.Action = "pull_request"
	c.NativeAction = string(pullRequest.Action)
	c.FullName = pullRequest.Repository.FullName
	c.RepoName = pullRequest.Repository.Name
	c.Username = pullRequest.Repository.Owner.UserName
	c.BaseRepo.RepoName = pullRequest.PullRequest.BaseRepo.Name
	c.BaseRepo.HTMLURL = pullRequest.PullRequest.BaseRepo.HTMLURL
	c.BaseRepo.CloneURL = pullRequest.PullRequest.BaseRepo.CloneURL
	c.BaseRepo.Branch = pullRequest.PullRequest.BaseBranch
	c.HeadRepo.RepoName = pullRequest.PullRequest.HeadRepo.Name
	c.HeadRepo.HTMLURL = pullRequest.PullRequest.HeadRepo.HTMLURL
	c.HeadRepo.CloneURL = pullRequest.PullRequest.HeadRepo.CloneURL
	c.HeadRepo.Branch = pullRequest.PullRequest.HeadBranch
	c.HTMLURL = pullRequest.PullRequest.HTMLURL
	c.CommentURL = fmt.Sprintf("%s/repos/%s/issues/%d/comments", c.APIURL, pullRequest.Repository.FullName, pullRequest.Index)
	c.Collaborator = pullRequest.PullRequest.HeadRepo.Owner.UserName
	return
}
