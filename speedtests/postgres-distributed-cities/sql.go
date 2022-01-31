package main

import (
	"fmt"
)

const SelectPostsTemplate = `
	SELECT 
		ID, Shortcode, ImageURL, IsVideo, Caption, CommentsCount, Timestamp, LikesCount, IsAd, AuthorID, LocationID, 
		ST_X(Location) as Lon, 
		ST_Y(Location) as Lat
	FROM posts
	WHERE Timestamp BETWEEN %v AND %v`

const SelectPostsCountTemplate = `
	SELECT 
		count(*) as count
	FROM posts
	WHERE Timestamp BETWEEN %v AND %v`

const SelectEventsTemplate = `
	SELECT 
		Title, Start, Finish, PostCodes, Tags,  
		ST_X(Center) as Lon, 
		ST_Y(Center) as Lat
	FROM events_6
	WHERE
		Start BETWEEN %v AND (%v - 1)`

const SelectEventsCountTemplate = `
	SELECT 
		COUNT(*) as count
	FROM events_6
	WHERE
		Start BETWEEN %v AND %v`

func buildSelectRequest(template string, startTimestamp, endTimestamp int64) string {
	statement := fmt.Sprintf(template, startTimestamp, endTimestamp)
	return statement
}

