package models

type Email struct {
	Message_ID string `json:"Message-ID"`
	Date       string `json:"Date"`
	From       string `json:"from"`
	To         string `json:"to"`
	Subject    string `json:"subject"`
	Body       string `json:"Body"`
}

type JsonBulk struct {
	Index   string  `json:"index"`
	Records []Email `json:"records"`
}
