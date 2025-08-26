package tui

import (
	"context"
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/sinfirst/GophKeeper/internal/auth"
)

// AuthApp представляет TUI приложение для аутентификации
type AuthApp struct {
	app        *tview.Application
	pages      *tview.Pages
	authClient *auth.Client
	statusText *tview.TextView
}

// NewAuthApp создает новое приложение
func NewAuthApp() *AuthApp {
	if os.Getenv("TERM") == "" {
		os.Setenv("TERM", "xterm-256color")
	}
	return &AuthApp{
		app: tview.NewApplication(),
	}
}

// Run запускает приложение
func (a *AuthApp) Run() error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Panic occurred: %v\n", r)
		}
	}()

	a.authClient = auth.NewClient("localhost:50051")
	a.setupUI()

	return a.app.SetRoot(a.pages, true).Run()
}

// setupUI настраивает пользовательский интерфейс
func (a *AuthApp) setupUI() {
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
func (a *AuthApp) createMainMenu() *tview.Flex {
	list := tview.NewList()
	list.SetBorder(true).SetTitle("GophKeeper - GRPC Server")

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

// createRegisterPage создает страницу регистрации
func (a *AuthApp) createRegisterPage() *tview.Flex {
	form := tview.NewForm()
	form.SetBorder(true).SetTitle(" Регистрация ")

	var username, password string

	form.AddInputField("Логин", "", 20, nil, func(text string) {
		username = text
	})

	form.AddPasswordField("Пароль", "", 20, '*', func(text string) {
		password = text
	})

	// Кнопки
	form.AddButton("Зарегистрировать", func() {
		if username == "" || password == "" {
			a.showStatus("Ошибка: заполните все поля")
			return
		}
		ctx := context.Background()
		response, err := a.authClient.Register(ctx, username, password)

		if err != nil {
			a.showStatus(fmt.Sprintf("Ошибка регистрации: %v", err))
			return
		}

		if response.Success {
			a.showStatus(fmt.Sprintf("Успешная регистрация! UserID: %s", response.UserId))
			a.pages.SwitchToPage("menu")
		} else {
			a.showStatus(fmt.Sprintf("Ошибка: %s", response.Message))
		}
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

// createHeader создает заголовок приложения
func (a *AuthApp) createHeader() *tview.TextView {
	header := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText("GophKeeper - GRPC Client").
		SetDynamicColors(true)

	header.SetBackgroundColor(tcell.ColorDarkBlue)
	return header
}

// createStatusBar создает строку статуса
func (a *AuthApp) createStatusBar() *tview.Flex {
	statusBar := tview.NewFlex().
		AddItem(a.statusText, 0, 1, false)

	statusBar.SetBackgroundColor(tcell.ColorDarkGray)
	return statusBar
}

// showStatus показывает статусное сообщение
func (a *AuthApp) showStatus(message string) {
	a.app.QueueUpdateDraw(func() {
		a.statusText.SetText(fmt.Sprintf(" Статус: %s", message))
	})
}

// globalInputHandler обрабатывает глобальные горячие клавиши
func (a *AuthApp) globalInputHandler(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyCtrlC:
		a.app.Stop()
		return nil
	case tcell.KeyEsc:
		// Возврат в главное меню с любой страницы
		if a.pages.HasPage("menu") {
			a.pages.SwitchToPage("menu")
		}
		return nil
	}
	return event
}
