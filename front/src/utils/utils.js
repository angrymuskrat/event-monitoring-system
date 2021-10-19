import axios from 'axios'
import { fromJS, Set, List } from 'immutable'
import moment from 'moment'
import { uuid } from 'uuidv4'
import distance from '@turf/distance'
import { point } from '@turf/helpers'

export const convertEventsToGeoJSON = async ({ data }) => {
  const events = await Promise.all(
    data.map(async d => {
      let lat = Number(d.Center.split(',')[0])
      let lon = Number(d.Center.split(',')[1])
      let postData
      try {
        postData = await axios.get(
          `https://www.instagram.com/p/${d.PostCodes[0]}/?__a=1`
        )
        return fromJS({
          properties: {
            tags: d.Tags,
            postcodes: d.PostCodes,
            title: d.Title,
            start: d.Start,
            finish: d.Finish,
            id: uuid(),
            photoUrl: `${postData.data.graphql.shortcode_media.display_url}`,
          },
          geometry: {
            coordinates: [lat, lon],
          },
        })
      } catch (error) {
        postData = await axios.get(
          `https://www.instagram.com/p/${d.PostCodes[1]}/?__a=1`
        )
        return fromJS({
          properties: {
            tags: d.Tags,
            postcodes: d.PostCodes,
            title: d.Title,
            start: d.Start,
            finish: d.Finish,
            id: uuid(),
            photoUrl: `${postData.data.graphql.shortcode_media.display_url}`,
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
          photoUrl: `https://instagram.com/p/${d.PostCodes[0]}/media/?size=m`,
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
export const convertPostData = post => {
  return {
    photoUrl: post.display_url,
    caption: post.edge_media_to_caption.edges[0].node.text,
    likes: post.edge_media_preview_like.count,
    location: post.location.name,
    profilePicUrl: post.owner.profile_pic_url,
    username: post.owner.username,
    profileLink: `https://www.instagram.com/${post.owner.username}`,
    postLink: `https://www.instagram.com/p/${post.shortcode}/`,
    id: post.id,
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
