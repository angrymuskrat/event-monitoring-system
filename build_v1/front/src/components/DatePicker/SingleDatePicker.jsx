import React, { useState, useEffect } from 'react'
import PropTypes from 'prop-types'
import { DateSingleInput } from '@datepicker-react/styled'
import { ThemeProvider } from 'styled-components'

// styles
import { datepickerStyles, xsDatepickerStyles } from '../../config/styles'

function SingleDatepicker({ selectedDate, setSelectedDate, stateDate }) {
  const [date, setDate] = useState(null)
  const [showDatepicker, setShowDatepicker] = useState(false)
  const [position, setPosition] = useState('top')
  const [styles, setStyles] = useState(datepickerStyles)

  useEffect(() => {
    if (window.innerWidth < 1180) {
      setPosition('bottom')
    }
    if (window.innerWidth < 500) {
      setStyles(xsDatepickerStyles)
    }
  }, [])
  useEffect(() => {
    if (stateDate) {
      setDate(new Date(stateDate * 1000))
    }
  }, [stateDate])

  const handleDateChange = ({ date, showDatepicker }) => {
    setDate(date)
    setShowDatepicker(showDatepicker)
    setSelectedDate(date)
  }
  return (
    <ThemeProvider theme={styles}>
      <DateSingleInput
        onDateChange={data => handleDateChange(data)}
        onFocusChange={focusedInput => {
          setShowDatepicker(focusedInput)
        }}
        date={date}
        showDatepicker={showDatepicker}
        displayFormat="dd.MM.yyyy"
        minBookingDate={new Date('January 01, 2019 00:00:00')}
        initialVisibleMonth={new Date(selectedDate * 1000)}
        placement={position}
      />
    </ThemeProvider>
  )
}
SingleDatepicker.propTypes = {
  selectedDate: PropTypes.number,
  setSelectedDates: PropTypes.func,
  stateDate: PropTypes.number,
}
export default SingleDatepicker
