import React, { useReducer, useEffect } from 'react'
import PropTypes from 'prop-types'
import { DateRangeInput } from '@datepicker-react/styled'
import { ThemeProvider } from 'styled-components'

// hooks
import { datepickerReducer } from '../../utils/hooks'

// styles
import { datepickerStyles } from '../../config/styles'

const initialState = {
  startDate: null,
  endDate: null,
  focusedInput: null,
}

function Datepicker({
  selectedDate,
  setSelectedDates,
  stateStartDate,
  stateEndDate,
}) {
  const [state, dispatch] = useReducer(datepickerReducer, initialState)
  useEffect(() => {
    if (stateStartDate || stateEndDate) {
      dispatch({
        type: 'dateChange',
        payload: {
          startDate: new Date(stateStartDate * 1000),
          endDate: new Date(stateEndDate * 1000),
          focusedInput: null,
        },
      })
    }
  }, [stateEndDate, stateStartDate])
  useEffect(() => {
    if (state.startDate && state.endDate && !state.focusedInput) {
      setSelectedDates(state)
    }
  }, [setSelectedDates, state])
  return (
    <ThemeProvider theme={datepickerStyles}>
      <DateRangeInput
        displayFormat="dd.MM.yyyy"
        onDatesChange={data => dispatch({ type: 'dateChange', payload: data })}
        onFocusChange={focusedInput =>
          dispatch({ type: 'focusChange', payload: focusedInput })
        }
        startDate={state.startDate} // Date or null
        endDate={state.endDate} // Date or null
        focusedInput={state.focusedInput} // START_DATE, END_DATE or null
        showSelectedDates={false}
        numberOfMonths={window.innerWidth > 1175 ? 2 : 1}
        minBookingDate={new Date('January 01, 2019 00:00:00')}
        initialVisibleMonth={new Date(selectedDate * 1000)}
      />
    </ThemeProvider>
  )
}
Datepicker.propTypes = {
  selectedDate: PropTypes.number,
  setSelectedDates: PropTypes.func.isRequired,
  stateStartDate: PropTypes.number,
  stateEndDate: PropTypes.number,
}
export default Datepicker
