import createImmutableSelector from 'create-immutable-selector'
import { connect } from 'react-redux'

// selectors
import {
  chartDataSelector,
  selectedHourSelector,
  selectedDateSelector,
} from '../../reducers/dataReducer'
import {
  isDemoPlaySelector,
  isShowAllEventsSelector,
} from '../../reducers/uiReducer'

// actions
import {
  toggleAllEvents,
  playDemo,
  stopDemo,
  setHighlightedHour,
  fetchData,
} from '../../actions/uiActions'
import { setSelectedHour } from '../../actions/dataActions'

import Chart from './Chart.jsx'

const mapStateToProps = createImmutableSelector(
  chartDataSelector,
  isDemoPlaySelector,
  isShowAllEventsSelector,
  selectedDateSelector,
  selectedHourSelector,
  (chartData, isDemoPlay, isShowAllEvents, selectedDate, selectedHour) => ({
    chartData,
    isDemoPlay,
    isShowAllEvents,
    selectedDate,
    selectedHour,
  })
)
const mapDispatchToProps = dispatch => ({
  fetchData: hour => dispatch(fetchData(hour)),
  setHighlightedHour: hour => dispatch(setHighlightedHour(hour)),
  setSelectedHour: hour => dispatch(setSelectedHour(hour)),
  playDemo: () => dispatch(playDemo()),
  stopDemo: () => dispatch(stopDemo()),
  toggleAllEvents: (chartData, bounds) =>
    dispatch(toggleAllEvents(chartData, bounds)),
})

export default connect(mapStateToProps, mapDispatchToProps)(Chart)
