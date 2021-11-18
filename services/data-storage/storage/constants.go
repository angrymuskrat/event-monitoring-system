package storage

import (
	"fmt"
	"strings"

	data "github.com/angrymuskrat/event-monitoring-system/services/proto"
)

const Hour = 60 * 60

const PostgresDBName = "postgres"
const GeneralDBName = "general"

const ExtensionTimescaleDBSQL = "CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;"
const ExtensionPostGISSQL = "CREATE EXTENSION IF NOT EXISTS postgis;"
const ExtensionPostGISTopologySQL = "CREATE EXTENSION IF NOT EXISTS postgis_topology;"
const CreateTimeFunctionSQL = "CREATE OR REPLACE FUNCTION unix_now() returns BIGINT LANGUAGE SQL STABLE as $$ SELECT extract(epoch from now())::BIGINT $$;"

const CreateDBTemplate = "CREATE DATABASE %v;"

func makeCreateDBSQL(dbname string) string {
	return fmt.Sprintf(CreateDBTemplate, dbname)
}

const SelectDBTemplate = "SELECT datname FROM pg_catalog.pg_database WHERE lower(datname) = lower('%v');"

func makeSelectDBSQL(dbname string) string {
	return fmt.Sprintf(SelectDBTemplate, dbname)
}

const CreateCitiesTableSQL = `
	CREATE TABLE IF NOT EXISTS cities(
		Title VARCHAR(100),
		Code VARCHAR(50) NOT NULL PRIMARY KEY,
		TopLeft geometry,
		BotRight geometry
	);
`
const InsertCitySQL = `
	INSERT INTO cities
		(Title, Code, TopLeft, BotRight)
	VALUES
		($1, $2, ST_SetSRID( ST_Point($3, $4), 4326), ST_SetSRID( ST_Point($5, $6), 4326));
`
const UpsertCitySQL = `
	INSERT INTO cities 
		(Title, Code, TopLeft, BotRight)
	VALUES 
		($1, $2, ST_SetSRID( ST_Point($3, $4), 4326), ST_SetSRID( ST_Point($5, $6), 4326))
 	ON CONFLICT (Code) DO UPDATE SET Title = EXCLUDED.Title, TopLeft = EXCLUDED.TopLeft, BotRight = EXCLUDED.BotRight;
`
const SelectCitiesSQL = `
	SELECT 
		Title,
		Code,
		ST_X(TopLeft) as tlLon,
		ST_Y(TopLeft) as tlLat,
		ST_X(BotRight) as brLon,
		ST_Y(BotRight) as brLat
	FROM cities;
`
const SelectCityTemplate = `
	SELECT 
		Title,
		Code,
		ST_X(TopLeft) as tlLon,
		ST_Y(TopLeft) as tlLat,
		ST_X(BotRight) as brLon,
		ST_Y(BotRight) as brLat
	FROM cities
	WHERE Code = '%v';
`

func makeSelectCitySQL(cityCode string) string {
	statement := fmt.Sprintf(SelectCityTemplate, cityCode)
	return statement
}

const CreateHyperTablePostsSQL = "SELECT create_hypertable('posts', 'timestamp', chunk_time_interval => 86400, if_not_exists => TRUE);"
const SetTimeFunctionForPostsSQL = "SELECT set_integer_now_func('posts', 'unix_now', replace_if_exists => true);"
const CreatePostsTableSQL = `
	CREATE TABLE IF NOT EXISTS posts(
		ID VARCHAR (30) NOT NULL,
		Shortcode VARCHAR (15) NOT NULL,
		ImageURL TEXT,
		IsVideo BOOLEAN NOT NULL,
		Caption TEXT, 
		CommentsCount BIGINT NOT NULL,
		Timestamp BIGINT NOT NULL,
		LikesCount BIGINT NOT NULL,
		IsAd BOOLEAN NOT NULL,
		AuthorID VARCHAR (15) NOT NULL,
		LocationID VARCHAR (20) NOT NULL,
		Location geometry,
		PRIMARY KEY (Shortcode, Timestamp)
	);
`
const CreatePostsIndexByShortcodeSQL = "CREATE INDEX IF NOT EXISTS shortcode_to_post ON posts (shortcode);"

const InsertPostSQL = `
	INSERT INTO posts
		(ID, Shortcode, ImageURL, IsVideo, Caption, CommentsCount, Timestamp, LikesCount, IsAd, AuthorID, LocationID, Location)
	VALUES 
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, ST_SetSRID( ST_Point($12, $13), 4326))
	ON CONFLICT (Shortcode, Timestamp) DO UPDATE SET Location = EXCLUDED.Location;
`
const SelectPostsTemplate = `
	SELECT 
		ID, Shortcode, ImageURL, IsVideo, Caption, CommentsCount, Timestamp, LikesCount, IsAd, AuthorID, LocationID, 
		ST_X(Location) as Lon, 
		ST_Y(Location) as Lat
	FROM posts
	WHERE Timestamp BETWEEN %v AND %v
`

func makeSelectPostsSQL(startTimestamp, endTimestamp int64) string {
	statement := fmt.Sprintf(SelectPostsTemplate, startTimestamp, endTimestamp)
	return statement
}

const CreateAggrPostsViewSQLTemplate = `
	CREATE MATERIALIZED VIEW aggr_posts
	WITH (timescaledb.continuous)
	AS
	SELECT
 		time_bucket('3600', timestamp) as hour,
 		COUNT(*) as count,
 		ST_Transform(ST_SnapToGrid(ST_Transform(location, 3857), %v, %v), 4326) as center
	FROM posts
	GROUP BY hour, center;
`

func makeCreateAggrPostsViewSQL(gridSize float64) string {
	statement := fmt.Sprintf(CreateAggrPostsViewSQLTemplate, gridSize, gridSize)
	return statement
}

const SelectAggrPostsTemplate = `
	SELECT
		count,
		ST_X(center) as Lon,
		ST_Y(center) as Lat
	FROM aggr_posts
	WHERE hour = %v AND ST_Contains(%v, center); 
`

func makeSelectAggrPostsSQL(interval data.SpatioHourInterval) string {
	poly := makePoly(interval.Area)
	statement := fmt.Sprintf(SelectAggrPostsTemplate, interval.Hour, poly)
	return statement
}

const CreateHyperTableEventsTemplate = "SELECT create_hypertable('%v', 'start', chunk_time_interval => 86400, if_not_exists => TRUE);"

func makeCreateHyperTableEventsSQL(eventTableName string) string {
	statement := fmt.Sprintf(CreateHyperTableEventsTemplate, eventTableName)
	return statement
}

const SetTimeFunctionForEventsTemplate = "SELECT set_integer_now_func('%v', 'unix_now', replace_if_exists => true);"

func makeSetTimeFunctionForEventsSQL(eventTableName string) string {
	statement := fmt.Sprintf(SetTimeFunctionForEventsTemplate, eventTableName)
	return statement
}

const CreateEventsTableTemplate = `
	CREATE TABLE IF NOT EXISTS %v (
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

func makeCreateEventsTableSQL(eventTableName string) string {
	statement := fmt.Sprintf(CreateEventsTableTemplate, eventTableName)
	return statement
}

const InsertEventTemplate = `
	INSERT INTO %v
		(Title, Start, Finish, Center, PostCodes, Tags)
	VALUES
		($1, $2, $3, ST_SetSRID( ST_Point($4, $5), 4326), $6, $7)
`

func makeInsertEventSQL(eventTableName string) string {
	statement := fmt.Sprintf(InsertEventTemplate, eventTableName)
	return statement
}

const SelectEventsTemplate = `
	SELECT 
		Title, Start, Finish, PostCodes, Tags,  
		ST_X(Center) as Lon, 
		ST_Y(Center) as Lat
	FROM %v
	WHERE 
		ST_Covers(%v, Center) 
		AND (Start BETWEEN %v AND (%v - 1))
`

func makeSelectEventsSQL(eventTableName string, interval data.SpatioHourInterval) string {
	poly := makePoly(interval.Area)
	statement := fmt.Sprintf(SelectEventsTemplate, eventTableName, poly, interval.Hour, interval.Hour+Hour)
	return statement
}

const SelectEventsTagsTemplate = `
	SELECT
		Title, Start, Finish, PostCodes, Tags,
		ST_X(Center) as Lon, 
		ST_Y(Center) as Lat
	FROM %v
	WHERE
		(%v <= Finish AND %v >= Start) %v;
`

func makeSelectEventsTagsSQL(eventsTableName string, tags []string, start, finish int64) string {
	tagsStr := ""
	if len(tags) > 0 {
		for _, tag := range tags {
			tagsStr += fmt.Sprintf("\n		AND '%v' = ANY (Tags)", tag)
		}
	}
	return fmt.Sprintf(SelectEventsTagsTemplate, eventsTableName, start, finish, tagsStr)
}

const CreatePostsTimelineViewSQL = `
	CREATE MATERIALIZED VIEW posts_timeline
	WITH (timescaledb.continuous)
	AS
	SELECT
	time_bucket('3600', timestamp) as time,
	COUNT(*) as count
	FROM posts 
	GROUP BY time;
`
const CreateEventsTimelineViewTemplate = `
	CREATE MATERIALIZED VIEW %v_timeline
	WITH (timescaledb.continuous)
	AS
	SELECT
	time_bucket('3600', start) as time,
	COUNT(*) as count
	FROM %v
	GROUP BY time;
`

func makeCreateEventsTimelineViewSQL(eventTableName string) string {
	statement := fmt.Sprintf(CreateEventsTimelineViewTemplate, eventTableName, eventTableName)
	return statement
}

const SelectTimelineTemplate = `
	SELECT
		SUM(posts) as posts,
		SUM(events) as events,
		time
	FROM (
 		SELECT count as posts, 0 as events, time
 		FROM posts_timeline
 		WHERE time BETWEEN %v AND %v
 		UNION
 		SELECT 0 as posts, count as events, time
 		FROM %v_timeline
		WHERE time BETWEEN %v AND %v
	) as tmp
	GROUP BY time;
`

func makeSelectTimelineSQL(startTimestamp, finishTimestamp int64, eventTableName string) string {
	statement := fmt.Sprintf(SelectTimelineTemplate, startTimestamp, finishTimestamp, eventTableName,
		startTimestamp, finishTimestamp)
	return statement
}

const CreateLocationsTableSQL = `
	CREATE TABLE IF NOT EXISTS locations (
		ID VARCHAR(20) NOT NULL PRIMARY KEY,
		Position geometry,
		Title TEXT,
		Slug TEXT
	);
`
const InsertLocationSQL = `
	INSERT INTO locations 
		(ID, Position, Title, Slug)
	VALUES
		($1, ST_SetSRID( ST_Point($2, $3), 4326), $4, $5)
	ON CONFLICT (ID) DO NOTHING;
`
const SelectLocationsSQL = `
	SELECT 
		ID, Title, Slug,
		ST_X(Position) as Lon,
		ST_Y(Position) as Lat
	FROM locations;
`

const CreateGridsTableSQL = `
	CREATE TABLE IF NOT EXISTS grids(
		ID BIGINT PRIMARY KEY,
		Blob BYTEA NOT NULL
	);
`
const InsertGridSQL = `
	INSERT INTO grids(id, blob)
	VALUES ($1, $2)
	ON CONFLICT (id) DO UPDATE SET blob = EXCLUDED.blob;
`

const SelectShortPostsInIntervalTemplate = `
	SELECT 
		Shortcode, Caption, CommentsCount, LikesCount, Timestamp, AuthorID, LocationID,
		ST_X(Location) as Lon, 
		ST_Y(Location) as Lat
	FROM posts
	WHERE Timestamp BETWEEN %v and %v
		AND Shortcode IN %v;
`

func makeSelectShortPostsInIntervalSQL(shortcodes []string, startTimestamp int64, endTimestamp int64) string {
	shortcodesSQL := "('"
	for _, code := range shortcodes {
		shortcodesSQL += code + "', '"
	}
	shortcodesSQL += "')"

	sqlRequest := fmt.Sprintf(SelectShortPostsInIntervalTemplate, startTimestamp, endTimestamp, shortcodesSQL)
	return sqlRequest
}

const makeSelectSinglePostTemplate = `
	SELECT 
		Shortcode, Caption, CommentsCount, LikesCount, Timestamp, AuthorID, LocationID,
		ST_X(Location) as Lon, 
		ST_Y(Location) as Lat
	FROM posts 
	WHERE shortcode = '%v';
`

func makeSelectSinglePostSQL(shortcode string) string {
	sqlRequest := fmt.Sprintf(makeSelectSinglePostTemplate, shortcode)
	return sqlRequest
}

func makePoly(area data.Area) string {
	return fmt.Sprintf("ST_Polygon('LINESTRING(%v %v, %v %v, %v %v, %v %v, %v %v)'::geometry, 4326)",
		area.TopLeft.Lon, area.TopLeft.Lat,
		area.TopLeft.Lon, area.BotRight.Lat,
		area.BotRight.Lon, area.BotRight.Lat,
		area.BotRight.Lon, area.TopLeft.Lat,
		area.TopLeft.Lon, area.TopLeft.Lat)
}

func isNotAlreadyExistsError(err error) bool {
	return !strings.Contains(err.Error(), "already exists")
}
