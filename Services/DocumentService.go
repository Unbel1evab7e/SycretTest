package Services

import (
	"SycretTest/Models/Document"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/antchfx/xmlquery"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
)

type DocumentService struct {
	ServiceUrl string
}

func NewDocumentService(serviceUrl string) *DocumentService {
	return &DocumentService{
		ServiceUrl: serviceUrl,
	}
}

func (s DocumentService) GetDocument(request *Document.DocumentRequest) (docx string, err error) {
	client := http.Client{}

	req, err := http.NewRequest(http.MethodGet, request.URLTemplate, strings.NewReader(""))
	if err != nil {
		return "", err
	}
	req.Header.Add("User-Agent", "Unbel1evab7e")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Connection", "keep-alive")

	docRes, err := client.Do(req)
	defer docRes.Body.Close()
	if err != nil {
		return "", err
	}

	if docRes.StatusCode != http.StatusOK {
		return "", errors.New("Bad response")
	}
	doc, err := xmlquery.Parse(docRes.Body)
	tags := xmlquery.Find(doc, "//ns1:text[@field != \"\"]")

	for _, tag := range tags {
		fieldName := tag.Attr[0]

		serviceDataUrl, err := url.Parse(s.ServiceUrl)
		if err != nil {
			return "", err
		}

		params := url.Values{}
		params.Add("text", fieldName.Value)
		params.Add("recordid", string(request.RecordID))

		serviceDataUrl.RawQuery = params.Encode()

		req, err := http.NewRequest(http.MethodGet, serviceDataUrl.String(), strings.NewReader(""))
		if err != nil {
			return "", err
		}
		req.Header.Add("User-Agent", "Unbel1evab7e")
		req.Header.Add("Accept", "*/*")
		req.Header.Add("Accept-Encoding", "gzip, deflate, br")
		req.Header.Add("Connection", "keep-alive")
		dataRes, err := client.Do(req)

		defer dataRes.Body.Close()
		if err != nil {
			return "", err
		}

		if dataRes.StatusCode != http.StatusOK {
			return "", errors.New(fmt.Sprintf("Bad response at word:%v", fieldName.Value))
		}

		b, err := io.ReadAll(dataRes.Body)
		if err != nil {
			return "", err
		}

		var body Document.DocumentDataResponse

		err = json.Unmarshal(b, &body)
		if err != nil {
			return "", err
		}

		textTag := xmlquery.Find(tag, "//w:t")

		textTag[0].FirstChild.Data = body.Resultdata
	}

	// Build fileName from fullPath
	fileURL, err := url.Parse(request.URLTemplate)
	if err != nil {
		log.Fatal(err)
	}
	link := fileURL.Path
	segments := strings.Split(link, "/")
	fileName := segments[len(segments)-1]

	dots := strings.Split(fileName, ".")

	folderPath := path.Join("temp", "documents", "response", strconv.Itoa(request.RecordID))

	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		mkErr := os.MkdirAll(folderPath, os.ModePerm)

		if mkErr != nil {
			return "", mkErr
		}
	}
	filePathDocx := path.Join(folderPath, fmt.Sprintf("%v.docx", strings.Join(dots[:len(dots)-1], "")))

	if _, err := os.Stat(filePathDocx); os.IsNotExist(err) {
		_, crErr := os.Create(filePathDocx)
		if crErr != nil {
			return "", crErr
		}
		os.WriteFile(filePathDocx, []byte(doc.OutputXML(true)), os.ModePerm)
	}

	return filePathDocx, nil

}
