import styled from 'styled-components'
import { darkGrey, lightGrey } from '../../config/styles'

export default styled.div`
  position: relative;
  height: 3.5rem;
  flex-basis: 100%;
  flex-grow: 1;
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  background: ${({ type }) =>
    type === 'starting-page' || type === 'sidebar' ? '#fff' : `${lightGrey}`};
  border: 1px solid ${lightGrey};
  box-shadow: 0px 2px 4px #bfc2c8;
  color: ${darkGrey};
  border-radius: ${({ showSuggestions }) =>
    showSuggestions ? '5px 5px 0 0' : `5px`};
  padding: 0.5rem 1rem 0.5rem;
  & input {
    width: 100%;
  }
  input::placeholder {
    padding-left: 0.5rem;
  }
`
