import React from 'react'
import PropTypes from 'prop-types'

// components
import SidebarCard from './SidebarCard.jsx'

// styled
import Container from './SidebarEventsList.styled'

// utils
import { sortEvents } from '../../../utils/utils'

function SidebarEventsList({
  chartData,
  currentEvents,
  eventFilters,
  events,
  isSearchingEvents,
  isSidebarOpen,
  isShowAllEvents,
  isLoading,
  searchQuery,
  searchQueryFilters,
  searchParameters,
  handleEventClick,
  handleEventHover,
  onCardClick,
  sidebarHeight,
  viewport,
}) {
  const renderEventsList = () => {
    if (isLoading) {
      return <p className="text text_p2 text__events">Loading...</p>
    }
    if (isSearchingEvents) {
      if (searchQuery && !isLoading) {
        searchQuery = searchQuery.toJS()
        if (
          searchQueryFilters.get('sortBy') ||
          searchQueryFilters.get('keyword')
        ) {
          searchQuery = sortEvents(
            searchQueryFilters.toJS(),
            searchQuery,
            viewport.toJS()
          )
        }

        return searchQuery.map((event, i) => (
          <SidebarCard
            key={i}
            event={event}
            handleEventHover={handleEventHover}
            handleEventClick={handleEventClick}
            handlePostsClick={onCardClick}
          />
        ))
      } else if (searchParameters) {
        return (
          <p className="text text_p2 text__events">
            No events found, try another search parameters
          </p>
        )
      } else {
        return (
          <p className="text text_p2 text__events">Enter search parameters</p>
        )
      }
    }
    if (isShowAllEvents) {
      let acc = []
      if (chartData && events) {
        chartData.forEach(data => {
          if (events.get(`${data.time}`) !== undefined) {
            events
              .get(`${data.time}`)
              .toJS()
              .forEach(d => acc.push(d))
          }
        })
      }
      if (eventFilters.get('sortBy') || eventFilters.get('keyword')) {
        acc = sortEvents(eventFilters.toJS(), acc, viewport.toJS())
        return acc.map((event, i) => (
          <SidebarCard
            key={i}
            event={event}
            handleEventHover={handleEventHover}
            handleEventClick={handleEventClick}
            handlePostsClick={onCardClick}
          />
        ))
      }
      return acc
        .sort((a, b) => {
          if (a.properties.start < b.properties.start) {
            return -1
          }
          if (a.properties.start < b.properties.start) {
            return 1
          }
          return 0
        })
        .map((event, i) => (
          <SidebarCard
            key={i}
            event={event}
            handleEventHover={handleEventHover}
            handleEventClick={handleEventClick}
            handlePostsClick={onCardClick}
          />
        ))
    }
    if (currentEvents && !isShowAllEvents && !isSearchingEvents) {
      currentEvents = currentEvents.toJS()
      if (eventFilters.get('sortBy') || eventFilters.get('keyword')) {
        currentEvents = sortEvents(
          eventFilters.toJS(),
          currentEvents,
          viewport.toJS()
        )
        if (currentEvents.length === 0) {
          return <p>No events found, try another search parameters</p>
        }
        return currentEvents.map((event, i) => (
          <SidebarCard
            key={i}
            event={event}
            handleEventHover={handleEventHover}
            handleEventClick={handleEventClick}
            handlePostsClick={onCardClick}
          />
        ))
      } else {
        return currentEvents
          .sort((a, b) => {
            if (a.properties.start < b.properties.start) {
              return -1
            }
            if (a.properties.start < b.properties.start) {
              return 1
            }
            return 0
          })
          .map((event, i) => (
            <SidebarCard
              key={i}
              event={event}
              handleEventHover={handleEventHover}
              handleEventClick={handleEventClick}
              handlePostsClick={onCardClick}
            />
          ))
      }
    } else {
      return <p className="text text_p2 text__events">No events found</p>
    }
  }
  return (
    <Container
      style={{
        top: `${sidebarHeight - 60}px`,
        minHeight: `${window.innerHeight - sidebarHeight}px`,
        maxWidth: `${window.innerWidth}px`,
      }}
    >
      {renderEventsList()}
    </Container>
  )
}
SidebarEventsList.propTypes = {
  chartData: PropTypes.array,
  currentevents: PropTypes.array,
  currentUserLocation: PropTypes.object,
  eventFilters: PropTypes.object,
  events: PropTypes.object,
  isSearchingEvents: PropTypes.bool.isRequired,
  isShowAllEvents: PropTypes.bool.isRequired,
  isLoading: PropTypes.bool.isRequired,
  searchQuery: PropTypes.object,
  searchQueryFilters: PropTypes.object,
  searchParameters: PropTypes.object,
  handleEventClick: PropTypes.func.isRequired,
  handleEventHover: PropTypes.func.isRequired,
  onCardClick: PropTypes.func.isRequired,
  sidebarHeight: PropTypes.number,
  viewport: PropTypes.object,
}
export default SidebarEventsList
