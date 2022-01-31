import React, { useState, useEffect } from 'react'
import PropTypes from 'prop-types'
import moment from 'moment'
import Slider from 'react-rangeslider'
import 'react-rangeslider/lib/index.css'

// styled
import Container from './RangeSlider.styled'

function RangeSlider({
  chartData,
  isLoading,
  setSelectedHour,
  selectedHour,
  fetchData,
}) {
  const [sliderPosition, setSliderPosition] = useState(selectedHour)

  useEffect(() => {
    setSliderPosition(selectedHour)
  }, [selectedHour])

  const handleChangeComplete = () => {
    setSliderPosition(sliderPosition)
    setSelectedHour(sliderPosition)
    fetchData(sliderPosition)
  }
  const handleChangeHorizontal = value => {
    setSliderPosition(value)
  }
  const formatTime = value => moment.unix(value).format('HH')

  if (isLoading) {
    return (
      <Container>
        <div className="rangeslider__placeholder" />
      </Container>
    )
  }

  return (
    <Container>
      <Slider
        min={chartData[0].time}
        max={chartData[23].time}
        step={3600}
        value={Number(sliderPosition)}
        format={formatTime}
        tooltip={true}
        onChange={handleChangeHorizontal}
        onChangeComplete={handleChangeComplete}
      />
    </Container>
  )
}
RangeSlider.propTypes = {
  chartData: PropTypes.array.isRequired,
  setSelectedHour: PropTypes.func.isRequired,
  fetchData: PropTypes.func.isRequired,
}

export default RangeSlider
