package main

import (
	"encoding/json"
	"flag"
	"log"
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
	flsgServer         = flag.String("s", "unifi", "Сервер")
	flsgPort           = flag.Int("p", 8443, "Порт сервера")
)

type Output struct {
	Success bool               `json:"success"`
	Message string             `json:"message,omitempty"`
	Error   string             `json:"error,omitempty"`
	Data    []vouchers.Voucher `json:"data,omitempty"`
}

func main() {
	// Перенаправляем логи в stderr
	log.SetOutput(os.Stderr)

	flag.Parse()

	vouchers.SetServerURL(*flsgServer, *flsgPort)

	// Проверка ключей -r и -d
	if *flagRegAuth && *flagDelAuth {
		log.Println("Ключи -r и -d не совместимы")
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

	vouchers.Login(creds.Username, creds.Password)

	nameVouch := vouchers.CreateVauchers(*flagCountVouch, *flagTTLVouch, *flagUpSpeedVouch, *flagDownSpeedVouch)
	// log.Printf("Искомое имя: %s\n", nameVouch)

	vouchersList, err := vouchers.GetFilterNoteVauchers(nameVouch)
	if err != nil {
		log.Println("Ошибка:", err)
		result := Output{
			Success: false,
			Error:   err.Error(),
		}
		json.NewEncoder(os.Stdout).Encode(result)
		os.Exit(1)
	} else {
		result := Output{
			Success: true,
			Data:    vouchersList,
		}
		if err := json.NewEncoder(os.Stdout).Encode(result); err != nil {
			log.Printf("Не удалось закодировать ответ в JSON: %v", err)
			os.Exit(1)
		}
	}
	// for _, v := range vouchersList {
	// 	log.Printf("Ваучер: %s, заметка: %s, статус: %s\n", v.Code, v.Note, v.Status)
	// }

	os.Exit(0)
	// ############## Убрать все panic, вывод ошибкок в json
}
