package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"syscall"
	"time"
	"unsafe"
)
var (
	user32 = syscall.MustLoadDLL("user32.dll")
)
const (
    coinbaseCryptoAPI = "https://api.coinbase.com/v2/prices/%s-%s/spot"
	currency          = "USD"
	retryDelay        = 30 // Delay between checking a tokens price
	MB_OK               = 0x00000000
	MB_ICONINFORMATION  = 0x00000040
	MB_SYSTEMMODAL      = 0x00001000
	WM_SETFOCUS         = 0x0007
	WM_ACTIVATEAPP      = 0x001C
)

type coinbasePriceResponse struct {
    Data struct {
        Amount string `json:"amount"`
    } `json:"data"`
}

type tokenConfig struct {
    Name      string  `json:"name"`
    Symbol    string  `json:"symbol"`
    Threshold float64 `json:"threshold"`
}

type config struct {
    Tokens []tokenConfig `json:"tokens"`
}

func main() {
    // Load configuration from file
    cfg, err := loadConfig("config.json")
    if err != nil {
        fmt.Println("Error loading configuration:", err)
        return
    }

    var wg sync.WaitGroup
    for _, token := range cfg.Tokens {
        wg.Add(1)
        go func(token tokenConfig) {
            defer wg.Done()
            for {
                price, err := getCoinbasePrice(token.Symbol)
                if err != nil {
                    fmt.Printf("Error retrieving %s price: %v\n", token.Name, err)
                    continue
                }
                if price > token.Threshold {
                    msg := fmt.Sprintf("%s price ($%.2f) is above threshold ($%.2f)!", token.Name, price, token.Threshold)
					fmt.Printf("%s\n", msg)
					MessageBox("Token has passed threshold", msg, MB_ICONINFORMATION|MB_OK|MB_SYSTEMMODAL)
                }

                time.Sleep(retryDelay * time.Second)
            }
        }(token)
    }

    wg.Wait()
}

func MessageBox(caption, message string, style uintptr) int {
	ret, _, _ := syscall.Syscall6(
		user32.MustFindProc("MessageBoxW").Addr(),
		4,
		0,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(message))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(caption))),
		style,
		0,
		0,
	)
	return int(ret)
}
func loadConfig(filename string) (*config, error) {
    data, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }

    var cfg config
    err = json.Unmarshal(data, &cfg)
    if err != nil {
        return nil, err
    }

    return &cfg, nil
}

func getCoinbasePrice(symbol string) (float64, error) {
    url := fmt.Sprintf(coinbaseCryptoAPI, symbol, currency)
    resp, err := http.Get(url)
    if err != nil {
        return 0, err
    }
    defer resp.Body.Close()

    var priceResp coinbasePriceResponse
    err = json.NewDecoder(resp.Body).Decode(&priceResp)
    if err != nil {
        return 0, err
    }

    price, err := strconv.ParseFloat(priceResp.Data.Amount, 64)
    if err != nil {
        return 0, err
    }

    return price, nil
}