import { fromJS } from 'immutable'
import createImmutableSelector from 'create-immutable-selector'

/******************************************************************************/
/******************************* TYPES ****************************************/
/******************************************************************************/

import { CLEAR_STORE, SET_CITY } from '../actions/types'

/******************************************************************************/
/******************************* INITIAL STATE ********************************/
/******************************************************************************/

const initialState = fromJS({
  cityId: null,
  city: null,
  country: null,
})

/******************************************************************************/
/******************************* SELECTORS ************************************/
/******************************************************************************/

const citySelector = createImmutableSelector(
  state => state.get('city'),
  city => city
)
export const currentCityIdSelector = createImmutableSelector(
  citySelector,
  city => city.get('cityId')
)
export const currentCitySelector = createImmutableSelector(citySelector, city =>
  city.get('city')
)
export const currentCityCountrySelector = createImmutableSelector(
  citySelector,
  city => city.get('country')
)

/******************************************************************************/
/******************************* REDUCERS *************************************/
/******************************************************************************/

export default function(state = initialState, { type, payload }) {
  switch (type) {
    case CLEAR_STORE:
      return state
        .set('cityId', null)
        .set('city', null)
        .set('country', null)
    case SET_CITY:
      return state
        .set('cityId', payload.id)
        .set('city', payload.city)
        .set('country', payload.country)
        .set('topLeft', String(payload.topLeft))
        .set('bottomRight', String(payload.bottomRight))
    default:
      return state
  }
}
