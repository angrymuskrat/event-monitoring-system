import React from 'react'

// components
import TeamCard from '../Card/TeamCard.jsx'
import Footer from '../Footer/Footer.jsx'

//styled
import Container from './TeamPage.styled'

// data
import { team } from '../../config/team'

function TeamPage() {
  return (
    <>
      <Container>
        <section className="page__section">
          <h1 className="title title_h1">Our Team</h1>
          <div className="team__list">
            {team.map((t, i) => (
              <TeamCard key={i} person={t} />
            ))}
          </div>
        </section>
      </Container>
      <Footer />
    </>
  )
}
export default TeamPage
