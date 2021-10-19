import React, { useState, useEffect } from 'react'
import PropTypes from 'prop-types'

// images
import Magnifier from '../../assets/svg/Magnifier'

// styled
import Container from './Input.styled'
import InputButton from './InputButton'
import InputField from './InputField'

function Input({ submitValue, type, value }) {
  const [inputValue, changeInputValue] = useState('')
  useEffect(() => {
    if (value) {
      changeInputValue(value)
    }
  }, [value])
  const handleChange = e => {
    changeInputValue(e.target.value)
  }
  const handleKeyDown = e => {
    switch (e.key) {
      case 'Enter':
        submitValue(inputValue)
        break
      default:
        break
    }
  }
  const handleButtonClick = () => {
    submitValue(inputValue)
  }
  return (
    <Container type={type} in={type}>
      <InputField
        placeholder="Search..."
        value={inputValue}
        onChange={handleChange}
        onKeyDown={handleKeyDown}
      />
      <InputButton onClick={handleButtonClick}>
        <Magnifier />
      </InputButton>
    </Container>
  )
}
Input.propTypes = {
  value: PropTypes.string,
  submitValue: PropTypes.func.isRequired,
  type: PropTypes.string,
}
export default Input
