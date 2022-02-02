package csvlib

import (
	"fmt"
)

const SelectAllEventsTemplate = `
	SELECT 
		Title, Start, Finish, PostCodes
	FROM %v
`

func makeSelectAllEventsSQL(eventTableName string) string {
	statement := fmt.Sprintf(SelectAllEventsTemplate, eventTableName)
	return statement
}

const SelectPostsTemplateByCodes = `
	SELECT 
		ID, Shortcode, ImageURL, IsVideo, Caption, CommentsCount, Timestamp, LikesCount, IsAd, AuthorID, LocationID, 
		ST_X(Location) as Lon, 
		ST_Y(Location) as Lat,
		NoiseProbability, EventUtility
	FROM posts
	WHERE 
		Timestamp BETWEEN %v AND %v
		AND Shortcode IN %v;
`

const SelectPostsTemplateByCodesOld = `
	SELECT 
		ID, Shortcode, ImageURL, IsVideo, Caption, CommentsCount, Timestamp, LikesCount, IsAd, AuthorID, LocationID, 
		ST_X(Location) as Lon, 
		ST_Y(Location) as Lat
	FROM posts
	WHERE 
		Timestamp BETWEEN %v AND %v
		AND Shortcode IN %v;
`

func makeSelectPostsByCodesSQL(shortcodes []string, startTimestamp, endTimestamp int64, hasNoiseAndUtility bool) string {
	shortcodesSQL := "('"
	for _, code := range shortcodes {
		shortcodesSQL += code + "', '"
	}
	shortcodesSQL += "')"
	var statement string
	if hasNoiseAndUtility {
		statement = fmt.Sprintf(SelectPostsTemplateByCodes, startTimestamp, endTimestamp, shortcodesSQL)
	} else {
		statement = fmt.Sprintf(SelectPostsTemplateByCodesOld, startTimestamp, endTimestamp, shortcodesSQL)
	}
	return statement
}
