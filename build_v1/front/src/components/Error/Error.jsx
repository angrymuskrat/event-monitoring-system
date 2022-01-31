import React from 'react'

// styled
import Container from './Error.styled'

// images
import ErrorIcon from '../../assets/svg/ErrorIcon'

export default function Error({ data }) {
  return (
    <Container>
      <div className="error__icon">
        <ErrorIcon />
      </div>
      <div>
        <p className="text text_p3">{data}</p>
      </div>
    </Container>
  )
}
