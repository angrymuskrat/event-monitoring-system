import { badgeOrange, orange } from '../../config/styles'

import styled from 'styled-components'

export default styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  & .chartstat__nav {
    display: flex;
    justify-content: space-between;
  }
  & .chartstat__label {
    display: inline-block;
    padding: 4px;
    border-radius: 0.5rem;
    color: ${orange};
    background-color: ${badgeOrange};
  }
`
