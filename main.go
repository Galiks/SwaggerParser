package main

import (
	"SwaggerParser/converter"
	"SwaggerParser/models"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

var (
	keyWordForAction = "permissions:"
)

func main() {

	var (
		document *openapi3.T
		err      error
	)

	url := flag.String("url", "", "url to doc.json")
	path := flag.String("path", "", "path to doc.json")
	flag.Parse()

	if *url != "" {
		resp, err := http.Get(*url)
		if err != nil {
			panic(err)
		}
		if resp.StatusCode != 200 {
			err := fmt.Errorf(resp.Status)
			panic(err)
		}
		defer resp.Body.Close()

		respBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		document, err = openapi3.NewLoader().LoadFromData(respBytes)
		if err != nil {
			panic(err)
		}
	}

	if *path != "" {
		document, err = openapi3.NewLoader().LoadFromFile(*path)
		if err != nil {
			panic(err)
		}
	}

	if *path == "" && *url == "" {
		log.Println("Args is empty")
		os.Exit(1)
	}

	if err = createHTML(document); err != nil {
		panic(err)
	}

}

func getActionFromDescription(description string) string {
	if description == "" {
		return ""
	}
	splitDescription := strings.Split(description, "\n")
	for _, desc := range splitDescription {
		if strings.Contains(desc, keyWordForAction) {
			return strings.TrimPrefix(desc, keyWordForAction)
		}
	}
	return ""
}

func getDescription(description string) string {
	var (
		result string = ""
	)
	if description == "" {
		return result
	}
	splitDescription := strings.Split(description, "\n")
	for _, desc := range splitDescription {
		if !strings.Contains(desc, keyWordForAction) {
			result += desc
		}
	}
	return result
}

func createHTML(document *openapi3.T) error {
	var (
		methods     map[string][]models.Method = make(map[string][]models.Method)
		dataSwagger []models.Swagger
		err         error
	)
	for key, value := range document.Paths {
		method := new(models.Method)
		method.Path = key
		var group string = ""

		if value.Get != nil {
			method.Action = getActionFromDescription(value.Get.Description)
			method.Description = getDescription(value.Get.Description)
			method.Summary = value.Get.Summary
			method.MethodName = "GET"
			group = value.Get.Tags[0]
			setSecurity(value.Get.Security, method)
		} else if value.Post != nil {
			method.Action = getActionFromDescription(value.Post.Description)
			method.Description = getDescription(value.Post.Description)
			method.Summary = value.Post.Summary
			method.MethodName = "POST"
			group = value.Post.Tags[0]
			setSecurity(value.Post.Security, method)
		} else if value.Put != nil {
			method.Action = getActionFromDescription(value.Put.Description)
			method.Description = getDescription(value.Put.Description)
			method.Summary = value.Put.Summary
			method.MethodName = "PUT"
			group = value.Put.Tags[0]
			setSecurity(value.Put.Security, method)
		} else if value.Delete != nil {
			method.Action = getActionFromDescription(value.Delete.Description)
			method.Description = getDescription(value.Delete.Description)
			method.Summary = value.Delete.Summary
			method.MethodName = "DELETE"
			group = value.Delete.Tags[0]
			setSecurity(value.Delete.Security, method)
		}

		methods[group] = append(methods[group], *method)
	}

	for k, v := range methods {
		var swag = new(models.Swagger)
		swag.Group = k
		swag.Methods = v
		dataSwagger = append(dataSwagger, *swag)
	}

	parser := converter.NewRequestPdf("")
	template := struct{ Data []models.Swagger }{
		Data: dataSwagger,
	}
	err = parser.ParseTemplate("index.html", template)
	if err != nil {
		return err
	}
	_, filename, err := parser.GeneratePDF()
	if err != nil {
		return err
	}
	fmt.Printf("filename: %v\n", filename)
	return nil
}

func setSecurity(security *openapi3.SecurityRequirements, method *models.Method) {
	if security != nil {
		method.IsJWT = "????"
	} else {
		method.IsJWT = "??????"
	}
}
