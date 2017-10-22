package main

import (
	"flag"
	"fmt"
	"runtime"
	"strconv"
	"strings"

	"github.com/ryanuber/columnize"
	"github.com/spf13/viper"
	"github.com/xanzy/go-gitlab"
)

type ProjectCommand struct {
}

func (c *ProjectCommand) Synopsis() string {
	return "Browse project"
}

func (c *ProjectCommand) Help() string {
	return "Usage: lab project [option]"
}

func (c *ProjectCommand) Run(args []string) int {
	var verbose bool

	// Set subcommand flags
	flags := flag.NewFlagSet("project", flag.ContinueOnError)
	flags.BoolVar(&verbose, "verbose", false, "Run as debug mode")
	flags.Usage = func() {}
	if err := flags.Parse(args); err != nil {
		return ExitCodeError
	}

	gitRemotes, err := GitRemotes()
	if err != nil {
		fmt.Println(err.Error())
		return ExitCodeError
	}

	gitlabRemote, err := FilterGitlabRemote(gitRemotes)
	if err != nil {
		fmt.Println(err.Error())
		return ExitCodeError
	}

	browser := SearchBrowserLauncher(runtime.GOOS)
	cmdOutput(browser, []string{gitlabRemote.RepositoryUrl()})

	return ExitCodeOK
}

type IssueCommand struct {
}

func (c *IssueCommand) Synopsis() string {
	return "Browse Issue"
}

func (c *IssueCommand) Help() string {
	return "Usage: lab issue [option]"
}

func (c *IssueCommand) Run(args []string) int {
	var verbose bool

	// Set subcommand flags
	flags := flag.NewFlagSet("project", flag.ContinueOnError)
	flags.BoolVar(&verbose, "verbose", false, "Run as debug mode")
	flags.Usage = func() {}
	if err := flags.Parse(args); err != nil {
		return ExitCodeError
	}

	gitRemotes, err := GitRemotes()
	if err != nil {
		fmt.Println(err.Error())
		return ExitCodeError
	}

	gitlabRemote, err := FilterGitlabRemote(gitRemotes)
	if err != nil {
		fmt.Println(err.Error())
		return ExitCodeError
	}

	// Read config file
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.lab")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
		return ExitCodeError
	}
	privateToken := viper.GetString("private_token")

	// Create client
	client := gitlab.NewClient(nil, privateToken)
	client.SetBaseURL(gitlabRemote.ApiUrl())

	listProjectOptions := &gitlab.ListProjectsOptions{Search: gitlab.String(gitlabRemote.Repository)}
	projects, _, err := client.Projects.ListProjects(listProjectOptions)

	if err != nil {
		fmt.Println(err)
		return ExitCodeError
	}

	// Get project id
	var projectId int
	for _, project := range projects {
		fullName := strings.Replace(project.NameWithNamespace, " ", "", -1)
		if fullName == gitlabRemote.FullName() {
			projectId = project.ID
		}
	}

	listOption := &gitlab.ListOptions{
		Page:    1,
		PerPage: 20,
	}
	listProjectIssuesOptions := &gitlab.ListProjectIssuesOptions{
		Scope:       gitlab.String("assigned-to-me"),
		OrderBy:     gitlab.String("updated_at"),
		Sort:        gitlab.String("desc"),
		ListOptions: *listOption,
	}
	issues, _, err := client.Issues.ListProjectIssues(projectId, listProjectIssuesOptions)

	if err != nil {
		fmt.Println(err)
		return ExitCodeError
	}

	var datas []string
	for _, issue := range issues {
		data := fmt.Sprint(issue.IID) + "|" + issue.Title
		datas = append(datas, data)
	}

	result := columnize.SimpleFormat(datas)
	fmt.Println(result)
	return ExitCodeOK
}

type MergeRequestCommand struct {
}

func (c *MergeRequestCommand) Synopsis() string {
	return "Browse merge request"
}

func (c *MergeRequestCommand) Help() string {
	return "Usage: lab merge-request [option]"
}

func (c *MergeRequestCommand) Run(args []string) int {
	var verbose bool

	// Set subcommand flags
	flags := flag.NewFlagSet("browse", flag.ContinueOnError)
	flags.BoolVar(&verbose, "verbose", false, "Run as debug mode")
	flags.Usage = func() {}
	if err := flags.Parse(args); err != nil {
		return ExitCodeError
	}

	gitRemotes, err := GitRemotes()
	if err != nil {
		fmt.Println(err.Error())
		return ExitCodeError
	}

	gitlabRemote, err := FilterGitlabRemote(gitRemotes)
	if err != nil {
		fmt.Println(err.Error())
		return ExitCodeError
	}

	browser := SearchBrowserLauncher(runtime.GOOS)

	if len(flags.Args()) > 0 {
		issueNo, err := strconv.Atoi(flags.Args()[0])
		if err != nil {
			fmt.Println(err.Error())
		}
		cmdOutput(browser, []string{gitlabRemote.MergeRequestDetailUrl(issueNo)})
	} else {
		cmdOutput(browser, []string{gitlabRemote.MergeRequestUrl()})
	}

	return ExitCodeOK
}