import styled from 'styled-components'
import { orange, badgeOrange, lightViolet, violet } from '../../config/styles'

export default styled.div`
  margin-top: 1rem;
  margin-bottom: 0.9rem;
  .input__tag {
    display: inline-block;
    background-color: ${({ type }) =>
      type === 'tags' ? `${badgeOrange}` : `${lightViolet}`};
    border-radius: 0.5rem;
    padding: 0.7rem 0.8rem 0.5rem 0.8rem;
    margin-right: 0.3rem;
    margin-top: 1rem;
    & p {
      display: inline-block;
      color: ${({ type }) => (type === 'tags' ? `${orange}` : `${violet}`)};
    }
  }
  .input__delete-button {
    display: inline-block;
    padding-left: 0.5rem;
    font-weight: 700;
    margin-block-end: 0;
    color: ${({ type }) => (type === 'tags' ? `${orange}` : `${violet}`)};
  }
`
