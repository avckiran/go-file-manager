package main

type Attributes struct {
	Color string `json:"color"`
	Size  string `json:"size"`
}

type InputData struct {
	ItemId     string     `json:"itemId"`
	ItemName   string     `json:"itemName"`
	Quantity   int        `json:"quantity"`
	Attributes Attributes `json:"attributes"`
}

type ApiResponseData struct {
	Fact   string `json:"fact"`
	Length int    `json:"length"`
}

type AggregatedRecord struct {
	ItemID         string `json:"itemId"`
	ItemName       string `json:"itemName"`
	Quantity       int    `json:"quantity"`
	Color          string `json:"color,omitempty"`
	Size           string `json:"size,omitempty"`
	FetchedAPIData string `json:"fetchedApiData"`
	APISource      string `json:"apiSource"`
	ProcessingDate string `json:"processingDate"`
}
