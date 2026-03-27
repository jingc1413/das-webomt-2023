package story

import (
	"crypto/md5"
	"fmt"

	"github.com/andygrunwald/go-jira"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var (
	JiraBase     = "http://20.0.0.70:8080"
	JiraUser     = "weiyu_chen"
	JiraPassword = "123321"
)

func getJiraClient() *jira.Client {
	tp := jira.BasicAuthTransport{
		Username: JiraUser,
		Password: JiraPassword,
	}
	cli, err := jira.NewClient(tp.Client(), JiraBase)
	if err != nil {
		logrus.Fatal(errors.Wrap(err, "create jira client"))
	}
	return cli
}

func SearchIssue(jql string, max int) ([]jira.Issue, error) {
	cli := getJiraClient()
	opts := jira.SearchOptions{
		MaxResults: max,
	}
	issues, _, err := cli.Issue.Search(jql, &opts)
	if err != nil {
		return issues, errors.Wrap(err, "search issues")
	}
	return issues, nil
}

func SearchProejctIssue(projectKey string, typ string, ext string) ([]jira.Issue, error) {
	jql := fmt.Sprintf("project = %v AND issuetype = %v", projectKey, typ)
	if ext != "" {
		jql = jql + " AND " + ext
	}
	return SearchIssue(jql, 10000)
}

func SearchProejctIssueBySummary(projectKey string, typ string, summary string, ext string) ([]jira.Issue, error) {
	jql := fmt.Sprintf(`project = %v AND issuetype = %v AND summary ~ "%v"`, projectKey, typ, summary)
	if ext != "" {
		jql = jql + " AND " + ext
	}
	return SearchIssue(jql, 10000)
}

func GetProject(projectID string) (*jira.Project, error) {
	cli := getJiraClient()
	project, _, err := cli.Project.Get(projectID)
	if err != nil {
		return nil, errors.Wrap(err, "get project")
	}
	if project == nil {
		return nil, errors.New("invalid project id")
	}
	return project, nil
}

func CreateProjectComponent(project *jira.Project, componentName string) (*jira.ProjectComponent, error) {
	cli := getJiraClient()
	opts := jira.CreateComponentOptions{
		Name:    componentName,
		Project: project.Key,
	}
	component, _, err := cli.Component.Create(&opts)
	if err != nil {
		return nil, errors.Wrap(err, "create component")
	}
	return component, nil
}

func GetProjectComponetsByNames(project *jira.Project, componentNames []string) ([]*jira.ProjectComponent, error) {
	components := []*jira.ProjectComponent{}
	for _, componentName := range componentNames {
		var component *jira.ProjectComponent = nil
		for _, v := range project.Components {
			if v.Name == componentName {
				tmp := v
				component = &tmp
				break
			}
		}
		if component == nil {
			v, err := CreateProjectComponent(project, componentName)
			if err != nil {
				return nil, errors.Wrap(err, "create component")
			}
			component = v
		}
		if component != nil {
			components = append(components, component)
		}
	}
	return components, nil
}

func UpdateOrCreateIssue(projectID string, componentName string, typeName string, summary string, description string) (*jira.Issue, error) {
	project, err := GetProject(projectID)
	if err != nil {
		return nil, errors.Wrap(err, "get project")
	}
	// projectComponents, err := GetProjectComponetsByNames(project, []string{componentName})
	// if err != nil {
	// 	return nil, errors.Wrap(err, "get project components by name")
	// }
	// components := []*jira.Component{}
	// for _, v := range projectComponents {
	// 	components = append(components, &jira.Component{
	// 		ID: v.ID,
	// 	})
	// }

	// ext := fmt.Sprintf(`resolution = Unresolved AND component = "%v"`, componentName)
	ext := fmt.Sprintf(`component = "%v"`, componentName)
	_find, err := SearchProejctIssueBySummary(project.Key, typeName, summary, ext)
	if err != nil {
		return nil, errors.Wrap(err, "search exists issues")
	}
	find := []jira.Issue{}
	for _, v := range _find {
		if v.Fields.Summary == summary {
			tmp := v
			find = append(find, tmp)
		}
	}
	if len(find) > 1 {
		return nil, errors.Errorf("update failed, mutiple issues with same summary, find=%v", find)
	}
	var issue *jira.Issue = nil
	if len(find) == 1 {
		issue = &find[0]
		if err := UpdateIssueDescription(issue, description); err != nil {
			return nil, errors.Wrap(err, "update issue description")
		}
	} else {
		tmp, err := CreateIssue(project.Key, typeName, summary, description)
		if err != nil {
			return nil, errors.Wrap(err, "crate issue")
		}
		issue = tmp
	}
	if err := UpdateIssueComponents(issue, []string{componentName}); err != nil {
		return nil, errors.Wrap(err, "update issue components")
	}
	return GetIssue(issue.ID)
}

func CreateIssue(projectKey string, typeName string, summary string, description string) (*jira.Issue, error) {
	cli := getJiraClient()
	_issue := jira.Issue{
		Fields: &jira.IssueFields{
			Project: jira.Project{
				Key: projectKey,
			},
			Type: jira.IssueType{
				Name: typeName,
			},
			Summary:     summary,
			Description: description,
		},
	}
	issue, _, err := cli.Issue.Create(&_issue)
	if err != nil {
		return issue, errors.Wrap(err, "create issue")
	}
	return issue, nil
}

func UpdateIssueComponents(issue *jira.Issue, components []string) error {
	cli := getJiraClient()
	setComponentNames := []map[string]string{}
	for _, v := range components {
		setComponentNames = append(setComponentNames, map[string]string{
			"name": v,
		})
	}
	data := map[string]any{}
	data["update"] = map[string]any{
		"components": []any{
			map[string]any{
				"set": setComponentNames,
			},
		},
	}
	_, err := cli.Issue.UpdateIssue(issue.ID, data)
	return err
}

func UpdateIssueDescription(issue *jira.Issue, description string) error {
	cli := getJiraClient()
	m1 := md5.Sum([]byte(issue.Fields.Description))
	m2 := md5.Sum([]byte(description))
	if m1 == m2 {
		return nil
	}
	data := map[string]any{}
	data["fields"] = map[string]any{
		"description": description,
	}
	issue.Fields.Description = description
	_, err := cli.Issue.UpdateIssue(issue.ID, data)
	return err
}

func GetIssue(issueID string) (*jira.Issue, error) {
	cli := getJiraClient()
	issue, _, err := cli.Issue.Get(issueID, nil)
	return issue, err
}

func RemoveIssue(issue *jira.Issue, description string) error {
	cli := getJiraClient()
	_, err := cli.Issue.Delete(issue.ID)
	if err != nil {
		return errors.Wrap(err, "update issue")
	}
	return nil
}
