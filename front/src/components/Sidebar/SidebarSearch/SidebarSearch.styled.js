import styled from 'styled-components'
import {
  lightGrey,
  badgeOrange,
  orange,
  errorRed,
} from '../../../config/styles'

export default styled.div`
  display: ${({ isShowSidebarSearch }) =>
    isShowSidebarSearch ? `block` : `none`};
  border-bottom: 1px solid ${lightGrey};
  padding-bottom: 2rem;
  & label {
    display: inline-block;
    width: 100%;
  }
  .button-container {
    display: flex;
    justify-content: center;
  }
  .datepicker__container {
    div[data-testid='DateRangeInputGrid'] {
      grid-template-columns: 45% 10% 45%;
    }
  }
  .sibebar-search__button {
    margin-top: 2rem;
    font-size: 1.2rem;
    background-color: ${orange};
    color: #fff;
    border-radius: 3rem;
    padding: 1rem;
    width: 70%;
    box-shadow: 0px 4px 3px ${lightGrey};
    cursor: pointer;
    transition-duration: 0.25s;
    &:hover {
      transform: translateY(-3px);
      box-shadow: 0 1px 4px rgba(0, 0, 0, 0.08), 0 2px 8px rgba(0, 0, 0, 0.06),
        0 0 16px rgba(0, 0, 0, 0.18);
    }
  }
  .sidebar-search__error {
    background: ${errorRed};
    padding: 0.5rem;
    height: 2.2rem;
    border-radius: 0.5rem;
    margin-bottom: 1rem;
    margin-top: 1.3rem;
    visibility: hidden;

    p {
      margin-block-end: 0rem;
    }
  }
  .sidebar-search__menu {
    display: flex;
    align-items: center;
    margin-top: 3rem;
  }
  .sidebar-search__menu-tab {
    flex-basis: 50%;
    min-width: 40%;
    display: flex;
    align-items: center;
    justify-content: center;
  }
  .sidebar-search__tab-button {
    cursor: pointer;
    border-bottom: 3px solid ${lightGrey};
    :hover {
      border-bottom: 3px solid ${badgeOrange};
    }
  }
  .sidebar-search__tab_filters {
    margin-top: 3rem;
    & > div {
      margin-right: 0;
      margin-top: 1rem;
      margin-bottom: 0.9rem;
    }
  }
`
