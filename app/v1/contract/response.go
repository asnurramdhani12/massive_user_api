package contract

type Response struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	MetaData   MetaData    `json:"meta_data"`
}

type MetaData struct {
	RequestID string `json:"request_id"`
}
