package container

import (
	"fmt"
	"testing"
)

func Test_Worker_Handler_Operator_Container_imaTag(t *testing.T) {
	testCases := []struct {
		str string
		tag string
		mat func(error) bool
	}{
		// Case 000
		{
			str: "foo/bar:v1.2.3",
			tag: "v1.2.3",
			mat: isNil,
		},
		// Case 001
		{
			str: "registry/repository:tag",
			tag: "tag",
			mat: isNil,
		},
		// Case 002
		{
			str: "ubuntu:22",
			tag: "22",
			mat: isNil,
		},
		// Case 003
		{
			str: "ubuntu",
			tag: "",
			mat: isNil,
		},
		// Case 004
		{
			str: "splits/server:latest",
			tag: "latest",
			mat: isNil,
		},
		// Case 005
		{
			str: "123.dkr.ecr.us-west-2.amazonaws.com/repo@sha256:94a00394bc5a8ef503fb59db0a7d0ae9e1110866e8aee8ba40cd864cea69ea1a",
			tag: "sha256:94a00394bc5a8ef503fb59db0a7d0ae9e1110866e8aee8ba40cd864cea69ea1a",
			mat: isNil,
		},
		// Case 006
		{
			str: "@digest",
			tag: "",
			mat: isErr,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%03d", i), func(t *testing.T) {
			tag, err := imaTag(tc.str)
			if !tc.mat(err) {
				t.Fatal("expected", true, "got", err)
			}

			if tag != tc.tag {
				t.Fatalf("expected %#v got %#v", tc.tag, tag)
			}
		})
	}
}
