import styled from 'styled-components'

import { lightGrey } from '../../config/styles'

export default styled.ul`
  will-change: display;
  display: ${({ showSuggestions }) =>
    showSuggestions !== '' ? 'block' : 'hidden'};
  position: absolute;
  width: 100%;
  left: 0;
  top: 3.5rem;
  box-shadow: 0px 5px 10px #bfc2c8;
  border-top-width: 0;
  list-style: none;
  margin-top: 0;
  max-height: 14.3rem;
  overflow-y: auto;
  padding-left: 0;
  z-index: 1;
  .suggestion-container {
    border-bottom: 1px solid ${lightGrey};
    border-left: 1px solid ${lightGrey};
    border-right: 1px solid ${lightGrey};
    padding: 1rem 1rem 0.5rem 1rem;
    background-color: #fff;
    p {
      font-size: 1.5rem;
      font-family: 'Montserrat-SemiBold', sans-serif;
    }
  }
  .suggestion-container:last-of-type {
    border-radius: 0 0 5px 5px;
  }
  .suggestion-active,
  & .suggestion-container:hover {
    background-color: ${lightGrey};
    cursor: pointer;
    font-weight: 700;
  }
`
