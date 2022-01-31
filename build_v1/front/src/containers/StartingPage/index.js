import { connect } from 'react-redux'

import { setCurrentCityAndViewport } from '../../actions/dataActions'

import StartingPage from './StartingPage.jsx'

const mapDispatchToProps = dispatch => ({
  setCurrentCityAndViewport: city => dispatch(setCurrentCityAndViewport(city)),
})

export default connect(null, mapDispatchToProps)(StartingPage)
