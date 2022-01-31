import { fromJS } from 'immutable'
import createImmutableSelector from 'create-immutable-selector'

/******************************************************************************/
/******************************* TYPES ****************************************/
/******************************************************************************/

import {
  CLEAR_STORE,
  TOGGLE_SIDEBAR,
  TOGGLE_ALL_EVENTS,
  PLAY_DEMO,
  STOP_DEMO,
  SET_HIGHLIGHTED_HOUR,
  FETCH_DATA,
  FETCH_ALL_DATA,
  FETCH_TIMELINE_SUCCESS,
  FETCH_SUCCESS,
  FETCH_FAILURE,
  SET_AUTH_ERROR,
  SET_HIGHTLIGHTED_EVENT,
  SET_SELECTED_EVENT,
  TOGGLE_POPUP,
  FETCH_TIMELINE_START,
  VIEWPORT_TRANSITION,
  SEARCH_EVENTS,
  LOADING_START,
  LOADING_FINISH,
} from '../actions/types'

/******************************************************************************/
/******************************* INITIAL STATE ********************************/
/******************************************************************************/

export const initialState = fromJS({
  isLoading: false,
  isTimelineLoading: false,
  isSidebarOpen: false,
  isDemoPlay: false,
  isShowAllEvents: false,
  isPopupOpen: false,
  isViewportTransition: false,
  highlightedHour: null,
  highlightedEvent: null,
  selectedEvent: null,
  errors: null,
  authError: null,
  isSearchingEvents: false,
  postcodes: null,
})

/******************************************************************************/
/******************************* SELECTORS ************************************/
/******************************************************************************/

const uiSelector = createImmutableSelector(
  state => state.get('ui'),
  ui => ui
)
export const isSidebarOpenSelector = createImmutableSelector(uiSelector, ui =>
  ui.get('isSidebarOpen')
)
export const isDemoPlaySelector = createImmutableSelector(uiSelector, ui =>
  ui.get('isDemoPlay')
)
export const isShowAllEventsSelector = createImmutableSelector(uiSelector, ui =>
  ui.get('isShowAllEvents')
)
export const isLoadingSelector = createImmutableSelector(uiSelector, ui =>
  ui.get('isLoading')
)
export const errorSelector = createImmutableSelector(uiSelector, ui =>
  ui.get('errors')
)
export const isTimelineLoadingSelector = createImmutableSelector(
  uiSelector,
  ui => ui.get('isTimelineLoading')
)
export const highlightedHourSelector = createImmutableSelector(uiSelector, ui =>
  ui.get('highlightedHour')
)
export const highlightedEventSelector = createImmutableSelector(
  uiSelector,
  ui => ui.get('highlightedEvent')
)
export const selectedEventSelector = createImmutableSelector(uiSelector, ui =>
  ui.get('selectedEvent')
)
export const isPopupOpenSelector = createImmutableSelector(uiSelector, ui =>
  ui.get('isPopupOpen')
)
export const isViewportTransitionSelector = createImmutableSelector(
  uiSelector,
  ui => ui.get('isViewportTransition')
)
export const postcodesSelector = createImmutableSelector(uiSelector, ui =>
  ui.get('postcodes')
)
export const isSearchingEventsSelector = createImmutableSelector(
  uiSelector,
  ui => ui.get('isSearchingEvents')
)
export const authErrorSelector = createImmutableSelector(uiSelector, ui =>
  ui.get('authError')
)

/******************************************************************************/
/******************************* REDUCERS *************************************/
/******************************************************************************/

export default function(state = initialState, { type, payload }) {
  switch (type) {
    case FETCH_DATA:
    case FETCH_ALL_DATA:
      return state.set('isLoading', true)
    case FETCH_SUCCESS:
      return state.set('isLoading', false)
    case FETCH_TIMELINE_START:
      return state.set('isTimelineLoading', true)
    case FETCH_TIMELINE_SUCCESS:
      return state.set('isTimelineLoading', false)
    case FETCH_FAILURE:
      return state.set('isLoading', false).set('errors', payload)
    case TOGGLE_SIDEBAR:
      return state.set('isSidebarOpen', !state.get('isSidebarOpen'))
    case TOGGLE_POPUP:
      return state
        .set('isPopupOpen', !state.get('isPopupOpen'))
        .set('postcodes', payload)
    case PLAY_DEMO:
    case STOP_DEMO:
      return state
        .set('isDemoPlay', !state.get('isDemoPlay'))
        .set('isShowAllEvents', false)
        .set('isLoading', false)
    case TOGGLE_ALL_EVENTS:
      return state
        .set('isShowAllEvents', !state.get('isShowAllEvents'))
        .set('isDemoPlay', false)
    case SET_AUTH_ERROR:
      return state.set('authError', payload)
    case SET_HIGHLIGHTED_HOUR:
      return state.set('highlightedHour', payload)
    case SET_HIGHTLIGHTED_EVENT:
      return state.set('highlightedEvent', payload)
    case SET_SELECTED_EVENT:
      return state.set('selectedEvent', payload)
    case SEARCH_EVENTS:
      return state.set('isSearchingEvents', !state.get('isSearchingEvents'))
    case VIEWPORT_TRANSITION:
      return state.set(
        'isViewportTransition',
        !state.get('isViewportTransition')
      )
    case LOADING_START:
      return state.set('isLoading', true)
    case LOADING_FINISH:
      return state.set('isLoading', false)
    case CLEAR_STORE:
      return state
        .set('isShowAllEvents', false)
        .set('isDemoPlay', false)
        .set('isSearchingEvents', false)
        .set('errors', null)
        .set('isPopupOpen', false)
        .set('isSidebarOpen', false)
        .set('selectedEvent', null)
    default:
      return state
  }
}
