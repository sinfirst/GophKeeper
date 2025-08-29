package tui

import (
	"context"
	"fmt"

	"github.com/rivo/tview"
)

// createRegisterPage создает страницу регистрации
func (a *TUI) createRegisterPage() *tview.Flex {
	form := tview.NewForm()
	form.SetBorder(true).SetTitle(" Регистрация ")

	var username, password string

	form.AddInputField("Логин", "", 20, nil, func(text string) {
		username = text
	})

	form.AddPasswordField("Пароль", "", 20, '*', func(text string) {
		password = text
	})

	form.AddButton("Зарегистрировать", func() {
		if username == "" || password == "" {
			a.showStatus("Ошибка: заполните все поля")
			return
		}
		ctx := context.Background()
		err := a.client.Register(ctx, username, password)

		if err != nil {
			a.showStatus(fmt.Sprintf("Ошибка регистрации: %v", err))
			return
		}

		a.showStatus("Успешная регистрация!")
		a.pages.SwitchToPage("menu")
	})

	form.AddButton("Назад", func() {
		a.pages.SwitchToPage("menu")
	})

	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(a.createHeader(), 3, 1, false).
		AddItem(form, 0, 1, true).
		AddItem(a.createStatusBar(), 2, 1, false)

	return layout
}
