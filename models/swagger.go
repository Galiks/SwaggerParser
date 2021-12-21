package models

type Swagger struct {
	Group   string
	Methods []Method
}

type Method struct {
	MethodName  string
	Path        string
	Summary     string
	IsJWT       string
	Action      string
	Description string
}
