import React from 'react'

// styled
import Container from './Loading.styled'

function Loading() {
  return (
    <Container>
      <p className="text text_p1">Loading...</p>
      <div className="spinner-container">
        <div className="spinner"></div>
      </div>
    </Container>
  )
}
export default Loading
