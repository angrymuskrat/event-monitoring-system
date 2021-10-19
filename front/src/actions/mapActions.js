import {
  CHANGE_VIEWPORT,
  SET_BOUNDS,
  SET_CURRENT_USER_LOCATION,
  SET_VIEWPORT,
} from './types'

export const changeViewport = config => ({
  type: CHANGE_VIEWPORT,
  payload: config,
})
export const setBounds = bounds => ({
  type: SET_BOUNDS,
  payload: bounds,
})
export const setCurrentUserLocation = location => ({
  type: SET_CURRENT_USER_LOCATION,
  payload: location,
})
export const setViewport = newViewport => ({
  type: SET_VIEWPORT,
  payload: newViewport,
})
