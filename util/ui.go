package util

import tea "github.com/charmbracelet/bubbletea"

type Updatable[T any] interface {
	TypedUpdate(message tea.Msg) (T, tea.Cmd)
}

func UpdateChild[T Updatable[T]](child *T, message tea.Msg, commands *[]tea.Cmd) {
	updatedChild, command := (*child).TypedUpdate(message)
	*child = updatedChild
	*commands = append(*commands, command)
}

func UpdateUntypedChild(child *tea.Model, message tea.Msg, commands *[]tea.Cmd) {
	updatedChild, command := (*child).Update(message)
	*child = updatedChild
	*commands = append(*commands, command)
}

func UpdateSingleChild[M any, T Updatable[T]](model *M, child *T, message tea.Msg) (M, tea.Cmd) {
	updatedChild, command := (*child).TypedUpdate(message)
	*child = updatedChild
	return *model, command
}

func UpdateSingleUntypedChild[M any](model *M, child *tea.Model, message tea.Msg) (M, tea.Cmd) {
	updatedChild, command := (*child).Update(message)
	*child = updatedChild
	return *model, command
}
