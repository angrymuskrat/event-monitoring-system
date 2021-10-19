import React from 'react'
import PropTypes from 'prop-types'
import moment from 'moment'
import {
  Bar,
  BarChart,
  Cell,
  ResponsiveContainer,
  Tooltip,
  XAxis,
} from 'recharts'

// components
import ChartStat from '../../components/ChartStat/ChartStat.jsx'
import Container from './Chart.styled'
import CustomTooltip from '../../components/CustomTooltip/CustomTooltip.jsx'

// styles
import {
  chartDarkGrey,
  chartLightGrey,
  grey,
  lightOrange,
  orange,
} from '../../config/styles'

function Chart({
  chartData,
  fetchData,
  isDemoPlay,
  isShowAllEvents,
  highlightedHour,
  playDemo,
  selectedDate,
  selectedHour,
  setHighlightedHour,
  setSelectedHour,
  stopDemo,
  toggleAllEvents,
}) {
  const handleEnter = d => {
    setHighlightedHour(d.time)
  }
  const handleMouseLeave = d => {
    setHighlightedHour(null)
  }
  const handleClick = d => {
    if (isDemoPlay) {
      stopDemo()
    }
    setSelectedHour(d.time === selectedHour ? null : d.time)
    fetchData(d.time)
  }
  const handleDemoPlay = () => {
    playDemo()
  }
  const handleDemoStop = () => {
    stopDemo()
  }
  const handleSelectAll = () => {
    toggleAllEvents(chartData)
  }
  const calculateHeight = e => {
    if (e > 0) {
      return e * 0.5
    }
  }
  return (
    <Container>
      {chartData ? (
        <>
          <ChartStat
            chartData={chartData}
            selectedHour={selectedHour}
            selectedDate={selectedDate}
            handleDemoStop={handleDemoStop}
            handleDemoPlay={handleDemoPlay}
            handleSelectAll={handleSelectAll}
            isDemoPlay={isDemoPlay}
          />
          <ResponsiveContainer width="100%" height={210}>
            <BarChart
              data={chartData}
              margin={{
                top: 20,
                right: 0,
                left: 0,
                bottom: 5,
              }}
              onMouseLeave={handleMouseLeave}
            >
              <Tooltip cursor={false} content={<CustomTooltip />} />
              <XAxis
                dataKey="time"
                stroke={grey}
                tickLine={false}
                interval={2}
                tickFormatter={tick => `${moment.unix(tick).format('HH')}`}
              />
              <Bar
                dataKey="posts"
                stackId="a"
                onMouseEnter={handleEnter}
                onClick={handleClick}
              >
                {chartData.map((d, index) => (
                  <Cell
                    cursor="pointer"
                    width={10}
                    fill={
                      isShowAllEvents
                        ? chartDarkGrey
                        : d.time === highlightedHour
                        ? chartDarkGrey
                        : d.time === Number(selectedHour)
                        ? chartDarkGrey
                        : chartLightGrey
                    }
                    key={`cell-${index}`}
                  />
                ))}
              </Bar>
              <Bar
                dataKey="events"
                stackId="a"
                onMouseEnter={handleEnter}
                onClick={handleClick}
                radius={[0, 0, 0, 0]}
              >
                {chartData.map((d, index) => (
                  <Cell
                    width={10}
                    height={calculateHeight(d.events)}
                    cursor="pointer"
                    fill={
                      isShowAllEvents
                        ? orange
                        : d.time === highlightedHour
                        ? orange
                        : d.time === Number(selectedHour)
                        ? orange
                        : lightOrange
                    }
                    key={`cell-${index}`}
                  />
                ))}
              </Bar>
            </BarChart>
          </ResponsiveContainer>
        </>
      ) : null}
    </Container>
  )
}
Chart.propTypes = {
  chartData: PropTypes.array,
  fetchData: PropTypes.func,
  isDemoPlay: PropTypes.bool,
  isShowAllEvents: PropTypes.bool,
  playDemo: PropTypes.func,
  setHighlightedHour: PropTypes.func,
  setSelectedHour: PropTypes.func,
  selectedDate: PropTypes.number,
  selectedHour: PropTypes.string,
  stopDemo: PropTypes.func,
  toggleAllEvents: PropTypes.func,
}
export default Chart
