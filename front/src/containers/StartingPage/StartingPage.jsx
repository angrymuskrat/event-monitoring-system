import React, { memo } from 'react'
import PropTypes from 'prop-types'

// components
import CardContainer from '../../components/Card/CardContainer'
import Card from '../../components/Card/StartingPageCard.jsx'
import MarkerStartingPage from '../../components/Marker/MarkerStartingPage'
import Footer from '../../components/Footer/Footer.jsx'

// styled
import Container from './StartingPage.styled'

// config
import { cities } from '../../config/cities'

function StartingPage({ setCurrentCityAndViewport }) {
  const handleClick = city => {
    setCurrentCityAndViewport(city)
  }
  return (
    <>
      <Container>
        <main className="page__main">
          <div className="page__first-screen">
            <section className="page_start">
              <h1 className="title title_h1 title__starting-page">
                Explore data-driven events in your town
              </h1>
              <p className="text text_subheading">
                Watch how Instagram geotagged posts with the power of data
                science are in real time transformed into events that you can
                view, follow and visit.
              </p>
            </section>
            <div className="starting-page__image">
              <MarkerStartingPage>
                <div className="marker marker_starting-page marker_starting-page-1">
                  <div className="marker__animation"></div>
                </div>
              </MarkerStartingPage>
              <MarkerStartingPage>
                <div className="marker marker_starting-page marker_starting-page-2">
                  <div className="marker__animation"></div>
                </div>
              </MarkerStartingPage>
              <MarkerStartingPage>
                <div className="marker marker_starting-page marker_starting-page-3">
                  <div className="marker__animation"></div>
                </div>
              </MarkerStartingPage>
              <MarkerStartingPage>
                <div className="marker marker_starting-page marker_starting-page-4">
                  <div className="marker__animation"></div>
                </div>
              </MarkerStartingPage>
            </div>
          </div>
          <section className="page__section">
            <h3 className="title title_h3 sidebar__title">Discover events</h3>
            <CardContainer className="page__cards">
              {cities.map(city => {
                if (city.avaliable) {
                  return (
                    <div onClick={() => handleClick(city)} key={city.id}>
                      <Card {...city} />
                    </div>
                  )
                } else {
                  return (
                    <div key={city.id} to={`/map/`} className="link__disabled">
                      <Card {...city} />
                    </div>
                  )
                }
              })}
            </CardContainer>
          </section>
        </main>
        <Footer />
      </Container>
    </>
  )
}

StartingPage.propTypes = {
  setCurrentCityAndViewport: PropTypes.func.isRequired,
}
export default memo(StartingPage)
