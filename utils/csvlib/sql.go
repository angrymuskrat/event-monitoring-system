package csvlib

import (
	"fmt"
)

const SelectEventsTemplate = `
	SELECT 
		Title, Start, Finish, PostCodes
	FROM %v %v
`

func makeSelectEventsSQL(eventTableName string, start, finish int64) string {
	timeBorder := "WHERE "
	if start != 0 && finish != 0 {
		timeBorder += fmt.Sprintf("Start >= %v AND Finish <= %v", start, finish)
	} else if start != 0 {
		timeBorder += fmt.Sprintf("Start >= %v", start)
	} else if finish != 0 {
		timeBorder += fmt.Sprintf("Finish <= %v", finish)
	} else {
		timeBorder = ""
	}
	statement := fmt.Sprintf(SelectEventsTemplate, eventTableName, timeBorder)
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

func makeSelectPostsByCodesSQL(shortcodes []string, startTimestamp, endTimestamp int64) string {
	shortcodesSQL := "('"
	for _, code := range shortcodes {
		shortcodesSQL += code + "', '"
	}
	shortcodesSQL += "')"
	statement := fmt.Sprintf(SelectPostsTemplateByCodes, startTimestamp, endTimestamp, shortcodesSQL)
	return statement
}

const UpdatePostsAdvTemplate = `UPDATE posts AS p SET IsAd = adv.IsAd FROM (values %v) as adv(Shortcode, IsAd) WHERE adv.Shortcode = p.Shortcode;`

func makeUpdatePostsAdvSQL(shortcodes []string, isAd []bool) string {
	valuesSQL := ""
	for ind, code := range shortcodes {
		valuesSQL += fmt.Sprintf("('%v', %v), ", code, isAd[ind])
	}
	valuesSQL = valuesSQL[:len(valuesSQL)-2]
	var statement string
	statement = fmt.Sprintf(UpdatePostsAdvTemplate, valuesSQL)
	return statement
}
