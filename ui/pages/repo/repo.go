package repo

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	ggit "github.com/charmbracelet/soft-serve/git"
	"github.com/charmbracelet/soft-serve/ui/common"
	"github.com/charmbracelet/soft-serve/ui/components/code"
	"github.com/charmbracelet/soft-serve/ui/components/selector"
	"github.com/charmbracelet/soft-serve/ui/components/statusbar"
	"github.com/charmbracelet/soft-serve/ui/components/tabs"
	"github.com/charmbracelet/soft-serve/ui/git"
)

type tab int

const (
	readmeTab tab = iota
	filesTab
	commitsTab
	branchesTab
	tagsTab
)

// RepoMsg is a message that contains a git.Repository.
type RepoMsg git.GitRepo

// RefMsg is a message that contains a git.Reference.
type RefMsg *ggit.Reference

// Repo is a view for a git repository.
type Repo struct {
	common       common.Common
	rs           git.GitRepoSource
	selectedRepo git.GitRepo
	activeTab    tab
	tabs         *tabs.Tabs
	statusbar    *statusbar.StatusBar
	readme       *code.Code
	log          *Log
	ref          *ggit.Reference
}

// New returns a new Repo.
func New(common common.Common, rs git.GitRepoSource) *Repo {
	sb := statusbar.New(common)
	tb := tabs.New(common, []string{"Readme", "Files", "Commits", "Branches", "Tags"})
	readme := code.New(common, "", "")
	readme.NoContentStyle = readme.NoContentStyle.SetString("No readme found.")
	log := NewLog(common)
	r := &Repo{
		common:    common,
		rs:        rs,
		tabs:      tb,
		statusbar: sb,
		readme:    readme,
		log:       log,
	}
	return r
}

// SetSize implements common.Component.
func (r *Repo) SetSize(width, height int) {
	r.common.SetSize(width, height)
	hm := r.common.Styles.RepoBody.GetVerticalFrameSize() +
		r.common.Styles.RepoHeader.GetHeight() +
		r.common.Styles.RepoHeader.GetVerticalFrameSize() +
		r.common.Styles.StatusBar.GetHeight() +
		r.common.Styles.Tabs.GetHeight()
	r.tabs.SetSize(width, height-hm)
	r.statusbar.SetSize(width, height-hm)
	r.readme.SetSize(width, height-hm)
	r.log.SetSize(width, height-hm)
}

// ShortHelp implements help.KeyMap.
func (r *Repo) ShortHelp() []key.Binding {
	b := make([]key.Binding, 0)
	tab := r.common.Keymap.Section
	tab.SetHelp("tab", "switch tab")
	b = append(b, r.common.Keymap.Back)
	b = append(b, tab)
	return b
}

// FullHelp implements help.KeyMap.
func (r *Repo) FullHelp() [][]key.Binding {
	b := make([][]key.Binding, 0)
	return b
}

// Init implements tea.View.
func (r *Repo) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (r *Repo) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)
	switch msg := msg.(type) {
	case selector.SelectMsg:
		r.activeTab = 0
		cmds = append(cmds, r.tabs.Init(), r.setRepoCmd(string(msg)))
	case RepoMsg:
		r.selectedRepo = git.GitRepo(msg)
		r.readme.GotoTop()
		cmds = append(cmds,
			r.updateReadmeCmd,
			r.updateRefCmd,
		)
	case RefMsg:
		r.ref = msg
		cmds = append(cmds,
			r.updateStatusBarCmd,
			r.log.Init(),
		)
	case tabs.ActiveTabMsg:
		r.activeTab = tab(msg)
	case tea.KeyMsg, tea.MouseMsg:
		if r.selectedRepo != nil {
			cmds = append(cmds, r.updateStatusBarCmd)
		}
	}
	t, cmd := r.tabs.Update(msg)
	r.tabs = t.(*tabs.Tabs)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}
	s, cmd := r.statusbar.Update(msg)
	r.statusbar = s.(*statusbar.StatusBar)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}
	switch r.activeTab {
	case readmeTab:
		b, cmd := r.readme.Update(msg)
		r.readme = b.(*code.Code)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	case filesTab:
	case commitsTab:
		l, cmd := r.log.Update(msg)
		r.log = l.(*Log)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	case branchesTab:
	case tagsTab:
	}
	return r, tea.Batch(cmds...)
}

// View implements tea.Model.
func (r *Repo) View() string {
	s := r.common.Styles.Repo.Copy().
		Width(r.common.Width).
		Height(r.common.Height)
	repoBodyStyle := r.common.Styles.RepoBody.Copy()
	hm := repoBodyStyle.GetVerticalFrameSize() +
		r.common.Styles.RepoHeader.GetHeight() +
		r.common.Styles.RepoHeader.GetVerticalFrameSize() +
		r.common.Styles.StatusBar.GetHeight() +
		r.common.Styles.Tabs.GetHeight()
	mainStyle := repoBodyStyle.
		Height(r.common.Height - hm)
	main := mainStyle.Render("")
	switch r.activeTab {
	case readmeTab:
		main = mainStyle.Render(r.readme.View())
	case filesTab:
	case commitsTab:
		main = mainStyle.Render(r.log.View())
	}
	view := lipgloss.JoinVertical(lipgloss.Top,
		r.headerView(),
		main,
		r.statusbar.View(),
	)
	return s.Render(view)
}

func (r *Repo) headerView() string {
	if r.selectedRepo == nil {
		return ""
	}
	name := r.common.Styles.RepoHeaderName.Render(r.selectedRepo.Name())
	style := r.common.Styles.RepoHeader.Copy().Width(r.common.Width)
	return style.Render(
		lipgloss.JoinVertical(lipgloss.Top,
			name,
			r.tabs.View(),
		),
	)
}

func (r *Repo) setRepoCmd(repo string) tea.Cmd {
	return func() tea.Msg {
		for _, r := range r.rs.AllRepos() {
			if r.Name() == repo {
				return RepoMsg(r)
			}
		}
		return common.ErrorMsg(git.ErrMissingRepo)
	}
}

func (r *Repo) updateStatusBarCmd() tea.Msg {
	info := ""
	switch r.activeTab {
	case readmeTab:
		info = fmt.Sprintf("%.f%%", r.readme.ScrollPercent()*100)
	}
	return statusbar.StatusBarMsg{
		Key:    r.selectedRepo.Name(),
		Value:  "",
		Info:   info,
		Branch: fmt.Sprintf(" %s", r.ref.Name().Short()),
	}
}

func (r *Repo) updateReadmeCmd() tea.Msg {
	if r.selectedRepo == nil {
		return common.ErrorCmd(git.ErrMissingRepo)
	}
	rm, rp := r.selectedRepo.Readme()
	return r.readme.SetContent(rm, rp)
}

func (r *Repo) updateRefCmd() tea.Msg {
	head, err := r.selectedRepo.HEAD()
	if err != nil {
		return common.ErrorMsg(err)
	}
	return RefMsg(head)
}
