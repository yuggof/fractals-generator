package main

import (
  "os/exec"
  "net/http"
  "os"
  "log"
  "fmt"
  "encoding/json"
  "io/ioutil"
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

func getTagName() string {
  tn := os.Getenv("CIRCLE_BUILD_NUM")

  if tn == "" {
    log.Fatal("CIRCLE_BUILD_NUM env variable is empty.")
  }

  return tn
}

func buildBinary() {
  err := exec.Command("go", "build", "-o", "generate-fractal", "generator.go").Run()
  if err != nil {
    log.Fatal(err)
  }
}

func createRelease(githubOauthToken string) *release {
  r := release{TagName: getTagName()}

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
    bs, err = ioutil.ReadAll(res.Body)
    if err != nil {
      log.Fatal(err)
    }

    log.Fatalf(
      "Could not create a relase.\nResponse status: %d.\nResponse body: %s",
      res.StatusCode,
      string(bs),
    )
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
  buildBinary()
  got := getGithubOauthToken()
  r := createRelease(got)
  uploadBinary(got, r)
}
