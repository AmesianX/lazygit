package gui

import (
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands"
	"github.com/jesseduffield/lazygit/pkg/gui/presentation"
)

// list panel functions

func (gui *Gui) getSelectedReflogCommit() *commands.Commit {
	selectedLine := gui.State.Panels.ReflogCommits.SelectedLine
	if selectedLine == -1 || len(gui.State.ReflogCommits) == 0 {
		return nil
	}

	return gui.State.ReflogCommits[selectedLine]
}

func (gui *Gui) handleReflogCommitSelect(g *gocui.Gui, v *gocui.View) error {
	if gui.popupPanelFocused() {
		return nil
	}

	gui.State.SplitMainPanel = false

	if _, err := gui.g.SetCurrentView(v.Name()); err != nil {
		return err
	}

	gui.getMainView().Title = "Reflog Entry"

	commit := gui.getSelectedReflogCommit()
	if commit == nil {
		return gui.newStringTask("main", "No reflog history")
	}
	if err := gui.focusPoint(0, gui.State.Panels.ReflogCommits.SelectedLine, len(gui.State.ReflogCommits), v); err != nil {
		return err
	}

	cmd := gui.OSCommand.ExecutableFromString(
		gui.GitCommand.ShowCmdStr(commit.Sha),
	)
	if err := gui.newPtyTask("main", cmd); err != nil {
		gui.Log.Error(err)
	}

	return nil
}

func (gui *Gui) refreshReflogCommits() error {
	commits, err := gui.GitCommand.GetReflogCommits()
	if err != nil {
		return gui.createErrorPanel(gui.g, err.Error())
	}

	gui.State.ReflogCommits = commits

	if gui.getCommitsView().Context == "reflog-commits" {
		return gui.renderReflogCommitsWithSelection()
	}

	return nil
}

func (gui *Gui) renderReflogCommitsWithSelection() error {
	commitsView := gui.getCommitsView()

	gui.refreshSelectedLine(&gui.State.Panels.ReflogCommits.SelectedLine, len(gui.State.ReflogCommits))
	displayStrings := presentation.GetCommitListDisplayStrings(gui.State.ReflogCommits, gui.State.ScreenMode != SCREEN_NORMAL)
	gui.renderDisplayStrings(commitsView, displayStrings)
	if gui.g.CurrentView() == commitsView && commitsView.Context == "reflog-commits" {
		if err := gui.handleReflogCommitSelect(gui.g, commitsView); err != nil {
			return err
		}
	}

	return nil
}

func (gui *Gui) handleCheckoutReflogCommit(g *gocui.Gui, v *gocui.View) error {
	commit := gui.getSelectedReflogCommit()
	if commit == nil {
		return nil
	}

	err := gui.createConfirmationPanel(g, gui.getCommitsView(), true, gui.Tr.SLocalize("checkoutCommit"), gui.Tr.SLocalize("SureCheckoutThisCommit"), func(g *gocui.Gui, v *gocui.View) error {
		return gui.handleCheckoutRef(commit.Sha)
	}, nil)
	if err != nil {
		return err
	}

	gui.State.Panels.ReflogCommits.SelectedLine = 0

	return nil
}

func (gui *Gui) handleCreateReflogResetMenu(g *gocui.Gui, v *gocui.View) error {
	commit := gui.getSelectedReflogCommit()

	return gui.createResetMenu(commit.Sha)
}
