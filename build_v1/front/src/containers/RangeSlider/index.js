import createImmutableSelector from 'create-immutable-selector'
import { connect } from 'react-redux'

// selectors
import {
  chartDataSelector,
  selectedHourSelector,
} from '../../reducers/dataReducer'
import { isLoadingSelector } from '../../reducers/uiReducer'

// actions
import { fetchData } from '../../actions/uiActions'
import { setSelectedHour } from '../../actions/dataActions'

import RangeSlider from './RangeSlider.jsx'

const mapStateToProps = createImmutableSelector(
  chartDataSelector,
  isLoadingSelector,
  selectedHourSelector,
  (chartData, isLoading, selectedHour) => ({
    chartData,
    isLoading,
    selectedHour,
  })
)
const mapDispatchToProps = dispatch => ({
  fetchData: hour => dispatch(fetchData(hour)),
  setSelectedHour: hour => dispatch(setSelectedHour(hour)),
})

export default connect(mapStateToProps, mapDispatchToProps)(RangeSlider)
