package core

type PostmanStructure struct {
	Info InfoStructure   `json:"info"`
	Item []ItemStructure `json:"item"`
}

type InfoStructure struct {
	PostmanID string `json:"_postman_id"`
	Name      string `json:"name"`
	Schema    string `json:"schema"`
	Export    string `json:"_export"`
}

type ItemStructure struct {
	Name    string                `json:"name"`
	Item    *[]ItemStructure      `json:"item,omitempty"`
	Request *ItemRequestStructure `json:"request,omitempty"`
}

type ItemRequestStructure struct {
	Method string                       `json:"method"`
	Header []ItemRequestHeaderStructure `json:"header"`
	Body   ItemRequestBodyStructure     `json:"body"`
	Url    ItemRequestURLStructure      `json:"url"`
	Auth   *ItemRequestAuth             `json:"auth,omitempty"`
}

type ItemRequestAuth struct {
	Type   string                  `json:"type"`
	Bearer []ItemRequestAuthBearer `json:"bearer"`
}

type ItemRequestAuthBearer struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

type ItemRequestHeaderStructure struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

type ItemRequestBodyStructure struct {
	Mode string `json:"mode"`
	Raw  string `json:"raw"`
}

type ItemRequestURLStructure struct {
	Raw      string   `json:"raw"`
	Protocol *string  `json:"protocol,omitempty"`
	Host     []string `json:"host"`
	Path     []string `json:"path"`
}
