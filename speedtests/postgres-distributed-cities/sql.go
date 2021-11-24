package main

import "fmt"

const SelectPostsSQLTemplate = `
	SELECT 
		ID, Shortcode, ImageURL, IsVideo, Caption, CommentsCount, Timestamp, LikesCount, IsAd, AuthorID, LocationID, 
		ST_X(Location) as Lon, 
		ST_Y(Location) as Lat
	FROM posts
	WHERE Timestamp BETWEEN %v AND %v`

const SelectPostsCountSQLTemplate = `
	SELECT 
		count(*) as count
	FROM posts
	WHERE Timestamp BETWEEN %v AND %v`

func makeSelectPostsSQL(startTimestamp, endTimestamp int64) string {
	request := fmt.Sprintf(SelectPostsSQLTemplate, startTimestamp, endTimestamp)
	return request
}

func makeSelectPostsCountSQL(startTimestamp, endTimestamp int64) string {
	request := fmt.Sprintf(SelectPostsCountSQLTemplate, startTimestamp, endTimestamp)
	return request
}
