import React, { useState, useEffect, memo } from 'react'
import PropTypes from 'prop-types'
import { cities } from '../../config/cities'

// images
import Magnifier from '../../assets/svg/Magnifier'

// styled
import InputStyled from './Input.styled'
import InputButton from './InputButton'
import InputField from './InputField'
import Suggestions from './Suggestions'

function AutocompleteInput({ type, defaultValue, handleClick }) {
  const [activeSuggestion, setActiveSuggestion] = useState(0)
  const [filteredSuggestions, setFilteredSuggestions] = useState([])
  const [showSuggestions, setShowSuggestions] = useState(false)
  const [userInput, setUserInput] = useState('')
  const suggestions = cities.map(c => {
    return {
      city: `${c.city}, ${c.country}`,
      avaliable: c.avaliable,
    }
  })

  useEffect(() => {
    if (defaultValue) {
      setUserInput(defaultValue)
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])
  const onChange = e => {
    const userInput = e.currentTarget.value

    const filteredSuggestions = suggestions.filter(
      suggestion =>
        suggestion.city.toLowerCase().indexOf(userInput.toLowerCase()) > -1
    )
    setActiveSuggestion(0)
    setFilteredSuggestions(filteredSuggestions)
    setShowSuggestions(true)
    setUserInput(e.currentTarget.value)
  }
  const onClick = e => {
    setActiveSuggestion(0)
    setFilteredSuggestions([])
    setShowSuggestions(false)
    setUserInput(e.currentTarget.innerText)
  }
  const onKeyDown = e => {
    // User pressed the enter key
    if (e.keyCode === 13) {
      setActiveSuggestion(0)
      setShowSuggestions(false)
      setUserInput(filteredSuggestions[activeSuggestion])
    }
    // User pressed the up arrow
    else if (e.keyCode === 38) {
      if (activeSuggestion === 0) {
        return
      }
      setActiveSuggestion(activeSuggestion - 1)
    }
    // User pressed the down arrow
    else if (e.keyCode === 40) {
      if (activeSuggestion - 1 === filteredSuggestions.length) {
        return
      }
      setActiveSuggestion(activeSuggestion + 1)
    }
  }
  const handleButtonClick = () => {
    handleClick(userInput)
  }

  return (
    <>
      <InputStyled type={type} showSuggestions={showSuggestions}>
        <InputField
          placeholder="Search..."
          onChange={onChange}
          onKeyDown={onKeyDown}
          value={userInput}
        />
        <InputButton
          style={{
            cursor: userInput ? 'pointer' : 'not-allowed',
          }}
          onClick={handleButtonClick}
        >
          <Magnifier />
        </InputButton>
        <Suggestions showSuggestions={showSuggestions}>
          {filteredSuggestions.map((suggestion, index) => {
            let className
            if (index === activeSuggestion) {
              className = 'suggestion-active'
            }
            if (suggestion.avaliable) {
              return (
                <li
                  className={className}
                  key={suggestion.city}
                  onClick={onClick}
                >
                  <div className="suggestion-container">
                    <p className="text">{suggestion.city}</p>
                  </div>
                </li>
              )
            } else {
              return (
                <li className={className} key={suggestion.city}>
                  <div
                    className="suggestion-container"
                    style={{ cursor: 'not-allowed' }}
                  >
                    <p>
                      <span className="text text_input text_crossed">
                        {suggestion.city}
                      </span>{' '}
                      <span className="text text_italic">avaliable soon</span>
                    </p>
                  </div>
                </li>
              )
            }
          })}
        </Suggestions>
      </InputStyled>
    </>
  )
}
AutocompleteInput.propTypes = {
  defaultValue: PropTypes.string,
  type: PropTypes.string,
  value: PropTypes.string,
}
export default memo(AutocompleteInput)
