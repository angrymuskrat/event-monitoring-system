import styled from 'styled-components'

export default styled.section`
  margin-top: 2rem;
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
  grid-gap: 2rem;
  align-items: center;
  .link__disabled {
    pointer-events: none;
  }
  @media (max-width: 800px) {
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  }
  @media (max-width: 490px) {
    grid-template-columns: repeat(auto-fill, minmax(1fr, 1fr));
  }
`
