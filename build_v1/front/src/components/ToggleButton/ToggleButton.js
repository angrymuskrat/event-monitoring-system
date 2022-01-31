import styled from 'styled-components'
import { lightGrey } from '../../config/styles'
export default styled.button`
  position: absolute;
  cursor: pointer;
  top: 50%;
  right: -3rem;
  z-index: 2;
  width: 3rem;
  height: 8rem;
  border-radius: 0px 10px 10px 0px;
  padding: 0;
  border: none;
  font: inherit;
  color: inherit;
  outline: 0;
  background-color: #ffffff;
  transition: all 0.3s;
  display: flex;
  align-items: center;
  justify-content: center;
  border-left: 1px solid ${lightGrey};
  & svg {
    width: 1.8rem;
    height: 1.8rem;
    padding-right: 2px;
  }
  &.toggle-button_closed {
    & svg {
      transform: scale(-1, 1);
    }
  }
  @media (max-width: 650px) {
    display: none;
    justify-content: center;
    align-items: center;
    top: 2.5%;
    right: -7.5rem;
    width: 5rem;
    height: 5rem;
    border-radius: 50%;
    box-shadow: 0 1px 4px rgba(0, 0, 0, 0.12), 0 2px 8px rgba(0, 0, 0, 0.06),
      0 0 16px rgba(0, 0, 0, 0.24);
    &.toggle-button_closed {
      display: flex;
    }
  }
`
