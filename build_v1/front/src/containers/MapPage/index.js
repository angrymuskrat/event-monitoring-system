import createImmutableSelector from 'create-immutable-selector'
import { connect } from 'react-redux'
import { replace } from 'connected-react-router'

// selectors
import {
  isLoadingSelector,
  isPopupOpenSelector,
  isShowAllEventsSelector,
  isSidebarOpenSelector,
  isSearchingEventsSelector,
} from '../../reducers/uiReducer'
import { viewportSelector } from '../../reducers/mapReducer'
import {
  allHeatmapDataSelector,
  chartDataSelector,
  currentHeatmapSelector,
  currentEventsSelector,
  currentPopupDataSelector,
  eventsSelector,
  selectedHourSelector,
  selectedDateSelector,
  searchQuerySelector,
} from '../../reducers/dataReducer'
import { currentCityIdSelector } from '../../reducers/cityReducer'

// actions
import {
  setViewport,
  setBounds,
  setCurrentUserLocation,
} from '../../actions/mapActions'
import { setupMapFromHistory } from '../../actions/dataActions'
import {
  fetchAllData,
  togglePopup,
  toggleSidebar,
  setSelectedEvent,
} from '../../actions/uiActions'

import MapPage from './MapPage.jsx'

const mapStateToProps = createImmutableSelector(
  allHeatmapDataSelector,
  chartDataSelector,
  currentCityIdSelector,
  currentHeatmapSelector,
  currentEventsSelector,
  eventsSelector,
  isLoadingSelector,
  isShowAllEventsSelector,
  isSidebarOpenSelector,
  selectedHourSelector,
  selectedDateSelector,
  viewportSelector,
  isPopupOpenSelector,
  isSearchingEventsSelector,
  searchQuerySelector,
  currentPopupDataSelector,
  (
    allHeatmapData,
    chartData,
    currentCityId,
    currentHeatmap,
    currentEvents,
    events,
    isLoading,
    isShowAllEvents,
    isSidebarOpen,
    selectedHour,
    selectedDate,
    viewport,
    isPopupOpen,
    isSearchingEvents,
    searchQuery,
    currentPopupData
  ) => ({
    allHeatmapData,
    chartData,
    currentCityId,
    currentHeatmap,
    currentEvents,
    events,
    isLoading,
    isShowAllEvents,
    isSidebarOpen,
    selectedHour,
    selectedDate,
    viewport,
    isPopupOpen,
    isSearchingEvents,
    searchQuery,
    currentPopupData,
  })
)
const mapDispatchToProps = dispatch => ({
  fetchAllData: time => dispatch(fetchAllData(time)),
  setupMapFromHistory: config => dispatch(setupMapFromHistory(config)),
  replace: path => dispatch(replace(path)),
  togglePopup: postcodes => dispatch(togglePopup(postcodes)),
  toggleSidebar: () => dispatch(toggleSidebar()),
  setBounds: bounds => dispatch(setBounds(bounds)),
  setCurrentUserLocation: location =>
    dispatch(setCurrentUserLocation(location)),
  setSelectedEvent: id => dispatch(setSelectedEvent(id)),
  setViewport: viewport => dispatch(setViewport(viewport)),
})

export default connect(mapStateToProps, mapDispatchToProps)(MapPage)
