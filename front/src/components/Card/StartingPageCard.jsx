import React from 'react'
import PropTypes from 'prop-types'
import Card from './StartingPageCard.styled'

function StartingPageCard({ avaliable, city, country, picture }) {
  return (
    <Card avaliable={avaliable} picture={picture}>
      <div className="card__text">
        <h4 className="title title_h4">{city}</h4>
        <p className="text text_subtitle">{country}</p>
        <p className="text text_p1-light">{!avaliable && 'coming soon'}</p>
      </div>
    </Card>
  )
}
StartingPageCard.propTypes = {
  avaliable: PropTypes.bool,
  city: PropTypes.string,
  country: PropTypes.string,
  picture: PropTypes.string,
}
export default StartingPageCard
