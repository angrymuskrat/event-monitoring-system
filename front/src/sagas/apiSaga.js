import {all, call, cancel, cancelled, fork, put, select, take, takeEvery, takeLatest,} from 'redux-saga/effects'
import {END, eventChannel} from 'redux-saga'
import {push, replace} from 'connected-react-router'

/******************************************************************************/
/******************************* UTILS ****************************************/
/******************************************************************************/
import moment from 'moment'
import {latLng, latLngBounds} from 'leaflet'
import {cities} from '../config/cities'

/******************************************************************************/
/******************************* TYPES ****************************************/
/******************************************************************************/
import {
    CHANGE_VIEWPORT,
    FETCH_ALL_DATA,
    FETCH_DATA,
    FETCH_FAILURE,
    FETCH_SUCCESS,
    INIT_SESSION,
    LOADING_FINISH,
    LOADING_START,
    PLAY_DEMO,
    SET_CURRENT_CITY_AND_VIEWPORT,
    SET_NEW_DATE,
    SET_SEARCH_PARAMETERS,
    SETUP_MAP_FROM_HISTORY,
    STOP_DEMO,
    TOGGLE_ALL_EVENTS,
    TOGGLE_POPUP,
} from '../actions/types'

/******************************************************************************/
/******************************* ACTIONS **************************************/
/******************************************************************************/
import {
    setAuth,
    setChartData,
    setEventsData,
    setHeatmapData,
    setPopupData,
    setSearchParameters,
    setSearchQuery,
    setSelectedDate,
    setSelectedHour,
} from '../actions/dataActions'
import {setCurrentCity} from '../actions/cityActions'
import {finishLoading, setAuthError, startLoading, stopDemo, toggleSidebar,} from '../actions/uiActions'
import {setBounds, setViewport} from '../actions/mapActions'

/******************************************************************************/
/******************************* HELPERS **************************************/
/******************************************************************************/
import {
    fetchChartTimeData,
    fetchEvents,
    fetchHeatmap,
    fetchPostData,
    fetchSearchEventsData,
    requestSessionCookie,
} from '../sagas/fetchData'
import {
    convertChartData,
    convertEventsToGeoJSON,
    convertHeatmapDataToLayer,
    convertPostData,
    convertSearchQueryToGeoJSON,
} from '../utils/utils'

/******************************************************************************/
/******************************* SAGAS *************************************/

/******************************************************************************/

/**
 * @func fetchSingleHeatmap - additional saga fn, that fetch each individual heatmap for one hour and put data into the store
 * @param {object} config - config object with parameters
 */
function* fetchSingleHeatmap(config) {
    try {
        const getIsCurrentHourExists = state =>
            state.getIn(['data', 'heatmap', `${config.time}`])
        const isCurrentHourExists = yield select(getIsCurrentHourExists)
        if (isCurrentHourExists !== undefined) {
            return
        }
        let heatmapData = yield call(fetchHeatmap, config)
        if (heatmapData.data) {
            heatmapData = yield call(convertHeatmapDataToLayer, heatmapData)
            const storeData = {
                [config.time]: heatmapData,
            }
            yield put(setHeatmapData(storeData))
        }
    } catch (err) {
        yield put({type: FETCH_FAILURE, payload: err})
    }
}

/**
 * @func fetchSingleEventLayer - additional saga fn, that fetch each individual events data for one hour and put data into the store
 * @param {object} config - config object with parameters
 */
function* fetchSingleEventLayer(config) {
    try {
        const getIsCurrentHourExists = state =>
            state.getIn(['data', 'events', `${config.time}`])
        const isCurrentHourExists = yield select(getIsCurrentHourExists)
        if (isCurrentHourExists !== undefined) {
            return
        }
        let eventsData = yield call(fetchEvents, config)
        if (eventsData.data) {
            eventsData = yield call(convertEventsToGeoJSON, eventsData)
            const storeData = {
                [config.time]: eventsData,
            }
            yield put(setEventsData(storeData))
        }
    } catch (err) {
        yield put({type: FETCH_FAILURE, payload: err})
    }
}

/**
 * @func handleToggleAllEvents - saga fn, that fetches heatmap data for all hours from timeline, when TOGGLE_ALL_EVENTS action is fired
 */
function* handleToggleAllEvents() {
    const getToggleAllEvents = state => state.getIn(['ui', 'isShowAllEvents'])
    const toggleActive = yield select(getToggleAllEvents)
    if (toggleActive) {
        try {
            yield put({type: LOADING_START})
            const getChartData = state => state.getIn(['data', 'chartData'])
            const hours = yield select(getChartData)
            const getCity = state => state.getIn(['city', 'cityId'])
            const city = yield select(getCity)
            const getTopLeft = state => state.getIn(['city', 'topLeft'])
            const topLeft = yield select(getTopLeft)
            const getBottomRight = state => state.getIn(['city', 'bottomRight'])
            const botRight = yield select(getBottomRight)
            yield all(
                hours.map(h => {
                    const config = {
                        city,
                        topLeft,
                        botRight,
                        time: h.time,
                    }
                    return call(fetchSingleHeatmap, config)
                })
            )
            yield all(
                hours.map(h => {
                    const config = {
                        city,
                        topLeft,
                        botRight,
                        time: h.time,
                    }
                    return call(fetchSingleEventLayer, config)
                })
            )
            yield put({type: LOADING_FINISH})
        } catch (err) {
            console.error(err)
        }
    }
}

/**
 * @func handleFetchData - saga fn, that handles fetching data when viewport settings changed
 * @param {Number} payload - time param
 */
function* handleFetchData({payload}) {
    try {
        const getIsCurrentHourExists = state =>
            state.getIn(['data', 'events', `${payload}`])
        const isCurrentHourExists = yield select(getIsCurrentHourExists)
        if (isCurrentHourExists !== undefined) {
            yield put({type: FETCH_SUCCESS})
            return
        }
        const getCity = state => state.getIn(['city', 'cityId'])
        const city = yield select(getCity)
        const getTopLeft = state => state.getIn(['city', 'topLeft'])
        const topLeft = yield select(getTopLeft)
        const getBottomRight = state => state.getIn(['city', 'bottomRight'])
        const botRight = yield select(getBottomRight)
        const time = payload
        const config = {
            city,
            topLeft,
            botRight,
            time,
        }
        let mapData = yield call(fetchHeatmap, config)

        if (mapData.data) {
            mapData = yield call(convertHeatmapDataToLayer, mapData)
            yield put(setHeatmapData(mapData))
        }
        let eventsData = yield call(fetchEvents, config)
        if (eventsData.data) {
            eventsData = yield call(convertEventsToGeoJSON, eventsData)
            yield put(setEventsData(eventsData))
        }
        yield put({type: FETCH_SUCCESS})
    } catch (err) {
        console.error(err)
        yield put({type: FETCH_FAILURE, payload: err})
    }
}

/**
 * @func handleFetchAllData - saga fn, that fetches all data for initial app setup
 * @param {Object} data - config object
 */
function* handleFetchAllData(data) {
    if (data.type) {
        data = data.payload
    }
    try {
        const getIsCurrentHourExists = state =>
            state.getIn([
                'data',
                'heatmap',
                `${data.selectedHour ? data.selectedHour : data.time}`,
            ])
        const isCurrentHourExists = yield select(getIsCurrentHourExists)
        if (isCurrentHourExists === undefined) {
            let mapData = yield call(fetchHeatmap, data)
            if (mapData.data) {
                mapData = yield call(convertHeatmapDataToLayer, mapData)
                yield put(setHeatmapData(mapData))
            }
            let eventsData = yield call(fetchEvents, data)
            if (eventsData.data) {
                eventsData = yield call(convertEventsToGeoJSON, eventsData)
                yield put(setEventsData(eventsData))
            }
        }
        let finish
        let start
        if (data.setup) {
            start = data.startDate
            finish = Number(data.startDate) + 82800
        } else {
            start = Number(data.time)
            finish = Number(data.time) + 82800
        }
        let chartData = yield call(fetchChartTimeData, data.city, start, finish)
        chartData = yield call(convertChartData, chartData)
        yield put(setChartData(chartData))
        yield put({type: FETCH_SUCCESS})
    } catch (err) {
        console.error(err)
        yield put({type: FETCH_FAILURE, payload: err})
    }
}

/**
 * @func countUpDemoPlay - helper fn, that counting time for demo play
 * @param {Number} start - count start
 * @param {Number} end - count end
 * @param {Array} array
 */
function countUpDemoPlay(start, end, array) {
    return eventChannel(emitter => {
        if (start === end) {
            emitter(END)
        }
        const iv = setInterval(() => {
            let time = array[start].time
            if (start < end) {
                start += 1
                emitter(time)
            }
            if (start === end) {
                emitter(END)
            }
        }, 1000)
        return () => {
            clearInterval(iv)
        }
    })
}

/**
 * @func handleDemoPlay - saga fn for playing timeline
 */
function* handlePlayDemo() {
    yield put({type: LOADING_START})
    const getChartData = state => state.getIn(['data', 'chartData'])
    const chartData = yield select(getChartData)
    const getSelectedHour = state => state.getIn(['data', 'selectedHour'])
    const selectedHour = yield select(getSelectedHour)
    let index = chartData.findIndex(h => h.time === Number(selectedHour)) + 1
    let end = chartData.length
    try {
        let hours = chartData.slice(index)
        const getCity = state => state.getIn(['city', 'cityId'])
        const city = yield select(getCity)
        const getTopLeft = state => state.getIn(['city', 'topLeft'])
        const topLeft = yield select(getTopLeft)
        const getBottomRight = state => state.getIn(['city', 'bottomRight'])
        const botRight = yield select(getBottomRight)
        yield all(
            hours.map(h => {
                const config = {
                    city,
                    topLeft,
                    botRight,
                    time: h.time,
                }
                return call(fetchSingleHeatmap, config)
            })
        )
        yield all(
            hours.map(h => {
                const config = {
                    city,
                    topLeft,
                    botRight,
                    time: h.time,
                }
                return call(fetchSingleEventLayer, config)
            })
        )
        const chan = yield call(countUpDemoPlay, index, end, chartData)
        yield put({type: LOADING_FINISH})
        while (true) {
            let seconds = yield take(chan)
            yield put(setSelectedHour(seconds))
        }
    } finally {
        yield put(stopDemo())
        if (yield cancelled()) {
            yield put(stopDemo())
        }
    }
}

/**
 * @func handleSetSelectedDate - saga fn, that set selected date && selected hour, and fethes data
 * @param {Object} payload - selected time
 */
function* handleSetSelectedDate({payload}) {
    try {
        yield put({type: LOADING_START})
        const date =
            Object.prototype.toString.call(payload) === '[object Date]'
                ? moment(payload).unix()
                : payload
        const getPrevSelectedHour = state => state.getIn(['data', 'selectedHour'])
        const prevSelectedHour = yield select(getPrevSelectedHour)
        const getPrevSelectedDate = state => state.getIn(['data', 'selectedDate'])
        const prevSelectedDate = yield select(getPrevSelectedDate)
        const nextSelectedHour = date + (prevSelectedHour - prevSelectedDate)

        const getCity = state => state.getIn(['city', 'cityId'])
        const city = yield select(getCity)
        const getTopLeft = state => state.getIn(['city', 'topLeft'])
        const topLeft = yield select(getTopLeft)
        const getBottomRight = state => state.getIn(['city', 'bottomRight'])
        const botRight = yield select(getBottomRight)
        const getCurrentBottomRight = state =>
            state.getIn(['map', 'bounds', 'bottomRight'])
        const currentBotRight = yield select(getCurrentBottomRight)
        const getCurrentTopLeft = state => state.getIn(['map', 'bounds', 'topLeft'])
        const currentTopLeft = yield select(getCurrentTopLeft)

        yield put(setSelectedDate(date))
        yield put(setSelectedHour(nextSelectedHour))
        const config = {
            city,
            topLeft,
            botRight,
            time: date,
            selectedHour: nextSelectedHour,
        }
        const getIsShowAllEvents = state => state.getIn(['ui', 'isShowAllEvents'])
        const isShowAllEvents = yield select(getIsShowAllEvents)

        const getSearchQuery = state => state.getIn(['data', 'searchQuery'])
        const searchQuery = yield select(getSearchQuery)

        if (searchQuery) {
            yield put(setSearchQuery(null))
            yield put(setSearchParameters(null))
        }

        if (isShowAllEvents) {
            const finish = date + 82800
            let chartData = yield call(fetchChartTimeData, city, date, finish)
            chartData = yield call(convertChartData, chartData)
            yield put(setChartData(chartData))
            yield call(handleToggleAllEvents)
        } else {
            yield call(handleFetchAllData, config)
        }
        yield put(
            replace(
                `/map/${city}/${currentTopLeft.join()}/${currentBotRight.join()}/${date}`
            )
        )
    } catch (err) {
        console.error(err)
    }
}

/**
 * @func handleSetCityAndViewport - saga fn for SET_CURRENT_CITY_AND_VIEWPORT action; sets initial data when city is selected
 * @param {Object} payload - config object with city parameters
 */
function* handleSetCityAndViewport({payload}) {
    // четверг, 9 мая 2019 г., 0:00:00 GMT+03:00
    const date = 1557349200
    // четверг, 9 мая 2019 г., 22:00:00 GMT+03:00
    const selectedHour = 1557428400
    if (typeof payload === 'string') {
        const cityName = payload.split(',')[0]
        payload = cities.find(c => c.city === cityName)
    }
    yield put(setSelectedDate(date))
    yield put(setSelectedHour(selectedHour))
    const city = {
        id: payload.id,
        city: payload.city,
        country: payload.country,
        topLeft: payload.topLeft,
        bottomRight: payload.bottomRight,
    }
    const viewport = {
        center: [payload.lat, payload.lng],
        zoom: 13,
    }
    yield put(setCurrentCity(city))
    yield put(setViewport(viewport))
    const bounds = {topLeft: payload.topLeft, bottomRight: payload.bottomRight}
    yield put(setBounds(bounds))
    yield put(
        push(`/map/${city.id}/${payload.topLeft}/${payload.bottomRight}/${date}`)
    )
    const config = {
        city: payload.id,
        topLeft: payload.topLeft,
        botRight: payload.bottomRight,
        time: selectedHour,
        setup: true,
        startDate: date,
    }
    yield call(handleFetchAllData, config)
    yield put(toggleSidebar())
}

/**
 * @func setInitialBounds - saga fn; sets initial bounds when app initialized from url history
 * @param {Object} payload - config object with bounds parameters
 */
function* setInitialBounds(payload) {
    let historyTopLeft = payload.topLeft.split(',').map(d => parseFloat(d))
    let historyBotRight = payload.botRight.split(',').map(d => parseFloat(d))

    const corner1 = latLng(historyTopLeft[0], historyTopLeft[1])
    const corner2 = latLng(historyBotRight[0], historyBotRight[1])
    const center = latLngBounds(corner1, corner2).getCenter()
    const newViewport = {
        zoom: 11,
        center: [center.lat, center.lng],
    }
    yield put(setViewport(newViewport))
    const newBounds = {topLeft: historyTopLeft, bottomRight: historyBotRight}
    yield put(setBounds(newBounds))
}

/**
 * @func handleSetupMapFromHistory - saga fn; fires when app is initialized from url history
 * @param {Object} payload - config object
 */
function* handleSetupMapFromHistory({payload}) {
    try {
        let currentHistoryCity = payload.city
        const findCity = configCity => {
            return configCity.id === currentHistoryCity
        }
        currentHistoryCity = cities.find(findCity)
        if (currentHistoryCity === undefined) {
            yield put(push('/'))
        }
        yield put(setCurrentCity(currentHistoryCity))
        yield put(setSelectedDate(Number(payload.time)))
        yield put(setSelectedHour(Number(payload.time)))
        yield call(setInitialBounds, payload)
        const config = {
            ...payload,
            topLeft: currentHistoryCity.topLeft,
            botRight: currentHistoryCity.bottomRight,
            time: Number(payload.time),
        }
        yield call(handleFetchAllData, config)
        yield put(toggleSidebar())
    } catch (error) {
        console.error(error)
    }
}

/**
 * @func handleChangeViewport - saga fn; fires when viewport of the map is changing
 * @param {Object} payload - config object with viewport settings
 */
function* handleChangeViewport({payload}) {
    yield put(setViewport(payload.viewport))
}

/**
 * @func handleSearchEvents - saga fn; fires when SET_SEARCH_PARAMETERS is dispatched; calculates config object to fetch events data in certain time limits
 * @param {Object} payload - config object with search tags, startDate, endDate for searching
 */
function* handleSearchEvents({payload}) {
    try {
        yield put(startLoading())
        const getCurrentCityId = state => state.getIn(['city', 'cityId'])
        const currentCityId = yield select(getCurrentCityId)
        const searchTags = payload.tags.map(tag => {
            if (tag[0] === '@') {
                return '%40'.concat(tag.slice(1))
            } else {
                return '%23'.concat(tag.slice(1))
            }
        })
        const config = {
            city: currentCityId,
            tags: searchTags.join(),
            start: payload.startDate,
            finish: payload.endDate,
        }
        let eventsData = yield call(fetchSearchEventsData, config)
        if (eventsData.data) {
            eventsData = yield call(convertSearchQueryToGeoJSON, eventsData)
            yield put(setSearchQuery(eventsData))
        } else {
            yield put(setSearchQuery(null))
        }
        yield put(finishLoading())
    } catch (error) {
        console.error(error)
    }
}

/**
 * @func fetchSingleInstagramPost - saga fn; fetches single instagram post
 * @param {String} postcode - instagram post postcode
 */
function* fetchSingleInstagramPost(postcode) {
    // console.log(`poascode`, postcode)
    try {
        yield put(startLoading())
        const getCurrentCityId = state => state.getIn(['city', 'cityId'])
        const currentCityId = yield select(getCurrentCityId)
        const config = {
            city: currentCityId,
            postcode: postcode
        }
        let postData = yield call(fetchPostData, config)
        if (!postData) {
            return
        } else {
            postData = yield call(
                convertPostData,
                postData.data
            )
            return postData
        }
    } catch (err) {
        yield put({type: FETCH_FAILURE, payload: err})
    }
}

/**
 * @func handleFetchPopupData - saga fn; fires when TOGGLE_POPUP action is dispatched; fetches instagram media for popup
 * @param {Array} payload - instagram postcodes
 */
function* handleFetchPopupData({payload}) {
    const getIsPopupOpen = state => state.getIn(['ui', 'isPopupOpen'])
    const isPopupOpen = yield select(getIsPopupOpen)
    const getPopupData = state => state.getIn(['data', 'currentPopupData'])
    const popupData = yield select(getPopupData)
    if (popupData) {
        yield put(setPopupData(null))
    }
    if (isPopupOpen) {
        const data = yield all(
            payload.map(p => {
                return call(fetchSingleInstagramPost, p)
            })
        )
        const filteredData = data.filter(d => d !== undefined)
        yield put(setPopupData(filteredData))
    }
}

const delay = time => new Promise(resolve => setTimeout(resolve, time))

function* handleInitSession({payload}) {
    try {
        const res = yield call(requestSessionCookie, payload)
        console.log('res', res
        )
        if (res.status === 200) {
            yield put(setAuth())
            yield put(replace(`/`))
        }
    } catch (error) {
        yield put(setAuthError('Incorrect login / password, please try again'))
        yield call(delay, 2000)
        yield put(setAuthError(''))
        console.error(error)
    }
}

/******************************************************************************/
/******************************* WATCHERS *************************************/

/******************************************************************************/

export function* fetchDataSaga() {
    yield takeLatest(FETCH_DATA, handleFetchData)
}

export function* fetchAllDataSaga() {
    yield takeLatest(FETCH_ALL_DATA, handleFetchAllData)
}

export function* toggleAllEventsSaga() {
    yield takeLatest(TOGGLE_ALL_EVENTS, handleToggleAllEvents)
}

export function* setSelectedDateSaga() {
    yield takeLatest(SET_NEW_DATE, handleSetSelectedDate)
}

export function* setCityAndViewportSaga() {
    yield takeLatest(SET_CURRENT_CITY_AND_VIEWPORT, handleSetCityAndViewport)
}

export function* setupMapFromHistorySaga() {
    yield takeLatest(SETUP_MAP_FROM_HISTORY, handleSetupMapFromHistory)
}

export function* changeViewportSaga() {
    yield takeLatest(CHANGE_VIEWPORT, handleChangeViewport)
}

export function* searchEventsSaga() {
    yield takeEvery(SET_SEARCH_PARAMETERS, handleSearchEvents)
}

export function* fetchPopupDataSaga() {
    yield takeLatest(TOGGLE_POPUP, handleFetchPopupData)
}

export function* initSessionSaga() {
    yield takeLatest(INIT_SESSION, handleInitSession)
}

export function* timerSaga() {
    while (yield take(PLAY_DEMO)) {
        // starts the task in the background
        const channel = yield fork(handlePlayDemo)
        // wait for the user stop action
        yield take(STOP_DEMO)
        // user clicked stop. cancel the background task
        // this will cause the forked bgSync task to jump into its finally block
        yield cancel(channel)
    }
}
