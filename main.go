package main

import (
	"flag"
	"fmt"
	"os"
	"syscall"

	"github.com/zalando/go-keyring" // для сохранения пароля в системном хранилище
	"golang.org/x/term"
)

var (
	flagRegAuth = flag.Bool("r", false, "Зарегистрировать учетные данные")
	flagDelAuth = flag.Bool("d", false, "Удалить учетные данные")
)

var (
	keyringNameService = "unifi-cli"
	keyringKeyUsername = "username"
	keyringKeyPasswrd  = "password"
)

type Credentials struct {
	username string
	password string
}

// Регистрация учетных данных в системе
func reg_auth() {
	if *flagRegAuth {
		fmt.Print("Username: ") // login
		var username string
		fmt.Scanln(&username)

		fmt.Print("Password: ") // password
		passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			panic(err)
		}

		fmt.Println()

		keyring.Set(keyringNameService, keyringKeyUsername, username)
		keyring.Set(keyringNameService, keyringKeyPasswrd, string(passwordBytes))
	}
}

// Удаление учетных данных
func unreg_auth() {
	errUsername := keyring.Delete(keyringNameService, keyringKeyUsername)
	errPassword := keyring.Delete(keyringNameService, keyringKeyPasswrd)
	if errUsername != nil && errPassword != nil {
		fmt.Fprintf(os.Stderr, "Ошибка удаления учетных данных из хранилища: \n\tЛогин(%v) \n\tПароль(%v)\n", errUsername, errPassword)
		os.Exit(1)
	} else {
		fmt.Println("Учетные данные удалены их хранилища системы")
	}
}

// Получение учетных данных
func get_auth() *Credentials {

	username, errUsername := keyring.Get(keyringNameService, keyringKeyUsername)

	password, errPassword := keyring.Get(keyringNameService, keyringKeyPasswrd)
	if errUsername != nil || errPassword != nil {
		fmt.Fprintf(os.Stderr, "Ошибка при получении учетных данных:\n")
		if errUsername != nil {
			fmt.Fprintf(os.Stderr, "\tЛогин: %v\n", errUsername)
		}
		if errPassword != nil {
			fmt.Fprintf(os.Stderr, "\tПароль: %v\n", errPassword)
		}
		os.Exit(1)
	}
	return &Credentials{
		username: username,
		password: password,
	}
}

func main() {

	flag.Parse()

	// Проверка ключей -r и -d
	if *flagRegAuth && *flagDelAuth {
		fmt.Fprintln(os.Stderr, "Ключи -r и -d не совместимы")
		os.Exit(1)
	}

	if *flagRegAuth {
		reg_auth()
		return
	}

	if *flagDelAuth {
		unreg_auth()
		return
	}

	cregentials := get_auth()
	login := cregentials.username
	password := cregentials.password

	// Основные действия
	fmt.Println(login)
	fmt.Println(password)

}
