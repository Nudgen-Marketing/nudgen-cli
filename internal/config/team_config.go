package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type CLIConfig struct {
	ActiveTeamID   string `json:"active_team_id"`
	ActiveTeamName string `json:"active_team_name"`
}

func GetConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".nudgen-config.json"
	}
	dir := filepath.Join(home, ".nudgen")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)
	}
	return filepath.Join(dir, "config.json")
}

func SaveActiveTeam(id, name string) error {
	cfg := CLIConfig{
		ActiveTeamID:   id,
		ActiveTeamName: name,
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(GetConfigPath(), data, 0644)
}

func GetActiveTeam() (*CLIConfig, error) {
	path := GetConfigPath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &CLIConfig{}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg CLIConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func ClearActiveTeam() error {
	return os.Remove(GetConfigPath())
}

func PrintActiveTeam() {
	cfg, err := GetActiveTeam()
	if err != nil || cfg.ActiveTeamID == "" {
		fmt.Println("No active team set. Run `nudgen teams switch` to set one.")
		return
	}
	fmt.Printf("Active Team: %s (%s)\n", cfg.ActiveTeamName, cfg.ActiveTeamID)
}
