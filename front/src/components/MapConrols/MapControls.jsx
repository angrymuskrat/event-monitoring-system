import React from 'react'
import PropTypes from 'prop-types'
import { GeolocateControl, NavigationControl } from 'react-map-gl'

// styles
import { geolocateControlStyle } from '../../config/styles'

// styled
import Container from './MapControls.styled'

function MapControls({ onGeolocate }) {
  return (
    <Container>
      <div
        className="control-group-container"
        style={{ position: 'absolute', right: '3rem', top: '12%' }}
      >
        <NavigationControl showCompass={false} />
      </div>
      <GeolocateControl
        showUserLocation={true}
        style={geolocateControlStyle}
        fitBoundsOptions={{ maxZoom: 30 }}
        onGeolocate={e => onGeolocate(e)}
        className="geolocate-control-container"
      />
    </Container>
  )
}
MapControls.propTypes = {
  onGeolocate: PropTypes.func.isRequired,
}
export default MapControls
