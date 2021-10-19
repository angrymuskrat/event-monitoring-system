import styled from 'styled-components'
import { lightGrey } from '../../../config/styles'

export default styled.div`
  .sidebar__filters {
    display: flex;
    padding-bottom: 2rem;
    border-bottom: 1px solid ${lightGrey};
    padding-top: 2rem;
    & > div {
      flex-basis: 50%;
    }
  }
  .sidebar__button-container {
    display: flex;
    justify-content: space-between;
    align-items: center;
    flex-direction: row;
    width: 100%;
    padding-bottom: 1.5rem;
  }
`
