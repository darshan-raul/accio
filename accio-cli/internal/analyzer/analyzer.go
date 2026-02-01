package analyzer

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type ProjectAnalysis struct {
	Language     string   `json:"language"`
	Framework    string   `json:"framework"`
	Dependencies []string `json:"dependencies"`
	Database     string   `json:"database"`
	Cache        string   `json:"cache"`
	Files        []string `json:"files"`
}

func Analyze(root string) (*ProjectAnalysis, error) {
	analysis := &ProjectAnalysis{
		Dependencies: []string{},
		Files:        []string{},
	}

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == ".git" || d.Name() == "node_modules" || d.Name() == "vendor" || d.Name() == "venv" {
				return filepath.SkipDir
			}
			return nil
		}

		relPath, _ := filepath.Rel(root, path)
		analysis.Files = append(analysis.Files, relPath)

		// Naive detection logic
		if d.Name() == "go.mod" {
			analysis.Language = "Go"
			// Parse go.mod for deps
		}
		if d.Name() == "package.json" {
			analysis.Language = "JavaScript/TypeScript"
		}
		if d.Name() == "requirements.txt" || d.Name() == "Pipfile" {
			analysis.Language = "Python"
		}
		if strings.HasSuffix(d.Name(), ".py") && analysis.Language == "" {
			analysis.Language = "Python"
		}
		if strings.HasSuffix(d.Name(), ".go") && analysis.Language == "Python" {
			// Mixed?
		}
		// Check for framework specific files
		if d.Name() == "main.go" {
			content, _ := os.ReadFile(path)
			s := string(content)
			if strings.Contains(s, "github.com/gin-gonic/gin") {
				analysis.Framework = "Gin"
			}
			if strings.Contains(s, "gorm.io/gorm") {
				analysis.Dependencies = append(analysis.Dependencies, "GORM")
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return analysis, nil
}

func (pa *ProjectAnalysis) ToJSON() (string, error) {
	b, err := json.MarshalIndent(pa, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
