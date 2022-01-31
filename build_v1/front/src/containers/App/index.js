import { connect } from 'react-redux'
import createImmutableSelector from 'create-immutable-selector'

import { authSelector } from '../../reducers/dataReducer'

import App from './App'

const mapStateToProps = createImmutableSelector(authSelector, auth => ({
  auth,
}))

const mapDispatchToProps = dispatch => ({})

export default connect(mapStateToProps, mapDispatchToProps)(App)
