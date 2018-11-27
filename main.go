package main

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
	"io/ioutil"
	"net/http"
)


type JiraAuthTransport struct {
	Token string

	Transport http.RoundTripper
}

func (t *JiraAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req2 := cloneRequest(req) // per RoundTripper contract

	//req2.SetBasicAuth(t.Username, t.Password)
	req2.Header.Set("Authorization", "Basic "+ t.Token)
	return t.transport().RoundTrip(req2)
}

func (t *JiraAuthTransport) Client() *http.Client {
	return &http.Client{Transport: t}
}

func (t *JiraAuthTransport) transport() http.RoundTripper {
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

func main() {

	tp := JiraAuthTransport{Token:""}
	jiraClient, err := jira.NewClient(tp.Client(), "https://jira.ozon.ru/")
	if err != nil {
		fmt.Println(err)
	} else {
		issue, res, err := jiraClient.Issue.Get("DEMO-10", nil)
		if err != nil {
			panic(err)

			data, err := ioutil.ReadAll(res.Body)
			if err != nil {
				fmt.Printf("Read Error: %s\n", err)
			}

			fmt.Printf("%+v\n", string(data))
		}else {
			fmt.Printf("%s: %+v\n", issue.Key, issue.Fields.Summary)
			fmt.Printf("Type: %s\n", issue.Fields.Type.Name)
			fmt.Printf("Priority: %s\n", issue.Fields.Priority.Name)
			fmt.Printf("%s\n", issue.Fields.Assignee.DisplayName)
			fmt.Printf("%+v\n", issue.Fields.Labels)
		}

		labels := append(issue.Fields.Labels, "test1")

		query := map[string]interface{}{
			"fields" : map[string]interface{}{
				"labels": labels,
			},
		}


		_, err = jiraClient.Issue.UpdateIssue("DEMO-10", query)
		if err != nil{
			panic(err)
		}

		query = map[string]interface{}{
			"fields" : map[string]interface{}{
				"description" : "new description",
			},
		}
		_, err = jiraClient.Issue.UpdateIssue("DEMO-10", query)
		if err != nil{
			panic(err)
		}
		}

	}
