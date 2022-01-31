import styled from 'styled-components'
import { lightGrey } from '../../config/styles'
export default styled.div`
  width: 100%;
  flex-basis: 49%;
  max-width: 54rem;
  min-width: 44rem;
  max-height: 35rem;
  display: flex;
  justify-content: space-between;
  flex-direction: row;
  flex-wrap: nowrap;
  background: #ffffff;
  border: 1px solid ${lightGrey};
  box-sizing: border-box;
  border-radius: 1rem;
  box-shadow: 0px 4px 4px rgba(0, 0, 0, 0.06);
  margin-bottom: 3rem;
  .team-card__image {
    flex-basis: 30%;
    width: 100%;
    max-width: 25.5rem;
    height: 100%;
    max-height: 35rem;
    @media (max-width: 940px) {
      max-width: 15.5rem;
    }
    @media (max-width: 940px) {
      max-width: 15.5rem;
    }
    @media (max-width: 400px) {
      max-width: 9.5rem;
    }
  }
  .team-card__content {
    padding: 2rem;
    @media (max-width: 940px) {
      padding: 1rem;
    }
  }
  @media (max-width: 940px) {
    min-width: 100%;
    margin-bottom: 2rem;
  }
`
