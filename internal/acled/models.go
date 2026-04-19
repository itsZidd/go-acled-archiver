package acled

type Event struct {
	EventID           string  `json:"event_id_cnty"`
	EventDate         string  `json:"event_date"`
	Year              int     `json:"year"`
	DisorderType      string  `json:"disorder_type"`
	EventType         string  `json:"event_type"`
	SubEventType      string  `json:"sub_event_type"`
	Actor1            string  `json:"actor1"`
	AssocActor1       string  `json:"assoc_actor_1"`
	Inter1            string  `json:"inter1"`
	Actor2            string  `json:"actor2"`
	AssocActor2       string  `json:"assoc_actor_2"`
	Inter2            string  `json:"inter2"`
	Interaction       string  `json:"interaction"`
	CivilianTargeting string  `json:"civilian_targeting"`
	ISO               int     `json:"iso"`
	Region            string  `json:"region"`
	Country           string  `json:"country"`
	Admin1            string  `json:"admin1"`
	Latitude          float64 `json:"latitude,string"`
	Longitude         float64 `json:"longitude,string"`
	Fatalities        int     `json:"fatalities"`
	Notes             string  `json:"notes"`
	PopulationBest    *int    `json:"population_best"`
}

type APIResponse struct {
	TotalCount int     `json:"total_count"`
	Data       []Event `json:"data"`
}
