import styled from 'styled-components'

export default styled.div`
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  align-items: left;
  width: 100%;
  max-width: 110rem;
  .page__section {
    padding: 0 3rem;
    margin-bottom: 10rem;
    margin-top: 11rem;
    display: flex;
    flex-direction: column;
    justify-content: flex-;
    @media (max-width: 800px) {
      margin-top: 4rem;
    }
    @media (max-width: 600px) {
      margin-top: 0rem;
      margin-bottom: 0rem;
    }
  }
`
