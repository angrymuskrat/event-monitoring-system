import styled from 'styled-components'
import { darkGrey, lightGrey, orange } from '../../config/styles'
export default styled.div`
  position: fixed;
  margin: 0 auto;
  width: 100%;
  max-width: 100%;
  z-index: 150;
  background-color: #ffffff;
  border-bottom: 1px solid ${lightGrey};

  .navigation {
    display: flex;
    justify-content: space-between;
    flex-direction: row;
    align-items: center;
    height: 7rem;
    width: 97%;
    margin: 0 auto;
    @media (max-width: 700px) {
      width: 93%;
    }
  }
  .navigation__container {
    flex-basis: 60%;
    display: flex;
    @media (max-width: 1180px) {
      flex-basis: 80%;
    }
    @media (max-width: 700px) {
      flex-basis: 100%;
    }
  }
  .navigation__input {
    width: 55%;
    & input {
      color: ${darkGrey};
      font-size: 1.5rem;
      font-family: 'Montserrat-SemiBold', sans-serif;
    }
    & svg {
      width: 2rem;
      height: 2rem;
    }
    @media (max-width: 1180px) {
      flex-basis: 80%;
    }
    @media (max-width: 1180px) {
      flex-basis: 100%;
    }
  }
  .navigation__button {
    cursor: pointer;
    & svg {
      transition: all 0.3s;
      width: 3rem;
      height: 3rem;
      & path {
        transition: all 0.3s;
        fill: ${darkGrey};
      }
      :hover {
        transform: translateX(-3px);
        & path {
          fill: ${orange};
        }
      }
    }
  }
  & .navigation__links-container {
    flex-basis: 40%;
    list-style-type: none;
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    width: 100%;
    max-width: 15rem;
    @media (max-width: 700px) {
      flex-basis: 0%;
      display: none;
    }
  }
  & .navigation__link {
    display: inline-block;
    &:last-of-type {
      margin-left: 3.5rem;
    }
  }
  & a {
    transition: all 0.3s;
  }
  & a:hover {
    color: ${orange};
  }
  .navigation__logo {
    margin-right: 4rem;

    & svg {
      width: 3.5rem;
      height: 3.5rem;
      & path {
        transition: all 0.3s;
      }
    }
    &:hover {
      & svg path {
        fill: ${darkGrey};
      }
    }
    @media (max-width: 1180px) {
      margin-right: 2rem;
    }
  }
`
