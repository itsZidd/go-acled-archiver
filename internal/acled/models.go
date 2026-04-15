package acled

type Event struct {
	EventID           string  `json:"event_id_cnty"`
	EventDate         string  `json:"event_date"`
	Year              int     `json:"year"` // No ,string
	TimePrecision     string  `json:"time_precision"`
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
	ISO               int     `json:"iso"` // No ,string
	Region            string  `json:"region"`
	Country           string  `json:"country"`
	Admin1            string  `json:"admin1"`
	Admin2            string  `json:"admin2"`
	Admin3            string  `json:"admin3"`
	Location          string  `json:"location"`
	GeoPrecision      int     `json:"geo_precision"`    // No ,string
	Latitude          float64 `json:"latitude,string"`  // KEEP ,string
	Longitude         float64 `json:"longitude,string"` // KEEP ,string
	Source            string  `json:"source"`
	SourceScale       string  `json:"source_scale"`
	Notes             string  `json:"notes"`
	Fatalities        int     `json:"fatalities"` // No ,string
	Timestamp         int64   `json:"timestamp"`  // No ,string
	Tags              *string `json:"tags"`
	Population1km     int     `json:"population_1km"`
	Population2km     int     `json:"population_2km"`
	Population5km     *int    `json:"population_5km"`
	PopulationBest    int     `json:"population_best"`
}

type APIResponse struct {
	TotalCount int     `json:"total_count"`
	Data       []Event `json:"data"`
}
