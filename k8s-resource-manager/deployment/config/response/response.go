package response

type ResponseErrorStatus struct {
	Name      string `json:"name"`
	NameSpace string `json:"namespace"`
	Reason    string `json:"reason"`
	Message   string `json:"message"`
	Age       string `json:"age"`
}

type ResponseDeleteErrorStatus struct {
	Name      string `json:"name"`
	NameSpace string `json:"namespace"`
	Reason    string `json:"reason"`
	Message   string `json:"message"`
	Age       string `json:"age"`
}
