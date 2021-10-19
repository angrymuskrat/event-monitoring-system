import React, { useState, useEffect } from 'react'
import PropTypes, { object } from 'prop-types'

// styled
import Container from './SelectInput.styled'

// images
import { ChevronDown, ChevronUp } from '../../assets/svg/Chevron'

function SelectInput({ options, placeholder, handleSelectOption, value }) {
  const [values, setValues] = useState([])
  const [focusedValue, setFocusedValue] = useState(-1)
  const [isOpen, setIsOpen] = useState(false)

  useEffect(() => {
    if (value) {
      setValues(value)
    }
  }, [value])

  const onBlur = () => {
    const value = values[0]
    let currentFocusedValue = -1

    if (value) {
      currentFocusedValue = options.findIndex(option => option.value === value)
    }
    setFocusedValue(currentFocusedValue)
    setIsOpen(false)
  }

  const onKeyDown = e => {
    let currentFocusedValue = focusedValue
    switch (e.key) {
      case ' ':
        e.preventDefault()
        if (isOpen) {
          setIsOpen(true)
        }
        break
      case 'Escape':
      case 'Tab':
        if (isOpen) {
          e.preventDefault()
          setIsOpen(false)
        }
        break
      case 'Enter':
        setIsOpen(!isOpen)
        break
      case 'ArrowDown':
        e.preventDefault()
        currentFocusedValue = focusedValue
        if (focusedValue < options.length - 1) {
          currentFocusedValue++
          setValues([options[currentFocusedValue].value])
          setFocusedValue(currentFocusedValue)
        }
        break
      case 'ArrowUp':
        e.preventDefault()
        currentFocusedValue = focusedValue
        if (focusedValue > 0) {
          currentFocusedValue--
          setValues([options[currentFocusedValue].value])
          setFocusedValue(currentFocusedValue)
        }
        break
      default:
        break
    }
  }

  const onClick = () => {
    setIsOpen(!isOpen)
  }

  const onHoverOption = e => {
    const { value } = e.currentTarget.dataset
    const index = options.findIndex(option => option.value === value)
    setFocusedValue(index)
  }

  const onClickOption = e => {
    const { value } = e.currentTarget.dataset
    const index = values.indexOf(value)
    if (index === -1) {
      setValues(value)
      handleSelectOption(value)
      setIsOpen(false)
    } else {
      const newValues = values.splice(index, 1)
      setValues(newValues)
      setIsOpen(false)
    }
  }

  const renderValues = () => {
    if (values.length === 0) {
      return <div className="select__placeholder text_p2">{placeholder}</div>
    }

    return <div className="select__value">{values}</div>
  }

  const renderOptions = () => {
    if (!isOpen) {
      return null
    }
    return <div className="select__options">{options.map(renderOption)}</div>
  }

  const renderOption = (option, index) => {
    const { value } = option
    const selected = values.includes(value)

    let className = 'select__option'
    if (selected) className += ' select__option_selected'
    if (index === focusedValue) className += ' select__option_focused'
    if (index === options.length - 1) className += ' select__option_last'
    return (
      <div
        key={value}
        data-value={value}
        className={className}
        onMouseOver={onHoverOption}
        onClick={onClickOption}
      >
        {value}
      </div>
    )
  }

  return (
    <Container>
      <div
        className="select"
        tabIndex="0"
        onBlur={onBlur}
        onKeyDown={onKeyDown}
      >
        <div
          className={
            isOpen ? 'select__input select__input_open' : 'select__input'
          }
          onClick={onClick}
        >
          {renderValues()}
          <span className="select__arrow">
            {isOpen ? <ChevronUp /> : <ChevronDown />}
          </span>
        </div>
        {renderOptions()}
      </div>
    </Container>
  )
}

SelectInput.propTypes = {
  value: PropTypes.string,
  placeholder: PropTypes.string,
  options: PropTypes.arrayOf(object).isRequired,
  handleSelectOption: PropTypes.func.isRequired,
}

export default SelectInput
