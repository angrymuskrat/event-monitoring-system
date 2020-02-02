package parser

import (
	"encoding/json"
	"fmt"
	"github.com/angrymuskrat/event-monitoring-system/services/insta-crawler/crawler/data"
)

func ParseFromProfileRequest(input []byte) ([]data.Post, data.Profile, string, bool, int64, error) {
	var d map[string]interface{}
	var timestamp int64
	err := json.Unmarshal(input, &d)
	if err != nil {
		return nil, data.Profile{}, "", false, timestamp, fmt.Errorf("Unable to parse response: %s \n %s", string(input), err)
	}
	hasProfileInfo := true
	root, succeed := d["graphql"]
	if !succeed || root == nil {
		hasProfileInfo = false
		root, succeed = d["data"]
		if !succeed || root == nil {
			return nil, data.Profile{}, "", false, timestamp, fmt.Errorf("Unable to get 'data' element: %s", string(input))
		}
	}
	var rawEntity interface{}
	rawEntity, succeed = root.(map[string]interface{})["user"]
	if !succeed || rawEntity == nil {
		return nil, data.Profile{}, "", false, timestamp, fmt.Errorf("Unable to get 'user' element: %s", string(input))
	}

	profile := data.Profile{}
	if hasProfileInfo {
		rawProfileID, succeed := rawEntity.(map[string]interface{})["id"]
		if !succeed || rawProfileID == nil {
			return nil, data.Profile{}, "", false, timestamp, fmt.Errorf("Unable to get 'id' element")
		}
		profile.ID = rawProfileID.(string)
		rawProfileBio, succeed := rawEntity.(map[string]interface{})["biography"]
		if !succeed || rawProfileBio == nil {
			return nil, data.Profile{}, "", false, timestamp, fmt.Errorf("Unable to get 'biography' element")
		}
		profile.Biography = rawProfileBio.(string)
		rawProfileUsername, succeed := rawEntity.(map[string]interface{})["username"]
		if !succeed || rawProfileUsername == nil {
			return nil, data.Profile{}, "", false, timestamp, fmt.Errorf("Unable to get 'username' element")
		}
		profile.Username = rawProfileUsername.(string)
		rawProfileFullname, succeed := rawEntity.(map[string]interface{})["full_name"]
		if !succeed || rawProfileFullname == nil {
			return nil, data.Profile{}, "", false, timestamp, fmt.Errorf("Unable to get 'full_name' element")
		}
		profile.FullName = rawProfileFullname.(string)
		rawProfileVerified, succeed := rawEntity.(map[string]interface{})["is_verified"]
		if !succeed || rawProfileVerified == nil {
			return nil, data.Profile{}, "", false, timestamp, fmt.Errorf("Unable to get 'is_verified' element")
		}
		profile.Verified = rawProfileVerified.(bool)
		rawProfilePrivate, succeed := rawEntity.(map[string]interface{})["is_private"]
		if !succeed || rawProfilePrivate == nil {
			return nil, data.Profile{}, "", false, timestamp, fmt.Errorf("Unable to get 'is_private' element")
		}
		profile.Private = rawProfilePrivate.(bool)
		rawFollowedBy, succeed := rawEntity.(map[string]interface{})["edge_followed_by"]
		if !succeed || rawFollowedBy == nil {
			return nil, data.Profile{}, "", false, timestamp, fmt.Errorf("Unable to get 'edge_followed_by' element")
		}
		rawCount, succeed := rawFollowedBy.(map[string]interface{})["count"]
		if !succeed || rawCount == nil {
			return nil, data.Profile{}, "", false, timestamp, fmt.Errorf("Unable to get 'count' element")
		}
		countString := rawCount.(float64)
		countInt := int(countString)
		profile.FollowersCount = countInt
		rawFollows, succeed := rawEntity.(map[string]interface{})["edge_follow"]
		if !succeed || rawFollows == nil {
			return nil, data.Profile{}, "", false, timestamp, fmt.Errorf("Unable to get 'edge_follow' element")
		}
		rawCount, succeed = rawFollows.(map[string]interface{})["count"]
		if !succeed || rawCount == nil {
			return nil, data.Profile{}, "", false, timestamp, fmt.Errorf("Unable to get 'count' element")
		}
		countString = rawCount.(float64)
		countInt = int(countString)
		profile.FollowsCount = countInt
	}

	var rawEdges interface{}
	rawEdges, succeed = rawEntity.(map[string]interface{})["edge_owner_to_timeline_media"]
	if !succeed || rawEdges == nil {
		return nil, data.Profile{}, "", false, timestamp, fmt.Errorf("Unable to get 'edge_owner_to_timeline_media' element")
	}
	rawPageInfo, succeed := rawEdges.(map[string]interface{})["page_info"]
	if !succeed || rawPageInfo == nil {
		return nil, data.Profile{}, "", false, timestamp, fmt.Errorf("Unable to get 'page_info' element")
	}
	// rawCursor, succeed := rawPageInfo.(map[string]interface{})["end_cursor"]
	// if !succeed {
	// 	return nil, data.Profile{}, "", false, timestamp, fmt.Errorf("Unable to get 'end_cursor' element: %s", string(input))
	// }
	rawNextPage, succeed := rawPageInfo.(map[string]interface{})["has_next_page"]
	if !succeed || rawNextPage == nil {
		return nil, data.Profile{}, "", false, timestamp, fmt.Errorf("Unable to get 'has_next_page' element")
	}
	nextPage := rawNextPage.(bool)
	endCursor := ""
	// endCursor := ""
	// if rawCursor != nil {
	// 	endCursor = rawCursor.(string)
	// }

	rawEdgesArray, succeed := rawEdges.(map[string]interface{})["edges"]
	if !succeed || rawEdgesArray == nil {
		return nil, data.Profile{}, "", false, timestamp, fmt.Errorf("Unable to get 'edges' element")
	}
	edgesArray := rawEdgesArray.([]interface{})
	posts := []data.Post{}
	if len(edgesArray) > 0 {
		for _, edge := range edgesArray {
			post := data.Post{}
			rawNode, succeed := edge.(map[string]interface{})["node"]
			if !succeed || rawNode == nil {
				return nil, data.Profile{}, "", false, timestamp, fmt.Errorf("Unable to get 'node' element")
			}
			rawID, succeed := rawNode.(map[string]interface{})["id"]
			if !succeed || rawID == nil {
				return nil, data.Profile{}, "", false, timestamp, fmt.Errorf("Unable to get 'id' element")
			}
			post.ID = rawID.(string)

			rawCode, succeed := rawNode.(map[string]interface{})["shortcode"]
			if !succeed || rawCode == nil {
				return nil, data.Profile{}, "", false, timestamp, fmt.Errorf("Unable to get 'shortcode' element")
			}
			post.Shortcode = rawCode.(string)

			rawURL, succeed := rawNode.(map[string]interface{})["display_url"]
			if !succeed || rawURL == nil {
				return nil, data.Profile{}, "", false, timestamp, fmt.Errorf("Unable to get 'display_url' element")
			}
			post.ImageURL = rawURL.(string)

			rawIsVideo, succeed := rawNode.(map[string]interface{})["is_video"]
			if !succeed || rawIsVideo == nil {
				return nil, data.Profile{}, "", false, timestamp, fmt.Errorf("Unable to get 'is_video' element")
			}
			post.IsVideo = rawIsVideo.(bool)

			rawCaptionNode, succeed := rawNode.(map[string]interface{})["edge_media_to_caption"]
			if !succeed || rawCaptionNode == nil {
				return nil, profile, "", false, timestamp, fmt.Errorf("Unable to get 'edge_media_to_caption' element")
			}
			rawCaptionEdges, succeed := rawCaptionNode.(map[string]interface{})["edges"]
			if !succeed || rawCaptionEdges == nil {
				return nil, profile, "", false, timestamp, fmt.Errorf("Unable to get 'edges' element")
			}
			captionEdgesArray := rawCaptionEdges.([]interface{})
			if len(captionEdgesArray) > 0 {
				rawCaptionEdge, succeed := captionEdgesArray[0].(map[string]interface{})["node"]
				if !succeed || rawCaptionEdge == nil {
					return nil, profile, "", false, timestamp, fmt.Errorf("Unable to get 'node' element")
				}
				rawCaption, succeed := rawCaptionEdge.(map[string]interface{})["text"]
				if !succeed || rawCaption == nil {
					return nil, profile, "", false, timestamp, fmt.Errorf("Unable to get 'text' element")
				}
				post.Caption = rawCaption.(string)
			}

			rawTimestamp, succeed := rawNode.(map[string]interface{})["taken_at_timestamp"]
			if !succeed || rawTimestamp == nil {
				return nil, data.Profile{}, "", false, timestamp, fmt.Errorf("Unable to get 'taken_at_timestamp' element")
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
				rawLikes, succeed = rawNode.(map[string]interface{})["edge_media_preview_like"]
				if !succeed || rawLikes == nil {
					continue
				}
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

			endCursor = post.ID

			posts = append(posts, post)
		}

	}
	return posts, profile, endCursor, nextPage, timestamp, nil
}
