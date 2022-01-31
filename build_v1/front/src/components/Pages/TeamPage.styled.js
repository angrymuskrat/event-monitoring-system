import styled from 'styled-components'
// styled
import CommonPage from './CommonPage.styled'

export default styled(CommonPage)`
  .team__list {
    display: flex;
    justify-content: space-between;
    flex-direction: row;
    flex-wrap: wrap;
    margin-top: 5rem;
    @media (max-width: 800px) {
      margin-top: 2rem;
    }
  }
`
