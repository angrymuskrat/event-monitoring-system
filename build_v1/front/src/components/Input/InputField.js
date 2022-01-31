import styled from 'styled-components'

import { darkGrey, lightGreyTransparent } from '../../config/styles'

export default styled.input`
  color: ${darkGrey};
  font-size: 1rem;
  line-height: 2rem;
  vertical-align: middle;
  &::placeholder {
    margin-right: -1rem;
  }
  &::-webkit-input-placeholder {
    color: ${lightGreyTransparent};
  }
`
