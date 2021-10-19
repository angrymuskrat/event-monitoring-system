import { connect } from 'react-redux'
import createImmutableSelector from 'create-immutable-selector'

// selectors
import { authErrorSelector } from '../../reducers/uiReducer'

// actions
import { initSession } from '../../actions/dataActions'

import LoginPage from './LoginPage.jsx'

const mapStateToProps = createImmutableSelector(
  authErrorSelector,
  authError => ({
    authError,
  })
)
const mapDispatchToProps = dispatch => ({
  initSession: config => dispatch(initSession(config)),
})

export default connect(mapStateToProps, mapDispatchToProps)(LoginPage)
