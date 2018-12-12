Server serve Japan's address by zipcode

## How to use

Download "KEN_ALL.CSV" from http://www.post.japanpost.jp/zipcode/dl/kogaki-zip.html
go run main.go -port 8888 -filename KEN_ALL.CSV

Browser:
http://localhost:8888?zipcode=1500031
http://localhost:8888?zipcode=351-0033

## How to build

### Mac
go build -o build/zip2address-linux-386 -v ./

### Linux 
env GOOS=linux GOARCH=386 go build -o build/zip2address-linux-386 -v ./