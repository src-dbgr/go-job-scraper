package openai

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

type FilePromptRepository struct {
	baseDir string
}

func NewFilePromptRepository() *FilePromptRepository {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}
	// Navigate up to the project root and then to the prompts directory
	baseDir := filepath.Join(filepath.Dir(filename), "..", "..", "..", "prompts")
	return &FilePromptRepository{baseDir: baseDir}
}

func (r *FilePromptRepository) GetPrompt(name string) (string, error) {
	filename := filepath.Join(r.baseDir, name+".txt")
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("error reading prompt file: %w", err)
	}
	return string(content), nil
}
