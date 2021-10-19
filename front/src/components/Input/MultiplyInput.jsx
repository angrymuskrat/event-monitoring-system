import React, { useState } from 'react'
import PropTypes from 'prop-types'

// images
import Plus from '../../assets/svg/Plus'

// styled
import Container from './Input.styled'
import InputButton from './InputButton'
import InputField from './InputField'
import MultiplyInputContainer from './MultiplyInput.styled'

function MultiplyInput({ type, handleChangeInputValues, inputValues }) {
  // local state
  const [currentValue, changeCurrentValue] = useState('')

  const handleChange = e => {
    changeCurrentValue(e.target.value)
  }
  const handleKeyDown = e => {
    switch (e.key) {
      case 'Enter':
        if (currentValue.length === 0 || /^\s+$/.test(currentValue)) {
          return
        }
        handleChangeInputValues(currentValue)
        changeCurrentValue('')
        break
      default:
        break
    }
  }
  const handleButtonClick = () => {
    if (currentValue.length === 0 || /^\s+$/.test(currentValue)) {
      return
    }
    handleChangeInputValues(currentValue)
    changeCurrentValue('')
  }
  const handleDelete = v => {
    const value = v.currentTarget.querySelector('.text').innerHTML
    const newValues = inputValues.filter(v => v !== value)
    handleChangeInputValues(newValues, true)
  }
  return (
    <MultiplyInputContainer type={type} in={type}>
      <Container>
        <InputField
          placeholder="Type here..."
          value={currentValue}
          onChange={handleChange}
          onKeyDown={handleKeyDown}
        />
        <InputButton onClick={handleButtonClick}>
          <Plus />
        </InputButton>
      </Container>
      {inputValues &&
        inputValues.map((v, i) => (
          <div
            className="input__tag"
            key={v + i}
            onClick={v => handleDelete(v)}
          >
            <p className="text text_subtitle">{v}</p>
            <button className="input__delete-button">x</button>
          </div>
        ))}
    </MultiplyInputContainer>
  )
}
MultiplyInput.propTypes = {
  handleChangeInputValues: PropTypes.func.isRequired,
  type: PropTypes.string,
  inputValues: PropTypes.array,
}
export default MultiplyInput
