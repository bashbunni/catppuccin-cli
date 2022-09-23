package ui

import (
	"fmt"
	"log"
	"time"

	"github.com/catppuccin/cli/internal/utils"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/go-git/go-git/v5"
)

/*
CURRENT STRUGGLE:
hard to get clone progress. sideband.Progress is mostly useless for this
https://github.com/src-d/go-git/issues/549
*/

type progressMsg float64 // Progress

type finishMsg string // Finish

type ProgressWriter struct {
	// total bytes downloaded
	downloaded int
	// total size of the file(s)
	total             int
	calculateProgress func(float64)
}

// make GitProgress implement io.Writer so we can store the git clone progress
// to this struct
func (g *ProgressWriter) Write(p []byte) (n int, err error) {
	log.Println(string(p))
	g.downloaded += len(p)
	if g.total > 0 {
		g.calculateProgress(float64(g.downloaded) / float64(g.total))
	}
	return len(p), nil
}

func NewProgressBar() ProgressWrapper {
	return ProgressWrapper{
		progress: progress.New(),
	}
}

// I feel like this should be running in a go routine and sending data to the
// BBT TUI with messages

// CloneRepo clones a repo into the specified location.
func CloneRepo(stagePath string, repo string) string {
	log.Print("in clone repo")
	org := utils.GetEnv("ORG_OVERRIDE", "catppuccin")
	gitProgress := &ProgressWriter{
		downloaded: 0,
		calculateProgress: func(ratio float64) {
			log.Printf("in calculate progress: %f", ratio)
			if p != nil {
				p.Send(progressMsg(ratio))
			}
		},
	}
	// idk how this works, I think it is writing to gitProgress each time...
	// maybe include p.Send in Write? already doing in calculateProgress
	_, err := git.PlainClone(stagePath, false, &git.CloneOptions{
		URL:      fmt.Sprintf("https://github.com/%s/%s.git", org, repo),
		Progress: gitProgress,
	})
	if err != nil {
		fmt.Println(err)
	}
	// TODO: I think this might need to be a command that can get reused
	return stagePath
}

func StartClone(repo string) tea.Cmd {
	return func() tea.Msg {
		log.Print("in start clone")
		CloneRepo(utils.GetTemplateDir(repo), "template")
		return nil
	}
}

func finalPause() tea.Cmd {
	return tea.Tick(time.Second*10, func(time.Time) tea.Msg {
		return tea.Quit
	})
}

// TODO: include this somewhere it will be seen
type ProgressWrapper struct {
	progress progress.Model
}

func (m ProgressWrapper) Init() tea.Cmd {
	return StartClone(RepoName)
}

func (m ProgressWrapper) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case progressMsg:
		var cmds []tea.Cmd
		if msg >= 1.0 {
			cmds = append(cmds, tea.Sequentially(finalPause(), tea.Quit))
		}
		cmds = append(cmds, m.progress.SetPercent(float64(msg)))
		return m, tea.Batch(cmds...)
	case tea.KeyMsg:
		return m, tea.Quit // Just quit
	case progress.FrameMsg:
		// Update bar
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd
	case finishMsg:
		return m, finalPause()
	}
	return m, nil
}

func (m ProgressWrapper) View() string {
	return lipgloss.JoinVertical(lipgloss.Left, "Downloading...", m.progress.View(), "Press any key to quit")
}
