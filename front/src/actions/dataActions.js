import {
  CLEAR_STORE,
  INIT_SESSION,
  SET_AUTH,
  SET_CHART_DATA,
  SET_CURRENT_CITY_AND_VIEWPORT,
  SET_EVENTS_DATA,
  SET_EVENT_FILTER,
  SET_HEATMAP_DATA,
  SET_NEW_DATE,
  SET_SEARCHQUERY_FILTER,
  SET_SELECTED_HOUR,
  SET_SELECTED_DATE,
  SET_SEARCH_PARAMETERS,
  SET_SEARCH_QUERY,
  SET_POPUP_DATA,
  SETUP_MAP_FROM_HISTORY,
} from './types'

export const clearStore = () => {
  return { type: CLEAR_STORE }
}
export const setChartData = data => {
  return {
    type: SET_CHART_DATA,
    payload: data,
  }
}
export const setCurrentCityAndViewport = city => {
  return {
    type: SET_CURRENT_CITY_AND_VIEWPORT,
    payload: city,
  }
}
export const setEventsData = data => {
  return {
    type: SET_EVENTS_DATA,
    payload: data,
  }
}
export const setEventFilter = filter => ({
  type: SET_EVENT_FILTER,
  payload: filter,
})
export const setHeatmapData = data => {
  return {
    type: SET_HEATMAP_DATA,
    payload: data,
  }
}
export const setNewDate = date => {
  return {
    type: SET_NEW_DATE,
    payload: date,
  }
}
export const setSelectedDate = date => {
  return {
    type: SET_SELECTED_DATE,
    payload: date,
  }
}
export const setSearchQuery = events => {
  return {
    type: SET_SEARCH_QUERY,
    payload: events,
  }
}
export const setSelectedHour = hour => {
  return {
    type: SET_SELECTED_HOUR,
    payload: hour,
  }
}
export const setSearchParameters = params => {
  return {
    type: SET_SEARCH_PARAMETERS,
    payload: params,
  }
}
export const setSearchQueryFilter = filter => ({
  type: SET_SEARCHQUERY_FILTER,
  payload: filter,
})
export const setupMapFromHistory = config => ({
  type: SETUP_MAP_FROM_HISTORY,
  payload: config,
})
export const setPopupData = data => {
  return {
    type: SET_POPUP_DATA,
    payload: data,
  }
}
export const initSession = config => {
  return {
    type: INIT_SESSION,
    payload: config,
  }
}
export const setAuth = () => {
  return {
    type: SET_AUTH,
  }
}
