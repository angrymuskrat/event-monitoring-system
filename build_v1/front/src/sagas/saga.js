import { all } from 'redux-saga/effects'
import {
  changeViewportSaga,
  fetchDataSaga,
  fetchAllDataSaga,
  toggleAllEventsSaga,
  setSelectedDateSaga,
  setCityAndViewportSaga,
  timerSaga,
  setupMapFromHistorySaga,
  searchEventsSaga,
  fetchPopupDataSaga,
  initSessionSaga,
} from './apiSaga'

function* rootSaga() {
  yield all([
    fetchDataSaga(),
    fetchAllDataSaga(),
    toggleAllEventsSaga(),
    timerSaga(),
    setSelectedDateSaga(),
    setCityAndViewportSaga(),
    setupMapFromHistorySaga(),
    changeViewportSaga(),
    searchEventsSaga(),
    fetchPopupDataSaga(),
    initSessionSaga(),
  ])
}

export default rootSaga
