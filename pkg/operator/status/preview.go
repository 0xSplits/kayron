package status

import (
	"context"
	"fmt"
	"strings"

	"github.com/0xSplits/kayron/pkg/cache"
	"github.com/google/go-github/v76/github"
	"github.com/xh3b4sd/choreo/parallel"
	"github.com/xh3b4sd/tracer"
)

const (
	// marker identifies those issue comments managed by the operator. This marker
	// is effectively a markdown comment that helps us find and update our status
	// updates.
	marker = "<!-- kayron:preview:status -->"
)

func (s *Status) preview() error {
	// Collect all injected preview releases, whether they have state drift or
	// not. Either of those cases requires us to manage the preview deployment
	// status.

	var dft []cache.Object
	for _, x := range s.cac.Releases() {
		if x.Preview() {
			dft = append(dft, x)
		}
	}

	fnc := func(_ int, o cache.Object) error {
		var err error

		// Try to find our marked issue comment for any given preview deployment.

		var com *github.IssueComment
		{
			com, err = s.issCom(o)
			if err != nil {
				return tracer.Mask(err)
			}
		}

		// Determine the status field of our status update. This also tells us
		// whether we should update the issue comment if it exists.

		var sta string
		var upd bool
		{
			sta, upd = s.comSta(o, com)
		}

		// Create or update the issue comment with any given status, if any.

		if com == nil {
			err = s.creCom(sta, o)
			if err != nil {
				return tracer.Mask(err)
			}
		} else if upd {
			err = s.updCom(sta, o, com.GetID())
			if err != nil {
				return tracer.Mask(err)
			}
		}

		return nil
	}

	{
		err := parallel.Slice(dft, fnc)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	return nil
}

func (s *Status) creCom(sta string, obj cache.Object) error {
	s.log.Log(
		"level", "info",
		"message", "creating status update",
		"pull", s.pulReq(obj),
		"status", sta,
	)

	var com *github.IssueComment
	{
		com = &github.IssueComment{
			Body: github.Ptr(s.comBod(sta, obj)),
		}
	}

	{
		_, _, err := s.git.Issues.CreateComment(context.Background(), s.own, obj.Release.Github.String(), obj.Release.Labels.Pull, com)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	return nil
}

// comBod returns the body content of an issue comment used to render a status
// update for preview deployments.
//
//	<!-- kayron:preview:status -->
//	Status | Commit | Endpoint
//	---|----|---
//	Updating | c8e29b2 | 1d0fd508.lite.testing.splits.org
func (s *Status) comBod(sta string, obj cache.Object) string {
	var dom string
	{
		dom = obj.Domain(s.env.Environment)
	}

	var end string
	{
		end = fmt.Sprintf("[%s](https://%s)", dom, dom)
	}

	var com string
	{
		com = shoStr(obj.Artifact.Scheduler.Current, 7)
	}

	return strings.Join(
		[]string{
			marker,
			fmt.Sprintf("%s | %s | %s", "Status", "Commit", "Endpoint"),
			/*********/ "---|----|---",
			fmt.Sprintf("%s | %s | %s", sta, com, end),
		},
		"\n")
}

func (s *Status) comSta(obj cache.Object, com *github.IssueComment) (string, bool) {
	// If the preview deployment has any form of state drift, and if the status
	// update is not marked as updating, then the preview status is "Updating".
	// This status may be applied to existing and new issue comments.

	if obj.Drift(false) && !strings.Contains(com.GetBody(), "Updating") {
		return "Updating", true
	}

	// If the preview deployment has no state drift at all, and if the status
	// update is not marked as ready, then the preview status is "Ready". This
	// status may be applied to existing and new issue comments.

	if !obj.Drift(false) && !strings.Contains(com.GetBody(), "Ready") {
		return "Ready", true
	}

	return "", false
}

func (s *Status) issCom(obj cache.Object) (*github.IssueComment, error) {
	var err error

	var opt *github.IssueListCommentsOptions
	{
		opt = &github.IssueListCommentsOptions{
			ListOptions: github.ListOptions{
				PerPage: 5,
			},
		}
	}

	var com []*github.IssueComment
	{
		com, _, err = s.git.Issues.ListComments(context.Background(), s.own, obj.Release.Github.String(), obj.Release.Labels.Pull, opt)
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	for _, x := range com {
		if strings.Contains(x.GetBody(), marker) {
			return x, nil
		}
	}

	return nil, nil
}

func (s *Status) pulReq(obj cache.Object) string {
	return fmt.Sprintf("https://github.com/%s/%s/pull/%d", s.own, obj.Release.Github.String(), obj.Release.Labels.Pull)
}

func (s *Status) updCom(sta string, obj cache.Object, cid int64) error {
	s.log.Log(
		"level", "info",
		"message", "updating status update",
		"pull", s.pulReq(obj),
		"status", sta,
	)

	var com *github.IssueComment
	{
		com = &github.IssueComment{
			Body: github.Ptr(s.comBod(sta, obj)),
		}
	}

	{
		_, _, err := s.git.Issues.EditComment(context.Background(), s.own, obj.Release.Github.String(), cid, com)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	return nil
}

// shoStr returns teh given string with a maximum length of max, or the provided
// string itself if that string is shorter than max.
func shoStr(str string, max int) string {
	if len(str) >= max {
		return str[:max]
	}

	return str
}
