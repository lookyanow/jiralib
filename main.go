package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/andygrunwald/go-jira"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
)

const JiraToken = ""

type jiraAuthTransport struct {
	Token string

	Transport http.RoundTripper
}

func (t *jiraAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req2 := cloneRequest(req) // per RoundTripper contract

	//req2.SetBasicAuth(t.Username, t.Password)
	req2.Header.Set("Authorization", "Basic "+ t.Token)
	return t.transport().RoundTrip(req2)
}

func (t *jiraAuthTransport) Client() *http.Client {
	return &http.Client{Transport: t}
}

func (t *jiraAuthTransport) transport() http.RoundTripper {
	if t.Transport != nil {
		return t.Transport
	}
	return http.DefaultTransport
}


func cloneRequest(r *http.Request) *http.Request {
	// shallow copy of the struct
	r2 := new(http.Request)
	*r2 = *r
	// deep copy of the Header
	r2.Header = make(http.Header, len(r.Header))
	for k, s := range r.Header {
		r2.Header[k] = append([]string(nil), s...)
	}
	return r2
}

type Issue struct {
	Id     string       `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Key    string       `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	Fields *IssueFields `protobuf:"bytes,3,opt,name=fields" json:"fields,omitempty"`
}

type IssueFields struct {
	Summary        string                `protobuf:"bytes,1,opt,name=summary,proto3" json:"summary,omitempty"`
	Description    string                `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	Assignee       *IssueFields_Assignee `protobuf:"bytes,3,opt,name=assignee" json:"assignee,omitempty"`
	Labels         []string              `protobuf:"bytes,4,rep,name=labels" json:"labels,omitempty"`
	Status         *IssueFields_Status   `protobuf:"bytes,5,opt,name=status" json:"status,omitempty"`
	Project        string                `protobuf:"bytes,6,opt,name=Project,json=customfield_18900,proto3" json:"customfield_18900"`
	BranchName     string                `protobuf:"bytes,7,opt,name=BranchName,json=customfield_18901,proto3" json:"customfield_18901"`
	UnitTestPassed string                `protobuf:"bytes,8,opt,name=UnitTestPassed,json=customfield_18807,proto3" json:"customfield_18807"`
}

type IssueFields_Assignee struct {
	Name         string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Key          string `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	EmailAddress string `protobuf:"bytes,3,opt,name=emailAddress,proto3" json:"emailAddress,omitempty"`
	DisplayName  string `protobuf:"bytes,4,opt,name=displayName,proto3" json:"displayName,omitempty"`
}

type IssueFields_Status struct {
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}


func JiraTokenToUserPass(token string) (user string, pass string){
	sDec, _ := base64.StdEncoding.DecodeString(token)
	s := strings.SplitN(string(sDec), ":",2)
	return s[0], s[1]
}

func GetIssue(issue string) (*Issue, error) {
	user, pass := JiraTokenToUserPass(JiraToken)
	tp := jira.BasicAuthTransport{Username: user, Password: pass}
	jiraClient, err := jira.NewClient(tp.Client(), "https://jit.ozon.ru/")
	if err != nil {
		fmt.Println(err)
	} else {

		issue1, _, err := jiraClient.Issue.Get(issue, nil)
		if err != nil {
			return nil, err
		} else {
			data, err := json.Marshal(issue1)
			o := &Issue{}
			err = json.Unmarshal(data, o)
			return o, err
		}
	}
	return nil, err
}

func SetIssueFields(issue string, fieldValues map[string]interface{}) error {

	issue = strings.ToUpper(issue)

	query := map[string]interface{}{
		"fields": fieldValues,
	}
	tp := jiraAuthTransport{Token: JiraToken}
	jiraClient, err := jira.NewClient(tp.Client(), "https://jit.ozon.ru/")
	if err != nil {
		return err
	}

	r, err := jiraClient.Issue.UpdateIssue(issue, query)
	if err != nil {
		body, err := ioutil.ReadAll(r.Body)
		fmt.Println(string(body))
		return err
	}
	return err

}

func SetIssueField(issue, field string, value interface{}) error {

	query := map[string]interface{}{
		"fields": map[string]interface{}{
			field: value,
		},
	}

	tp := jiraAuthTransport{Token: JiraToken}
	jiraClient, err := jira.NewClient(tp.Client(), "https://jira.ozon.ru/")
	if err != nil {
		return err
	} else {

		_, err := jiraClient.Issue.UpdateIssue(issue, query)
		if err != nil {

			return err
		}

	}
	return err

}

func getUnitFieldNames() string {
	t := reflect.TypeOf(IssueFields{})
	field1, _ := t.FieldByName("UnitTestPassed")
	tag1 := field1.Tag.Get("json")
	return tag1
}

func getProjectBranchFieldNames() (string, string) {
	t := reflect.TypeOf(IssueFields{})
	field1, _ := t.FieldByName("Project")
	field2, _ := t.FieldByName("BranchName")
	tag1 := field1.Tag.Get("json")
	tag2 := field2.Tag.Get("json")
	return tag1, tag2
}

func main() {
	issue, err := GetIssue("SF-1355")
	if err != nil{
		panic(err)
	} else {
		//data, err := json.Marshal(issue)
		//fmt.Printf("%s\n", string(data))
		fmt.Printf("%s: %+v\n", issue.Key, issue.Fields.Summary)
		//fmt.Printf("Type: %s\n", issue.Fields.Type.Name)
		//fmt.Printf("Priority: %s\n", issue.Fields.Priority.Name)
		fmt.Printf("%s\n", issue.Fields.Assignee.DisplayName)
		fmt.Printf("%+v\n", issue.Fields.Project)
		fmt.Printf("%+v\n", issue.Fields.Status.Name)


		tp := jiraAuthTransport{Token: JiraToken}
		jiraClient, err := jira.NewClient(tp.Client(), "https://jira.ozon.ru/")
		if err != nil {
			panic(err)
		}


		trans, _, err := jiraClient.Issue.GetTransitions("RE-350")
		if err != nil {
			panic(err)
		}
		transition := ""
		for _, r := range trans {
			if r.Name == "Start"{
				transition = r.ID
			}
			fmt.Printf("%s %s\n", r.ID, r.Name)
		}
		fmt.Println(transition)
/*
		comment := &jira.Comment{
			Body: "test failed",
		}
		comm, res, err := jiraClient.Issue.AddComment("DEMO-10", comment)
		if err != nil{
			panic(err)
			fmt.Printf("%d responce code", res.StatusCode)
			fmt.Println(comm.Author.Name)
		}
		fmt.Println(comm.Author.Name)
		fmt.Println("Issue updated")*/
	}
}
