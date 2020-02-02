package parser

import (
	"encoding/json"
	"fmt"
	"github.com/angrymuskrat/event-monitoring-system/services/insta-crawler/crawler/data"
)

func ParseFromPostRequest(input []byte) (data.Post, error) {
	var d map[string]interface{}
	err := json.Unmarshal(input, &d)
	if err != nil {
		return data.Post{}, fmt.Errorf("Unable to parse response: %s \n %s", string(input), err)
	}
	root, succeed := d["graphql"]
	if !succeed || root == nil {
		return data.Post{}, fmt.Errorf("Unable to get 'graphql' element: %s", string(input))
	}
	rawNode, succeed := root.(map[string]interface{})["shortcode_media"]
	if !succeed || rawNode == nil {
		return data.Post{}, fmt.Errorf("Unable to get 'shortcode_media' element: %s", string(input))
	}
	post := data.Post{}
	rawID, succeed := rawNode.(map[string]interface{})["id"]
	if !succeed || rawID == nil {
		return data.Post{}, fmt.Errorf("Unable to get 'id' element")
	}
	post.ID = rawID.(string)

	rawCode, succeed := rawNode.(map[string]interface{})["shortcode"]
	if !succeed || rawCode == nil {
		return data.Post{}, fmt.Errorf("Unable to get 'shortcode' element")
	}
	post.Shortcode = rawCode.(string)

	rawURL, succeed := rawNode.(map[string]interface{})["display_url"]
	if !succeed || rawURL == nil {
		return data.Post{}, fmt.Errorf("Unable to get 'display_url' element")
	}
	post.ImageURL = rawURL.(string)

	rawIsVideo, succeed := rawNode.(map[string]interface{})["is_video"]
	if !succeed || rawIsVideo == nil {
		return data.Post{}, fmt.Errorf("Unable to get 'is_video' element")
	}
	post.IsVideo = rawIsVideo.(bool)

	rawCaptionNode, succeed := rawNode.(map[string]interface{})["edge_media_to_caption"]
	if !succeed || rawCaptionNode == nil {
		return data.Post{}, fmt.Errorf("Unable to get 'edge_media_to_caption' element")
	}
	rawCaptionEdges, succeed := rawCaptionNode.(map[string]interface{})["edges"]
	if !succeed || rawCaptionEdges == nil {
		return data.Post{}, fmt.Errorf("Unable to get 'edges' element")
	}
	captionEdgesArray := rawCaptionEdges.([]interface{})
	if len(captionEdgesArray) > 0 {
		rawCaptionEdge, succeed := captionEdgesArray[0].(map[string]interface{})["node"]
		if !succeed || rawCaptionEdge == nil {
			return data.Post{}, fmt.Errorf("Unable to get 'node' element")
		}
		rawCaption, succeed := rawCaptionEdge.(map[string]interface{})["text"]
		if !succeed || rawCaption == nil {
			return data.Post{}, fmt.Errorf("Unable to get 'text' element")
		}
		post.Caption = rawCaption.(string)
	}

	rawTimestamp, succeed := rawNode.(map[string]interface{})["taken_at_timestamp"]
	if !succeed || rawTimestamp == nil {
		return data.Post{}, fmt.Errorf("Unable to get 'taken_at_timestamp' element")
	}
	timeString := rawTimestamp.(float64)
	timeInt := int64(timeString)
	post.Timestamp = timeInt

	rawComment, succeed := rawNode.(map[string]interface{})["edge_media_to_comment"]
	if !succeed || rawComment == nil {
		return data.Post{}, fmt.Errorf("Unable to get 'edge_media_to_comment' element")
	}
	rawCount, succeed := rawComment.(map[string]interface{})["count"]
	if !succeed || rawCount == nil {
		return data.Post{}, fmt.Errorf("Unable to get 'count' element")
	}
	countString := rawCount.(float64)
	countInt := int(countString)
	post.CommentsCount = countInt

	rawLikes, succeed := rawNode.(map[string]interface{})["edge_media_preview_like"]
	if !succeed || rawLikes == nil {
		return data.Post{}, fmt.Errorf("Unable to get 'edge_media_preview_like' element")
	}
	rawCount, succeed = rawLikes.(map[string]interface{})["count"]
	if !succeed || rawCount == nil {
		return data.Post{}, fmt.Errorf("Unable to get 'count' element")
	}
	countString = rawCount.(float64)
	countInt = int(countString)
	post.LikesCount = countInt

	rawOwner, succeed := rawNode.(map[string]interface{})["owner"]
	if !succeed || rawOwner == nil {
		return data.Post{}, fmt.Errorf("Unable to get 'owner' element")
	}
	rawOwnerID, succeed := rawOwner.(map[string]interface{})["id"]
	if !succeed || rawOwnerID == nil {
		return data.Post{}, fmt.Errorf("Unable to get 'id' element")
	}
	post.AuthorID = rawOwnerID.(string)

	rawLocation, succeed := rawNode.(map[string]interface{})["location"]
	if !succeed || rawLocation == nil {
		return data.Post{}, fmt.Errorf("Unable to get 'location' element")
	}
	rawLocationID, succeed := rawOwner.(map[string]interface{})["id"]
	if !succeed || rawLocationID == nil {
		return data.Post{}, fmt.Errorf("Unable to get 'id' element")
	}
	post.LocationID = rawLocationID.(string)

	rawTaggedNode, succeed := rawNode.(map[string]interface{})["edge_media_to_tagged_user"]
	if !succeed || rawTaggedNode == nil {
		return data.Post{}, fmt.Errorf("Unable to get 'edge_media_to_tagged_user' element")
	}
	rawTaggedEdges, succeed := rawTaggedNode.(map[string]interface{})["edges"]
	if !succeed || rawTaggedEdges == nil {
		return data.Post{}, fmt.Errorf("Unable to get 'edges' element")
	}
	taggedEdgesArray := rawTaggedEdges.([]interface{})
	if len(taggedEdgesArray) > 0 {
		post.TaggedUsers = []data.UserTag{}
		for i := 0; i < len(taggedEdgesArray); i++ {
			taggedUser := data.UserTag{}
			rawTaggedEdge, succeed := taggedEdgesArray[i].(map[string]interface{})["node"]
			if !succeed || rawTaggedEdge == nil {
				return data.Post{}, fmt.Errorf("Unable to get 'node' element")
			}
			rawUser, succeed := rawTaggedEdge.(map[string]interface{})["user"]
			if !succeed || rawUser == nil {
				return data.Post{}, fmt.Errorf("Unable to get 'user' element")
			}
			rawUsername, succeed := rawUser.(map[string]interface{})["username"]
			if !succeed || rawUsername == nil {
				return data.Post{}, fmt.Errorf("Unable to get 'username' element")
			}
			taggedUser.Username = rawUsername.(string)

			rawX, succeed := rawTaggedEdge.(map[string]interface{})["x"]
			if !succeed || rawX == nil {
				return data.Post{}, fmt.Errorf("Unable to get 'x' element")
			}
			taggedUser.X = rawX.(float64)
			rawY, succeed := rawTaggedEdge.(map[string]interface{})["y"]
			if !succeed || rawY == nil {
				return data.Post{}, fmt.Errorf("Unable to get 'y' element")
			}
			taggedUser.Y = rawY.(float64)
			post.TaggedUsers = append(post.TaggedUsers, taggedUser)
		}
	}
	return post, nil
}
