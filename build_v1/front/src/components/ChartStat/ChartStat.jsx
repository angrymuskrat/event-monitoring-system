import ChartStatStyled from './ChartStat.styled'
import PropTypes from 'prop-types'
import React from 'react'
import moment from 'moment'

// components
import ChartStatButton from '../Button/ChartStatButton.jsx'

function ChartStat({
  isDemoPlay,
  handleDemoPlay,
  handleDemoStop,
  handleSelectAll,
  selectedHour,
}) {
  return (
    <ChartStatStyled>
      <div>
        <p className="title title_h5">Timeline</p>
      </div>
      <div className="chartstat__nav">
        <button className="chartstat__nav-item">
          Selected time:{' '}
          <span className="chartstat__label">
            {selectedHour ? moment.unix(selectedHour).format('HH:mm') : `none`}
          </span>
        </button>
        {isDemoPlay ? (
          <ChartStatButton type="play" onClick={handleDemoStop}>
            Stop timeline
          </ChartStatButton>
        ) : (
          <ChartStatButton onClick={handleDemoPlay}>
            Play timeline
          </ChartStatButton>
        )}
        <ChartStatButton
          type="period"
          isDemoPlay={isDemoPlay}
          onClick={handleSelectAll}
        >
          All period
        </ChartStatButton>
      </div>
    </ChartStatStyled>
  )
}
ChartStat.propTypes = {
  isDemoPlay: PropTypes.bool,
  handleDemoPlay: PropTypes.func,
  handleDemoStop: PropTypes.func,
  handleSelectAll: PropTypes.func,
  selectedHour: PropTypes.string,
}
export default ChartStat
