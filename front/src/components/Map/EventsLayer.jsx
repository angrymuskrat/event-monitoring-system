import React from "react";
import L from "leaflet";
import MarkerClusterGroup from "react-leaflet-markercluster";
import "react-leaflet-markercluster/dist/styles.min.css";
import MarkerContainer from "../Marker/Marker.jsx";

function EventsLayer(props) {
  const { events, onEventClick } = props;

  const createClusterCustomIcon = (cluster) => {
    const children = cluster.getAllChildMarkers();
    const firstChildHtml = children[0].options.icon.options.html;
    const coverPictureUrl = firstChildHtml.substring(
      firstChildHtml.indexOf('data-photoUrl="') + 15,
      firstChildHtml.indexOf('class="marker-inner"')
    );
    const firstChildTitle = firstChildHtml.substring(
      firstChildHtml.indexOf('<span class="text">') + 19,
      firstChildHtml.indexOf("</span>")
    );
    const count = cluster.getChildCount();

    // console.log(`coverPictureUrl`, coverPictureUrl);

    return L.divIcon({
      html: `<div class="marker-inner" 
                style="
                  background: url(${coverPictureUrl.replace(
                    '"',
                    ""
                  )}) no-repeat center center;
                  background-color: #fff;
                  background-size: cover;
                "
                >
                  <div class="marker__description" style="width:${30 +
                    firstChildTitle.length * 6}px">
                    <span class="text">${firstChildTitle} + ${count} events</span>
                  </div>
              </div>`,
      className: `marker`,
    });
  };

  const Markers = events
    ? events.map((p) => (
        <MarkerContainer
          key={`${p.properties.id}`}
          center={p.geometry.coordinates}
          onEventClick={onEventClick}
          {...p}
        />
      ))
    : null;

  return (
    <MarkerClusterGroup
      iconCreateFunction={createClusterCustomIcon}
      showCoverageOnHover={false}
      spiderfyOnMaxZoom={true}
      spiderLegPolylineOptions={{
        weight: 0,
        opacity: 0,
      }}
      removeOutsideVisibleBounds={true}
    >
      {Markers}
    </MarkerClusterGroup>
  );
}

export default EventsLayer;
