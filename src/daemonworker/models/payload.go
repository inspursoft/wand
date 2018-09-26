package models

import (
	"encoding/json"
	"errors"
	"strings"
	"time"
)

var (
	ErrInvalidReceiveHook = errors.New("Invalid JSON payload received over webhook")
)

type PayloadUser struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	UserName string `json:"username"`
}

// FIXME: consider use same format as API when commits API are added.
type PayloadCommit struct {
	ID        string       `json:"id"`
	Message   string       `json:"message"`
	URL       string       `json:"url"`
	Author    *PayloadUser `json:"author"`
	Committer *PayloadUser `json:"committer"`

	Added    []string `json:"added"`
	Removed  []string `json:"removed"`
	Modified []string `json:"modified"`

	Timestamp time.Time `json:"timestamp"`
}

// _________                        __
// \_   ___ \_______   ____ _____ _/  |_  ____
// /    \  \/\_  __ \_/ __ \\__  \\   __\/ __ \
// \     \____|  | \/\  ___/ / __ \|  | \  ___/
//  \______  /|__|    \___  >____  /__|  \___  >
//         \/             \/     \/          \/

type CreatePayload struct {
	Ref           string      `json:"ref"`
	RefType       string      `json:"ref_type"`
	DefaultBranch string      `json:"default_branch"`
	Repo          *Repository `json:"repository"`
	Sender        *User       `json:"sender"`
}

func (p *CreatePayload) JSONPayload() ([]byte, error) {
	return json.MarshalIndent(p, "", "  ")
}

// ParseCreateHook parses create event hook content.
func ParseCreateHook(raw []byte) (*CreatePayload, error) {
	hook := new(CreatePayload)
	if err := json.Unmarshal(raw, hook); err != nil {
		return nil, err
	}

	// it is possible the JSON was parsed, however,
	// was not from Gogs (maybe was from Bitbucket)
	// So we'll check to be sure certain key fields
	// were populated
	switch {
	case hook.Repo == nil:
		return nil, ErrInvalidReceiveHook
	case len(hook.Ref) == 0:
		return nil, ErrInvalidReceiveHook
	}
	return hook, nil
}

// ________         .__          __
// \______ \   ____ |  |   _____/  |_  ____
//  |    |  \_/ __ \|  | _/ __ \   __\/ __ \
//  |    `   \  ___/|  |_\  ___/|  | \  ___/
// /_______  /\___  >____/\___  >__|  \___  >
//         \/     \/          \/          \/

type PusherType string

const (
	PUSHER_TYPE_USER PusherType = "user"
)

type DeletePayload struct {
	Ref        string      `json:"ref"`
	RefType    string      `json:"ref_type"`
	PusherType PusherType  `json:"pusher_type"`
	Repo       *Repository `json:"repository"`
	Sender     *User       `json:"sender"`
}

func (p *DeletePayload) JSONPayload() ([]byte, error) {
	return json.MarshalIndent(p, "", "  ")
}

// ___________           __
// \_   _____/__________|  | __
//  |    __)/  _ \_  __ \  |/ /
//  |     \(  <_> )  | \/    <
//  \___  / \____/|__|  |__|_ \
//      \/                   \/

type ForkPayload struct {
	Forkee *Repository `json:"forkee"`
	Repo   *Repository `json:"repository"`
	Sender *User       `json:"sender"`
}

func (p *ForkPayload) JSONPayload() ([]byte, error) {
	return json.MarshalIndent(p, "", "  ")
}

// __________             .__
// \______   \__ __  _____|  |__
//  |     ___/  |  \/  ___/  |  \
//  |    |   |  |  /\___ \|   Y  \
//  |____|   |____//____  >___|  /
//                      \/     \/

// PushPayload represents a payload information of push event.
type PushPayload struct {
	Ref        string           `json:"ref"`
	Before     string           `json:"before"`
	After      string           `json:"after"`
	CompareURL string           `json:"compare_url"`
	Commits    []*PayloadCommit `json:"commits"`
	Repo       *Repository      `json:"repository"`
	Pusher     *User            `json:"pusher"`
	Sender     *User            `json:"sender"`
}

func (p *PushPayload) JSONPayload() ([]byte, error) {
	return json.MarshalIndent(p, "", "  ")
}

// ParsePushHook parses push event hook content.
func ParsePushHook(raw []byte) (*PushPayload, error) {
	hook := new(PushPayload)
	if err := json.Unmarshal(raw, hook); err != nil {
		return nil, err
	}

	switch {
	case hook.Repo == nil:
		return nil, ErrInvalidReceiveHook
	case len(hook.Ref) == 0:
		return nil, ErrInvalidReceiveHook
	}
	return hook, nil
}

// Branch returns branch name from a payload
func (p *PushPayload) Branch() string {
	return strings.Replace(p.Ref, "refs/heads/", "", -1)
}

// .___
// |   | ______ ________ __   ____
// |   |/  ___//  ___/  |  \_/ __ \
// |   |\___ \ \___ \|  |  /\  ___/
// |___/____  >____  >____/  \___  >
//          \/     \/            \/

type HookIssueAction string

const (
	HOOK_ISSUE_OPENED        HookIssueAction = "opened"
	HOOK_ISSUE_CLOSED        HookIssueAction = "closed"
	HOOK_ISSUE_REOPENED      HookIssueAction = "reopened"
	HOOK_ISSUE_EDITED        HookIssueAction = "edited"
	HOOK_ISSUE_ASSIGNED      HookIssueAction = "assigned"
	HOOK_ISSUE_UNASSIGNED    HookIssueAction = "unassigned"
	HOOK_ISSUE_LABEL_UPDATED HookIssueAction = "label_updated"
	HOOK_ISSUE_LABEL_CLEARED HookIssueAction = "label_cleared"
	HOOK_ISSUE_MILESTONED    HookIssueAction = "milestoned"
	HOOK_ISSUE_DEMILESTONED  HookIssueAction = "demilestoned"
	HOOK_ISSUE_SYNCHRONIZED  HookIssueAction = "synchronized"
)

type ChangesFromPayload struct {
	From string `json:"from"`
}

type ChangesPayload struct {
	Title *ChangesFromPayload `json:"title,omitempty"`
	Body  *ChangesFromPayload `json:"body,omitempty"`
}

// IssuesPayload represents a payload information of issues event.
type IssuesPayload struct {
	Action     HookIssueAction `json:"action"`
	Index      int64           `json:"number"`
	Issue      *Issue          `json:"issue"`
	Changes    *ChangesPayload `json:"changes,omitempty"`
	Repository *Repository     `json:"repository"`
	Sender     *User           `json:"sender"`
}

func (p *IssuesPayload) JSONPayload() ([]byte, error) {
	return json.MarshalIndent(p, "", "  ")
}

type HookIssueCommentAction string

const (
	HOOK_ISSUE_COMMENT_CREATED HookIssueCommentAction = "created"
	HOOK_ISSUE_COMMENT_EDITED  HookIssueCommentAction = "edited"
	HOOK_ISSUE_COMMENT_DELETED HookIssueCommentAction = "deleted"
)

// IssueCommentPayload represents a payload information of issue comment event.
type IssueCommentPayload struct {
	Action     HookIssueCommentAction `json:"action"`
	Issue      *Issue                 `json:"issue"`
	Comment    *Comment               `json:"comment"`
	Changes    *ChangesPayload        `json:"changes,omitempty"`
	Repository *Repository            `json:"repository"`
	Sender     *User                  `json:"sender"`
}

func (p *IssueCommentPayload) JSONPayload() ([]byte, error) {
	return json.MarshalIndent(p, "", "  ")
}

// __________      .__  .__    __________                                     __
// \______   \__ __|  | |  |   \______   \ ____  ________ __   ____   _______/  |_
//  |     ___/  |  \  | |  |    |       _// __ \/ ____/  |  \_/ __ \ /  ___/\   __\
//  |    |   |  |  /  |_|  |__  |    |   \  ___< <_|  |  |  /\  ___/ \___ \  |  |
//  |____|   |____/|____/____/  |____|_  /\___  >__   |____/  \___  >____  > |__|
//                                     \/     \/   |__|           \/     \/

// PullRequestPayload represents a payload information of pull request event.
type PullRequestPayload struct {
	Action      HookIssueAction `json:"action"`
	Index       int64           `json:"number"`
	PullRequest *PullRequest    `json:"pull_request"`
	Changes     *ChangesPayload `json:"changes,omitempty"`
	Repository  *Repository     `json:"repository"`
	Sender      *User           `json:"sender"`
}

func (p *PullRequestPayload) JSONPayload() ([]byte, error) {
	return json.MarshalIndent(p, "", "  ")
}

// __________       .__
// \______   \ ____ |  |   ____ _____    ______ ____
//  |       _// __ \|  | _/ __ \\__  \  /  ___// __ \
//  |    |   \  ___/|  |_\  ___/ / __ \_\___ \\  ___/
//  |____|_  /\___  >____/\___  >____  /____  >\___  >
//         \/     \/          \/     \/     \/     \/

type HookReleaseAction string

const (
	HOOK_RELEASE_PUBLISHED HookReleaseAction = "published"
)

// ReleasePayload represents a payload information of release event.
type ReleasePayload struct {
	Action     HookReleaseAction `json:"action"`
	Release    *Release          `json:"release"`
	Repository *Repository       `json:"repository"`
	Sender     *User             `json:"sender"`
}

func (p *ReleasePayload) JSONPayload() ([]byte, error) {
	return json.MarshalIndent(p, "", "  ")
}
