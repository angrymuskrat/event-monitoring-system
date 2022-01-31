import React from 'react'
import PropTypes from 'prop-types'

// styled
import Container from './SidebarTitle.styled'

function SidebarTitle({ time, events }) {
  return (
    <Container>
      <h2 className="title title_h2 title_sidebar">
        {`${events === null ? 0 : events} ${events === 1 ? 'event' : 'events'}`}{' '}
        <span className="title title_h2 title_light title_sidebar">
          for {time}
        </span>
      </h2>
    </Container>
  )
}
SidebarTitle.propTypes = {
  time: PropTypes.string,
  events: PropTypes.number,
}
export default SidebarTitle
