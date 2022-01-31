import { fromJS, Set, List, Map } from 'immutable'
import createImmutableSelector from 'create-immutable-selector'

/******************************************************************************/
/******************************* TYPES ****************************************/
/******************************************************************************/

import {
  CLEAR_STORE,
  SET_AUTH,
  SET_SELECTED_HOUR,
  SET_EVENTS_DATA,
  SET_HEATMAP_DATA,
  SET_CHART_DATA,
  SET_NEW_DATE,
  SET_CURRENT_CITY_AND_VIEWPORT,
  SETUP_MAP_FROM_HISTORY,
  SET_EVENT_FILTER,
  SET_SEARCHQUERY_FILTER,
  SET_SELECTED_DATE,
  SET_SEARCH_PARAMETERS,
  SET_SEARCH_QUERY,
  SET_POPUP_DATA,
} from '../actions/types'

/******************************************************************************/
/******************************* INITIAL STATE ********************************/
/******************************************************************************/

const initialState = fromJS({
  selectedHour: null,
  selectedDate: null,
  heatmap: {},
  events: {},
  searchQuery: null,
  searchParameters: null,
  chartData: null,
  eventFilters: {
    sortBy: null,
    keyword: null,
  },
  searchQueryFilters: {
    sortBy: null,
    keyword: null,
  },
  currentPopupData: null,
  auth: false,
})

/******************************************************************************/
/******************************* SELECTORS ************************************/
/******************************************************************************/

const dataSelector = createImmutableSelector(
  state => state.get('data'),
  data => data
)
export const selectedHourSelector = createImmutableSelector(
  dataSelector,
  data => data.get('selectedHour')
)
export const selectedDateSelector = createImmutableSelector(
  dataSelector,
  data => data.get('selectedDate')
)
export const heatmapSelector = createImmutableSelector(dataSelector, data =>
  data.get('heatmap')
)
export const currentHeatmapSelector = createImmutableSelector(
  selectedHourSelector,
  dataSelector,
  (selectedHour, data) => {
    return data.getIn(['heatmap', selectedHour])
  }
)
export const eventsSelector = createImmutableSelector(dataSelector, data =>
  data.get('events')
)
export const currentEventsSelector = createImmutableSelector(
  selectedHourSelector,
  dataSelector,
  (selectedHour, data) => {
    return data.getIn(['events', selectedHour])
  }
)
export const chartDataSelector = createImmutableSelector(dataSelector, data =>
  data.get('chartData')
)
export const allChartEventsSelector = createImmutableSelector(
  dataSelector,
  data =>
    data
      .get('chartData')
      .reduce((accum, data) => (accum += data.get('events')), 0)
)

export const allHeatmapDataSelector = createImmutableSelector(
  dataSelector,
  data => {
    return data.get('heatmap')
  }
)
export const allEventsDataSelector = createImmutableSelector(
  dataSelector,
  data => {
    if (JSON.stringify(data.get('events')) !== '{}') {
      return data
        .get('events')
        .valueSeq()
        .flatten(1)
    }
  }
)
export const filteredAllEventsDataSelector = createImmutableSelector(
  dataSelector,
  data => {
    return data.get('events')
  }
)
export const eventFiltersSelector = createImmutableSelector(
  dataSelector,
  data => data.get('eventFilters')
)
export const searchQueryFiltersSelector = createImmutableSelector(
  dataSelector,
  data => data.get('searchQueryFilters')
)
export const searchQuerySelector = createImmutableSelector(dataSelector, data =>
  data.get('searchQuery')
)
export const searchParametersSelector = createImmutableSelector(
  dataSelector,
  data => data.get('searchParameters')
)
export const searchParametersStartDateSelector = createImmutableSelector(
  dataSelector,
  data => data.getIn(['searchParameters', 'startDate'])
)
export const searchParametersEndDateSelector = createImmutableSelector(
  dataSelector,
  data => data.getIn(['searchParameters', 'endDate'])
)
export const currentPopupDataSelector = createImmutableSelector(
  dataSelector,
  data => data.get('currentPopupData')
)
export const authSelector = createImmutableSelector(dataSelector, data =>
  data.get('auth')
)

/******************************************************************************/
/******************************* REDUCERS *************************************/
/******************************************************************************/

export default (state = initialState, { type, payload }) => {
  switch (type) {
    case CLEAR_STORE:
      return state
        .set('events', Map())
        .set('heatmap', Map())
        .set('chartData', null)
        .set('searchQuery', null)
        .set('searchParameters', null)
        .set(
          'eventFilters',
          fromJS({
            sortBy: null,
            keyword: null,
          })
        )
        .set(
          'searchQueryFilters',
          fromJS({
            sortBy: null,
            keyword: null,
          })
        )
    case SET_SELECTED_HOUR:
      return state.set('selectedHour', payload.toString())
    case SET_SELECTED_DATE:
      return state.set('selectedDate', payload)
    case SET_HEATMAP_DATA:
      if (Set.isSet(payload)) {
        const selectedHour = state.get('selectedHour')
        if (state.getIn(['heatmap', selectedHour]) !== undefined) {
          let prevHeatmap = state.getIn(['heatmap', selectedHour])
          return state.setIn(
            ['heatmap', selectedHour],
            Set(prevHeatmap.union(payload))
          )
        } else {
          return state.setIn(['heatmap', selectedHour], payload)
        }
      } else {
        const selectedHour = Object.keys(payload)[0]
        const data = payload[selectedHour]
        if (state.getIn(['heatmap', selectedHour]) !== undefined) {
          let prevHeatmap = state.getIn(['heatmap', selectedHour])
          return state.setIn(
            ['heatmap', selectedHour],
            Set(prevHeatmap.union(data))
          )
        } else {
          return state.setIn(['heatmap', selectedHour], data)
        }
      }
    case SET_EVENTS_DATA:
      if (Set.isSet(payload)) {
        const selectedHour = state.get('selectedHour')
        if (state.getIn(['events', selectedHour]) !== undefined) {
          let prevEvents = state.getIn(['events', selectedHour])
          return state.setIn(
            ['events', selectedHour],
            Set(prevEvents.union(payload))
          )
        } else {
          return state.setIn(['events', selectedHour], payload)
        }
      } else {
        const selectedHour = Object.keys(payload)[0]
        const data = payload[selectedHour]
        if (state.getIn(['events', selectedHour]) !== undefined) {
          let prevEvents = state.getIn(['events', selectedHour])
          return state.setIn(
            ['events', selectedHour],
            Set(prevEvents.union(data))
          )
        } else {
          return state.setIn(['events', selectedHour], data)
        }
      }
    case SET_SEARCH_QUERY:
      if (payload === null) {
        return state.set('searchQuery', null)
      } else {
        return state.set('searchQuery', List(payload))
      }
    case SET_CHART_DATA:
      return state.set('chartData', payload)
    case SET_POPUP_DATA:
      return state.set('currentPopupData', payload)
    case SET_CURRENT_CITY_AND_VIEWPORT:
    case SETUP_MAP_FROM_HISTORY:
    case SET_NEW_DATE:
      return state
    case SET_EVENT_FILTER:
      if (payload.type === 'sort') {
        return state.setIn(['eventFilters', 'sortBy'], payload.parameter)
      } else if (payload.type === 'keyword') {
        return state.setIn(['eventFilters', 'keyword'], payload.parameter)
      }
      break
    case SET_SEARCHQUERY_FILTER:
      if (payload.type === 'sort') {
        return state.setIn(['searchQueryFilters', 'sortBy'], payload.parameter)
      } else if (payload.type === 'keyword') {
        return state.setIn(['searchQueryFilters', 'keyword'], payload.parameter)
      }
      break
    case SET_SEARCH_PARAMETERS:
      return state.set('searchParameters', payload)
    case SET_AUTH: {
      return state.set('auth', !state.get('auth'))
    }
    default:
      return state
  }
}
