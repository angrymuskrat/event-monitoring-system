import PropTypes from 'prop-types'
import React from 'react'
import moment from 'moment'

// styled
import Container from './CustomTooltip.styled'

const CustomTooltip = ({ active, label, payload }) => {
  if (active) {
    return (
      <Container>
        <p className="text text_p1 chart__tooltip-text_main">{`Time: ${moment
          .unix(label)
          .format('HH:mm')}`}</p>
        <p className="text text_p2 chart__tooltip-text">
          Events for hour:{'  '}
          <span className="text_p2-bold chart__tooltip-label chart__tooltip-label_events">{`${payload[1].value}`}</span>
        </p>
        <p className="text text_p2 chart__tooltip-text">
          Posts for hour:{'  '}
          <span className="text_p2-bold chart__tooltip-label chart__tooltip-label_posts">{`${payload[0].value}`}</span>
        </p>
      </Container>
    )
  }
  return null
}
CustomTooltip.propTypes = {
  active: PropTypes.bool,
  label: PropTypes.number,
  payload: PropTypes.array,
}
export default CustomTooltip
