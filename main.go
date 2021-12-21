package main

import (
	"SwaggerParser/converter"
	"SwaggerParser/models"
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
)

func main() {

	var (
		methods     map[string][]models.Method = make(map[string][]models.Method)
		dataSwagger []models.Swagger
	)

	// swaggerUrl := "http://10.250.232.104:8001/swagger/doc.json"
	// resp, err := http.Get(swaggerUrl)
	// if err != nil {
	// 	panic(err)
	// }
	// if resp.StatusCode != 200 {
	// 	err := fmt.Errorf(resp.Status)
	// 	panic(err)
	// }
	// defer resp.Body.Close()

	// respBytes, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	panic(err)
	// }
	document, err := openapi3.NewLoader().LoadFromFile("doc.json") //LoadFromData(respBytes)
	if err != nil {
		panic(err)
	}
	for key, value := range document.Paths {
		method := new(models.Method)
		method.Path = key
		var group string = ""

		if value.Get != nil {
			method.Description = value.Get.Description
			method.Summary = value.Get.Summary
			method.MethodName = "GET"
			group = value.Get.Tags[0]
			fmt.Printf("value.Get.Security: %+v\n", value.Get.Security)
			if value.Get.Security != nil {
				method.IsJWT = "Да"
			} else {
				method.IsJWT = "Нет"
			}
		} else if value.Post != nil {
			method.Description = value.Post.Description
			method.Summary = value.Post.Summary
			method.MethodName = "POST"
			group = value.Post.Tags[0]
			fmt.Printf("value.Post.Security: %+v\n", value.Post.Security)
			if value.Post.Security != nil {
				method.IsJWT = "Да"
			} else {
				method.IsJWT = "Нет"
			}
		} else if value.Put != nil {
			method.Description = value.Put.Description
			method.Summary = value.Put.Summary
			method.MethodName = "PUT"
			group = value.Put.Tags[0]
			fmt.Printf("value.Put.Security: %+v\n", value.Put.Security)
			if value.Put.Security != nil {
				method.IsJWT = "Да"
			} else {
				method.IsJWT = "Нет"
			}
		} else if value.Delete != nil {
			method.Description = value.Delete.Description
			method.Summary = value.Delete.Summary
			method.MethodName = "DELETE"
			group = value.Delete.Tags[0]
			fmt.Printf("value.Delete.Security: %+v\n", value.Delete.Security)
			if value.Delete.Security != nil {
				method.IsJWT = "Да"
			} else {
				method.IsJWT = "Нет"
			}
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
		panic(err)
	}
	_, err = parser.GeneratePDF()
	if err != nil {
		panic(err)
	}
}
