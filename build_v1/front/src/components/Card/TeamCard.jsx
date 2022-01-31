import React from 'react'
import PropTypes from 'prop-types'

//styled
import Container from './TeamCard.styled'

function TeamCard({ person: { photo, name, title, description } }) {
  return (
    <Container data-test="team-card-container">
      <img className="team-card__image" src={photo} alt={name} />
      <div className="team-card__content" data-test="team-card__content">
        <h2 className="title title_h2" data-test="team-card__name">
          {name}
        </h2>
        <h4 className="title title_h4" data-test="team-card__title">
          {title}
        </h4>
        <p className="text text_p2" data-test="team-card__descr">
          {description}
        </p>
      </div>
    </Container>
  )
}
TeamCard.propTypes = {
  person: PropTypes.object.isRequired,
}
export default TeamCard
