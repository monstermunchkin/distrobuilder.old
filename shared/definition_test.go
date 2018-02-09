package shared

import (
	"bytes"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"testing"

	yaml "gopkg.in/yaml.v2"
)

func TestSetDefinitionDefaults(t *testing.T) {
	def := getTestDefinition(t, "defaults.yaml")

	SetDefinitionDefaults(def)

	if def.Image.Arch != runtime.GOARCH {
		t.Fatalf("Expected image.arch to be '%s', got '%s'", runtime.GOARCH, def.Image.Arch)
	}

	if def.Image.Expiry != "30d" {
		t.Fatalf("Expected image.expiry to be '%s', got '%s'", "30d", def.Image.Expiry)
	}
}

func TestValidateDefinition(t *testing.T) {
	tests := []struct {
		name       string
		definition Definition
		expected   string
		shouldFail bool
	}{
		{
			"valid Definition",
			Definition{
				Image: DefinitionImage{
					Distribution: "ubuntu",
					Release:      "artful",
				},
				Source: DefinitionSource{
					Downloader: "debootstrap",
					URL:        "https://ubuntu.com",
				},
				Packages: DefinitionPackages{
					Manager: "apt",
				},
			},
			"",
			false,
		},
		{
			"empty image.distribution",
			Definition{},
			"image.distribution may not be empty",
			true,
		},
		{
			"empty image.release",
			Definition{
				Image: DefinitionImage{
					Distribution: "ubuntu",
				},
			},
			"image.release may not be empty",
			true,
		},
		{
			"invalid source.downloader",
			Definition{
				Image: DefinitionImage{
					Distribution: "ubuntu",
					Release:      "artful",
				},
				Source: DefinitionSource{
					Downloader: "foo",
				},
			},
			"source.downloader must be one of .+",
			true,
		},
		{
			"empty source.url",
			Definition{
				Image: DefinitionImage{
					Distribution: "ubuntu",
					Release:      "artful",
				},
				Source: DefinitionSource{
					Downloader: "debootstrap",
				},
			},
			"source.url may not be empty",
			true,
		},
		{
			"invalid package.manager",
			Definition{
				Image: DefinitionImage{
					Distribution: "ubuntu",
					Release:      "artful",
				},
				Source: DefinitionSource{
					Downloader: "debootstrap",
					URL:        "https://ubuntu.com",
				},
				Packages: DefinitionPackages{
					Manager: "foo",
				},
			},
			"packages.manager must be one of .+",
			true,
		},
	}

	for i, tt := range tests {
		log.Printf("Running test #%d: %s", i, tt.name)
		err := ValidateDefinition(tt.definition)
		if !tt.shouldFail && err != nil {
			t.Fatalf("Validation failed: %s", err)
		} else if tt.shouldFail {
			match, _ := regexp.MatchString(tt.expected, err.Error())
			if !match {
				t.Fatalf("Validation failed: Expected '%s', got '%s'", tt.expected, err.Error())
			}
		}
	}
}

func getTestDefinition(t *testing.T, fname string) *Definition {
	var (
		buf bytes.Buffer
		def Definition
	)

	f, err := os.Open(filepath.Join("..", "testdata", fname))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	_, err = io.Copy(&buf, f)
	if err != nil {
		t.Fatal(err)
	}

	err = yaml.Unmarshal(buf.Bytes(), &def)
	if err != nil {
		t.Fatal(err)
	}

	return &def
}
