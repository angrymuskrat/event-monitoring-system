import React, { useState } from 'react'
import PropTypes from 'prop-types'
import moment from 'moment'

// components
import Input from '../../Input/Input.jsx'
import SelectInput from '../../Select/SelectInput.jsx'
import MultiplyInput from '../../Input/MultiplyInput.jsx'
import DatePicker from '../../DatePicker/Datepicker.jsx'

// styled
import Container from './SidebarSearch.styled'

// colors
import { orange } from '../../../config/styles'

function SidebarSearch({
  currentTab,
  setCurrentTab,
  display,
  handleSearchSubmit,
  handleSelectSearchQueryOption,
  handleSubmitSearchQueryValue,
  selectedDate,
  searchQueryFilters,
  searchParametersStartDate,
  searchParametersEndDate,
}) {
  const [dayValues, setDayValues] = useState({
    startDate: null,
    endDate: null,
  })
  const [inputTagValues, changeInputTagValues] = useState([])
  const [inputMentionValues, changeInputMentionValues] = useState([])
  const [error, setError] = useState('')
  const displayErrors = error => {
    setError(error)
    setTimeout(() => {
      setError('')
    }, 2000)
  }
  const handleChangeTagInputValues = (value, replace) => {
    if (replace) {
      changeInputTagValues([...value])
    } else {
      if (value[0] !== '#') {
        value = '#'.concat(value)
      } else {
        value = '#'.concat(value.slice(1))
      }
      changeInputTagValues([...inputTagValues, value])
    }
  }
  const handleChangeMentionInputValues = (value, replace) => {
    if (replace) {
      changeInputMentionValues([...value])
    } else {
      if (value[0] !== '@') {
        value = '@'.concat(value)
      } else {
        value = '@'.concat(value.slice(1))
      }
      changeInputMentionValues([...inputMentionValues, value])
    }
  }
  const handleSubmit = () => {
    const params = {
      tags: [...inputTagValues, ...inputMentionValues],
      startDate: moment(dayValues.startDate).unix(),
      endDate: moment(dayValues.endDate)
        .add({ hours: 23, minutes: 59 })
        .unix(),
    }
    if (inputTagValues.length === 0 && inputMentionValues.length === 0) {
      displayErrors('Please, enter at least one tag or mention')
      return
    }
    if (!dayValues.startDate || !dayValues.endDate) {
      displayErrors('Please, select date range')
      return
    }
    handleSearchSubmit(params)
  }
  const handleSetSelectedDate = data => {
    setDayValues(data)
  }
  const handleTabClick = tabName => {
    setCurrentTab(tabName)
  }
  return (
    <Container isShowSidebarSearch={display}>
      <div className="sidebar-search__menu">
        <div className="sidebar-search__menu-tab">
          <button
            className="sidebar-search__tab-button"
            onClick={() => handleTabClick('Search')}
            style={{
              borderBottom: currentTab === 'Search' && `3px solid ${orange}`,
            }}
          >
            <p className="text text_p2">Search parameters</p>
          </button>
        </div>
        <div className="sidebar-search__menu-tab">
          <button
            className="sidebar-search__tab-button"
            onClick={() => handleTabClick('Filters')}
            style={{
              borderBottom: currentTab === 'Filters' && `3px solid ${orange}`,
            }}
          >
            <p className="text text_p2">Filters</p>
          </button>
        </div>
      </div>
      {currentTab === 'Search' && (
        <div className="sidebar-search__tab">
          <div
            className="sidebar-search__error"
            style={{
              visibility: error ? `inherit` : 'hidden',
            }}
          >
            <p className="text text_p2">{error}</p>
          </div>
          <p className="text text_p2">Search by tags</p>
          <MultiplyInput
            handleChangeInputValues={handleChangeTagInputValues}
            inputValues={inputTagValues}
            type="tags"
          />
          <br />
          <p className="text text_p2">Search by mentions</p>
          <MultiplyInput
            handleChangeInputValues={handleChangeMentionInputValues}
            inputValues={inputMentionValues}
            type="mentions"
          />
          <br />
          <p className="text text_p2">Select date range</p>
          <div className="datepicker__container">
            <DatePicker
              selectedDate={selectedDate}
              setSelectedDates={handleSetSelectedDate}
              stateStartDate={searchParametersStartDate}
              stateEndDate={searchParametersEndDate}
            />
          </div>
          <div className="button-container">
            <button className="sibebar-search__button" onClick={handleSubmit}>
              Search events
            </button>
          </div>
        </div>
      )}
      {currentTab === 'Filters' && (
        <div className="sidebar-search__tab sidebar-search__tab_filters">
          <div>
            <p className="text text_p2">Search by word</p>
            <Input
              value={searchQueryFilters.get('keyword')}
              type="multiply"
              submitValue={handleSubmitSearchQueryValue}
            />
          </div>
          <SelectInput
            placeholder="Sort by:"
            value={searchQueryFilters.get('sortBy')}
            options={[
              { value: 'A - Z' },
              { value: 'Popular' },
              { value: 'Nearby' },
              { value: 'By time' },
            ]}
            handleSelectOption={handleSelectSearchQueryOption}
          />
        </div>
      )}
    </Container>
  )
}
SidebarSearch.propTypes = {
  currentTab: PropTypes.string.isRequired,
  setCurrentTab: PropTypes.func.isRequired,
  display: PropTypes.bool,
  handleSearchSubmit: PropTypes.func.isRequired,
  handleSelectSearchQueryOption: PropTypes.func.isRequired,
  handleSubmitSearchQueryValue: PropTypes.func.isRequired,
  selectedDate: PropTypes.number,
  searchQueryFilters: PropTypes.object,
  searchParametersStartDate: PropTypes.number,
  searchParametersEndDate: PropTypes.number,
}
export default SidebarSearch
