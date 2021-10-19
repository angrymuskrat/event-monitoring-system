import styled from 'styled-components'
import { darkGrey, lightGrey } from '../../config/styles'
export default styled.div`
  max-width: 520px;
  max-height: 50%;
  display: flex;
  justify-content: center;
  align-items: center;
  flex-direction: column;
  position: relative;
  overflow: hidden;
  @media (max-width: 650px) {
    height: 100%;
  }
  .popup__counter {
    width: 100%;
    z-index: 150;
    background-color: #ffffff;
    border-top: 1px solid ${lightGrey};
    display: flex;
    justify-content: center;
    box-shadow: 0px -6px 8px -1px rgba(0, 0, 0, 0.01);
    margin-top: -2rem;
    & .text {
      display: block;
      margin-top: 2rem;
      margin-bottom: 2rem;
    }
  }
  .popup__loading {
    min-height: 70rem;
    display: flex !important;
    justify-content: center;
    align-items: center;
    flex-direction: column;
    & p {
      display: block;
      color: ${darkGrey};
      font-weight: 700;
      font-size: 1.4rem;
    }
    .spinner-container {
      margin-top: 3rem;
    }
    .spinner {
      background-color: transparent;
      width: 1rem;
      height: 1rem;
      border-radius: 100%;
      box-shadow: 12px -12px 0 hsla(208, 16%, 35%, 0.125),
        17px 0 0 -1px hsla(208, 16%, 35%, 0.25),
        12px 12px 0 -2px hsla(208, 16%, 35%, 0.375),
        0 17px 0 -3px hsla(208, 16%, 35%, 0.5),
        -12px 12px 0 -4px hsla(208, 16%, 35%, 0.625),
        -17px 0 0 -5px hsla(208, 16%, 35%, 0.75),
        -12px -12px 0 -6px hsla(208, 16%, 35%, 0.875),
        0 -17px 0 -7px hsla(208, 16%, 35%, 1);
      animation: clockwise 0.75s steps(8, end) infinite;
    }
    @keyframes clockwise {
      to {
        transform: rotate(360deg) translatez(0);
      }
    }
  }
  .slick-slider {
    max-width: 52rem;
    width: 100%;
    scrollbar-width: none;
    @media (max-width: 650px) {
      max-width: 100vw;
      max-height: 75vh;
    }
  }
  .slick-slide {
    max-height: 70rem;
    overflow: scroll;
    display: flex;
    justify-content: center;
    overflow: scroll;
    scrollbar-color: white;
    scrollbar-width: none;
    @media (max-width: 650px) {
      height: 90%;
    }
  }
  .slick-list {
    overflow: hidden;
    max-width: 52rem;
    margin: 0 auto;
    @media (max-width: 650px) {
      max-width: 100%;
    }
  }
  .slick-track {
    display: flex;
    align-items: center;
  }
  .slider__button {
    position: absolute;
    top: 45%;
    display: inline-flex;
    color: rgb(34, 34, 34);
    background-color: rgba(255, 255, 255, 0.9);
    cursor: pointer;
    align-items: center;
    justify-content: center;
    background-clip: padding-box;
    box-shadow: transparent 0px 0px 0px 1px, transparent 0px 0px 0px 4px,
      rgba(0, 0, 0, 0.18) 0px 2px 4px;
    width: 3.2rem;
    height: 3.2rem;
    border-radius: 50%;
    border-image: initial;
    outline: 0;
    margin: 0;
    padding: 0;
    border-style: solid;
    border-width: 1px;
    border-color: rgba(0, 0, 0, 0.08);
    transition: box-shadow 0.2s ease 0s, -ms-transform 0.25s ease 0s,
      -webkit-transform 0.25s ease 0s, transform 0.25s ease 0s;
    &:active {
      background-color: rgb(255, 255, 255);
      color: rgb(0, 0, 0);
      box-shadow: none;
      transform: scale(1);
      border-color: rgba(0, 0, 0, 0.08);
    }
    &:hover {
      background-color: rgb(255, 255, 255);
      color: rgb(0, 0, 0);
      box-shadow: transparent 0px 0px 0px 1px, transparent 0px 0px 0px 4px,
        rgba(0, 0, 0, 0.12) 0px 6px 16px;
      transform: scale(1.04);
      border-color: rgba(0, 0, 0, 0.08);
    }
    & svg {
      width: 1rem;
      height: 1rem;
    }
  }
  .slider__button_next {
    right: 3px;
  }
  .slider__button_prev {
    left: 3px;
    & svg {
      transform: scale(-1, 1);
    }
  }
  .slick-prev {
    display: none !important;
  }
  .slick-next {
    display: none !important;
  }
`
