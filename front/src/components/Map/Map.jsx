import React, { useEffect, useState } from 'react'
import PropTypes from 'prop-types'
// leaflet
import { Map, TileLayer } from 'react-leaflet'
import 'leaflet/dist/leaflet.css'
import HeatmapLayer from 'react-leaflet-heatmap-layer'

// styles
import { heatmapGradient } from '../../config/styles'

// components
import EventsLayer from './EventsLayer'
import Loading from '../Loading/Loading.jsx'

function MapContainer(props) {
  const {
    isShowAllEvents,
    isSearchingEvents,
    isLoading,
    chartData,
    allHeatmapData,
    currentHeatmap,
    currentEvents,
    events,
    mapRef,
    searchQuery,
    selectedHour,
    handleViewportChangeEnd,
    viewport,
    onEventClick,
  } = props
  const [heatmapData, setHeatmapData] = useState(null)
  const [eventsData, setEventsData] = useState(null)
  const [innerViewport, setInnerViewport] = useState(null)
  const [currentZoom, setCurrentZoom] = useState()

  useEffect(() => {
    setInnerViewport(viewport.toJS())
  }, [viewport])
  const handleViewportChanged = viewport => {
    if (mapRef.current) {
      handleViewportChangeEnd(viewport)
    }
  }
  const handleViewportChange = e => {
    setCurrentZoom(e.zoom)
  }
  useEffect(() => {
    setHeatmapData(calculateHeatmap())
    setEventsData(calculateEventsData())
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [
    currentHeatmap,
    currentEvents,
    events,
    searchQuery,
    allHeatmapData,
    isShowAllEvents,
    isSearchingEvents,
    selectedHour,
  ])

  useEffect(() => {
    setHeatmapData(calculateHeatmap())
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [currentZoom])

  const calculateHeatmap = () => {
    if (
      mapRef.current &&
      (currentHeatmap || (isShowAllEvents && allHeatmapData)) &&
      selectedHour &&
      !isSearchingEvents
    ) {
      let acc = []
      if (isShowAllEvents && chartData && allHeatmapData) {
        chartData.forEach(data => {
          if (allHeatmapData.get(`${data.time}`) !== undefined) {
            allHeatmapData
              .get(`${data.time}`)
              .toJS()
              .forEach(d => acc.push(d))
          }
        })
      }
      return isShowAllEvents
        ? acc
        : currentHeatmap
        ? currentHeatmap.toJS()
        : null
    } else {
      return null
    }
  }
  const calculateEventsData = () => {
    if (isSearchingEvents && searchQuery) {
      return searchQuery.toJS()
    }
    if (
      mapRef.current &&
      (currentEvents || (isShowAllEvents && events)) &&
      selectedHour &&
      !isSearchingEvents
    ) {
      let acc = []
      if (isShowAllEvents && chartData && events) {
        chartData.forEach(data => {
          if (events.get(`${data.time}`) !== undefined) {
            events
              .get(`${data.time}`)
              .toJS()
              .forEach(d => acc.push(d))
          }
        })
      }
      return isShowAllEvents ? acc : currentEvents ? currentEvents.toJS() : null
    } else {
      return null
    }
  }
  return (
    <>
      {innerViewport ? (
        <Map
          id="mapId"
          ref={mapRef}
          animate={true}
          viewport={innerViewport}
          onViewportChanged={handleViewportChanged}
          onViewportChange={handleViewportChange}
        >
          <TileLayer
            attribution='&amp;copy <a href="http://osm.org/copyright">OpenStreetMap</a> contributors'
            url="https://{s}.basemaps.cartocdn.com/rastertiles/light_all/{z}/{x}/{y}.png"
          />
          {isLoading ? (
            <Loading />
          ) : (
            <>
              {!isSearchingEvents && (
                <HeatmapLayer
                  points={heatmapData}
                  longitudeExtractor={m => m[1]}
                  latitudeExtractor={m => m[0]}
                  gradient={heatmapGradient}
                  minOpacity={0.25}
                  intensityExtractor={m => parseFloat(m[2] / 3)}
                  radius={30}
                  blur={10}
                  max={3}
                />
              )}

              <EventsLayer events={eventsData} onEventClick={onEventClick} />
            </>
          )}
        </Map>
      ) : (
        <Loading />
      )}
    </>
  )
}
MapContainer.propTypes = {
  isShowAllEvents: PropTypes.bool.isRequired,
  isSearchingEvents: PropTypes.bool.isRequired,
  isLoading: PropTypes.bool.isRequired,
  chartData: PropTypes.array,
  allHeatmapData: PropTypes.object,
  currentHeatmap: PropTypes.object,
  currentEvents: PropTypes.object,
  events: PropTypes.object,
  mapRef: PropTypes.object,
  searchQuery: PropTypes.object,
  selectedHour: PropTypes.string.isRequired,
  handleViewportChangeEnd: PropTypes.func.isRequired,
  viewport: PropTypes.object,
  onEventClick: PropTypes.func.isRequired,
}
export default MapContainer
