import styled from 'styled-components'
// styles
import { darkGrey } from '../../config/styles'
// styled
import CommonPage from './CommonPage.styled'

export default styled(CommonPage)`
  height: calc(100vh - 17.5rem);
  .about-page__section {
    display: flex;
    flex-direction: row;
    flex-wrap: wrap;
    justify-content: space-between;
    margin-top: 6rem;
    :first-of-type {
      margin-top: 4rem;
    }
  }
  .about-page__image {
    flex-basis: 25%;
    width: 100%;
    background-color: ${darkGrey};
  }
  .about-page__content {
    flex-basis: 70%;
  }
`
