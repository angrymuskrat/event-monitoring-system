import styled from 'styled-components'
import { errorRed } from '../../config/styles'
export default styled.div`
  position: relative;
  top: 10%;
  left: 2%;
  background: ${errorRed};
  width: 96%;
  padding: 1rem 1.5rem;
  border-radius: 1rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
  p {
    margin-block-start: 0.5em;
    margin-block-end: 0.5rem;
  }
  .error__icon {
    margin-right: 1rem;
    & svg {
      width: 4rem;
      height: 4rem;
    }
  }
`
