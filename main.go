package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

/* CSV format (Shift JIS), http://www.post.japanpost.jp/zipcode/dl/kogaki-zip.html

全国地方公共団体コード（JIS X0401、X0402）………　半角数字
（旧）郵便番号（5桁）………………………………………　半角数字
郵便番号（7桁）………………………………………　半角数字
都道府県名　…………　半角カタカナ（コード順に掲載）　（注1）
市区町村名　…………　半角カタカナ（コード順に掲載）　（注1）
町域名　………………　半角カタカナ（五十音順に掲載）　（注1）
都道府県名　…………　漢字（コード順に掲載）　（注1,2）
市区町村名　…………　漢字（コード順に掲載）　（注1,2）
町域名　………………　漢字（五十音順に掲載）　（注1,2）
一町域が二以上の郵便番号で表される場合の表示　（注3）　（「1」は該当、「0」は該当せず）
小字毎に番地が起番されている町域の表示　（注4）　（「1」は該当、「0」は該当せず）
丁目を有する町域の場合の表示　（「1」は該当、「0」は該当せず）
一つの郵便番号で二以上の町域を表す場合の表示　（注5）　（「1」は該当、「0」は該当せず）
更新の表示（注6）（「0」は変更なし、「1」は変更あり、「2」廃止（廃止データのみ使用））
変更理由　（「0」は変更なし、「1」市政・区政・町政・分区・政令指定都市施行、「2」住居表示の実施、「3」区画整理、「4」郵便区調整等、「5」訂正、「6」廃止（廃止データのみ使用））
================
Response
pref: 都道府県の文字列
city: 市区町村の文字列
town: 町域名の文字列
address: 市区町村の文字列（cityとtownを結合したもの）
fullAddress: 都道府県+市区町村+町域名の結合文字列
*/
type Zipcode struct {
	Prefecture  string `json:"pref"`
	City        string `json:"city"`
	Town        string `json:"town"`
	Address     string `json:"address"`
	FullAddress string `json:"fullAddress"`
}
type ZipcodeResponse struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Data    *Zipcode `json:"data,omitempty"`
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	port := flag.Int("port", 80, "Port")
	filename := flag.String("filename", "KEN_ALL.CSV", "Path to csv file (ShiftJIS)")
	flag.Parse()

	csvFile, err := os.Open(*filename)
	check(err)
	// Parse SHIFTJIS
	reader := csv.NewReader(transform.NewReader(csvFile, japanese.ShiftJIS.NewDecoder()))
	zipcode := make(map[string]Zipcode)
	// Assign data to map
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		zipcode[line[2]] = Zipcode{
			Prefecture:  line[6],
			City:        line[7],
			Town:        line[8],
			Address:     line[7] + line[8],
			FullAddress: line[6] + line[7] + line[8],
		}
	}
	// Only accept xxx-yyyy or xxxyyyy
	var validZipcode = regexp.MustCompile(`^(\d{3})[-]?(\d{4})$`)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		q := r.URL.Query()

		// Check zipcode is valid
		if !validZipcode.MatchString(q.Get("zipcode")) {
			log.Println(q.Get("zipcode"))
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ZipcodeResponse{Code: 400, Message: "Invalid zipcode format."})
			return
		}

		z := validZipcode.ReplaceAllString(q.Get("zipcode"), "${1}${2}")
		log.Println(z)

		// Get zipcode
		value, ok := zipcode[z]
		if ok {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(ZipcodeResponse{Code: 200, Message: "Success", Data: &value})
		} else {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ZipcodeResponse{Code: 404, Message: "Zipcode not found"})
			return
		}

	})
	log.Println("Server serve at :" + strconv.Itoa(*port))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), nil))

}
