package llm

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type OllamaClient struct {
	host       string
	model      string
	httpClient *http.Client
}

type ollamaGenerateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type ollamaGenerateResponse struct {
	Response string `json:"response"`
	Error    string `json:"error"`
}

type TemplateSelection struct {
	Templates []string          `json:"templates"`
	Params    map[string]string `json:"params"`
}

func NewOllamaClient(host, model string, timeout time.Duration) *OllamaClient {
	if host == "" {
		host = "http://localhost:11434"
	}
	if model == "" {
		model = "llama3.2"
	}
	return &OllamaClient{
		host:       strings.TrimRight(host, "/"),
		model:      model,
		httpClient: &http.Client{Timeout: timeout},
	}
}

func (c *OllamaClient) Name() string {
	return "ollama"
}

func (c *OllamaClient) IsAvailable() bool {
	request, err := http.NewRequest(http.MethodGet, c.host+"/api/tags", nil)
	if err != nil {
		return false
	}
	response, err := c.httpClient.Do(request)
	if err != nil {
		return false
	}
	defer func() { _ = response.Body.Close() }()
	return response.StatusCode >= 200 && response.StatusCode < 500
}

func (c *OllamaClient) Generate(prompt string) (string, error) {
	body, err := json.Marshal(ollamaGenerateRequest{
		Model:  c.model,
		Prompt: prompt,
		Stream: false,
	})
	if err != nil {
		return "", err
	}

	request, err := http.NewRequest(http.MethodPost, c.host+"/api/generate", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	request.Header.Set("Content-Type", "application/json")

	response, err := c.httpClient.Do(request)
	if err != nil {
		return "", err
	}
	defer func() { _ = response.Body.Close() }()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return "", fmt.Errorf("ollama returned %d: %s", response.StatusCode, strings.TrimSpace(string(responseBody)))
	}

	var decoded ollamaGenerateResponse
	if err := json.Unmarshal(responseBody, &decoded); err != nil {
		return "", err
	}
	if decoded.Error != "" {
		return "", errors.New(decoded.Error)
	}
	return decoded.Response, nil
}

func ParseTemplateSelection(raw string) (TemplateSelection, error) {
	var selection TemplateSelection
	if err := json.Unmarshal([]byte(raw), &selection); err != nil {
		return TemplateSelection{}, err
	}
	if len(selection.Templates) == 0 {
		return TemplateSelection{}, fmt.Errorf("no templates returned")
	}
	if selection.Params == nil {
		selection.Params = map[string]string{}
	}
	return selection, nil
}

func GenerateTemplateSelection(client Client, prompt string, allowedTemplates []string) (TemplateSelection, error) {
	if client == nil || !client.IsAvailable() {
		return TemplateSelection{}, fmt.Errorf("LLM backend not available")
	}

	requestPrompt := prompt
	for attempt := 0; attempt < 2; attempt++ {
		raw, err := client.Generate(requestPrompt)
		if err != nil {
			return TemplateSelection{}, err
		}

		selection, parseErr := ParseTemplateSelection(strings.TrimSpace(raw))
		if parseErr == nil {
			if err := validateSelection(selection, allowedTemplates); err != nil {
				return TemplateSelection{}, err
			}
			return selection, nil
		}

		requestPrompt = prompt + "\nReturn valid JSON only with templates and params."
		if attempt == 1 {
			return TemplateSelection{}, parseErr
		}
	}

	return TemplateSelection{}, fmt.Errorf("unable to generate template selection")
}

func validateSelection(selection TemplateSelection, allowedTemplates []string) error {
	allowed := make(map[string]struct{}, len(allowedTemplates))
	for _, template := range allowedTemplates {
		allowed[template] = struct{}{}
	}
	for _, template := range selection.Templates {
		if _, ok := allowed[template]; !ok {
			return fmt.Errorf("unknown template %q", template)
		}
	}
	return nil
}
