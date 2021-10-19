import React from 'react'
// styled
import Container from './Footer.styled'

function Footer() {
  return (
    <Container>
      <p className="text text_p2">
        Made with ❤ in{' '}
        <a
          className="text text__link"
          href="http://www.ifmo.ru/ru/"
          traget="_blank"
        >
          ITMO
        </a>
      </p>
    </Container>
  )
}
export default Footer
