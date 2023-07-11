package services

import (
	//"bufio"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"

	"encoding/base64"
	//"io/ioutil"
	"path"
	//"strings"

	"github.com/joho/godotenv"

	"EnronEmailApi/models"
	"encoding/json"
	"log"
	"net/http"
	"net/mail"
	"os"
	"sync"
)

var MutexJson sync.Mutex

func ListAllFolders(folderName string) ([]string, []string) {
	files, err := os.ReadDir(folderName)
	if err != nil {
		log.Println(err)
	}
	var listFolders []string
	var listFiles []string
	for _, f := range files {
		if f.IsDir() {
			listFiles = append(listFiles, f.Name())
		} else {
			listFolders = append(listFolders, f.Name())
		}
	}
	return listFiles, listFolders
}

func DivideFolders(list []string, numParts int) [][]string {
	var divided [][]string
	partSize := len(list) / numParts
	remaining := len(list) % numParts

	index := 0
	for i := 0; i < numParts; i++ {
		size := partSize
		if i < remaining {
			size++
		}
		divided = append(divided, list[index:index+size])
		index += size
	}

	return divided
}

func ListFiles(folderName string) []string {
	files, err := os.ReadDir(folderName)
	if err != nil {
		log.Fatal(err)
	}
	var filesNames []string
	for _, file := range files {
		filesNames = append(filesNames, file.Name())
	}
	return filesNames
}

func parseData(dataLines *bufio.Scanner) models.Email {
	var data models.Email
	for dataLines.Scan() {
		line := dataLines.Text()
		switch {
		case strings.Contains(line, "Message-ID:"):
			data.Message_ID = line[11:]
		case strings.Contains(line, "Date:"):
			data.Date = line[5:]
		case strings.Contains(line, "From:"):
			data.From = line[5:]
		case strings.Contains(line, "To:"):
			data.To = line[3:]
		case strings.Contains(line, "Subject:"):
			data.Subject = line[8:]
		default:
			data.Body += line
		}
	}
	return data
}

func Algodeaca(folderList []string, path string, wg *sync.WaitGroup) {

	defer wg.Done()

	for _, user := range folderList {

		var jsonForBulk models.JsonBulk

		fmt.Println(user)

		processDir(path+user, &jsonForBulk)

		IndexDataBulk(jsonForBulk)
		jsonForBulk = models.JsonBulk{}

	}

}

func ProcessMailFile(path string, jsonForBulk *models.JsonBulk, wg *sync.WaitGroup) {
	defer wg.Done()

	sysFile, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Error opening file: %s\n", err)
		return
	}

	r := bytes.NewReader(sysFile)
	m, err := mail.ReadMessage(r)
	if err != nil {
		fmt.Printf("ISSUE WHILE READING EMAIL: %s\n", err)
		fmt.Printf("PATH TO MAIL: %s\n", path)
		data := bytes.NewReader(sysFile)
		lines := bufio.NewScanner(data)
		newEmail := parseData(lines)
		jsonForBulk.Records = append(jsonForBulk.Records, newEmail)
		return
	}
	body, err := io.ReadAll(m.Body)
	if err != nil {
		fmt.Println("aca se rompe")
	}

	newEmail := models.Email{
		Message_ID: path,
		From:       m.Header.Get("From"),
		To:         m.Header.Get("To"),
		Subject:    m.Header.Get("Subject"),
		Body:       string(body),
	}
	jsonForBulk.Records = append(jsonForBulk.Records, newEmail)

}

func IndexDataBulk(jsonForBulk models.JsonBulk) {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error al cargar el archivo .env")
	}

	var (
		user     = os.Getenv("USER")
		password = os.Getenv("PASSWORD")
	)

	auth := user + ":" + password
	bas64encoded_creds := base64.StdEncoding.EncodeToString([]byte(auth))

	jsonForBulk.Index = "emails"
	zinc_host := "http://localhost:4080"
	zinc_url := zinc_host + "/api/_bulkv2"

	body, _ := json.Marshal(jsonForBulk)

	req, err := http.NewRequest("POST", zinc_url, bytes.NewBuffer(body))

	if err != nil {
		log.Fatal("Error reading request. ", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+bas64encoded_creds)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
}

func processDir(name string, jsonForBulk *models.JsonBulk) {
	d, err := os.Open(name)
	if err != nil {
		log.Println(err.Error())
	}

	defer d.Close()

	files, err := d.ReadDir(-1)
	if err != nil {
		log.Println(err.Error())
	}

	for _, f := range files {
		if f.IsDir() {
			processDir(path.Join(name, f.Name()), jsonForBulk)
		} else {
			var wg sync.WaitGroup
			wg.Add(1)
			go ProcessMailFile(path.Join(name, f.Name()), jsonForBulk, &wg)
			wg.Wait()
		}
	}

}
