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
	flagTTLVouch       = flag.Int("t", 60, "Время дуйствия ваучера (0 - Unlimited)")
	flagUpSpeedVouch   = flag.Int("up", 1024, "Скорость отдачи для ваучера (0 - Unlimited)")
	flagDownSpeedVouch = flag.Int("down", 1024, "Скорость загрузки для ваучера (0 - Unlimited)")
	flsgServer         = flag.String("s", "unifi", "Сервер")
	flsgPort           = flag.Int("p", 8443, "Порт сервера")
)

type Output struct {
	Success bool               `json:"success"`
	Message string             `json:"message,omitempty"`
	Error   string             `json:"error,omitempty"`
	Data    []vouchers.Voucher `json:"data,omitempty"`
}

func result(vouchersList []vouchers.Voucher, err error) {
	if err != nil {
		result := Output{
			Success: false,
			Error:   err.Error(),
		}
		json.NewEncoder(os.Stdout).Encode(result)
		return
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
}

func main() {
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

	// Login
	err := vouchers.Login(creds.Username, creds.Password)
	if err != nil {
		result(nil, err)
		return
	}

	// Create
	nameVouch, err := vouchers.CreateVauchers(*flagCountVouch, *flagTTLVouch, *flagUpSpeedVouch, *flagDownSpeedVouch)
	if err != nil {
		result(nil, err)
		return
	}

	// Filter
	vouchersList, err := vouchers.GetFilterNoteVauchers(nameVouch)
	result(vouchersList, err)

}
