Server serve Japan's address by zipcode

## How to use

Download "KEN_ALL.CSV" from http://www.post.japanpost.jp/zipcode/dl/kogaki-zip.html

`go run main.go -port 8888 -filename KEN_ALL.CSV`

Open browser:

`http://localhost:8888?zipcode=1500031`
```json
{
  "code":200,
  "message":"Success",
  "data":{
    "pref":"東京都",
    "city":"渋谷区",
    "town":"桜丘町",
    "address":"渋谷区桜丘町",
    "fullAddress":"東京都渋谷区桜丘町"
  }
}
```

`http://localhost:8888?zipcode=351-0033`
```json
{
  "code":200,
  "message":"Success",
  "data":{
    "pref":"埼玉県",
    "city":"朝霞市",
    "town":"浜崎",
    "address":"朝霞市浜崎",
    "fullAddress":"埼玉県朝霞市浜崎"
  }
}
```


## How to build

### Mac
go build -o build/zip2address-linux-386 -v ./

### Linux 
env GOOS=linux GOARCH=386 go build -o build/zip2address-linux-386 -v ./
