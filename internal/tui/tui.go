package tui

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/sinfirst/GophKeeper/internal/client"
	"github.com/sinfirst/GophKeeper/internal/config"
)

// AuthApp представляет TUI приложение для аутентификации
type TUI struct {
	app        *tview.Application
	pages      *tview.Pages
	client     *client.Client
	statusText *tview.TextView
}

// NewAuthApp создает новое приложение
func NewApp() *TUI {
	if os.Getenv("TERM") == "" {
		os.Setenv("TERM", "xterm-256color")
	}
	return &TUI{
		app: tview.NewApplication(),
	}
}

// Run запускает приложение
func (a *TUI) Run() error {
	defer func() {
		if r := recover(); r != nil {
			panic(r)
		}
	}()
	config := config.NewConfig()
	a.client = client.NewClient(config.Host)
	a.setupUI()

	return a.app.SetRoot(a.pages, true).Run()
}

// setupUI настраивает пользовательский интерфейс
func (a *TUI) setupUI() {
	a.pages = tview.NewPages()
	a.statusText = tview.NewTextView().SetTextAlign(tview.AlignLeft)

	mainMenu := a.createMainMenu()

	registerPage := a.createRegisterPage()

	a.pages.AddPage("menu", mainMenu, true, true)
	a.pages.AddPage("register", registerPage, true, false)

	a.app.SetRoot(a.pages, true)

	a.app.SetInputCapture(a.globalInputHandler)
}

// createMainMenu создает главное меню
func (a *TUI) createMainMenu() *tview.Flex {
	list := tview.NewList()
	list.SetBorder(true).SetTitle(" Главное меню ")

	list.AddItem("Регистрация", "Создать новый аккаунт", '1', func() {
		a.pages.SwitchToPage("register")
	})

	list.AddItem("Вход", "Войти в существующий аккаунт", '2', func() {
		a.showStatus("Функция входа будет реализована позже")
	})

	list.AddItem("Выход", "Завершить программу", 'q', func() {
		a.app.Stop()
	})

	layout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(a.createHeader(), 3, 1, false).
		AddItem(list, 0, 1, true).
		AddItem(a.createStatusBar(), 2, 1, false)

	return layout
}

// createHeader создает заголовок приложения
func (a *TUI) createHeader() *tview.TextView {
	header := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText("GophKeeper - GRPC Client").
		SetDynamicColors(true)

	header.SetBackgroundColor(tcell.ColorDarkBlue)
	return header
}

// createStatusBar создает строку статуса
func (a *TUI) createStatusBar() *tview.Flex {
	statusBar := tview.NewFlex().
		AddItem(a.statusText, 0, 1, false)

	statusBar.SetBackgroundColor(tcell.ColorDarkGray)
	return statusBar
}

// showStatus показывает статусное сообщение
func (a *TUI) showStatus(message string) {
	a.app.QueueUpdateDraw(func() {
		a.statusText.SetText(fmt.Sprintf(" Статус: %s", message))
	})
}

// globalInputHandler обрабатывает глобальные горячие клавиши
func (a *TUI) globalInputHandler(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyCtrlC:
		a.app.Stop()
		return nil
	case tcell.KeyEsc:
		if a.pages.HasPage("menu") {
			a.pages.SwitchToPage("menu")
		}
		return nil
	}
	return event
}
