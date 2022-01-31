import styled from 'styled-components'
export default styled.div`
  .arrow {
    display: block;
  }
  .menu {
    display: none;
  }
  @media (max-width: 650px) {
    .arrow {
      display: none !important;
    }
    .menu {
      display: block;
    }
  }
`
