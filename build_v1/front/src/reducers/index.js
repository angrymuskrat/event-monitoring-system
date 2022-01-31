import { combineReducers } from 'redux-immutable'
import { connectRouter } from 'connected-react-router/immutable'

import dataReducer from './dataReducer'
import cityReducer from './cityReducer'
import mapReducer from './mapReducer'
import uiReducer from './uiReducer'

const createRootReducer = history =>
  combineReducers({
    router: connectRouter(history),
    city: cityReducer,
    ui: uiReducer,
    map: mapReducer,
    data: dataReducer,
  })
export default createRootReducer
