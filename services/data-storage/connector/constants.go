package connector

import (
	"fmt"
	data "github.com/angrymuskrat/event-monitoring-system/services/proto"
)

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

const GridTable = `
	CREATE TABLE IF NOT EXISTS grids(
		ID VARCHAR (100) NOT NULL PRIMARY KEY,
		Blob BYTEA NOT NULL
	);`

const SelectAggrPostsTemplate  = `
	SELECT
		count,
		ST_X(center) as lat,
		ST_Y(center) as lon
	FROM aggr_posts
	WHERE hour = %v AND ST_Contains(%v, center); 
`

func makePoly(topLeft, botRight *data.Point) string {
	return fmt.Sprintf("ST_Polygon('LINESTRING(%v %v, %v %v, %v %v, %v %v, %v %v)'::geometry, 4326)",
		topLeft.Lat, topLeft.Lon,
		topLeft.Lat, botRight.Lon,
		botRight.Lat, botRight.Lon,
		botRight.Lat, topLeft.Lon,
		topLeft.Lat, topLeft.Lon)
}