package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

type GoScrumClient struct {
	URL   string
	Token string
}

func NewGoScrumClient(URL, token string) GoScrumClient {
	return GoScrumClient{URL: URL, Token: token}
}

func (g *GoScrumClient) GetWorkspaceByToken() (*Workspace, error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	fmt.Println("Making call to GoScrum service")

	req, err := http.NewRequest("GET", g.URL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", g.Token))

	fmt.Println("Making request", g.URL)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return nil, err

	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return nil, err
	}
	workspace := Workspace{}

	err = json.Unmarshal(body, &workspace)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return nil, err
	}
	return &workspace, nil
}

func (g *GoScrumClient) GetParticipantQuestion(projectId, participantId string) (*Question, error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	url := fmt.Sprintf("%s/%s/%s/question", g.URL, projectId, participantId)
	fmt.Println("Making call to GoScrum service", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", g.Token))

	fmt.Println("Making request", url)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return nil, err

	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return nil, err
	}
	question := Question{}

	err = json.Unmarshal(body, &question)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return nil, err
	}
	return &question, nil
}
