package gist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os/exec"
	"path"
	"runtime"

	"github.com/zbindenren/gost/configuration"
)

const (
	url = "https://api.github.com/gists"
)

//GetResponse is used to unmarshal the get response
type GetResponse struct {
	Comments    int              `json:"comments"`
	CommentsURL string           `json:"comments_url"`
	CommitsURL  string           `json:"commits_url"`
	CreatedAt   string           `json:"created_at"`
	Description string           `json:"description"`
	Files       map[string]*file `json:"files"`
	ForksURL    string           `json:"forks_url"`
	GitPullURL  string           `json:"git_pull_url"`
	GitPushURL  string           `json:"git_push_url"`
	HTMLURL     string           `json:"html_url"`
	ID          string           `json:"id"`
	Owner       owner            `json:"owner"`
	Public      bool             `json:"public"`
	UpdatedAt   string           `json:"updated_at"`
	URL         string           `json:"url"`
	User        interface{}      `json:"user"`
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

//Client is used to interact with the github api for gists
type Client struct {
	cfg *configuration.Configuration
	cl  *http.Client
}

//New creates a new Client
func New(cfg *configuration.Configuration) Client {
	return Client{cfg: cfg, cl: new(http.Client)}
}

//List gets all the users gists
func (g *Client) List() ([]GetResponse, error) {
	req, err := g.newRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	gists := []GetResponse{}
	err = g.do(req, &gists)
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

	gist := postData{Desc: description, Public: !g.cfg.Private, Files: files}

	// post json
	req, err := g.newRequest("POST", url, gist)
	if err != nil {
		return err
	}
	err = g.do(req, nil)
	if err != nil {
		return err
	}
	return nil
}

//Update updates a gist
func (g *Client) Update(id string, description string, filesPath []string) error {
	gist, err := g.Get(id)
	if err != nil {
		return err
	}
	files, err := createFiles(filesPath)
	if err != nil {
		return err
	}

	d := gist.Description
	if len(description) > 0 {
		d = description
	}
	patchdata := postData{
		Desc: d,
	}
	if len(filesPath) > 0 {
		patchdata.Files = files
	}
	return g.patch(id, patchdata)
}

//Delete deletes a gist
func (g *Client) Delete(id string, fileName string) error {

	if len(fileName) == 0 {
		req, err := g.newRequest("DELETE", url+"/"+id, nil)
		if err != nil {
			return err
		}
		err = g.do(req, nil)
		if err != nil {
			return err
		}
	} else {
		gist, err := g.Get(id)
		if err != nil {
			return err
		}
		patchdata := postData{}
		patchdata.Desc = gist.Description
		patchdata.Files = make(map[string]*file)
		for key, val := range gist.Files {
			if key == fileName {
				patchdata.Files[fileName] = nil
			} else {
				patchdata.Files[key] = new(file)
				patchdata.Files[key].Content = val.Content
			}
		}
		return g.patch(id, patchdata)
	}
	return nil
}

//Download a gist
func (g *Client) Download(id string, name string) error {
	r, err := g.Get(id)
	if err != nil {
		return err
	}
	for _, file := range r.Files {
		if len(name) == 0 {
			err := ioutil.WriteFile(file.FileName, []byte(file.Content), 0660)
			if err != nil {
				return err
			}
		} else {
			if file.FileName == name {
				err := ioutil.WriteFile(file.FileName, []byte(file.Content), 0660)
				if err != nil {
					return err
				}
			}
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
				fmt.Println(file.FileName + ":")
				fmt.Println(file.Content)
			}
		}
	}
	return nil
}

//Get fetches information for a gist with an id
func (g *Client) Get(id string) (*GetResponse, error) {
	req, err := g.newRequest("GET", url+"/"+id, nil)
	if err != nil {
		return nil, err
	}
	r := new(GetResponse)
	err = g.do(req, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

type file struct {
	FileName string `json:"filename,omitempty"`
	Type     string `json:"type,omitempty"`
	Language string `json:"language,omitempty"`
	RawURL   string `json:"raw_url,omitempty"`
	Size     int    `json:"size,omitempty"`
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

type postData struct {
	Desc   string           `json:"description"`
	Public bool             `json:"public,omitempty"`
	Files  map[string]*file `json:"files"`
}

func (g *Client) newRequest(method string, url string, postData interface{}) (*http.Request, error) {
	var bu io.Reader
	if postData != nil {
		b, err := json.Marshal(postData)
		if err != nil {
			return nil, err
		}
		bu = bytes.NewBuffer(b)
	}
	req, err := http.NewRequest(method, url, bu)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "token "+g.cfg.Token)
	return req, nil
}

func (g *Client) do(req *http.Request, responseData interface{}) error {
	res, err := g.cl.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode >= 300 {
		return fmt.Errorf("request failed, code=%d, message=%s", res.StatusCode, res.Body)
	}
	if responseData != nil {
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		err = json.Unmarshal(body, responseData)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *Client) patch(id string, gist postData) error {
	req, err := g.newRequest("PATCH", url+"/"+id, gist)
	if err != nil {
		return err
	}
	err = g.do(req, nil)
	if err != nil {
		return err
	}
	return nil
}

func createFiles(filesPath []string) (map[string]*file, error) {
	files := make(map[string]*file)
	for _, f := range filesPath {
		content, err := ioutil.ReadFile(f)
		if err != nil {
			return files, err
		}
		fileName := path.Base(f)
		files[fileName] = &file{Content: string(content)}
	}
	return files, nil
}
