package acled

type Event struct {
	EventID      string `json:"event_id_cnty"`
	EventDate    string `json:"event_date"`
	Year         int    `json:"year"` // Removed ,string
	DisorderType string `json:"disorder_type"`
	EventType    string `json:"event_type"`
	SubEventType string `json:"sub_event_type"`
	Actor1       string `json:"actor1"`
	Actor2       string `json:"actor2"`
	Inter1       string `json:"inter1"`
	Inter2       string `json:"inter2"`
	Interaction  string `json:"interaction"`
	Country      string `json:"country"`
	ISO          int    `json:"iso"`        // Removed ,string
	Fatalities   int    `json:"fatalities"` // Removed ,string
	Notes        string `json:"notes"`
	Timestamp    int64  `json:"timestamp"`
}

type APIResponse struct {
	Status  int     `json:"status"`
	Count   int     `json:"count"`
	Message string  `json:"message"`
	Data    []Event `json:"data"`
}
