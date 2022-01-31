import {
  FETCH_ALL_DATA,
  FETCH_DATA,
  FETCH_DATA_AFTER_VIEWPORT_CHANGE,
  FETCH_FAILURE,
  FETCH_SUCCESS,
  PLAY_DEMO,
  SET_HIGHLIGHTED_HOUR,
  SET_HIGHTLIGHTED_EVENT,
  SET_SELECTED_EVENT,
  STOP_DEMO,
  TOGGLE_ALL_EVENTS,
  TOGGLE_POPUP,
  TOGGLE_SIDEBAR,
  VIEWPORT_TRANSITION,
  SEARCH_EVENTS,
  LOADING_START,
  LOADING_FINISH,
  SET_AUTH_ERROR,
} from './types'

export const toggleSidebar = () => ({ type: TOGGLE_SIDEBAR })
export const toggleAllEvents = () => ({ type: TOGGLE_ALL_EVENTS })
export const togglePopup = postcodes => ({
  type: TOGGLE_POPUP,
  payload: postcodes,
})
export const playDemo = () => ({ type: PLAY_DEMO })
export const startLoading = () => ({ type: LOADING_START })
export const finishLoading = () => ({ type: LOADING_FINISH })
export const stopDemo = () => ({ type: STOP_DEMO })
export const fetchData = data => ({ type: FETCH_DATA, payload: data })
export const fetchAllData = data => ({ type: FETCH_ALL_DATA, payload: data })
export const fetchSuccess = () => ({ type: FETCH_SUCCESS })
export const fetchFailure = error => ({ type: FETCH_FAILURE, payload: error })
export const setSearchEvents = () => ({ type: SEARCH_EVENTS })
export const setAuthError = error => {
  return { type: SET_AUTH_ERROR, payload: error }
}
export const setHighlightedHour = hour => {
  return {
    type: SET_HIGHLIGHTED_HOUR,
    payload: hour,
  }
}
export const setHighlightedEvent = id => {
  return {
    type: SET_HIGHTLIGHTED_EVENT,
    payload: id,
  }
}
export const setSelectedEvent = id => {
  return {
    type: SET_SELECTED_EVENT,
    payload: id,
  }
}
export const fetchDataAfterViewportChange = data => {
  return {
    type: FETCH_DATA_AFTER_VIEWPORT_CHANGE,
    payload: data,
  }
}
export const viewportTransition = () => ({ type: VIEWPORT_TRANSITION })
