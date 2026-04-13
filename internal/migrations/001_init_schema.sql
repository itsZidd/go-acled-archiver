DROP TABLE IF EXISTS acled_events;

CREATE TABLE acled_events (
    event_id_cnty VARCHAR(255) PRIMARY KEY,
    event_date DATE,
    year INTEGER,
    disorder_type VARCHAR(100),
    event_type VARCHAR(100),
    sub_event_type VARCHAR(100),
    actor1 TEXT,
    actor2 TEXT,
    inter1 TEXT,
    inter2 TEXT,
    interaction TEXT,
    country VARCHAR(255),
    iso INTEGER,
    fatalities INTEGER,
    notes TEXT,
    timestamp BIGINT,
    last_updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_acled_year_country ON acled_events(year, country);
