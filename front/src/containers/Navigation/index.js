import createImmutableSelector from 'create-immutable-selector'
import { connect } from 'react-redux'
import { replace } from 'react-router-redux'

// selectors
import {
  currentCityIdSelector,
  currentCitySelector,
  currentCityCountrySelector,
} from '../../reducers/cityReducer'
import { boundsSelector } from '../../reducers/mapReducer'
import { selectedDateSelector } from '../../reducers/dataReducer'

// actions
import { setCurrentCityAndViewport } from '../../actions/dataActions'

// actions
import { clearStore } from '../../actions/dataActions'

import Navigation from './Navigation.jsx'

const mapStateToProps = createImmutableSelector(
  boundsSelector,
  currentCityIdSelector,
  currentCitySelector,
  currentCityCountrySelector,
  selectedDateSelector,
  (bounds, currentCityId, currentCity, currentCityCountry, selectedDate) => ({
    bounds,
    currentCityId,
    currentCity,
    currentCityCountry,
    selectedDate,
  })
)
const mapDispatchToProps = dispatch => ({
  clearStore: () => dispatch(clearStore()),
  setCurrentCityAndViewport: city => dispatch(setCurrentCityAndViewport(city)),
  replace: path => dispatch(replace(path)),
})

export default connect(mapStateToProps, mapDispatchToProps)(Navigation)
