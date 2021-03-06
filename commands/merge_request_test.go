package commands

import (
	"testing"
	"time"

	"github.com/lighttiger2505/lab/git"
	lab "github.com/lighttiger2505/lab/gitlab"
	"github.com/lighttiger2505/lab/ui"
	gitlab "github.com/xanzy/go-gitlab"
)

var mergeRequest = &gitlab.MergeRequest{
	IID:   12,
	Title: "Title12",
	Assignee: struct {
		ID        int        `json:"id"`
		Username  string     `json:"username"`
		Name      string     `json:"name"`
		State     string     `json:"state"`
		CreatedAt *time.Time `json:"created_at"`
	}{
		Name: "AssigneeName",
	},
	Author: struct {
		ID        int        `json:"id"`
		Username  string     `json:"username"`
		Name      string     `json:"name"`
		State     string     `json:"state"`
		CreatedAt *time.Time `json:"created_at"`
	}{
		Name: "AuthorName",
	},
	WebURL:      "http://gitlab.jp/namespace/repo12",
	CreatedAt:   &createdAt,
	UpdatedAt:   &updatedAt,
	Description: "Description",
}

var mergeRequests []*gitlab.MergeRequest = []*gitlab.MergeRequest{
	&gitlab.MergeRequest{IID: 12, Title: "Title12", WebURL: "http://gitlab.jp/namespace/repo12"},
	&gitlab.MergeRequest{IID: 13, Title: "Title13", WebURL: "http://gitlab.jp/namespace/repo13"},
}

var mockGitlabMergeRequestClient = &lab.MockLabMergeRequestClient{
	MockGetMergeRequest: func(pid int, repositoryName string) (*gitlab.MergeRequest, error) {
		return mergeRequest, nil
	},
	MockGetAllProjectMergeRequest: func(opt *gitlab.ListMergeRequestsOptions) ([]*gitlab.MergeRequest, error) {
		return mergeRequests, nil
	},
	MockGetProjectMargeRequest: func(opt *gitlab.ListProjectMergeRequestsOptions, repositoryName string) ([]*gitlab.MergeRequest, error) {
		return mergeRequests, nil
	},
	MockCreateMergeRequest: func(opt *gitlab.CreateMergeRequestOptions, repositoryName string) (*gitlab.MergeRequest, error) {
		return mergeRequest, nil
	},
	MockUpdateMergeRequest: func(opt *gitlab.UpdateMergeRequestOptions, pid int, repositoryName string) (*gitlab.MergeRequest, error) {
		return mergeRequest, nil
	}}

var mockMergeRequestProvider = &lab.MockProvider{
	MockInit: func() error { return nil },
	MockGetSpecificRemote: func(namespace, project string) *git.RemoteInfo {
		return &git.RemoteInfo{
			Domain:     "domain",
			NameSpace:  "namespace",
			Repository: "repository",
		}
	},
	MockGetCurrentRemote: func() (*git.RemoteInfo, error) {
		return &git.RemoteInfo{
			Domain:     "domain",
			NameSpace:  "namespace",
			Repository: "repository",
		}, nil
	},
	MockGetMergeRequestClient: func(remote *git.RemoteInfo) (lab.MergeRequest, error) {
		return mockGitlabMergeRequestClient, nil
	},
}

func TestMergeRequestCommandRun_List(t *testing.T) {
	mockUI := ui.NewMockUi()
	c := MergeRequestCommand{
		Ui:       mockUI,
		Provider: mockMergeRequestProvider,
	}

	args := []string{}
	if code := c.Run(args); code != 0 {
		t.Fatalf("wrong exit code. errors: \n%s", mockUI.ErrorWriter.String())
	}

	got := mockUI.Writer.String()
	want := "!12  Title12\n!13  Title13\n"
	if want != got {
		t.Fatalf("bad output value \nwant %#v \ngot  %#v", want, got)
	}
}

func TestMergeRequestCommandRun_ListAll(t *testing.T) {
	mockUI := ui.NewMockUi()
	c := MergeRequestCommand{
		Ui:       mockUI,
		Provider: mockMergeRequestProvider,
	}

	args := []string{"--all-project"}
	if code := c.Run(args); code != 0 {
		t.Fatalf("wrong exit code. errors: \n%s", mockUI.ErrorWriter.String())
	}

	got := mockUI.Writer.String()
	want := "namespace/repo12  !12  Title12\nnamespace/repo13  !13  Title13\n"
	if want != got {
		t.Fatalf("bad output value \nwant %#v \ngot  %#v", want, got)
	}
}

func TestMergeRequestCommandRun_Create(t *testing.T) {
	mockUI := ui.NewMockUi()
	c := MergeRequestCommand{
		Ui:       mockUI,
		Provider: mockMergeRequestProvider,
	}

	args := []string{"-i", "title", "-m", "message"}
	if code := c.Run(args); code != 0 {
		t.Fatalf("wrong exit code. errors: \n%s", mockUI.ErrorWriter.String())
	}

	got := mockUI.Writer.String()
	want := "!12\n"
	if want != got {
		t.Fatalf("bad output value \nwant %#v \ngot  %#v", want, got)
	}
}

func TestMergeRequestCommandRun_CreateOnEditor(t *testing.T) {
	mockUI := ui.NewMockUi()
	c := MergeRequestCommand{
		Ui:       mockUI,
		Provider: mockMergeRequestProvider,
		EditFunc: func(program, file string) error {
			return nil
		},
	}

	args := []string{"-e", "-i", "title", "-m", "message"}
	if code := c.Run(args); code != 0 {
		t.Fatalf("wrong exit code. errors: \n%s", mockUI.ErrorWriter.String())
	}

	got := mockUI.Writer.String()
	want := "!12\n"
	if want != got {
		t.Fatalf("bad output value \nwant %#v \ngot  %#v", want, got)
	}
}

func TestMergeRequestCommandRun_Update(t *testing.T) {
	mockUI := ui.NewMockUi()
	c := MergeRequestCommand{
		Ui:       mockUI,
		Provider: mockMergeRequestProvider,
	}

	args := []string{"-i", "title", "-m", "message", "12"}
	if code := c.Run(args); code != 0 {
		t.Fatalf("wrong exit code. errors: \n%s", mockUI.ErrorWriter.String())
	}

	got := mockUI.Writer.String()
	want := "!12\n"
	if want != got {
		t.Fatalf("bad output value \nwant %#v \ngot  %#v", want, got)
	}
}

func TestMergeRequestCommandRun_UpdateOnEditor(t *testing.T) {
	mockUI := ui.NewMockUi()
	c := MergeRequestCommand{
		Ui:       mockUI,
		Provider: mockMergeRequestProvider,
		EditFunc: func(program, file string) error {
			return nil
		},
	}

	args := []string{"-e", "-i", "title", "-m", "message", "12"}
	if code := c.Run(args); code != 0 {
		t.Fatalf("wrong exit code. errors: \n%s", mockUI.ErrorWriter.String())
	}

	got := mockUI.Writer.String()
	want := "!12\n"
	if want != got {
		t.Fatalf("bad output value \nwant %#v \ngot  %#v", want, got)
	}
}

func TestMergeRequestCommandRun_Show(t *testing.T) {
	mockUI := ui.NewMockUi()
	c := MergeRequestCommand{
		Ui:       mockUI,
		Provider: mockMergeRequestProvider,
	}

	args := []string{"12"}
	if code := c.Run(args); code != 0 {
		t.Fatalf("wrong exit code. errors: \n%s", mockUI.ErrorWriter.String())
	}

	got := mockUI.Writer.String()
	want := `!12
Title: Title12
Assignee: AssigneeName
Author: AuthorName
CreatedAt: 2018-02-14 00:00:00 +0000 UTC
UpdatedAt: 2018-03-14 00:00:00 +0000 UTC

Description
`
	if want != got {
		t.Fatalf("bad output value \nwant %#v \ngot  %#v", want, got)
	}
}
