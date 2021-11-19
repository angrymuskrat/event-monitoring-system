# Backend API

## Content

### Objects

* [Event](#event-object) - contains main info about an event
* [Timeline unit](#timeline-object) - is used to display a timeline of the number of publications and events
* [Heatmap unit](#heatmap-object) - is used to display a heatmap of the number of publications
* [Short Instagram post](#shortpost-object) - contains main info about Instagram post

### Requests
* [login](#/login) - get session cookie
* [heatmap](#heatmap) - get heatmap of posts counts for a rectangle
* [timeline](#timeline) - get city timelines
* [events](#events) - search events
* [search](#search) - search events by hashtags

## Objects

### Event object
A city event is several posts close in location and time and linked by a common theme or hashtag.
```
{
  "Center": "{float64},{float64}", // concatinating of latitude and longitude of a event
  "PostCodes": array of string, // shortcodes of Instagram posts related with the event  
  "Tags": array of string, // hashtags related with the events (with symbol '#')
  "Title": string, // title of the event (most popular hashtag from related hashtags)
  "Start": int // unix event begin timestamp
  "Finish": int // unix event end timestamp 
}
```

### Timeline object
Timeline object characterizes the number of publications and events in a given hour
```
{
    "time": int, // unix timestamp of begining the hour
    "posts": int, // count of posts in the hour 
    "events": int // count of events in the hour
}
```

### Heatmap object
Heatmap object is a map area of size 50 by 50 meters for which the number of publications is calculated, the cell is
set by the coordinates of its center.
```
{
  "c": "{float64},{float64}", // concatinating of latitude and longitude - center the cell
  "n": int // count of posts in this cell
}
```

### ShortPost object
ShortPost contains basic information about the Instagram post.
```
{
  "Shortcode": string, // unic id of the post
  "Caption": string, // text of the post
  "LikesCount": int, // amount of likes
  "Timestamp": int, // unix timestamp 
  "Lat": float64, // lattitude
  "Lon": float64 // longitude
}
```
## Requests

### login
Request: /login <br>
Type: POST <br>
Description: For an existing user it gives the "session" cookie needed for all other requests. <br>
Input: body - raw object:
```
{
  "login": string,  // login 
  "password": string // password
}
```
Output: set cookie with name "session"
### heatmap
Request: /heatmap/city/topLeftLat,topLeftLon/botRightLat,botRightLon/hourTimestamp <br>
Type: GET <br>
Description: For a city and a rectangle given in two corners gives out a grid with the size of cells x by y meters (the default
value is 50x50 meters, the value is set in the service storage), with the number of Instagram posts for a given hour
in each cell.<br>
Input:
* city: string - code of the city
* topLeftLat, topLeftLon: float64 - the latitude and longitude of top left corner of the rectangle
* botRightLat, botRightLon: float64 - the latitude and longitude of bottom right corner of the rectangle
* hourTimestamp - Unix timestamp of the beginning of the hour <br>

Cookie: session <br>
Output: JSON array of Heatmap objects (response contains only cells with one or more posts).<br>

Example: <br>
&nbsp;&nbsp;&nbsp;request: /heatmap/spb/59.945257,30.305417/59.938587,30.319471/1559404800 <br>
&nbsp;&nbsp;&nbsp;response: JSON array
```json
[
  {
    "c": "59.9428,30.3068",
    "n": 12
  },
  {
    "c": "59.9414,30.3126",
    "n": 2
  },
  {
    "c": "59.9401,30.3146",
    "n": 24
  },
  {
    "c": "59.9410,30.3150",
    "n": 7
  },
  {
    "c": "59.9419,30.3157",
    "n": 1
  }
]
```
### timeline
Request: /timeline/spb/startTimestamp/endTimestamp <br>
Type: GET <br>
Description: For city and time interval gives the number of posts and events for each hour of the given time interval.<br>
Input:
* city: string - code of the city
* startTimestamp - Unix timestamp of the beginning of the *first* hour from the time interval.
* endTimestamp - Unix timestamp of the ending of the *last* hour from the time interval.<br>

Cookie: session <br>
Output: JSON array of Timeline objects. <br>

Example: <br>
&nbsp;&nbsp;&nbsp;request: /timeline/spb/1559408400/1559415600<br>
&nbsp;&nbsp;&nbsp;response: JSON array
```json
[
  {
    "time": 1559408400,
    "posts": 1845,
    "events": 5
  },
  {
    "time": 1559412000,
    "posts": 1967,
    "events": 3
  },
  {
    "time": 1559415600,
    "posts": 1927,
    "events": 7
  }
]

```

### events
Request: /events/city/topLeftLat,topLeftLon/botRightLat,botRightLon/hourTimestamp
Type: GET <br>
Description: For a city and a rectangle given in two corners gives out events for a given hour.<br>
Input:
* city: string - code of the city
* topLeftLat, topLeftLon: float64 - the latitude and longitude of top left corner of the rectangle
* botRightLat, botRightLon: float64 - the latitude and longitude of bottom right corner of the rectangle
* hourTimestamp - Unix timestamp of the beginning of the hour <br>

Cookie: session <br>
Output: JSON array of Event objects<br>

Example: <br>
&nbsp;&nbsp;&nbsp; request: /events/spb/60.12,30.11/59.84,30.69/1557442800 <br>
&nbsp;&nbsp;&nbsp; response: JSON array
```json
[
  {
    "Center": "59.9382,30.3203",
    "PostCodes": ["BxQnlpAgyG9", "BxQohGujwPW", "BxQqAldAyNF", "BxQogYthDEq", "BxQomc5j7K7", "BxQsls3DbHT", "BxQoarUBbnY"],
    "Tags": ["#деньпобеды", "#9мая", "#салют", "#9мая2019", "#мирногонеба"],
    "Title": "#деньпобеды",
    "Start": 1557442800,
    "Finish": 1557446400
  }
]
```
### search
Request: /search/spb/tags/startTimestamp/endTimestamp <br>
Type: GET <br>
Description:  For a city and set of hashtags gives out events for a given time interval.<br>
Input:
* city: string - code of the city
* tags: string - comma-separated hashtags concatenation in UTF-8 format(with '#' symbol, without any space symbols)
* startTimestamp - Unix timestamp of the beginning of the *first* hour from the time interval.
* endTimestamp - Unix timestamp of the ending of the *last* hour from the time interval.<br>
  Output: JSON array of Event objects (see [events](#events) request) <br>
  Example: <br>
  &nbsp;&nbsp;&nbsp; request: search/spb/%239%D0%BC%D0%B0%D1%8F%2C%23%D1%81%D0%B0%D0%BB%D1%8E%D1%82/1557432000/1557442800 <br>
  &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; Note: tags - string "#9мая,#салют" in UTF-8 format<br>
  &nbsp;&nbsp;&nbsp; response: JSON array
```json
[
  {
    "Center": "59.9404,30.3113",
    "PostCodes": [
      "BxQPLSVFJOE",
      "BxQOaMZBsJc",
      ...
    ],
    "Tags": [
      "#9мая",
      "#деньпобеды",
      "#салют",
      "#9мая2019",
      "#сднемпобеды"
    ],
    "Title": "#9мая",
    "Start": 1557428400,
    "Finish": 1557432000
  },
  {
    "Center": "59.9408,30.3319",
    "PostCodes": [
      "BxQRtcrDgx6",
      "BxQRN7Hnn0A",
      ...
    ],
    "Tags": [
      "#9мая",
      "#кириллгордеев",
      "#деньпобеды",
      "#игорькроль",
      "#салют"
    ],
    "Title": "#9мая",
    "Start": 1557428400,
    "Finish": 1557432000
  }, ...
]
```
### shortPosts
Request: /shortPosts/spb/startTimestamp/endTimestamp/shortcodes <br>
Type: GET <br>
Description: Get Instagram posts by their shortcodes and time interval. The time interval is needed to optimize the
query of a group of posts, if the time is unknown use the query [/singlePost](#singleShortPost) <br>
Input:
* city: string - code of the city
* shortcodes: string - comma-separated shortcodes of needed posts
* startTimestamp - Unix timestamp of the beginning of the time interval.
* endTimestamp - Unix timestamp of the ending of the time interval.<br>

Output: JSON array of Shortpost objects <br>
Example: <br>
&nbsp;&nbsp;&nbsp; request: /shortPosts/spb/1572615416/1572619016/B4U0-0Ygoc3,B4U0-NaAMex <br>
&nbsp;&nbsp;&nbsp; response:
```json
[
    {
        "Shortcode": "B4U0-0Ygoc3",
        "Caption": "Кресла для тех, кому не сидится на месте.\n⠀\nДетское сидение из коллекции оригинальных аксессуаров Audi создано для детей от первых дней жизни до достижения роста 105 см. Сидение прочно соединено с базой и не вынимается. Кроме того, его можно поворачивать на 360 градусов не отсоединяя от сидения.\n⠀\nНаличие уточняйте по телефону:\n+7 (812) 561-08-67\nили в шоуруме по адресу:\nВитебский просп., 17/2.\n⠀\n#Audi #AudiВитебский",
        "LikesCount": 74,
        "Timestamp": 1572616940,
        "Lat": 59.8729024727,
        "Lon": 30.3523912688
    },
    {
        "Shortcode": "B4U0-NaAMex",
        "CommentsCount": 2,
        "LikesCount": 41,
        "Timestamp": 1572616920,
        "Lat": 59.9341,
        "Lon": 30.3062
    }
]
```

### singleShortPost
Request: /singleShortPost/spb/shortcode <br>
Type: GET <br>
Description: Get basic info about Instagram post by its shortcode. <br>
Input:
* city: string - code of the city
* shortcode: string - shortcode of the needed post

Output: JSON object - ShortPost <br>
Example: <br>
&nbsp;&nbsp;&nbsp; request: /singleShortPost/spb/B4U01n5I0H5 <br>
&nbsp;&nbsp;&nbsp; response:
```json
{
  "Shortcode": "B4U01n5I0H5",
  "Caption": "#зениттомь💪⚽ 4:0",
  "LikesCount": 30,
  "Timestamp": 1572616850,
  "Lat": 59.97301135029,
  "Lon": 30.220320224762
}
```

### image
Request: /image/shortcode <br>
Type: GET <br>
Description: Makes a request for the first image of a post, bypassing Instagram's strict CORS policy <br>
Input: shortcode: string - shortcode of the needed post <br>
Output: image/jpeg <br>
Example: <br>
&nbsp;&nbsp;&nbsp; request: /image/B4U01n5I0H5 <br>
&nbsp;&nbsp;&nbsp; response:
![Image for B4U01n5I0H5](https://instagram.frix7-1.fna.fbcdn.net/v/t51.2885-15/e35/s1080x1080/73393262_149724593086603_8485685819203930514_n.jpg?_nc_ht=instagram.frix7-1.fna.fbcdn.net&_nc_cat=107&_nc_ohc=ODsbYDeSGUoAX9GhvfF&edm=AGenrX8BAAAA&ccb=7-4&oh=ee431fc84637bce882951cc56bf861a8&oe=6195B7CF&_nc_sid=5eceaa)