import { SET_CITY } from './types'

export const setCurrentCity = city => {
  return {
    type: SET_CITY,
    payload: city,
  }
}
