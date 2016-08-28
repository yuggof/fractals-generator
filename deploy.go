package main

import (
  "net/http"
  "os"
  "log"
  "fmt"
  "encoding/json"
  "io/ioutil"
  "strings"
  "bytes"
)

type release struct {
  ID int64 `json:"id"`
  TagName string `json:"tag_name"`
  TargetCommitish string `json:"target_commitish"`
}

func getGithubOauthToken() string {
  got := os.Getenv("GITHUB_OAUTH_TOKEN")

  if got == "" {
    log.Fatal("GITHUB_OAUTH_TOKEN env variable is empty.")
  }

  return got
}

func getVersion() string {
  bs, err := ioutil.ReadFile("VERSION")
  if err != nil {
    log.Fatal(err)
  }

  return strings.TrimSpace(string(bs))
}

func createRelease(githubOauthToken, version string) *release {
  r := release{
    TagName: fmt.Sprintf("staging-%s", version),
    TargetCommitish: "staging-release",
  }

  bs, err := json.Marshal(r)
  if err != nil {
    log.Fatal(err)
  }

  req, err := http.NewRequest("POST", "https://api.github.com/repos/yuggof/fractals-generator/releases", bytes.NewBuffer(bs))
  if err != nil {
    log.Fatal(err)
  }

  req.Header.Add("Authorization", fmt.Sprintf("token %s", githubOauthToken))

  res, err := http.DefaultClient.Do(req)
  if err != nil {
    log.Fatal(err)
  }
  defer res.Body.Close()

  if res.StatusCode != 201 {
    log.Fatal("Could not create a release.")
  }

  d := json.NewDecoder(res.Body)
  d.Decode(&r)

  return &r
}

func uploadBinary(githubOauthToken string, release *release) {
  bs, err := ioutil.ReadFile("generate-fractal")
  if err != nil {
    log.Fatal(err)
  }

  req, err := http.NewRequest(
    "POST",
    fmt.Sprintf("https://uploads.github.com/repos/yuggof/fractals-generator/releases/%d/assets?name=generate-fractal", release.ID),
    bytes.NewBuffer(bs),
  )
  if err != nil {
    log.Fatal(err)
  }

  req.Header.Add("Authorization", fmt.Sprintf("token %s", githubOauthToken))
  req.Header.Add("Content-Type", "text/plain")

  res, err := http.DefaultClient.Do(req)
  if err != nil {
    log.Fatal(err)
  }
  res.Body.Close()
}

func main() {
  got := getGithubOauthToken()
  v := getVersion()
  r := createRelease(got, v)
  uploadBinary(got, r)
}
