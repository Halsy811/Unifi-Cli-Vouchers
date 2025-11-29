package main

import (
	"flag"
	"fmt"
	"os"
	"unifi-cli-vouchers/auth"
	"unifi-cli-vouchers/vouchers"
)

var (
	flagRegAuth        = flag.Bool("r", false, "Зарегистрировать учетные данные в системном хранилище")
	flagDelAuth        = flag.Bool("d", false, "Удалить учетные данные в системном хранилище")
	flagCountVouch     = flag.Int("c", 1, "Количество создаваемых ваучеров")
	flagTTLVouch       = flag.Int("t", 60, "Время дуйствия ваучера")
	flagUpSpeedVouch   = flag.Int("up", 1024, "Скорость отдачи для ваучера")
	flagDownSpeedVouch = flag.Int("down", 1024, "Скорость загрузки для ваучера")
)

func main() {

	flag.Parse()

	// Проверка ключей -r и -d
	if *flagRegAuth && *flagDelAuth {
		fmt.Fprintln(os.Stderr, "Ключи -r и -d не совместимы")
		os.Exit(1)
	}

	if *flagRegAuth {
		auth.Reg_auth()
		return
	}

	if *flagDelAuth {
		auth.Unreg_auth()
		return
	}

	creds := auth.Get_auth()

	// Основные действия
	// fmt.Println(login)
	// fmt.Println(password)

	vouchers.Login(creds.Username, creds.Password)
	// vouchers.GetVauchers()
	nameVouch := vouchers.CreateVauchers(*flagCountVouch, *flagTTLVouch, *flagUpSpeedVouch, *flagDownSpeedVouch)
	fmt.Printf("Искомое имя: %s\n", nameVouch)

	vouchers := vouchers.GetFilterNoteVauchers(nameVouch)

	for _, v := range vouchers {
		fmt.Printf("Ваучер: %s, заметка: %s, статус: %s\n", v.Code, v.Note, v.Status)
	}
}
