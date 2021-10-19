import createImmutableSelector from 'create-immutable-selector'
import { connect } from 'react-redux'

import {
  isPopupOpenSelector,
  isSidebarOpenSelector,
  isShowAllEventsSelector,
  isSearchingEventsSelector,
  isLoadingSelector,
  errorSelector,
} from '../../reducers/uiReducer'
import {
  allEventsDataSelector,
  chartDataSelector,
  currentEventsSelector,
  eventsSelector,
  selectedHourSelector,
  selectedDateSelector,
  eventFiltersSelector,
  searchQueryFiltersSelector,
  searchQuerySelector,
  searchParametersSelector,
  searchParametersStartDateSelector,
  searchParametersEndDateSelector,
} from '../../reducers/dataReducer'
import {
  currentUserLocationSelector,
  viewportSelector,
} from '../../reducers/mapReducer'

// action
import {
  toggleSidebar,
  setHighlightedEvent,
  setSelectedEvent,
  setSearchEvents,
} from '../../actions/uiActions'
import {
  setEventFilter,
  setNewDate,
  setSearchParameters,
  setSearchQueryFilter,
} from '../../actions/dataActions'

import Sidebar from './Sidebar.jsx'

const mapStateToProps = createImmutableSelector(
  allEventsDataSelector,
  chartDataSelector,
  currentEventsSelector,
  currentUserLocationSelector,
  eventsSelector,
  isPopupOpenSelector,
  isSidebarOpenSelector,
  isShowAllEventsSelector,
  selectedHourSelector,
  selectedDateSelector,
  eventFiltersSelector,
  isSearchingEventsSelector,
  searchQuerySelector,
  searchParametersSelector,
  isLoadingSelector,
  viewportSelector,
  errorSelector,
  searchQueryFiltersSelector,
  searchParametersStartDateSelector,
  searchParametersEndDateSelector,
  (
    allEventsData,
    chartData,
    currentEvents,
    currentUserLocation,
    events,
    isPopupOpen,
    isSidebarOpen,
    isShowAllEvents,
    selectedHour,
    selectedDate,
    eventFilters,
    isSearchingEvents,
    searchQuery,
    searchParameters,
    isLoading,
    viewport,
    errors,
    searchQueryFilters,
    searchParametersStartDate,
    searchParametersEndDate
  ) => ({
    allEventsData,
    chartData,
    currentEvents,
    currentUserLocation,
    events,
    isPopupOpen,
    isSidebarOpen,
    isShowAllEvents,
    selectedHour,
    selectedDate,
    eventFilters,
    isSearchingEvents,
    searchQuery,
    searchParameters,
    isLoading,
    viewport,
    errors,
    searchQueryFilters,
    searchParametersStartDate,
    searchParametersEndDate,
  })
)
const mapDispatchToProps = dispatch => ({
  toggleSidebar: () => dispatch(toggleSidebar()),
  setNewDate: date => dispatch(setNewDate(date)),
  setHighlightedEvent: id => dispatch(setHighlightedEvent(id)),
  setSelectedEvent: id => dispatch(setSelectedEvent(id)),
  setEventFilter: filter => dispatch(setEventFilter(filter)),
  setSearchQueryFilter: filter => dispatch(setSearchQueryFilter(filter)),
  setSearchParameters: params => dispatch(setSearchParameters(params)),
  setSearchEvents: () => dispatch(setSearchEvents()),
})

export default connect(mapStateToProps, mapDispatchToProps)(Sidebar)
