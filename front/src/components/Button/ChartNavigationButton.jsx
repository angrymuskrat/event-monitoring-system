import React from 'react'
import PropTypes from 'prop-types'

//styled
import Button from './ChartNavigationButton.styled'

export default function ChartStatButton({ onClick, type, children }) {
  return (
    <Button type={type} onClick={onClick}>
      <span className="text text_s1">{children}</span>
    </Button>
  )
}
ChartStatButton.propTypes = {
  onClick: PropTypes.func.isRequired,
  type: PropTypes.string,
  children: PropTypes.string.isRequired,
}
