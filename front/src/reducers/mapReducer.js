import { fromJS } from 'immutable'
import createImmutableSelector from 'create-immutable-selector'

/******************************************************************************/
/******************************* TYPES ****************************************/
/******************************************************************************/

import {
  SET_BOUNDS,
  SET_CURRENT_USER_LOCATION,
  SET_VIEWPORT,
  SETUP_MAP_FROM_HISTORY,
  CHANGE_VIEWPORT,
} from '../actions/types'

/******************************************************************************/
/******************************* INITIAL STATE ********************************/
/******************************************************************************/

const initialState = fromJS({
  bounds: {
    topLeft: [0, 0],
    bottomRight: [0, 0],
  },
  viewport: {
    center: [],
    zoom: 13,
  },
  currentUserLocation: null,
})

/******************************************************************************/
/******************************* SELECTORS ************************************/
/******************************************************************************/

const mapSelector = createImmutableSelector(
  state => state.get('map'),
  map => map
)
export const viewportSelector = createImmutableSelector(mapSelector, map =>
  map.get('viewport')
)
export const boundsSelector = createImmutableSelector(mapSelector, map =>
  map.get('bounds')
)
export const currentUserLocationSelector = createImmutableSelector(
  mapSelector,
  map => map.get('currentUserLocation')
)

/******************************************************************************/
/******************************* REDUCERS *************************************/
/******************************************************************************/

export default (state = initialState, { type, payload }) => {
  switch (type) {
    case SET_VIEWPORT:
      return state.update('viewport', previousViewport =>
        previousViewport.merge(payload)
      )
    case SET_BOUNDS:
      return state.update('bounds', previousBounds =>
        previousBounds.merge(payload)
      )
    case SET_CURRENT_USER_LOCATION:
      return state.set('currentUserLocation', payload)
    case SETUP_MAP_FROM_HISTORY:
    case CHANGE_VIEWPORT:
      return state
    default:
      return state
  }
}
