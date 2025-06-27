package activity

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/charmbracelet/lipgloss"
)

type UserActivity struct {
	Type      string `json:"type"`
	Repo      Repo   `json:"repo"`
	CreatedAt string `json:"created_at"`
	Payload   struct {
		Action  string `json:"action"`
		Ref     string `json:"ref"`
		RefType string `json:"ref_type"`
		Commits []struct {
			Message string `json:"message"`
		} `json:"commits"`
	} `json:"payload"`
}

type Repo struct {
	Name string `json:"name"`
}

func FetchActivity(username string) ([]UserActivity, error) {
	res, err := http.Get(fmt.Sprintf("https://api.github.com/users/%s/events", username))
	if err != nil {
		return nil, err
	}

	if res.StatusCode == 404 {
		return nil, fmt.Errorf("no se encontró el usuario: %s", username)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error en el fetching de datos : %d", res.StatusCode)
	}

	var activities []UserActivity
	if err := json.NewDecoder(res.Body).Decode(&activities); err != nil {
		return nil, err
	}

	return activities, nil
}

func DisplayActivity(username string, events []UserActivity) error {
	if len(events) == 0 {
		return fmt.Errorf("no se encontraron eventos para el usuario: %s", username)
	}
	fmt.Println(
		lipgloss.NewStyle().
			Bold(true).
			Padding(2).
			MarginBottom(2).
			Foreground(lipgloss.Color("#0cad85")).
			Border(lipgloss.NormalBorder(), true, true, true, true).
			Render(fmt.Sprintf("🚀 %s's recent activities on Github", username)),
	)
	for _, event := range events {
		parsedDate, err := time.Parse(time.RFC3339, event.CreatedAt)
		if err != nil {
			return fmt.Errorf("error al parsear la fecha: %v", err)
		}

		var action string
		switch event.Type {
		case "PushEvent":
			commitCount := len(event.Payload.Commits)
			action = fmt.Sprintf("Pushed %d commit(s) to %s - %s", commitCount, event.Repo.Name, parsedDate)
		case "IssuesEvent":
			action = fmt.Sprintf("%s an issue in %s - %s", event.Payload.Action, event.Repo.Name, parsedDate)
		case "WatchEvent":
			action = fmt.Sprintf("Starred %s - %s", event.Repo.Name, parsedDate)
		case "ForkEvent":
			action = fmt.Sprintf("Forked %s - %s", event.Repo.Name, parsedDate)
		case "CreateEvent":
			action = fmt.Sprintf("Created %s in %s - %s", event.Payload.RefType, event.Repo.Name, parsedDate)
		default:
			action = fmt.Sprintf("%s in %s - %s", event.Type, event.Repo.Name, parsedDate)
		}

		actionStyle := lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(lipgloss.Color("#3C3C3C")).
			Render(fmt.Sprintf("✅ %s", action))
		fmt.Println(actionStyle)
	}

	return nil
}
