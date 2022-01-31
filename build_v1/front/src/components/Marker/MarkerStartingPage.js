import styled from 'styled-components'
import pic1 from '../../assets/img/markerPictures/pic1.png'
import pic2 from '../../assets/img/markerPictures/pic2.png'
import pic3 from '../../assets/img/markerPictures/pic3.png'
import pic4 from '../../assets/img/markerPictures/pic4.png'

export default styled.div`
  .marker_starting-page {
    position: relative;
    width: 5.5rem !important;
    height: 5.5rem !important;
    background-size: cover;
    background-position: top;
    background: #999;
    border: 3px solid hsla(0, 0%, 100%, 0.95);
    border-radius: 50%;
    padding: 1rem;
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 100;
  }
  .marker_starting-page-1 {
    left: 60%;
    top: 18rem;
    background: ${`url(${pic4}) no-repeat`};
    background-size: cover;
    @media (max-width: 800px) {
      left: 57%;
      top: 6.4rem;
    }
  }
  .marker_starting-page-2 {
    left: 60%;
    top: 38.7rem;
    background: ${`url(${pic2}) no-repeat`};
    background-size: cover;
    @media (max-width: 800px) {
      left: 73%;
      top: 29.7rem;
    }
    @media (max-width: 600px) {
      display: none;
    }
  }
  .marker_starting-page-3 {
    left: 32%;
    top: 7.2rem;
    background: ${`url(${pic3}) no-repeat`};
    background-size: cover;
    @media (max-width: 800px) {
      left: 30%;
      top: -2.4rem;
    }
  }
  .marker_starting-page-4 {
    left: 63%;
    top: 15.4rem;
    background: ${`url(${pic1}) no-repeat`};
    background-size: cover;
    @media (max-width: 800px) {
      left: 48%;
      top: 8rem;
    }
  }
  .marker__animation {
    position: absolute;
    left: 5px;
    z-index: -1;
    &:before,
    &:after {
      content: '';
      width: 4rem;
      height: 4rem;
      border-radius: 50%;
      position: absolute;
      top: 0;
      right: 0;
      bottom: 0;
      left: 0;
      margin: auto;
      transform: scale(0.5);
      transform-origin: center center;
      animation: pulse-me 3s linear infinite;
    }
    &:after {
      animation-delay: 2s;
    }
  }
  @keyframes pulse-me {
    0% {
      transform: scale(0.5);
      opacity: 0;
    }
    50% {
      opacity: 0.1;
      background-color: tomato;
    }
    70% {
      opacity: 0.09;
      background-color: tomato;
    }
    100% {
      transform: scale(5);
      opacity: 0;
    }
  }
`
