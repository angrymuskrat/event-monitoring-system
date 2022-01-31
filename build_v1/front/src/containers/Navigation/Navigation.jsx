import React from 'react'
import PropTypes from 'prop-types'
import { Link, useLocation } from 'react-router-dom'

// components
import AutocompleteInput from '../../components/Input/AutocompleteInput.jsx'

// styled
import Container from './Navigation.styled'

// images
import { BackButton } from '../../assets/svg/BackButton'
import EMlogo from '../../assets/svg/EMlogo'

function Navigation({
  bounds,
  clearStore,
  currentCityId,
  currentCity,
  currentCityCountry,
  replace,
  selectedDate,
  setCurrentCityAndViewport,
}) {
  const location = useLocation()

  const handleGoToMainClick = () => {
    clearStore()
  }
  const onButtonClick = () => {
    setCurrentCityAndViewport()
  }
  const handleGoBack = () => {
    const topLeft = bounds.toJS().topLeft.join()
    const bottomRight = bounds.toJS().bottomRight.join()
    replace(`/map/${currentCityId}/${topLeft}/${bottomRight}/${selectedDate}`)
  }
  return (
    <Container>
      <nav className="navigation">
        {(location.pathname === '/about' || location.pathname === '/team') &&
        currentCity ? (
          <button className="navigation__button" onClick={handleGoBack}>
            <BackButton />
          </button>
        ) : (
          <div className="navigation__container">
            <div className="navigation__logo">
              <Link
                className="text text_p3 text_link"
                onClick={handleGoToMainClick}
                to="/"
              >
                <EMlogo />
              </Link>
            </div>
            <div className="navigation__input">
              {currentCity && currentCityCountry && (
                <AutocompleteInput
                  defaultValue={`${currentCity}, ${currentCityCountry}`}
                  type="sidebar"
                  handleClick={onButtonClick}
                />
              )}
              {!currentCity && !currentCityCountry && (
                <AutocompleteInput type="sidebar" handleClick={onButtonClick} />
              )}
            </div>
          </div>
        )}

        <div className="navigation__links">
          <ul className="navigation__links-container">
            <li className="navigation__link">
              <Link to="/team" className="text text_p3">
                Team
              </Link>
            </li>
            <li className="navigation__link">
              <Link to="/about" className="text text_p3">
                About
              </Link>
            </li>
          </ul>
        </div>
      </nav>
    </Container>
  )
}
Navigation.propTypes = {
  bounds: PropTypes.object,
  clearStore: PropTypes.func.isRequired,
  currentCityId: PropTypes.string,
  currentCity: PropTypes.string,
  currentCityCountry: PropTypes.string,
  replace: PropTypes.func.isRequired,
  selectedDate: PropTypes.number,
  setCurrentCityAndViewport: PropTypes.func.isRequired,
}
export default Navigation
