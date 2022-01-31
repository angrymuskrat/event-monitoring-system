import React, { Suspense, lazy, useState, useEffect, useRef } from 'react'
import PropTypes from 'prop-types'
import moment from 'moment'

// components
import SidebarNav from '../../components/Sidebar/SidebarNav/SidebarNav.jsx'
import SidebarBody from '../../components/Sidebar/SidebarBody/SidebarBody.jsx'

// images
import { ToggleMenu } from '../../assets/svg/ToggleMenu'

// styled
import SidebarContainer from './SidebarContainer.styled'
import Container from './Sidebar.styled'
import ToggleButton from '../../components/ToggleButton/ToggleButton'

// utils
import { useWindowWidth } from '../../utils/hooks'

// lazy
const EventsList = lazy(() =>
  import('../../components/Sidebar/SidebarEventsList/SidebarEventsList.jsx')
)

function Sidebar({
  allEventsData,
  chartData,
  currentEvents,
  currentUserLocation,
  events,
  isLoading,
  isPopupOpen,
  isSearchingEvents,
  isSidebarOpen,
  isShowAllEvents,
  toggleSidebar,
  searchQuery,
  searchQueryFilters,
  searchParameters,
  searchParametersStartDate,
  searchParametersEndDate,
  selectedHour,
  selectedDate,
  setNewDate,
  setHighlightedEvent,
  setSelectedEvent,
  setSearchParameters,
  setSearchEvents,
  onCardClick,
  setEventFilter,
  setSearchQueryFilter,
  eventFilters,
  viewport,
  mapRef,
}) {
  const [currentTab, setCurrentTab] = useState('Search')
  const fixedSidebarRef = useRef()
  const [sidebarHeight, setSidebarHeight] = useState()
  const width = useWindowWidth()
  useEffect(() => {
    if (fixedSidebarRef.current) {
      setSidebarHeight(fixedSidebarRef.current.offsetHeight)
    }
  }, [chartData, currentTab, isSearchingEvents, searchParameters, width])

  const handleDateChange = data => {
    if (data.startDate) {
      setNewDate(moment(data.startDate).unix())
    } else {
      setNewDate(data)
    }
  }
  const handleCheckboxChange = () => {
    setSearchEvents()
  }
  const handleEventHover = id => {
    setHighlightedEvent(id)
  }
  const handleEventClick = id => {
    let event
    if (isShowAllEvents) {
      event = allEventsData.find(obj => {
        return obj.getIn(['properties', 'id']) === id
      })
    } else if (isSearchingEvents) {
      event = searchQuery.find(obj => {
        return obj.getIn(['properties', 'id']) === id
      })
    } else {
      event = currentEvents.find(obj => {
        return obj.getIn(['properties', 'id']) === id
      })
    }
    const eventCoordinates = event.getIn(['geometry', 'coordinates']).toJS()
    mapRef.current.leafletElement.flyTo(
      [eventCoordinates[0], eventCoordinates[1]],
      16
    )
    if (window.innerWidth < 600) {
      toggleSidebar()
    }
    setSelectedEvent(id)
  }
  const handleSelectOption = option => {
    const parameter = {
      type: 'sort',
      parameter: option,
    }
    setEventFilter(parameter)
  }
  const handleSubmitValue = value => {
    const parameter = {
      type: 'keyword',
      parameter: value,
    }
    setEventFilter(parameter)
  }
  const handleSelectSearchQueryOption = option => {
    const parameter = {
      type: 'sort',
      parameter: option,
    }
    setSearchQueryFilter(parameter)
  }
  const handleSubmitSearchQueryValue = value => {
    const parameter = {
      type: 'keyword',
      parameter: value,
    }
    setSearchQueryFilter(parameter)
  }
  const handleSearchSubmit = params => {
    setSearchParameters(params)
  }
  return (
    <SidebarContainer
      className={
        isSidebarOpen ? 'sidebar-container_active' : 'sidebar-container_closed'
      }
    >
      <Container>
        {window.innerWidth < 600 && isPopupOpen ? null : (
          <ToggleButton
            className={
              isSidebarOpen ? 'toggle-button_active' : 'toggle-button_closed'
            }
            in={isSidebarOpen.toString()}
            onClick={() => toggleSidebar()}
          >
            <ToggleMenu />
          </ToggleButton>
        )}
        <div
          className={
            isSidebarOpen
              ? 'sidebar-fixed sidebar-fixed_active'
              : 'sidebar-fixed sidebar-fixed_closed'
          }
          ref={fixedSidebarRef}
        >
          <SidebarNav
            toggleSidebar={toggleSidebar}
            handleCheckboxChange={handleCheckboxChange}
          />
          <SidebarBody
            chartData={chartData}
            currentTab={currentTab}
            setCurrentTab={setCurrentTab}
            eventFilters={eventFilters}
            isSearchingEvents={isSearchingEvents}
            isShowAllEvents={isShowAllEvents}
            handleDateChange={handleDateChange}
            handleSearchSubmit={handleSearchSubmit}
            handleSubmitValue={handleSubmitValue}
            handleSelectOption={handleSelectOption}
            handleSelectSearchQueryOption={handleSelectSearchQueryOption}
            handleSubmitSearchQueryValue={handleSubmitSearchQueryValue}
            selectedDate={Number(selectedDate)}
            selectedHour={Number(selectedHour)}
            searchQueryFilters={searchQueryFilters}
            searchParametersStartDate={searchParametersStartDate}
            searchParametersEndDate={searchParametersEndDate}
          />
        </div>

        {fixedSidebarRef.current && (
          <Suspense fallback={<div>Loading...</div>}>
            <EventsList
              chartData={chartData}
              currentEvents={currentEvents}
              currentUserLocation={currentUserLocation}
              events={events}
              eventFilters={eventFilters}
              isSearchingEvents={isSearchingEvents}
              isSidebarOpen={isSidebarOpen}
              isShowAllEvents={isShowAllEvents}
              isLoading={isLoading}
              searchQuery={searchQuery}
              searchQueryFilters={searchQueryFilters}
              searchParameters={searchParameters}
              handleEventClick={handleEventClick}
              handleEventHover={handleEventHover}
              onCardClick={onCardClick}
              mapRef={mapRef}
              sidebarHeight={sidebarHeight}
              viewport={viewport}
            />
          </Suspense>
        )}
      </Container>
    </SidebarContainer>
  )
}

Sidebar.propTypes = {
  allEventsData: PropTypes.object,
  chartData: PropTypes.array,
  currentEvents: PropTypes.object,
  currentUserLocation: PropTypes.object,
  events: PropTypes.object,
  isLoading: PropTypes.bool.isRequired,
  isPopupOpen: PropTypes.bool.isRequired,
  isSearchingEvents: PropTypes.bool.isRequired,
  isSidebarOpen: PropTypes.bool.isRequired,
  isShowAllEvents: PropTypes.bool.isRequired,
  toggleSidebar: PropTypes.func.isRequired,
  searchQuery: PropTypes.object,
  searchQueryFilters: PropTypes.object,
  searchParameters: PropTypes.object,
  setSearchParameters: PropTypes.func.isRequired,
  searchParametersStartDate: PropTypes.number,
  searchParametersEndDate: PropTypes.number,
  selectedHour: PropTypes.string,
  selectedDate: PropTypes.number,
  setNewDate: PropTypes.func.isRequired,
  setHighlightedEvent: PropTypes.func.isRequired,
  setSelectedEvent: PropTypes.func.isRequired,
  setSearchEvents: PropTypes.func.isRequired,
  onCardClick: PropTypes.func.isRequired,
  setEventFilter: PropTypes.func.isRequired,
  setSearchQueryFilter: PropTypes.func.isRequired,
  eventFilters: PropTypes.object,
  viewport: PropTypes.object,
  mapRef: PropTypes.object,
}

export default Sidebar
