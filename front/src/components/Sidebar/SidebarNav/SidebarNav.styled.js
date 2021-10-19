import styled from 'styled-components'
import {
  darkGrey,
  lightGreyTransparent,
  lightGrey,
  orange,
} from '../../../config/styles'
export default styled.nav`
  display: flex;
  flex-direction: column;
  width: 100%;
  @media (max-width: 650px) {
    padding-top: 0;
  }
  .sidebar__close-button {
    position: absolute;
    top: 8rem;
    cursor: pointer;
    display: none;
    align-self: flex-end;
    width: 4rem;
    height: 4rem;
    margin-bottom: 2rem;
    font-size: 1rem;
    line-height: 1rem;
    vertical-align: middle;
    & svg path {
      transition: fill 0.25s;
      fill: ${lightGreyTransparent};
    }
    &:hover {
      & svg path {
        fill: ${darkGrey};
      }
    }
    @media (max-width: 650px) {
      display: block;
    }
  }
  .sidebar__filters-toggle {
    display: flex;
    justify-content: flex-end;
    align-items: flex-end;
    height: 11rem;
    @media (max-width: 650px) {
      height: 15rem;
    }
  }
  /* Toggle button styles */
  .sidebar__toggle-active {
    position: relative;
    display: inline-block;
    width: 4rem;
    height: 2.2rem;
    margin-left: 2rem;
  }
  .sidebar__toggle-active input {
    opacity: 0;
    width: 0;
    height: 0;
  }
  .slider {
    position: absolute;
    cursor: pointer;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-color: ${lightGreyTransparent};
    box-shadow: 0px 2px 4px #bfc2c8;
    -webkit-transition: 0.4s;
    transition: 0.4s;
  }
  .slider:before {
    position: absolute;
    content: '';
    height: 1.5rem;
    width: 1.5rem;
    left: 0.3rem;
    bottom: 0.4rem;
    background-color: ${lightGrey};
    -webkit-transition: 0.4s;
    transition: 0.4s;
  }
  input:checked + .slider {
    background-color: ${orange};
  }
  input:focus + .slider {
    box-shadow: 0px 2px 4px #bfc2c8;
  }
  input:checked + .slider:before {
    -webkit-transform: translateX(1.9rem);
    -ms-transform: translateX(1.9rem);
    transform: translateX(1.9rem);
  }
  /* Rounded slider */
  .slider.round {
    border-radius: 3.4rem;
  }
  .slider.round:before {
    border-radius: 50%;
  }
  & .text_link {
    padding-bottom: 1rem;
    padding-left: 0.5rem;
  }
`
