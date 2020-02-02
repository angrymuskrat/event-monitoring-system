package parser

import (
	"encoding/json"
	"fmt"
	"github.com/angrymuskrat/event-monitoring-system/services/insta-crawler/crawler/data"
)

func ParseFromLocationRequest(input []byte) ([]data.Post, data.Location, string, bool, int64, error) {
	var d map[string]interface{}
	var timestamp int64
	err := json.Unmarshal(input, &d)
	if err != nil {
		return nil, data.Location{}, "", false, timestamp, fmt.Errorf("Unable to parse response: %s", err)
	}
	root, succeed := d["data"]
	if !succeed || root == nil {
		root, succeed = d["graphql"]
		if !succeed || root == nil {
			return nil, data.Location{}, "", false, timestamp, fmt.Errorf("Unable to get 'graphql' element: %s", string(input))
		}
	}
	var rawEntity interface{}
	rawEntity, succeed = root.(map[string]interface{})["location"]
	if !succeed || rawEntity == nil {
		return nil, data.Location{}, "", false, timestamp, fmt.Errorf("Unable to get 'location' element: %s", string(input))
	}

	location := data.Location{}
	rawLocID, succeed := rawEntity.(map[string]interface{})["id"]
	if !succeed || rawLocID == nil {
		return nil, data.Location{}, "", false, timestamp, fmt.Errorf("Unable to get 'id' element")
	}
	location.ID = rawLocID.(string)
	rawLocName, succeed := rawEntity.(map[string]interface{})["name"]
	if !succeed || rawLocName == nil {
		return nil, data.Location{}, "", false, timestamp, fmt.Errorf("Unable to get 'name' element")
	}
	location.Title = rawLocName.(string)
	rawLocSlug, succeed := rawEntity.(map[string]interface{})["slug"]
	if !succeed || rawLocSlug == nil {
		return nil, data.Location{}, "", false, timestamp, fmt.Errorf("Unable to get 'slug' element")
	}
	location.Slug = rawLocSlug.(string)
	rawLocLat, succeed := rawEntity.(map[string]interface{})["lat"]
	if !succeed || rawLocLat == nil {
		location.Lat = 0
	} else {
		location.Lat = rawLocLat.(float64)
	}
	rawLocLon, succeed := rawEntity.(map[string]interface{})["lng"]
	if !succeed || rawLocLon == nil {
		location.Lon = 0
	} else {
		location.Lon = rawLocLon.(float64)
	}

	var rawEdges interface{}
	rawEdges, succeed = rawEntity.(map[string]interface{})["edge_location_to_media"]
	if !succeed || rawEdges == nil {
		return nil, location, "", false, timestamp, fmt.Errorf("Unable to get 'edge_location_to_media' element")
	}
	rawPageInfo, succeed := rawEdges.(map[string]interface{})["page_info"]
	if !succeed || rawPageInfo == nil {
		return nil, location, "", false, timestamp, fmt.Errorf("Unable to get 'page_info' element")
	}
	rawNextPage, succeed := rawPageInfo.(map[string]interface{})["has_next_page"]
	if !succeed || rawNextPage == nil {
		return nil, location, "", false, timestamp, fmt.Errorf("Unable to get 'has_next_page' element")
	}
	nextPage := rawNextPage.(bool)
	endCursor := ""

	rawEdgesArray, succeed := rawEdges.(map[string]interface{})["edges"]
	if !succeed || rawEdgesArray == nil {
		return nil, location, "", false, timestamp, fmt.Errorf("Unable to get 'edges' element")
	}
	edgesArray := rawEdgesArray.([]interface{})
	posts := []data.Post{}
	if len(edgesArray) > 0 {
		for _, edge := range edgesArray {
			post := data.Post{}
			rawNode, succeed := edge.(map[string]interface{})["node"]
			if !succeed || rawNode == nil {
				return nil, location, "", false, timestamp, fmt.Errorf("Unable to get 'node' element")
			}
			rawID, succeed := rawNode.(map[string]interface{})["id"]
			if !succeed || rawID == nil {
				return nil, location, "", false, timestamp, fmt.Errorf("Unable to get 'id' element")
			}
			post.ID = rawID.(string)

			rawCode, succeed := rawNode.(map[string]interface{})["shortcode"]
			if !succeed || rawCode == nil {
				return nil, location, "", false, timestamp, fmt.Errorf("Unable to get 'shortcode' element")
			}
			post.Shortcode = rawCode.(string)

			rawURL, succeed := rawNode.(map[string]interface{})["display_url"]
			if !succeed || rawURL == nil {
				return nil, location, "", false, timestamp, fmt.Errorf("Unable to get 'display_url' element")
			}
			post.ImageURL = rawURL.(string)

			rawIsVideo, succeed := rawNode.(map[string]interface{})["is_video"]
			if !succeed || rawIsVideo == nil {
				return nil, location, "", false, timestamp, fmt.Errorf("Unable to get 'is_video' element")
			}
			post.IsVideo = rawIsVideo.(bool)

			rawCaptionNode, succeed := rawNode.(map[string]interface{})["edge_media_to_caption"]
			if !succeed || rawCaptionNode == nil {
				return nil, location, "", false, timestamp, fmt.Errorf("Unable to get 'edge_media_to_caption' element")
			}
			rawCaptionEdges, succeed := rawCaptionNode.(map[string]interface{})["edges"]
			if !succeed || rawCaptionEdges == nil {
				return nil, location, "", false, timestamp, fmt.Errorf("Unable to get 'edges' element")
			}
			captionEdgesArray := rawCaptionEdges.([]interface{})
			if len(captionEdgesArray) > 0 {
				rawCaptionEdge, succeed := captionEdgesArray[0].(map[string]interface{})["node"]
				if !succeed || rawCaptionEdge == nil {
					return nil, location, "", false, timestamp, fmt.Errorf("Unable to get 'node' element")
				}
				rawCaption, succeed := rawCaptionEdge.(map[string]interface{})["text"]
				if !succeed || rawCaption == nil {
					return nil, location, "", false, timestamp, fmt.Errorf("Unable to get 'text' element")
				}
				post.Caption = rawCaption.(string)
			}

			rawTimestamp, succeed := rawNode.(map[string]interface{})["taken_at_timestamp"]
			if !succeed || rawTimestamp == nil {
				return nil, location, "", false, timestamp, fmt.Errorf("Unable to get 'taken_at_timestamp' element")
			}
			timeString := rawTimestamp.(float64)
			timeInt := int64(timeString)
			timestamp = timeInt
			post.Timestamp = timeInt

			rawComment, succeed := rawNode.(map[string]interface{})["edge_media_to_comment"]
			if !succeed || rawComment == nil {
				continue
			}
			rawCount, succeed := rawComment.(map[string]interface{})["count"]
			if !succeed || rawCount == nil {
				continue
			}
			countString := rawCount.(float64)
			countInt := int(countString)
			post.CommentsCount = countInt

			rawLikes, succeed := rawNode.(map[string]interface{})["edge_liked_by"]
			if !succeed || rawLikes == nil {
				continue
			}
			rawCount, succeed = rawLikes.(map[string]interface{})["count"]
			if !succeed || rawCount == nil {
				continue
			}
			countString = rawCount.(float64)
			countInt = int(countString)
			post.LikesCount = countInt

			rawOwner, succeed := rawNode.(map[string]interface{})["owner"]
			if !succeed || rawOwner == nil {
				continue
			}
			rawOwnerID, succeed := rawOwner.(map[string]interface{})["id"]
			if !succeed || rawOwnerID == nil {
				continue
			}
			post.AuthorID = rawOwnerID.(string)

			post.LocationID = location.ID
			post.Lat = location.Lat
			post.Lon = location.Lon

			endCursor = post.ID
			posts = append(posts, post)
		}

	}
	return posts, location, endCursor, nextPage, timestamp, nil
}
