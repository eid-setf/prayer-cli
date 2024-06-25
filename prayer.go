package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

var (
	latitude  = 30.983334
	longitude = 41.016666
)

func DownloadTimesForMonth(monthPath string) {
	fmt.Printf("Downloading %v...\n", monthPath)
	year, month, _ := time.Now().Date()

	apiUrl := "https://api.aladhan.com/v1/calendar"
	method := 4

	requestUrl := fmt.Sprintf("%v/%v/%v?latitude=%v&longitude=%v&method=%v",
		apiUrl, year, int(month), latitude, longitude, method)

	resp, err := http.Get(requestUrl)
	if err != nil {
		panic(err)
	}

	f, err := os.Create(monthPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	io.Copy(f, resp.Body)
}

func GetTodayTimes() []time.Time {
	times := make([]time.Time, 5)

	_, month, day := time.Now().Date()

	f, err := os.Open(fmt.Sprintf("%v/%v.json", timesDir, int(month)))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	fb, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	var v struct {
		Data []map[string]map[string]interface{}
	}
	if err := json.Unmarshal(fb, &v); err != nil {
		panic(err)
	}

	todayTimes := v.Data[day-1]["timings"]
	for i, p := range []string{"Fajr", "Dhuhr", "Asr", "Maghrib", "Isha"} {
		t, err := time.Parse("15:04 (-07)", todayTimes[p].(string))
		if err != nil {
			panic(err)
		}
		times[i] = t
	}

	return times
}

// TODO: continue this function
func main() {
	os.MkdirAll(timesDir, 0644)
	for i := 1; i <= 12; i++ {
		monthPath := fmt.Sprintf("%v/%v.json", timesDir, i)
		if _, err := os.Stat(monthPath); os.IsNotExist(err) {
			DownloadTimesForMonth(monthPath)
		}
	}

	times := GetTodayTimes()

	fmt.Println(times)
}
