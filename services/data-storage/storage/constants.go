package storage

import (
	"fmt"
	data "github.com/angrymuskrat/event-monitoring-system/services/proto"
)

const ExtensionTimescaleDB = "CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;"
const ExtensionPostGIS = "CREATE EXTENSION IF NOT EXISTS postgis;"
const ExtensionPostGISTopology = "CREATE EXTENSION IF NOT EXISTS postgis_topology;"

const PostTable = `
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
	);`

const CreateHyperTablePosts = "SELECT create_hypertable('posts', 'timestamp', chunk_time_interval => 86400, if_not_exists => TRUE);"
const CreateTimeFunction = "CREATE OR REPLACE FUNCTION unix_now() returns BIGINT LANGUAGE SQL STABLE as $$ SELECT extract(epoch from now())::BIGINT $$;"
const SetTimeFunctionForPosts = "SELECT set_integer_now_func('posts', 'unix_now', replace_if_exists => true);"

const DropAggregationPosts = "DROP VIEW IF EXISTS aggr_posts CASCADE;"
const AggregationPosts = `
	CREATE VIEW aggr_posts
	WITH (timescaledb.continuous)
	AS
	SELECT
 		time_bucket('3600', timestamp) as hour,
 		COUNT(*) as count,
 		ST_Transform(ST_SnapToGrid(ST_Transform(location, 3857), %v, %v), 4326) as center
	FROM posts
	GROUP BY hour, center;`

const EventsTable = `
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
const CreateHyperTableEvents = "SELECT create_hypertable('events', 'start', chunk_time_interval => 86400, if_not_exists => TRUE);"
const SetTimeFunctionForEvents = "SELECT set_integer_now_func('events', 'unix_now', replace_if_exists => true);"

const CitiesTable = `
	CREATE TABLE IF NOT EXISTS cities (
		Title VARCHAR(50),
		ID VARCHAR(20) NOT NULL PRIMARY KEY,
		Location geometry
	);
`
const LocationsTable = `
	CREATE TABLE IF NOT EXISTS locations (
		ID VARCHAR(20) NOT NULL PRIMARY KEY,
		CityId VARCHAR(20),
		Position geometry,
		Title VARCHAR(100),
		Slug VARCHAR(100),
		FOREIGN KEY (CityId) REFERENCES cities (ID)
	);`

const GridTable = `
	CREATE TABLE IF NOT EXISTS grids(
		ID VARCHAR (100) NOT NULL PRIMARY KEY,
		Blob BYTEA NOT NULL
	);`

const PushPostsTemplate = `
	INSERT INTO posts
		(ID, Shortcode, ImageURL, IsVideo, Caption, CommentsCount, Timestamp, LikesCount, IsAd, AuthorID, LocationID, Location)
	VALUES 
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, ST_SetSRID( ST_Point($12, $13), 4326))
	ON CONFLICT (Shortcode, Timestamp) DO NOTHING;
`

const SelectPostsTemplate = `
	SELECT 
		ID, Shortcode, ImageURL, IsVideo, Caption, CommentsCount, Timestamp, LikesCount, IsAd, AuthorID, LocationID, 
		ST_X(Location) as Lat, 
		ST_Y(Location) as Lon
	FROM posts
	WHERE 
		ST_Covers(%v, Location) 
		AND (Timestamp BETWEEN %v AND %v)
	`

const SelectAggrPostsTemplate = `
	SELECT
		count,
		ST_X(center) as lat,
		ST_Y(center) as lon
	FROM aggr_posts
	WHERE hour = %v AND ST_Contains(%v, center); 
`

const PushGridTemplate = `
	INSERT INTO grids(ID, Blob)
	VALUES ($1, $2);
`

const PullGridTemplate = `SELECT Blob FROM grids WHERE '%v' = Id;`

const PushEventsTemplate = `
	INSERT INTO events
		(Title, Start, Finish, Center, PostCodes, Tags)
	VALUES
		($1, $2, $3, ST_SetSRID( ST_Point($4, $5), 4326), $6, $7)
`
const Hour = 60 * 60
const SelectEventsTemplate = `
	SELECT 
		Title, Start, Finish, PostCodes, Tags,  
		ST_X(Center) as Lat, 
		ST_Y(Center) as Lon
	FROM events
	WHERE 
		ST_Covers(%v, Center) 
		AND (Start BETWEEN %v AND %v)
	`

const PushCityIfNotExists = `
	INSERT INTO cities
		(Title, Id)
	VALUES
		($1, $2)
	ON CONFLICT (Id) DO NOTHING;
`

const PushLocationTemplate = `
	INSERT INTO locations 
		(ID, CityId, Position, Title, Slug)
	VALUES
		($1, $2, ST_SetSRID( ST_Point($3, $4), 4326), $5, $6)
	ON CONFLICT (ID) DO NOTHING;
`

const SelectLocationsTemplate = `
	SELECT 
		ID, Title, Slug,
		ST_X(Position) as Lat,
		ST_Y(Position) as Lon
	FROM locations
	WHERE CityId = '%v';
`

func makePoly(topLeft, botRight data.Point) string {
	return fmt.Sprintf("ST_Polygon('LINESTRING(%v %v, %v %v, %v %v, %v %v, %v %v)'::geometry, 4326)",
		topLeft.Lat, topLeft.Lon,
		topLeft.Lat, botRight.Lon,
		botRight.Lat, botRight.Lon,
		botRight.Lat, topLeft.Lon,
		topLeft.Lat, topLeft.Lon)
}
