import React from 'react'
import PropTypes from 'prop-types'
//styled
import Button from './ChartStatButton.styled'

export default function ChartStatButton({
  isDemoPlay,
  onClick,
  type,
  children,
}) {
  return (
    <Button type={type} onClick={onClick} isDemoPlay={isDemoPlay}>
      <span className="text text_light text_s2">{children}</span>
    </Button>
  )
}
ChartStatButton.propTypes = {
  isDemoPlay: PropTypes.bool,
  onClick: PropTypes.func.isRequired,
  type: PropTypes.string,
  children: PropTypes.string.isRequired,
}
