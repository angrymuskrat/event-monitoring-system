import React, { useEffect, useState } from "react";
import PropTypes from "prop-types";
import moment from "moment";

// styled
import Container from "./SidebarCard.styled";

// styles
import { lightGrey } from "../../../config/styles";
import noPhoto from "../../../assets/img/noPhotoEvent.png";
import { makeInstagramImageUrl } from "../../../utils/utils";

function SidebarCard({
  event,
  handleEventHover,
  handleEventClick,
  handlePostsClick,
}) {
  const {
    title,
    tags,
    start,
    finish,
    photoUrl,
    id,
    postcodes,
  } = event.properties;

  const [isPhotoValid, setIsPhotoValid] = useState(noPhoto);

  const handleMouseEnter = () => {
    handleEventHover(id);
  };
  const handleMouseLeave = () => {
    handleEventHover(null);
  };
  const handleClick = () => {
    handleEventClick(id);
  };
  const handleLinkClick = () => {
    handlePostsClick(id, postcodes);
  };

  useEffect(() => {
    const imagesUrl = postcodes.map((i) => makeInstagramImageUrl(i));

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
  }, [postcodes]);

  return (
    <Container
      key={title}
      onMouseEnter={handleMouseEnter}
      onMouseLeave={handleMouseLeave}
    >
      <div
        className="event-card__image"
        style={{
          background: `url(${isPhotoValid}) ${lightGrey} no-repeat center`,
          backgroundSize: "cover",
        }}
        onClick={handleClick}
      />
      <div className="event-card__content">
        <h4 className="title title_h4 event-card__title" onClick={handleClick}>
          {title}
        </h4>
        <p>Start: {moment.unix(start).format("HH:mm, DD.MM.YYYY")}</p>
        <p>Finish: {moment.unix(finish).format("HH:mm, DD.MM.YYYY")}</p>
        {tags.map((tag, i) => (
          <div className="event-card__tag" key={`${tag + i}`}>
            {tag}
          </div>
        ))}
        <br />
        <button className="event-card__button" onClick={handleLinkClick}>
          View all posts
        </button>
      </div>
    </Container>
  );
}
SidebarCard.propTypes = {
  event: PropTypes.object,
  handleEventClick: PropTypes.func.isRequired,
  handleEventHover: PropTypes.func.isRequired,
  handlePostsClick: PropTypes.func.isRequired,
};
export default SidebarCard;
