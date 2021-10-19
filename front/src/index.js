import React from 'react'
import ReactDOM from 'react-dom'

import { createStore, applyMiddleware } from 'redux'
import { Provider } from 'react-redux'
import createSagaMiddleware from 'redux-saga'

import { createBrowserHistory } from 'history'
import { routerMiddleware } from 'connected-react-router'
// import { BrowserRouter as Router } from 'react-router-dom'
import { ConnectedRouter } from 'connected-react-router/immutable'
import { composeWithDevTools } from 'redux-devtools-extension'
import rootSaga from './sagas/saga'

// root reducer
import createRootReducer from './reducers'

import App from './containers/App'

export const history = createBrowserHistory()

const sagaMiddleware = createSagaMiddleware()

// all middlewares with redux dev tools
//const middlewares = [sagaMiddleware]

// const composeEnhancers = composeWithDevTools(applyMiddleware(...middlewares))

const store = createStore(
  createRootReducer(history),
  composeWithDevTools(
    applyMiddleware(
      routerMiddleware(history), // for dispatching history actions
      sagaMiddleware
    )
  )
)
sagaMiddleware.run(rootSaga)

ReactDOM.render(
  <Provider store={store}>
    <ConnectedRouter history={history}>
      <App />
    </ConnectedRouter>
  </Provider>,
  document.getElementById('root')
)
