import React from 'react'

// components
import Footer from '../Footer/Footer.jsx'

// styled
import Container from './AboutPage.styled'

function AboutPage() {
  return (
    <>
      <Container>
        <main className="page__section">
          <h1 className="title title_h1">How this works</h1>
          <section className="about-page__section">
            <div className="about-page__content">
              <h4 className="title title_h4">1. Select a city.</h4>
              <h4 className="title title_h4">
                2. Use the timeline to select a particular hour.
              </h4>
              <h4 className="title title_h4">
                3. Click on a pop-up to explore events!
              </h4>
              <h4 className="title title_h4">
                4. You can use a calendar and timeline to switch date.
              </h4>
              <h4 className="title title_h4">
                5. Events can be sorted by name, popularity, and your location.
              </h4>
              <h4 className="title title_h4">
                6. To search for a profile, use Advanced search and type
                username.
              </h4>
              <h4 className="title title_h4">
                7. To search for a hashtag, type the hashtag in the form.
              </h4>
            </div>
          </section>
        </main>
      </Container>
      <Footer />
    </>
  )
}
export default AboutPage
