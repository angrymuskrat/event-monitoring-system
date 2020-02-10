package storage

import (
	"fmt"
	data "github.com/angrymuskrat/event-monitoring-system/services/proto"
	"strings"
)

const Hour = 60 * 60



const PostgresDBName = "postgres"
const GeneralDBName = "general"



const ExtensionTimescaleDB = "CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;"
const ExtensionPostGIS = "CREATE EXTENSION IF NOT EXISTS postgis;"
const ExtensionPostGISTopology = "CREATE EXTENSION IF NOT EXISTS postgis_topology;"
const CreateTimeFunction = "CREATE OR REPLACE FUNCTION unix_now() returns BIGINT LANGUAGE SQL STABLE as $$ SELECT extract(epoch from now())::BIGINT $$;"



const CreateDB = "CREATE DATABASE %v;"
const SelectDB = "SELECT datname FROM pg_catalog.pg_database WHERE lower(datname) = lower('%v');"



const CreateCitiesTable = `
	CREATE TABLE IF NOT EXISTS cities(
		Title VARCHAR(100),
		Code VARCHAR(50) NOT NULL PRIMARY KEY,
		TopLeft geometry,
		BotRight geometry
	);
`
const InsertCity = `
	INSERT INTO cities
		(Title, Code, TopLeft, BotRight)
	VALUES
		($1, $2, ST_SetSRID( ST_Point($3, $4), 4326), ST_SetSRID( ST_Point($5, $6), 4326));
`
const UpsertCity = `
	INSERT INTO cities 
		(Title, Code, TopLeft, BotRight)
	VALUES 
		($1, $2, ST_SetSRID( ST_Point($3, $4), 4326), ST_SetSRID( ST_Point($5, $6), 4326))
 	ON CONFLICT (Code) DO UPDATE SET Title = EXCLUDED.Title, TopLeft = EXCLUDED.TopLeft, BotRight = EXCLUDED.BotRight;
`
const SelectCities = `
	SELECT 
		Title,
		Code,
		ST_X(TopLeft) as tlLat,
		ST_Y(TopLeft) as tlLon,
		ST_X(BotRight) as brLat,
		ST_Y(BotRight) as brLon
	FROM cities;
`
const SelectCity = `
	SELECT 
		Title,
		Code,
		ST_X(TopLeft) as tlLat,
		ST_Y(TopLeft) as tlLon,
		ST_X(BotRight) as brLat,
		ST_Y(BotRight) as brLon
	FROM cities
	WHERE Code = '%v';
`


const CreateHyperTablePosts = "SELECT create_hypertable('posts', 'timestamp', chunk_time_interval => 86400, if_not_exists => TRUE);"
const SetTimeFunctionForPosts = "SELECT set_integer_now_func('posts', 'unix_now', replace_if_exists => true);"
const CreatePostsTable = `
	CREATE TABLE IF NOT EXISTS posts(
		ID VARCHAR (30) NOT NULL,
		Shortcode VARCHAR (15) NOT NULL,
		ImageURL TEXT,
		IsVideo BOOLEAN NOT NULL,
		Caption TEXT, -- max size of text in Instagram
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
const InsertPost = `
	INSERT INTO posts
		(ID, Shortcode, ImageURL, IsVideo, Caption, CommentsCount, Timestamp, LikesCount, IsAd, AuthorID, LocationID, Location)
	VALUES 
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, ST_SetSRID( ST_Point($12, $13), 4326))
	ON CONFLICT (Shortcode, Timestamp) DO NOTHING;
`
const SelectPosts = `
	SELECT 
		ID, Shortcode, ImageURL, IsVideo, Caption, CommentsCount, Timestamp, LikesCount, IsAd, AuthorID, LocationID, 
		ST_X(Location) as Lat, 
		ST_Y(Location) as Lon
	FROM posts
	WHERE Timestamp BETWEEN %v AND %v
`



const CreateAggrPostsView = `
	CREATE VIEW aggr_posts
	WITH (timescaledb.continuous)
	AS
	SELECT
 		time_bucket('3600', timestamp) as hour,
 		COUNT(*) as count,
 		ST_Transform(ST_SnapToGrid(ST_Transform(location, 3857), %v, %v), 4326) as center
	FROM posts
	GROUP BY hour, center;
`
const SelectAggrPosts = `
	SELECT
		count,
		ST_X(center) as lat,
		ST_Y(center) as lon
	FROM aggr_posts
	WHERE hour = %v AND ST_Contains(%v, center); 
`


const CreateHyperTableEvents = "SELECT create_hypertable('events', 'start', chunk_time_interval => 86400, if_not_exists => TRUE);"
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
`
const SelectEvents = `
	SELECT 
		Title, Start, Finish, PostCodes, Tags,  
		ST_X(Center) as Lat, 
		ST_Y(Center) as Lon
	FROM events
	WHERE 
		ST_Covers(%v, Center) 
		AND (Start BETWEEN %v AND %v)
`


const CreatePostsTimelineView = `
	CREATE VIEW posts_timeline
	WITH (timescaledb.continuous)
	AS
	SELECT
	time_bucket('3600', timestamp) as time,
	COUNT(*) as count
	FROM posts 
	GROUP BY time;
`
const CreateEventsTimelineView = `
	CREATE VIEW events_timeline
	WITH (timescaledb.continuous)
	AS
	SELECT
	time_bucket('3600', start) as time,
	COUNT(*) as count
	FROM events 
	GROUP BY time;
`
const SelectTimeline = `
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
 		FROM events_timeline
		WHERE time BETWEEN %v AND %v
	) as tmp
	GROUP BY time;
`


const CreateLocationsTable = `
	CREATE TABLE IF NOT EXISTS locations (
		ID VARCHAR(20) NOT NULL PRIMARY KEY,
		Position geometry,
		Title VARCHAR(100),
		Slug VARCHAR(100)
	);
`
const InsertLocation = `
	INSERT INTO locations 
		(ID, Position, Title, Slug)
	VALUES
		($1, ST_SetSRID( ST_Point($2, $3), 4326), $4, $5)
	ON CONFLICT (ID) DO NOTHING;
`
const SelectLocations = `
	SELECT 
		ID, Title, Slug,
		ST_X(Position) as Lat,
		ST_Y(Position) as Lon
	FROM locations;
`



const CreateGridsTable = `
	CREATE TABLE IF NOT EXISTS grids(
		ID BIGINT PRIMARY KEY,
		Blob BYTEA NOT NULL
	);
`
const InsertGrid = `
	INSERT INTO grids(id, blob)
	VALUES ($1, $2)
	ON CONFLICT (id) DO UPDATE SET blob = EXCLUDED.blob;
`
const SelectGrid = `
	SELECT id, blob 
	FROM grids 
	WHERE id BETWEEN %v AND %v;
`




func MakePoly(area data.Area) string {
	return fmt.Sprintf("ST_Polygon('LINESTRING(%v %v, %v %v, %v %v, %v %v, %v %v)'::geometry, 4326)",
		area.TopLeft.Lat, area.TopLeft.Lon,
		area.TopLeft.Lat, area.BotRight.Lon,
		area.BotRight.Lat, area.BotRight.Lon,
		area.BotRight.Lat, area.TopLeft.Lon,
		area.TopLeft.Lat, area.TopLeft.Lon)
}

func MakeCreateDB(name string) string {
	return fmt.Sprintf(CreateDB, name)
}

func MakeSelectDB(name string) string {
	return fmt.Sprintf(SelectDB, name)
}

func IsNotAlreadyExistsError(err error) bool {
	return !strings.Contains(err.Error(), "already exists")
}
