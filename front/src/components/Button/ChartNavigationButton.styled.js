import styled from 'styled-components'
import { darkGrey, transparentBackground } from '../../config/styles'

export default styled.button`
  display: block;
  border-radius: 2rem;
  letter-spacing: 0.5px;
  padding: 0.4rem 0.8rem;
  background: ${transparentBackground};
  cursor: pointer;
  transition-duration: 0.25s;
  &:hover {
    transform: translateY(-1.5px);
    box-shadow: 0 1px 1px rgba(0, 0, 0, 0.02), 0 2px 2px rgba(0, 0, 0, 0.0015),
      0 0 4px rgba(0, 0, 0, 0.045);
  }
  .text {
    color: ${darkGrey};
  }
`
