package status

import (
	"context"
	"fmt"
	"strings"

	"github.com/0xSplits/kayron/pkg/cache"
	"github.com/google/go-github/v73/github"
	"github.com/xh3b4sd/choreo/parallel"
	"github.com/xh3b4sd/tracer"
)

const (
	marker = "<!-- kayron:preview:status -->"
)

func (s *Status) preview() error {
	// Collect all injected preview releases that have state drift.

	var dft []cache.Object
	for _, x := range s.pol.Drift() {
		if x.Preview() {
			dft = append(dft, x)
		}
	}

	fnc := func(_ int, o cache.Object) error {
		var err error

		var com *github.IssueComment
		{
			com, err = s.issCom(o)
			if err != nil {
				return tracer.Mask(err)
			}
		}

		// TODO comment

		var cre bool
		{
			cre = com == nil
		}

		if cre {
			err = s.creCom("Creating", o)
			if err != nil {
				return tracer.Mask(err)
			}
		}

		// TODO comment

		var upd bool
		{
			upd = !cre && !s.pol.Cancel() && !strings.Contains(com.GetBody(), "Ready")
		}

		if upd {
			err = s.updCom("Ready", o, com.GetID())
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
//	Status | Hash | Endpoint
//	---|----|---
//	Creating | 1D0FD508 | https://1d0fd508.lite.testing.splits.org
func (s *Status) comBod(sta string, obj cache.Object) string {
	return strings.Join(
		[]string{
			marker,
			fmt.Sprintf("%s | %s | %s", "Status", "Hash", "Endpoint"),
			/*********/ "---|----|---",
			fmt.Sprintf("%s | %s | %s", sta, obj.Release.Labels.Hash.Upper(), obj.Domain(s.env.Environment)),
		},
		"\n")
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

func (s *Status) updCom(sta string, obj cache.Object, cid int64) error {
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
