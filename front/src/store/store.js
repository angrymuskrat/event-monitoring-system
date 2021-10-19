import { createStore, applyMiddleware } from 'redux'
import createSagaMiddleware from 'redux-saga'
import { composeWithDevTools } from 'redux-devtools-extension'

import rootReducer from '../reducers'

import rootSaga from '../sagas/saga'

const configureStore = () => {
  const sagaMiddleware = createSagaMiddleware()
  const middlewares = [sagaMiddleware]
  const composeEnhancers = composeWithDevTools(applyMiddleware(...middlewares))
  const store = createStore(rootReducer, composeEnhancers)
  sagaMiddleware.run(rootSaga)
  return store
}

export default configureStore
