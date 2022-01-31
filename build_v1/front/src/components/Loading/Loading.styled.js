import styled from 'styled-components'
import { darkGrey, transparentBackground } from '../../config/styles'

export default styled.div`
  position: absolute;
  width: 70%;
  height: 50%;
  top: 50%;
  right: -20%;
  transform: translate(-50%, -50%);
  background: ${transparentBackground};
  border-radius: 1rem;
  z-index: 20;
  display: flex;
  align-items: center;
  flex-direction: column;
  justify-content: center;
  & .text {
    color: ${darkGrey};
    font-weight: 700;
    font-size: 1.4rem;
  }
  .spinner-container {
    margin-top: 3rem;
  }
  .spinner {
    background-color: transparent;
    width: 1rem;
    height: 1rem;
    border-radius: 100%;
    box-shadow: 12px -12px 0 hsla(208, 16%, 35%, 0.125),
      17px 0 0 -1px hsla(208, 16%, 35%, 0.25),
      12px 12px 0 -2px hsla(208, 16%, 35%, 0.375),
      0 17px 0 -3px hsla(208, 16%, 35%, 0.5),
      -12px 12px 0 -4px hsla(208, 16%, 35%, 0.625),
      -17px 0 0 -5px hsla(208, 16%, 35%, 0.75),
      -12px -12px 0 -6px hsla(208, 16%, 35%, 0.875),
      0 -17px 0 -7px hsla(208, 16%, 35%, 1);
    animation: clockwise 0.75s steps(8, end) infinite;
  }
  @keyframes clockwise {
    to {
      transform: rotate(360deg) translatez(0);
    }
  }
`
