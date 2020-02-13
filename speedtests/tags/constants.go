package main

import "fmt"

const CreateHyperTableEvents = "SELECT create_hypertable('events', 'start', 'id', 2, chunk_time_interval => 86400, if_not_exists => TRUE);"
const SetTimeFunctionForEvents = "SELECT set_integer_now_func('events', 'unix_now', replace_if_exists => true);"

const CreateEventsTable = `
	CREATE TABLE IF NOT EXISTS events (
		Id SERIAL,
		Title VARCHAR (100),
		Start BIGINT,
		Finish BIGINT,
		Center geometry,
		PostCodes VARCHAR(15)[],
		Tags TEXT[],
		PRIMARY KEY (Id, Start)
	);
`
const InsertEvent = `
	INSERT INTO events
		(Title, Start, Finish, Center, PostCodes, Tags)
	VALUES
		($1, $2, $3, ST_SetSRID( ST_Point($4, $5), 4326), $6, $7)
	RETURNING 
		(Id)
`

const CreateEventTagsTables1 = `
	CREATE TABLE IF NOT EXISTS event_tags1 (
		Text TEXT,
		EventId INTEGER,
		Start BIGINT,
		Finish BIGINT,
		PRIMARY KEY (EventId, Text)
	)
`

const CreateTagsTables2 = `
	CREATE TABLE IF NOT EXISTS  tags2 (
		Text Text UNIQUE,
		Id SERIAL PRIMARY KEY
	)
`

const CreateEventTagsTables2 = `
	CREATE TABLE IF NOT EXISTS  event_tags2 (
		TagId SERIAL REFERENCES tags2 (Id),
		EventId SERIAL,
		Start INTEGER,
		FINISH INTEGER
	)
`

const InsertEventTags1 = `
	INSERT INTO event_tags1 
		(EventId, Start, Finish, Text)
	VALUES
		($1, $2, $3, $4)
	ON CONFLICT DO NOTHING;
`

const InsertTags2 = `
	INSERT INTO tags2 
		(Text)
	VALUES
		($1)
	ON CONFLICT (Text) DO UPDATE SET Text = EXCLUDED.Text
	RETURNING (Id);
`

const InsertEventTags2 = `
	INSERT INTO event_tags2 
		(TagId, EventId, Start, Finish)
	VALUES
		($1, $2, $3, $4);
`

const SelectEvents1 = `
	SELECT
		Title, Start, Finish, PostCodes, Tags,
		ST_X(Center) as Lat, 
		ST_Y(Center) as Lon
	FROM events
	WHERE
		((%v <= Start AND %v > Start) OR (%v BETWEEN Start AND Finish)) %v;
`
func makeSelectEvents1(tags []string, start, finish int64) string {
	tagsStr := ""
	if len(tags) > 0 {
		for _, tag := range tags {
			tagsStr += fmt.Sprintf("\n		AND '%v' = ANY (Tags)", tag)
		}
	}
	return fmt.Sprintf(SelectEvents1, start, finish, start, tagsStr)
}


const SelectEvents2 = `
	SELECT  
  		e.Title, e.Start, e.Finish, e.PostCodes, e.Tags as Tags,
  		ST_X(e.Center) as Lat, 
  		ST_Y(e.Center) as Lon
	FROM events e                                                                                 
	INNER JOIN (
		SELECT EventId, array_agg(Text)
		FROM event_tags1 
		WHERE 
			((%v <= Start AND %v > Start) OR (%v BETWEEN Start AND Finish))
			AND Text = ANY ('%v')
		GROUP BY EventId                                                                              
		HAVING true
			%v
	) t ON e.Id = t.EventId;
	
`

func makeSelectEvents2(tags []string, start, finish int64) string {
	tagsArray := `{""}`
	tagsStr := ""
	if len(tags) > 0 {
		tagsStr = fmt.Sprintf("\n		AND '%v' = ANY (array_agg(Text))", tags[0])
		tagsArray = fmt.Sprintf(`"%v"`, tags[0])
		for i := 1; i < len(tags); i++ {
			tagsStr += fmt.Sprintf("\n		AND '%v' = ANY (array_agg(Text))", tags[i])
			tagsArray += fmt.Sprintf(`, "%v"`, tags[i])
		}
		tagsArray = fmt.Sprintf(`{%v}`, tagsArray)
	}
	return fmt.Sprintf(SelectEvents2, start, finish, start, tagsArray, tagsStr)
}
