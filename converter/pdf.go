package converter

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"text/template"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

//pdf requestpdf struct
type RequestPdf struct {
	body string
}

//new request to pdf function
func NewRequestPdf(body string) *RequestPdf {
	return &RequestPdf{
		body: body,
	}
}

//parsing template function
func (r *RequestPdf) ParseTemplate(templateFileName string, data interface{}) error {
	file, err := os.Open(templateFileName)
	if err != nil {
		return err
	}
	document, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		return err
	}

	document.Find(".field15").Each(func(i int, s *goquery.Selection) {
		page := s.Find("p")
		page.SetText("QWERRTY")
	})
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		err = fmt.Errorf("31: %s", err)
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		err = fmt.Errorf("37: %s", err)
		return err
	}
	r.body = buf.String()
	return nil
}

//generate pdf function
func (r *RequestPdf) GeneratePDF() ([]byte, string, error) {
	t := time.Now().Unix()
	// write whole the body

	if _, err := os.Stat("cloneTemplate/"); os.IsNotExist(err) {
		err := os.Mkdir("cloneTemplate/", 0777)
		if err != nil {
			return nil, "", err
		}
	}
	fileName := "cloneTemplate/" + strconv.FormatInt(int64(t), 10) + ".html"
	err := ioutil.WriteFile(fileName, []byte(r.body), 0644)
	if err != nil {
		return nil, "", err
	}

	f, err := os.Open(fileName)
	if f != nil {
		defer f.Close()
	}
	if err != nil {
		return nil, "", err
	}

	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return nil, "", err
	}

	pdfg.AddPage(wkhtmltopdf.NewPageReader(f))

	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)

	pdfg.Dpi.Set(300)

	err = pdfg.Create()
	if err != nil {
		return nil, "", err
	}

	return pdfg.Bytes(), f.Name(), nil
}

// GetTextLength возвращает длину строки
//
// Если соединить массив текста в одну строчку, например, "Самостоятельной постройки", и взять длину, то получим 48
// Это происходит потому что, например, букву Я он принимает за 2 символа
// Чтобы этого избежать нужно преобразовать строку в срез рун. Длина такого среза будет действительной длинной строки
func GetTextLength(allTextField string) int {
	var converter []rune = make([]rune, 0)
	for _, runeValue := range allTextField {
		converter = append(converter, runeValue)
	}
	return len(converter)
}
