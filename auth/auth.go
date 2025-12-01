package auth

import (
	"fmt"
	"log"
	"os"
	"syscall"

	"github.com/zalando/go-keyring" // для сохранения пароля в системном хранилище
	"golang.org/x/term"
)

var (
	keyringNameService = "unifi-cli"
	keyringKeyUsername = "username"
	keyringKeyPasswrd  = "password"
)

type Credentials struct {
	Username string
	Password string
}

// Регистрация учетных данных в системе
func Reg_auth() {
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

// Удаление учетных данных
func Unreg_auth() {
	errUsername := keyring.Delete(keyringNameService, keyringKeyUsername)
	errPassword := keyring.Delete(keyringNameService, keyringKeyPasswrd)
	if errUsername != nil && errPassword != nil {
		log.Printf("Ошибка удаления учетных данных из хранилища: \n\tЛогин(%v) \n\tПароль(%v)\n", errUsername, errPassword)
		os.Exit(1)
	} else {
		log.Println("Учетные данные удалены их хранилища системы")
	}
}

// Получение учетных данных
func Get_auth() *Credentials {

	username, errUsername := keyring.Get(keyringNameService, keyringKeyUsername)

	password, errPassword := keyring.Get(keyringNameService, keyringKeyPasswrd)
	if errUsername != nil || errPassword != nil {
		log.Printf("Ошибка при получении учетных данных:\n")
		if errUsername != nil {
			log.Printf("\tЛогин: %v\n", errUsername)
		}
		if errPassword != nil {
			log.Printf("\tПароль: %v\n", errPassword)
		}
		os.Exit(1)
	}
	return &Credentials{
		Username: username,
		Password: password,
	}
}
