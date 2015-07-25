package gist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"path"
	"runtime"

	"github.com/zbindenren/gost/configuration"
)

var (
	url    = "https://api.github.com/gists"
	client = new(http.Client)
)

//GetResponse is used to unmarshal the get response
type GetResponse struct {
	Comments    int                    `json:"comments"`
	CommentsURL string                 `json:"comments_url"`
	CommitsURL  string                 `json:"commits_url"`
	CreatedAt   string                 `json:"created_at"`
	Description string                 `json:"description"`
	Files       map[string]fileDetails `json:"files"`
	ForksURL    string                 `json:"forks_url"`
	GitPullURL  string                 `json:"git_pull_url"`
	GitPushURL  string                 `json:"git_push_url"`
	HTMLURL     string                 `json:"html_url"`
	ID          string                 `json:"id"`
	Owner       owner                  `json:"owner"`
	Public      bool                   `json:"public"`
	UpdatedAt   string                 `json:"updated_at"`
	URL         string                 `json:"url"`
	User        interface{}            `json:"user"`
}

// OpenInBrowser opens the current gist response in a browser
func (g *GetResponse) OpenInBrowser() error {
	os := runtime.GOOS
	url := g.HTMLURL
	var err error
	switch {
	case os == "windows":
		err = exec.Command("cmd", "/c", "start", url).Run()
	case os == "darwin":
		err = exec.Command("open", url).Run()
	case os == "linux":
		err = exec.Command("xdg-open", url).Run()
	}
	return err
}

type fileDetails struct {
	FileName string `json:"filename"`
	Type     string `json:"type"`
	Language string `json:"language"`
	RawURL   string `json:"raw_url"`
	Size     int    `json:"size"`
	Content  string `json:"content"`
}
type owner struct {
	AvatarURL         string `json:"avatar_url"`
	EventsURL         string `json:"events_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	GravatarID        string `json:"gravatar_id"`
	HTMLURL           string `json:"html_url"`
	ID                int    `json:"id"`
	Login             string `json:"login"`
	OrganizationsURL  string `json:"organizations_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	ReposURL          string `json:"repos_url"`
	SiteAdmin         bool   `json:"site_admin"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	Type              string `json:"type"`
	URL               string `json:"url"`
}

//PostData is used to marshal the post request
type PostData struct {
	Desc   string          `json:"description"`
	Public bool            `json:"public"`
	Files  map[string]file `json:"files"`
}

//PatchData is used to marshal the patch request
type PatchData struct {
	Desc  string          `json:"description"`
	Files map[string]file `json:"files"`
}

type file struct {
	Content string `json:"content"`
}

//Client is used to interact with the github api for gists
type Client struct {
	cfg *configuration.Configuration
}

//New creates a new Client
func New(cfg *configuration.Configuration) Client {
	return Client{cfg: cfg}
}

//List gets all the users gists
func (g *Client) List() ([]GetResponse, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "token "+g.cfg.Token)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	gists := []GetResponse{}
	err = json.Unmarshal(body, &gists)
	if err != nil {
		return nil, err
	}
	return gists, nil
}

//Post posts a new gist
func (g *Client) Post(description string, filesPath []string) error {
	files, err := createFiles(filesPath)
	if err != nil {
		return err
	}

	gist := PostData{
		Desc:   description,
		Public: !g.cfg.Private,
		Files:  files,
	}

	// encode json
	b, err := json.Marshal(gist)
	if err != nil {
		return err
	}

	// post json
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "token "+g.cfg.Token)
	_, err = client.Do(req)
	if err != nil {
		return err
	}
	return nil
}

//Update updates a gist
func (g *Client) Update(id string, filesPath []string) error {
	files, err := createFiles(filesPath)
	if err != nil {
		return err
	}

	gist := PatchData{
		Files: files,
	}

	// encode json
	b, err := json.Marshal(gist)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PATCH", url+"/"+id, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "token "+g.cfg.Token)
	_, err = client.Do(req)
	if err != nil {
		return err
	}
	return nil
}

//Delete deletes a gist
func (g *Client) Delete(id string) error {

	req, err := http.NewRequest("DELETE", url+"/"+id, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "token "+g.cfg.Token)
	_, err = client.Do(req)
	if err != nil {
		return err
	}
	return nil
}

//Download a gist
func (g *Client) Download(id string) error {
	r, err := g.Get(id)
	if err != nil {
		return err
	}
	for _, file := range r.Files {
		err := ioutil.WriteFile(file.FileName, []byte(file.Content), 0660)
		if err != nil {
			return err
		}
	}
	return nil
}

//ViewBrowser opens a browser for the gist
func (g *Client) ViewBrowser(id string) error {
	r, err := g.Get(id)
	if err != nil {
		return err
	}
	return r.OpenInBrowser()
}

//View prints the contents of a gist to stdout
func (g *Client) View(id string, name string) error {
	r, err := g.Get(id)
	if err != nil {
		return err
	}
	for _, file := range r.Files {
		if len(name) == 0 {
			fmt.Println(file.FileName + ":")
			fmt.Println(file.Content)
		} else {
			if file.FileName == name {
				fmt.Println(file.Content)
			}
		}
	}
	return nil
}

func (g *Client) Get(id string) (*GetResponse, error) {
	req, err := http.NewRequest("GET", url+"/"+id, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "token "+g.cfg.Token)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("could not download gists, code: %d, body: %s", res.StatusCode, body)
	}

	var r GetResponse
	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil

}

func createFiles(filesPath []string) (map[string]file, error) {
	files := make(map[string]file)
	for _, f := range filesPath {
		content, err := ioutil.ReadFile(f)
		if err != nil {
			return files, err
		}
		fileName := path.Base(f)
		files[fileName] = file{Content: string(content)}
	}
	return files, nil
}
