package tui

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/sinfirst/GophKeeper/internal/client"
	"github.com/sinfirst/GophKeeper/internal/models"
)

const (
	menu     string = "1. Регистрация \n2. Вход в аккаунт \n3. Сохранение данных \n4. Извлечение данных \n5. Лист данных \n6. Обновление данных \n7. Удаление данных \n8. Получение версии программы \n0. Выход из программы\n"
	typeData string = "1. Пара логин-пароль \n2. Текстовые данные \n3. Банковская карта \n4. Бинарные данные \n 0. Назад \n"
)

type TUI struct {
	Client *client.Client
}

func StartTUI(tui TUI) {
	var chooseInt int = 1
	var choose string
	for chooseInt != 0 {
		fmt.Print(menu)
		fmt.Print("Введите число: ")
		fmt.Scan(&choose)
		chooseInt, err := strconv.Atoi(choose)
		if err != nil {
			fmt.Println("Введите число, а не строку!")
			continue
		}

		switch chooseInt {
		case 1:
			tui.auth(1)
		case 2:
			tui.auth(2)
		case 3:
			tui.store()
		case 4:
			tui.retrieve()
		case 5:
			tui.listData()
		case 6:
			tui.updateData()
		case 7:
			tui.deleteData()
		case 8:
			tui.getVersion()
		default:
			fmt.Println("Число не входит в пункты меню!")
		}
	}
	fmt.Println("Счастливо!")
	tui.Client.Close()
}

func (t *TUI) auth(typeAuth int) {
	var username, password string
	var err error
	fmt.Print("Введите логин:")
	fmt.Scan(&username)
	fmt.Print("Введите пароль:")
	fmt.Scan(&password)

	if typeAuth == 1 {
		err = t.Client.Register(context.Background(), username, password)
	} else {
		err = t.Client.Login(context.Background(), username, password)
	}

	if err != nil {
		fmt.Println("Ошибка: ", err)
		return
	}

	fmt.Println("Успешно!")
}

func (t *TUI) store() {
	var chooseInt int = 1
	var choose string
	for chooseInt != 0 {
		fmt.Print(typeData)
		fmt.Print("Введите число: ")
		fmt.Scan(&choose)
		chooseInt, err := strconv.Atoi(choose)
		if err != nil {
			fmt.Println("Введите число, а не строку!")
		}
		switch chooseInt {
		case 1:
			jsonReq, meta, err := separateDataByTypeToInput(models.Login)
			if err != nil {

			}
			id, err := t.Client.StoreData(context.Background(), models.Login, meta, jsonReq)
			if err != nil {
				fmt.Println("Ошибка: ", err)
				continue
			}
			fmt.Println("Успешено ID сохраненных данных: ", id)
			return
		case 2, 4:
			var id int
			var err error
			req, meta, err := separateDataByTypeToInput(models.Text)
			if err != nil {
				continue
			}
			if chooseInt == 2 {
				id, err = t.Client.StoreData(context.Background(), models.Text, meta, req)
			} else {
				id, err = t.Client.StoreData(context.Background(), models.Binary, meta, req)
			}
			if err != nil {
				fmt.Println("Ошибка: ", err)
				continue
			}
			fmt.Println("Успешено ID сохраненных данных: ", id)
			return
		case 3:
			jsonReq, meta, err := separateDataByTypeToInput(models.Card)
			if err != nil {
				continue
			}
			id, err := t.Client.StoreData(context.Background(), models.Card, meta, jsonReq)
			if err != nil {
				fmt.Println("Ошибка: ", err)
				continue
			}
			fmt.Println("Успешено ID сохраненных данных: ", id)
			return
		default:
			fmt.Println("Число не входит в пункты типа данных!")
		}
	}
}

func (t *TUI) retrieve() {
	var id string
	fmt.Print("Введите id данных: ")
	fmt.Scan(&id)
	record, err := t.Client.RetrieveData(context.Background(), id)
	if err != nil {
		fmt.Println("Ошибка: ", err)
		return
	}
	err = separateDataByTypeToOutput(record)
	if err != nil {
		fmt.Println("Ошибка, попробуйте еще раз")
		return
	}
}

func (t *TUI) listData() {
	records, err := t.Client.ListData(context.Background())
	if err != nil {
		fmt.Println("Ошибка: ", err)
		return
	}

	for index, value := range records {
		fmt.Printf("Номер: %d\n", index)
		err = separateDataByTypeToOutput(value)
		if err != nil {
			fmt.Println("Ошибка, попробуйте еще раз")
			return
		}
		fmt.Print("\n")
	}
}

func (t *TUI) updateData() {
	var id string
	fmt.Print("Введите id данных: ")
	fmt.Scan(&id)
	record, err := t.Client.RetrieveData(context.Background(), id)
	if err != nil {
		fmt.Println("Ошибка: ", err)
		return
	}

	req, meta, err := separateDataByTypeToInput(record.TypeRecord)
	if err != nil {
		fmt.Println("Ошибка, попробуйте еще раз")
		return
	}

	err = t.Client.UpdateData(context.Background(), id, meta, req)
	if err != nil {
		fmt.Println("Ошибка: ", err)
		return
	}
	fmt.Println("Успешно!")
}

func (t *TUI) deleteData() {
	var id string
	fmt.Print("Введите id данных: ")
	fmt.Scan(&id)
	err := t.Client.DeleteData(context.Background(), id)
	if err != nil {
		fmt.Println("Ошибка: ", err)
		return
	}
	fmt.Println("Успешно!")

}

func (t *TUI) getVersion() {
	ver, err := t.Client.GetVersion(context.Background())
	if err != nil {
		fmt.Println("Ошибка, попробуйте еще раз")
		return
	}
	fmt.Printf("Версия сборки: %s\n Дата: %s\n", ver.Version, ver.Date)
}

func separateDataByTypeToOutput(record models.Record) error {
	switch record.TypeRecord {
	case models.Login:
		var jsonResp models.LoginJSON
		err := json.Unmarshal(record.Data, &jsonResp)
		if err != nil {
			return err
		}
		fmt.Printf("ID: %d\nЛогин: %s\nПароль: %s\nЗаметка: %s\n", record.Id, jsonResp.Login, jsonResp.Password, record.Meta)

	case models.Text, models.Binary:
		fmt.Printf("ID: %d\nДанные: %s\nЗаметка: %s\n", record.Id, string(record.Data), record.Meta)

	case models.Card:
		var jsonResp models.CardJSON
		err := json.Unmarshal(record.Data, &jsonResp)
		if err != nil {
			return err
		}
		fmt.Printf("ID: %d\nНомер карты: %s\nСрок действия: %s\nCVV: %s\nЗаметка: %s\n", record.Id, jsonResp.Number, jsonResp.Date, jsonResp.CVV, record.Meta)
	}
	return nil
}

func separateDataByTypeToInput(dataType string) ([]byte, string, error) {
	switch dataType {
	case models.Login:
		var username, password, meta string
		fmt.Print("Введите логин:")
		fmt.Scan(&username)
		fmt.Print("Введите пароль:")
		fmt.Scan(&password)
		fmt.Print("Введите заметку к данным:")
		fmt.Scan(&meta)
		jsonReq, err := json.Marshal(models.LoginJSON{Login: username, Password: password})
		if err != nil {
			fmt.Println("Ошибка, попробуйте еще раз")
			return nil, "", err
		}
		return jsonReq, meta, err

	case models.Text, models.Binary:
		var text, meta string
		fmt.Print("Введите данные:")
		fmt.Scan(&text)
		fmt.Print("Введите заметку к данным:")
		fmt.Scan(&meta)
		return []byte(text), meta, nil
	case models.Card:
		var number, date, cvv, meta string
		fmt.Print("Введите номер карты:")
		fmt.Scan(&number)
		fmt.Print("Введите срок действия карты:")
		fmt.Scan(&date)
		fmt.Print("Введите cvv:")
		fmt.Scan(&cvv)
		fmt.Print("Введите заметку к данным:")
		fmt.Scan(&meta)
		jsonReq, err := json.Marshal(models.CardJSON{Number: number, Date: date, CVV: cvv})
		if err != nil {
			fmt.Println("Ошибка, попробуйте еще раз")
			return nil, "", err
		}
		return jsonReq, meta, err
	}
	return nil, "", fmt.Errorf("not found data type")
}
