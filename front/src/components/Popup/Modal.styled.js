import styled from 'styled-components'
import { darkGrey, lightGreyTransparent } from '../../config/styles'
export default styled.div`
  & .modal__button {
    position: absolute;
    top: 1.2rem;
    right: 1rem;
    width: 2rem;
    height: 2rem;
    cursor: pointer;
    & svg path {
      transition: fill 0.25s;
      fill: ${lightGreyTransparent};
    }
    &:hover {
      & svg path {
        fill: ${darkGrey};
      }
    }
  }
`
