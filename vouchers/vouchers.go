package vouchers

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"
)

var (
	httpClient http.Client
	baseURL    string
)

var (
	NameNoteVouchers string
)

var (
	loginURL          = baseURL + "/api/login"
	getVoucherURL     = baseURL + "/api/s/default/stat/voucher"
	CreateVauchersURL = baseURL + "/api/s/default/cmd/hotspot"
)

type Voucher struct {
	Duration       int    `json:"duration"`
	Note           string `json:"note"`
	Qos_overwrite  bool   `json:"qos_overwrite"`
	For_hotspot    bool   `json:"for_hotspot"`
	Code           string `json:"code"`
	Create_time    int64  `json:"create_time"`
	Quota          int    `json:"quota"`
	Site_id        string `json:"site_id"`
	External_id    string `json:"external_id"`
	Id             string `json:"_id"`
	Admin_name     string `json:"admin_name"`
	Used           int    `json:"used"`
	Status         string `json:"status"`
	Status_expires int    `json:"status_expires"`
}

type VoucherResponse struct {
	Meta struct {
		RC string `json:"rc"`
	} `json:"meta"`
	Data []Voucher `json:"data"`
}

// Специальная функция которая выполняется перед main
func init() {
	// Игнорировать ошибки SSL (для самоподписанных сертификатов UniFi)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	jar, _ := cookiejar.New(nil)
	httpClient = http.Client{
		Transport: tr,
		Jar:       jar, // создаёт "банку" (jar) для cookies
	}
}

func SetServerURL(server string, port int) {
	baseURL = fmt.Sprintf("https://%s:%d", server, port)
}

// Создать сессию. Login
func Login(login, password string) error {

	loginData := map[string]string{
		"username": login,
		"password": password,
	}

	loginBody, _ := json.Marshal(loginData) // Преобразование в JSON

	req, err := http.NewRequest("POST", loginURL, bytes.NewBuffer(loginBody))
	if err != nil {
		log.Printf("Не удалось сформировать запрос авторизации. Ошибка: %s\n", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Printf("Ошибка выполнения запроса на сервер: %s\n", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("Ошибка авторизации: %s\n", resp.Status)
		log.Println(err)
		return err
	}

	// (опционально) вывести ответ
	// body, _ := io.ReadAll(resp.Body)
	// log.Println(string(body))
	log.Println("Успешный вход в UniFi Controller")

	return nil
}

// Запросить список всех ваучеров
// func GetVauchers() {
// 	req, err := http.NewRequest("GET", getVoucherURL, nil)
// 	if err != nil {
// 		panic(err)
// 	}

// 	resp, err := httpClient.Do(req)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer resp.Body.Close()

// 	// (опционально) прочитать и вывести ответ
// 	// body, _ := io.ReadAll(resp.Body)
// 	// log.Println(string(body))
// }

// GetAPIVouchers возвращает только ваучеры с note, начинающейся с "API-created-*"
func GetFilterNoteVauchers(NameNoteVouchers string) ([]Voucher, error) {
	req, err := http.NewRequest("GET", getVoucherURL, nil)
	if err != nil {
		return nil, fmt.Errorf("Не удалось создать запрос: %w", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Ошибка HTTP-запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Сервер вернул статус: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Ошибка чтения тела ответа: %w", err)
	}

	var vr VoucherResponse
	if err := json.Unmarshal(body, &vr); err != nil {
		return nil, fmt.Errorf("Ошибка парсинга JSON: %w", err)
	}

	if vr.Meta.RC != "ok" {
		return nil, fmt.Errorf("Ошибка UniFi API: %s", string(body))
	}

	// Фильтрация
	var filtered []Voucher
	for _, v := range vr.Data {
		if strings.HasPrefix(v.Note, NameNoteVouchers) {
			filtered = append(filtered, v)
		}
	}

	return filtered, nil

	// (опционально) прочитать и вывести ответ
	// body, _ := io.ReadAll(resp.Body)
	// log.Println(string(body))
}

func CreateVauchers(count, ttl, uploadSpeed, downloadSpeed int) (string, error) {

	now := time.Now()
	dateTime := now.Format("2006-01-02-15-04-05-")
	nanosecStr := fmt.Sprintf("%09d", now.Nanosecond())
	uniqueStr := strings.ReplaceAll(dateTime+nanosecStr, "-", "")

	// Тело запроса
	NameNoteVouchers := "API-created-" + uniqueStr
	payload := map[string]interface{}{
		"cmd":    "create-voucher",
		"expire": ttl,
		"n":      count,
		"quota":  1,
		"note":   NameNoteVouchers,
		"up":     uploadSpeed,
		"down":   downloadSpeed,
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", CreateVauchersURL, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Ошибка при формировании запроса на создание ваучеров: %s\n", err)
		return "", err
	}
	// Заголовки
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Printf("Ошибка при запросе на сервер: %s\n", err)
		return "", err
	}
	defer resp.Body.Close()

	return NameNoteVouchers, nil
}
