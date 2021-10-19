import styled from 'styled-components'
import { darkGrey, orange } from '../../config/styles'
export default styled.div`
  @media (min-width: 1175px) {
    display: none;
  }
  .rangeslider-horizontal {
    height: 1rem;
    outline: 0;
  }
  .rangeslider,
  .rangeslider .rangeslider__fill {
    display: block;
    box-shadow: none;
    outline: 0;
  }
  .rangeslider-horizontal .rangeslider__fill {
    height: 100%;
    background-color: ${orange};
    outline: 0;
  }
  .rangeslider-horizontal .rangeslider__handle:after {
    content: ' ';
    position: absolute;
    width: 0;
    height: 0;
    box-shadow: 0;
    outline: 0;
  }
  .rangeslider .rangeslider__handle {
    box-shadow: none;
  }
  .rangeslider__handle-label {
    position: relative;
  }
  .rangeslider-horizontal .rangeslider__handle {
    width: 1.5rem;
    height: 1.5rem;
    border-radius: 50%;
    top: 50%;
    transform: translate3d(-50%, -50%, 0);
    outline: 0;
  }
  .rangeslider-horizontal .rangeslider__handle-tooltip {
    top: -3.5rem;
    background-color: transparent;
    color: ${darkGrey};
  }
  .rangeslider-horizontal .rangeslider__handle-tooltip {
    background-color: none !important;
    &:after {
      border-left: 8px solid transparent;
      border-right: 8px solid transparent;
      border-top: 8px solid transparent !important;
      left: 50%;
      bottom: -7px;
      transform: translate3d(-50%, 0, 0);
    }
  }
  .rangeslider__placeholder {
    display: block;
    box-shadow: none;
    outline: 0;
    width: 100%;
    height: 10px;
    border-radius: 10px;
    margin: 2rem 0;
    position: relative;
    background: #e6e6e6;
    touch-action: none;
    animation-duration: 1.5s;
    animation-fill-mode: forwards;
    animation-iteration-count: infinite;
    animation-timing-function: linear;
    animation-name: placeholderAnimate;
    background: linear-gradient(to right, #feefec 2%, #ff8c69 18%, #feefec 33%);
    background-size: 1300px; /* Animation Area */
  }
  @keyframes placeholderAnimate {
    0% {
      background-position: -650px 0;
    }
    100% {
      background-position: 650px 0;
    }
  }
`
