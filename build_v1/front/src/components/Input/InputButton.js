import styled from 'styled-components'

import { darkGrey, grey } from '../../config/styles'

export default styled.button`
  cursor: pointer;
  & svg path {
    transition: fill 0.25s;
    fill: ${grey};
  }
  font-size: 1rem;
  line-height: 1rem;
  vertical-align: middle;
  &:hover {
    & svg path {
      fill: ${darkGrey};
    }
  }
`
