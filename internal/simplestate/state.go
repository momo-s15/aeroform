package simplestate

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

const statePath = ".aeroform/simple_state.json"

type Project struct {
	Name            string    `json:"name"`
	Provider        string    `json:"provider"`
	Template        string    `json:"template"`
	Prompt          string    `json:"prompt"`
	CustomDomain    string    `json:"custom_domain,omitempty"`
	MonthlyEstimate float64   `json:"monthly_estimate"`
	WorkDir         string    `json:"work_dir,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

type State struct {
	Projects []Project `json:"projects"`
}

func Load() (State, error) {
	data, err := os.ReadFile(statePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return State{}, nil
		}
		return State{}, err
	}

	var st State
	if len(data) == 0 {
		return State{}, nil
	}
	if err := json.Unmarshal(data, &st); err != nil {
		return State{}, err
	}
	return st, nil
}

func Save(st State) error {
	if err := os.MkdirAll(filepath.Dir(statePath), 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(st, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(statePath, data, 0o644)
}

func AddProject(project Project) error {
	st, err := Load()
	if err != nil {
		return err
	}

	for i := range st.Projects {
		if st.Projects[i].Name == project.Name {
			st.Projects[i] = project
			return Save(st)
		}
	}

	st.Projects = append(st.Projects, project)
	return Save(st)
}

func Clear() error {
	return Save(State{})
}
