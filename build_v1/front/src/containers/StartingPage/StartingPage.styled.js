import styled from 'styled-components'
import map from './img/map.png'
export default styled.div`
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  align-items: center;
  & .page__first-screen {
    height: 72vh;
    margin-top: 10rem;
    padding: 0 30px;
    display: flex;
    justify-content: space-between;
    position: relative;
    @media (max-width: 800px) {
      flex-direction: column;
    }
    @media (max-width: 1000px) {
      height: 80vh;
    }
    @media (max-width: 600px) {
      height: 25vh;
    }
  }
  & .title__starting-page {
    margin-bottom: 3rem;
  }
  & .page__title {
    max-width: 55rem;
    margin-bottom: 2rem;
  }
  & .page__text {
    max-width: 55rem;
    margin: 2rem 0;
  }
  & .page_start {
    width: 50%;
    margin: 10rem 0;
    align-self: flex-start;
    @media (max-width: 800px) {
      width: 100%;
      margin: 3rem 0;
    }
  }
  & .starting-page__image {
    position: absolute;
    top: 0;
    right: 0;
    z-index: -1;
    background-image: url(${map});
    background-position: right center;
    background-size: cover;
    background-repeat: no-repeat;
    background-clip: border-box;
    width: 100%;
    max-width: 62rem;
    height: 100%;
    max-height: 62rem;
    @media (max-width: 800px) {
      top: 34vh;
      left: 8%;
    }
    @media (max-width: 600px) {
      display: none;
    }
  }
  & .page__section {
    padding: 0 3rem;
    margin-bottom: 10rem;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    @media (max-width: 800px) {
      margin-top: 20rem;
    }
    @media (max-width: 600px) {
      margin-top: 0rem;
      margin-bottom: 0rem;
    }
  }
  & .page__cards {
    margin-bottom: 5rem;
  }
  & .page__main {
    width: 100%;
    max-width: 110rem;
    display: flex;
    flex-direction: column;
    justify-content: center;
  }
  .page__input-group {
    position: relative;
  }
  & input {
    font-size: 1.4rem;
  }
`
