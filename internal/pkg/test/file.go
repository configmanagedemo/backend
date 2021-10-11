package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

var (
	path = "http://localhost:8080/"
)

var (
	client  *http.Client
	jar     *cookiejar.Jar
	cookies []*http.Cookie
)

func init() {
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err.Error())
	}

	cookies = []*http.Cookie{}
	url, err := url.Parse(path)
	if err != nil {
		panic(err.Error())
	}
	jar.SetCookies(url, cookies)
	client = &http.Client{Jar: jar}

}

func login() {
	jsonData, _ := json.Marshal(map[string]interface{}{
		"username": "admin",
		"password": "admin",
	})
	r, err := client.Post(path+"api/v1/login", "application/json", strings.NewReader(string(jsonData)))
	if err != nil {
		panic(err.Error())
	}
	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)
	fmt.Printf("%s\n", body)
}

func getFile(ch chan int, fileID int) {
	r, err := client.Get(fmt.Sprintf("%sapi/v1/bfile/%d/download", path, fileID))
	if err != nil {
		panic(err.Error())
	}
	defer r.Body.Close()
	ioutil.ReadAll(r.Body)
	ch <- 1
	// fmt.Printf("%s\n", body)
}

// FileInfoRsp
type FileInfoRsp struct {
	Data struct {
		CreatedAt time.Time `json:"created_at"`
		Desc      string    `json:"desc"`
		FileID    uint      `json:"file_id"`
		Filename  string    `json:"filename"`
		FileSize  uint      `json:"filesize"`
		IsUse     bool      `json:"is_use"`
		Uploader  string    `json:"uploader"`
	} `json:"data"`
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func getFilename(fileID int) FileInfoRsp {
	r, err := client.Get(fmt.Sprintf("%sapi/v1/bfile/%d", path, fileID))
	if err != nil {
		panic(err.Error())
	}
	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)
	data := &FileInfoRsp{}
	if err := json.Unmarshal(body, data); err != nil {
		panic(err.Error())
	}
	return *data
}

func DoGetFileWithLoop(loop, fileID int) {
	ch := make(chan int, loop)
	t := time.Now()
	for i := 0; i < loop; i++ {
		go getFile(ch, fileID)
	}

	for i := 0; i < loop; i++ {
		<-ch
		i++
	}

	total := time.Since(t)
	fmt.Printf("fileID[%d], get[%d], total cost[%f], cost per time[%f]\n",
		fileID, loop, total.Seconds(), total.Seconds()/float64(loop))
}

var (
	OneHundredTimes  = 100
	TwoHundredTimes  = 200
	OneThousandTimes = 1000
	KB               = 10
	MB               = 20
)

func DoGetFile() {
	login()
	fileID := 18
	fileData := getFilename(fileID)
	fmt.Printf("fileID[%d], filename[%s], size[%dB][%dK][%dM]\n",
		fileID, fileData.Data.Filename, fileData.Data.FileSize, fileData.Data.FileSize>>uint(KB), fileData.Data.FileSize>>uint(MB))
	DoGetFileWithLoop(OneHundredTimes, fileID)
	DoGetFileWithLoop(TwoHundredTimes, fileID)
	DoGetFileWithLoop(OneThousandTimes, fileID)
}
