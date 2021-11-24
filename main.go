package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/Hudayberdyyev/online_translate/constants"
	"github.com/Hudayberdyyev/online_translate/model"
)

func readCsvFile(filePath string) [][]string {
	// Load a csv file.
	f, _ := os.Open(filePath)

	// Create a new reader.
	r := csv.NewReader(f)
	r.FieldsPerRecord = -1
	r.LazyQuotes = true
	r.Comma = ';'

	records, err := r.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	return records
}

func translateFromGoogle(from model.TextStruct, toLang string) string {
	var resStr string
	// fmt.Println(url.QueryEscape("هی"))
	// fmt.Println(url.QueryEscape(from.Val))

	spaceClient := http.Client{
		Timeout: time.Second * 5, // Timeout after 2 seconds
	}

	url := fmt.Sprintf("https://translate.googleapis.com/translate_a/single?client=gtx&sl=%s&tl=%s&dt=t&q=%s", from.Lang, toLang, url.QueryEscape(from.Val))
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36")

	res, getErr := spaceClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	var resByte []byte
	resByte, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var jsonBody interface{}
	if err := json.Unmarshal(resByte, &jsonBody); err != nil {
		log.Fatal(err)
	}

	convertArr, _ := jsonBody.([]interface{})
	doubleArr, _ := convertArr[0].([]interface{})
	thirdArray, _ := doubleArr[0].([]interface{})

	resStr, _ = thirdArray[0].(string)
	return resStr
}

func main() {
	records := readCsvFile("/Users/ahmet/Desktop/all_data.csv")
	fmt.Println(cap(records), " records")

	file, err := os.Create("korey.csv")
	defer file.Close()
	if err != nil {
		log.Fatalln("failed to open file", err)
	}

	w := csv.NewWriter(file)
	w.Comma = ';'
	defer w.Flush()

	for index, v := range records {
		if index < 112 {
			continue
		}
		srcText := model.TextStruct{
			Val:  v[1],
			Lang: constants.RUSSIAN,
		}

		tl := translateFromGoogle(srcText, constants.KOREY)
		v[1] = tl
		if err := w.Write(v); err != nil {
			log.Fatalln("error writing record to file", err)
		}
		fmt.Println(index)
	}

}
