import React, { Suspense, lazy } from 'react'
import PropTypes from 'prop-types'
import moment from 'moment'

// containers
import Chart from '../../../containers/Chart'
import RangeSlider from '../../../containers/RangeSlider'

// components
import Button from '../../Button/ChartNavigationButton.jsx'
import Datepicker from '../../DatePicker/SingleDatePicker.jsx'
import Input from '../../Input/Input.jsx'
import SidebarTitle from '../SidebarTitle/SidebarTitle.jsx'
import SelectInput from '../../Select/SelectInput.jsx'

// styled
import Container from './SidebarBody.styled'

// lazy
const SidebarSearch = lazy(() => import('../SidebarSearch/SidebarSearch.jsx'))

function SidebarBody({
  chartData,
  currentTab,
  setCurrentTab,
  eventFilters,
  isSearchingEvents,
  isShowAllEvents,
  handleDateChange,
  handleSubmitValue,
  handleSearchSubmit,
  handleSelectOption,
  handleSelectSearchQueryOption,
  handleSubmitSearchQueryValue,
  selectedDate,
  selectedHour,
  searchQueryFilters,
  searchParametersStartDate,
  searchParametersEndDate,
}) {
  const calculateEvents = () => {
    if (chartData) {
      if (isShowAllEvents) {
        return chartData.reduce((accum, data) => (accum += data.events), 0)
      } else {
        const currentEventsNumber = chartData.find(
          d => d.time === Number(selectedHour)
        )
        if (!currentEventsNumber) {
          return null
        } else {
          return currentEventsNumber.events
        }
      }
    } else {
      return null
    }
  }
  const handleBackClick = () => {
    const prevDay = Number(parseInt(selectedDate) - 86400)
    handleDateChange(prevDay)
  }
  const handleNextClick = () => {
    const nextDay = Number(parseInt(selectedDate) + 86400)
    handleDateChange(nextDay)
  }
  return (
    <Container>
      {!isSearchingEvents && chartData !== null && (
        <>
          <SidebarTitle
            time={
              isShowAllEvents
                ? moment.unix(selectedDate).format('DD.MM.YYYY')
                : moment.unix(selectedHour).format('HH:mm, DD.MM.YYYY')
            }
            events={calculateEvents()}
          />
          {window.innerWidth > 1175 ? <Chart /> : <RangeSlider />}
          <div className="sidebar__button-container">
            <Button onClick={handleBackClick} type="back">
              {`<  ${moment
                .unix(Number(selectedDate) - 86400)
                .format('DD.MM')}`}
            </Button>
            <Button onClick={handleNextClick} type="back">
              {`${moment
                .unix(Number(selectedDate) + 86400)
                .format('DD.MM')}  >`}
            </Button>
          </div>
          <Input
            value={eventFilters.get('keyword')}
            type="multiply"
            submitValue={handleSubmitValue}
          />
          <div className="sidebar__filters">
            <SelectInput
              placeholder="Sort by:"
              value={eventFilters.get('sortBy')}
              options={[
                { value: 'A - Z' },
                { value: 'Popular' },
                { value: 'Nearby' },
                { value: 'By time' },
              ]}
              handleSelectOption={handleSelectOption}
            />
            <Datepicker
              stateDate={selectedDate}
              selectedDate={selectedDate}
              setSelectedDate={handleDateChange}
            />
          </div>
        </>
      )}
      <Suspense fallback={<div>Loading...</div>}>
        <SidebarSearch
          currentTab={currentTab}
          setCurrentTab={setCurrentTab}
          display={isSearchingEvents}
          selectedDate={selectedDate}
          handleSearchSubmit={handleSearchSubmit}
          handleSelectSearchQueryOption={handleSelectSearchQueryOption}
          handleSubmitSearchQueryValue={handleSubmitSearchQueryValue}
          searchQueryFilters={searchQueryFilters}
          searchParametersStartDate={searchParametersStartDate}
          searchParametersEndDate={searchParametersEndDate}
        />
      </Suspense>
    </Container>
  )
}
SidebarBody.propTypes = {
  chartData: PropTypes.array,
  currentTab: PropTypes.string.isRequired,
  setCurrentTab: PropTypes.func.isRequired,
  eventFilters: PropTypes.object,
  isSearchingEvents: PropTypes.bool.isRequired,
  isShowAllEvents: PropTypes.bool.isRequired,
  handleDateChange: PropTypes.func.isRequired,
  handleSubmitValue: PropTypes.func.isRequired,
  handleSearchSubmit: PropTypes.func.isRequired,
  handleSelectOption: PropTypes.func.isRequired,
  handleSelectSearchQueryOption: PropTypes.func.isRequired,
  handleSubmitSearchQueryValue: PropTypes.func.isRequired,
  selectedDate: PropTypes.number,
  selectedHour: PropTypes.number,
  searchQueryFilters: PropTypes.object.isRequired,
  searchParametersStartDate: PropTypes.number,
  searchParametersEndDate: PropTypes.number,
}
export default SidebarBody
