import React, { useEffect, useRef, useState, lazy, Suspense } from 'react'
import PropTypes from 'prop-types'
import { useParams } from 'react-router-dom'

// containers
import MapContainer from '../../components/Map/Map.jsx'
import Sidebar from '../Sidebar'

// components
import Loading from '../../components/Loading/Loading.jsx'

// styled
import Container from './MapPage.styled'

// lazy
const Popup = lazy(() => import('../../components/Popup/Popup.jsx'))

function MapPage(props) {
  const {
    setViewport,
    currentCityId,
    currentPopupData,
    currentUserLocation,
    isSidebarOpen,
    isPopupOpen,
    replace,
    selectedDate,
    setupMapFromHistory,
    setBounds,
    setCurrentUserLocation,
    setSelectedEvent,
    togglePopup,
    toggleSidebar,
    viewport,
  } = props
  // use refs to get mapbox bounds
  const mapRef = useRef()
  // react router hooks
  let { city, topLeft, botRight, time } = useParams()

  const [isMapStoreLoaded, setIsMapStoreLoaded] = useState(false)

  useEffect(() => {
    const config = { city, topLeft, botRight, time }
    if (currentCityId === null) {
      setupMapFromHistory(config)
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])
  useEffect(() => {
    if (viewport.toJS().center.length > 0) {
      setIsMapStoreLoaded(true)
    }
    if (mapRef.current && !currentUserLocation) {
      const pos = mapRef.current.leafletElement.locate()._renderer
      const location = pos && pos._center
      setCurrentUserLocation(location)
    }
  }, [currentUserLocation, setCurrentUserLocation, viewport])

  const handleViewportChangeEnd = viewport => {
    let bottomRight = [
      Number(
        mapRef.current.leafletElement.getBounds()._southWest.lat.toFixed(3)
      ),
      Number(
        mapRef.current.leafletElement.getBounds()._southWest.lng.toFixed(3)
      ),
    ]
    let topLeft = [
      Number(
        mapRef.current.leafletElement.getBounds()._northEast.lat.toFixed(3)
      ),
      Number(
        mapRef.current.leafletElement.getBounds()._northEast.lng.toFixed(3)
      ),
    ]
    setViewport(viewport)
    setBounds({ topLeft, bottomRight })
    replace(`/map/${currentCityId}/${topLeft}/${bottomRight}/${selectedDate}`)
  }
  const onEventClick = (id, postcodes) => {
    if (window.innerWidth < 600) {
      isSidebarOpen && toggleSidebar()
    }
    setSelectedEvent(id)
    togglePopup(postcodes)
  }
  const closePopup = () => {
    togglePopup()
  }
  return (
    <Container>
      {isMapStoreLoaded ? (
        <>
          <Sidebar
            onCardClick={onEventClick}
            viewport={viewport}
            mapRef={mapRef}
          />
          <MapContainer
            {...props}
            mapRef={mapRef}
            handleViewportChangeEnd={handleViewportChangeEnd}
            onEventClick={onEventClick}
          />
        </>
      ) : (
        <Loading />
      )}
      {isPopupOpen ? (
        <Suspense fallback={<Loading />}>
          <Popup
            isPopupOpen={isPopupOpen}
            closePopup={closePopup}
            data={currentPopupData}
          />
        </Suspense>
      ) : null}
    </Container>
  )
}
MapPage.propTypes = {
  allHeatmapData: PropTypes.object,
  currentCityId: PropTypes.any,
  currentHeatmap: PropTypes.object,
  currentEvents: PropTypes.object,
  currentPopupData: PropTypes.array,
  events: PropTypes.object,
  isSidebarOpen: PropTypes.bool.isRequired,
  isShowAllEvents: PropTypes.bool.isRequired,
  isPopupOpen: PropTypes.bool.isRequired,
  isSearchingEvents: PropTypes.bool.isRequired,
  replace: PropTypes.func.isRequired,
  searchQuery: PropTypes.object,
  selectedDate: PropTypes.any,
  selectedHour: PropTypes.any,
  selectedEvent: PropTypes.string,
  setupMapFromHistory: PropTypes.func.isRequired,
  setBounds: PropTypes.func.isRequired,
  setSelectedEvent: PropTypes.func.isRequired,
  togglePopup: PropTypes.func.isRequired,
  toggleSidebar: PropTypes.func.isRequired,
  viewport: PropTypes.object,
}

export default MapPage
