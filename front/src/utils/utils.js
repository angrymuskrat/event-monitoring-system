import axios from 'axios'
import { fromJS, Set, List } from 'immutable'
import moment from 'moment'
import { uuid } from 'uuidv4'
import distance from '@turf/distance'
import { point } from '@turf/helpers'

//const makeInstagramImageUrl = code => `https://www.instagram.com/p/${code}/media/?size=m\``
export const makeInstagramImageUrl = (code) => `http://10.64.0.206:17112/image/${code}`

export const convertEventsToGeoJSON = async ({ data }) => {

  const events = await Promise.all(
    data.map(async d => {
      let lat = Number(d.Center.split(',')[0])
      let lon = Number(d.Center.split(',')[1])
      let postData
      try {
        return fromJS({
          properties: {
            tags: d.Tags,
            postcodes: d.PostCodes,
            title: d.Title,
            start: d.Start,
            finish: d.Finish,
            id: uuid(),
            photoUrl: makeInstagramImageUrl(d.PostCodes[0]),
          },
          geometry: {
            coordinates: [lat, lon],
          },
        })
      } catch (error) {
        return fromJS({
          properties: {
            tags: d.Tags,
            postcodes: d.PostCodes,
            title: d.Title,
            start: d.Start,
            finish: d.Finish,
            id: uuid(),
            photoUrl: makeInstagramImageUrl(d.PostCodes[1]),
          },
          geometry: {
            coordinates: [lat, lon],
          },
        })
      }
    })
  )
  return Set(events)
}
export const convertSearchQueryToGeoJSON = ({ data }) => {
  return List(
    data.map(d => {
      let lat = Number(d.Center.split(',')[0])
      let lon = Number(d.Center.split(',')[1])
      return fromJS({
        type: 'Feature',
        properties: {
          cluster: false,
          tags: d.Tags,
          postcodes: d.PostCodes,
          title: d.Title,
          start: d.Start,
          finish: d.Finish,
          id: uuid(),
          photoUrl: makeInstagramImageUrl(d.PostCodes[0]),
        },
        geometry: {
          type: 'Point',
          coordinates: [lat, lon],
        },
      })
    })
  )
}
export const convertHeatmapDataToLayer = ({ data }) => {
  return Set(
    data.map(d => {
      let coordinates = d.c.split(',')
      return fromJS([Number(coordinates[0]), Number(coordinates[1]), `${d.n}`])
    })
  )
}
export const convertChartData = ({ data }) => {
  data.sort((a, b) => {
    if (a.time > b.time) {
      return 1
    }
    if (a.time < b.time) {
      return -1
    }
    return 0
  })
  return data.map(d => {
    return {
      localTime: moment.unix(d.time).format('HH'),
      posts: d.posts,
      events: d.events,
      time: d.time,
    }
  })
}
export const convertPostData = (post) => {
  return {
    photoUrl: makeInstagramImageUrl(post.Shortcode),
    caption: post.Caption,
    likes: post.LikesCount,
    location: post.LocationID,
    locationLink: `https://www.instagram.com/explore/locations/${post.LocationID}`,
    profilePicUrl: `https://www.instagram.com/p/${post.Shortcode}/media/?size=l`,
    username: `user id: ${post.AuthorID}`,
    profileLink: `https://www.instagram.com/p/${post.Shortcode}/`,
    postLink: `https://www.instagram.com/p/${post.Shortcode}/`,
    id: post.Shortcode,
    comments: post.CommentsCount || 0,
  }
}
export const calculateDistance = (
  viewportCenterCoordinates,
  eventCoordinates
) => {
  const from = point(viewportCenterCoordinates)
  const to = point(eventCoordinates)
  const kilometers = distance(from, to, { units: 'kilometers' })
  return kilometers
}
export const filterListByKeyword = (str, arr) => {
  if (Set.isSet(arr)) {
    return arr.toJS().filter(function(el) {
      return str.split(/[-.,\s!]+/).some(function(word) {
        return (
          el.properties.title.toLowerCase().indexOf(word.toLowerCase()) !==
            -1 ||
          el.properties.tags.some(tag => {
            return tag.toLowerCase().indexOf(word.toLowerCase()) !== -1
          })
        )
      })
    })
  } else {
    return arr.filter(function(el) {
      return str.split(/[-.,\s!]+/).some(function(word) {
        return (
          el.properties.title.toLowerCase().indexOf(word.toLowerCase()) !==
            -1 ||
          el.properties.tags.some(tag => {
            return tag.toLowerCase().indexOf(word.toLowerCase()) !== -1
          })
        )
      })
    })
  }
}
export const sortEvents = (filters, events, viewport) => {
  if (filters.keyword) {
    events = filterListByKeyword(filters.keyword, events)
  }
  switch (filters.sortBy) {
    case 'A - Z':
      return events.sort((a, b) => {
        if (a.properties.title < b.properties.title) {
          return -1
        }
        if (a.properties.title > b.properties.title) {
          return 1
        }
        return 0
      })
    case 'Popular':
      return events.sort((a, b) => {
        if (a.properties.postcodes.length > b.properties.postcodes.length) {
          return -1
        }
        if (a.properties.postcodes.length < b.properties.postcodes.length) {
          return 1
        }
        return 0
      })
    case 'Nearby':
      const location = [viewport.center[0], viewport.center[1]]
      const eventsWithDistance = events.map(x => {
        x.properties.distance = calculateDistance(
          location,
          x.geometry.coordinates
        )
        return x
      })
      return eventsWithDistance.sort((a, b) => {
        if (a.properties.distance < b.properties.distance) {
          return -1
        }
        if (a.properties.distance > b.properties.distance) {
          return 1
        }
        return 0
      })
    case 'By time':
      return events.sort((a, b) => {
        if (a.properties.start < b.properties.start) {
          return -1
        }
        if (a.properties.start < b.properties.start) {
          return 1
        }
        return 0
      })
    default:
      return events
  }
}
export const getColor = (
  isSelected,
  isSelectedStartOrEnd,
  isWithinHoverRange,
  isDisabled
) => {
  return ({
    selectedFirstOrLastColor,
    normalColor,
    selectedColor,
    rangeHoverColor,
    disabledColor,
  }) => {
    if (isSelectedStartOrEnd) {
      return selectedFirstOrLastColor
    } else if (isSelected) {
      return selectedColor
    } else if (isWithinHoverRange) {
      return rangeHoverColor
    } else if (isDisabled) {
      return disabledColor
    } else {
      return normalColor
    }
  }
}
export const checkIfPhotoExists = async photoUrl => {
  try {
    let photo = await axios.get(`${photoUrl}`)
    if (photo.data) {
      return true
    }
    return false
  } catch (err) {
    console.error(err)
  }
}
