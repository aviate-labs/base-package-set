package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

const (
	BASE_URL   = "https://api.github.com/repos/dfinity/motoko-base"
	DFX_URL    = "https://api.github.com/repos/dfinity/sdk"
	VESSEL_URL = "https://api.github.com/repos/dfinity/vessel"
	MOC_PREFIX = "moc-"
)

type Tag struct {
	Name   string
	Commit struct {
		Sha string
	}
}

type Tags []Tag

func (t Tag) MocVersion() string {
	return strings.TrimPrefix(t.Name, MOC_PREFIX)
}

func (t Tag) PackageName() string {
	if strings.HasPrefix(t.Name, MOC_PREFIX) {
		return fmt.Sprintf("base-%s", t.MocVersion())
	} else {
		return fmt.Sprintf("base-%s", t.Name)
	}
}

func getTags(url string) Tags {
	resp, err := http.Get(fmt.Sprintf("%s/tags", url))
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var tags Tags
	if err := json.Unmarshal(body, &tags); err != nil {
		log.Fatal(err)
	}
	if len(tags) == 0 {
		log.Fatal("No tags found.")
	}
	return tags
}

type Release struct {
	Url     string
	TagName string `json:"tag_name"`
}

type Releases []Release

func (rs Releases) getLatestStable() Release {
	for _, r := range rs {
		if !strings.Contains(r.TagName, "-") {
			return r
		}
	}
	return rs[0]
}

func getReleases(url string) Releases {
	resp, err := http.Get(fmt.Sprintf("%s/releases", url))
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var releases Releases
	if err := json.Unmarshal(body, &releases); err != nil {
		log.Fatal(err)
	}
	if len(releases) == 0 {
		log.Fatal("No releases found.")
	}
	return releases
}

func main() {
	tags := getTags(BASE_URL)
	t := template.Must(template.ParseGlob("templates/*"))
	{
		fileName := "package-set.dhall"
		file, err := os.OpenFile(fileName, os.O_RDWR, 0644)
		if err != nil {
			log.Fatal(err)
		}
		if err := t.ExecuteTemplate(file, fmt.Sprintf("%s.tmpl", fileName), struct {
			Tags []Tag
		}{
			Tags: tags,
		}); err != nil {
			log.Fatal(err)
		}
	}

	latest := tags[0]
	if !strings.HasPrefix(latest.Name, MOC_PREFIX) {
		log.Fatalf("script is outdated, expected tag to start with %q", MOC_PREFIX)
	}

	vesselCmd := exec.Command("vessel", "verify", "--version", latest.MocVersion())
	if msg, err := vesselCmd.CombinedOutput(); err != nil {
		log.Fatalf("vessel: %s %s %s", vesselCmd, msg, err)
	}

	{
		dhallCmd := exec.Command("dhall", "hash", "--file", "package-set.dhall")
		rawHash, err := dhallCmd.Output()
		if err != nil {
			log.Fatalf("dhall: %s %s", dhallCmd, err)
		}

		fileName := "README.md"
		file, err := os.OpenFile(fileName, os.O_RDWR, 0644)
		if err != nil {
			log.Fatal(err)
		}
		if err := t.ExecuteTemplate(file, fmt.Sprintf("%s.tmpl", fileName), struct {
			Latest Tag
			Hash   string
		}{
			Latest: latest,
			Hash:   strings.TrimPrefix(strings.TrimSpace(string(rawHash)), "sha256:"),
		}); err != nil {
			log.Fatal(err)
		}
	}
	{
		fileName := "package-set.yml"
		file, err := os.OpenFile(fmt.Sprintf(".github/workflows/%s", fileName), os.O_RDWR, 0644)
		if err != nil {
			log.Fatal(err)
		}

		dfxLatest := getReleases(DFX_URL).getLatestStable()
		vesselLatest := getReleases(VESSEL_URL).getLatestStable()

		if err := t.ExecuteTemplate(file, fmt.Sprintf("%s.tmpl", fileName), struct {
			Latest        Tag
			DfxVersion    string
			VesselVersion string
		}{
			Latest:        latest,
			DfxVersion:    dfxLatest.TagName,
			VesselVersion: strings.TrimPrefix(vesselLatest.TagName, "v"),
		}); err != nil {
			log.Fatal(err)
		}

		vesselCmd := exec.Command("vessel", "verify", "--version", latest.MocVersion())
		if msg, err := vesselCmd.CombinedOutput(); err != nil {
			log.Fatalf("vessel: %s %s %s", vesselCmd, msg, err)
		}
	}
}
