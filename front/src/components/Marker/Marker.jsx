import React, { useEffect, useState } from "react";
import { divIcon } from "leaflet";
import { Marker } from "react-leaflet";

import noPhotoMarker from "../../assets/img/noPhotoMarker.png";

import { makeInstagramImageUrl } from "../../utils/utils";

function MarkerContainer(props) {
  const { properties, center, onEventClick } = props;
  const [isPhotoValid, setIsPhotoValid] = useState(noPhotoMarker);

  useEffect(() => {
    const imagesUrl = properties.postcodes.map((i) => makeInstagramImageUrl(i));

    const getImage = (index) => {
      if (index === imagesUrl.length) return;
      const image = new Image();
      image.onload = () => {
        setIsPhotoValid(image.src);
      };
      image.onerror = () => {
        getImage(index + 1);
      };
      image.src = imagesUrl[index];
    };

    getImage(0);
  }, [properties.postcodes]);

  const handleClick = () => {
    onEventClick(properties.id, properties.postcodes);
  };
  const createIcon = () => {
    const iconMarkup = `
      <div data-photoUrl="${properties.photoUrl}"
          class="marker-inner" 
          style="
            background: url(${isPhotoValid}) no-repeat center center;
            background-color: #fff;
            background-size: cover;
          "
          >
          <div class="marker__description" style="width:${25 +
            properties.title.length * 6}px">
            <span class="text">${properties.title}</span>
          </div>
      </div>
    `;
    return divIcon({
      className: "marker",
      html: iconMarkup,
    });
  };
  return (
    <Marker
      className="circle"
      icon={createIcon()}
      position={center}
      onClick={handleClick}
    />
  );
}

export default MarkerContainer;
