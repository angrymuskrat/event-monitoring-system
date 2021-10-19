import 'normalize.css/normalize.css'
import { createGlobalStyle } from 'styled-components'
import {
  fontSize,
  darkGrey,
  grey,
  defaultBackground,
  orange,
} from '../../config/styles'
import openSansRegular from '../../assets/fonts/Open_Sans/OpenSans-Regular.ttf'
import openSansLight from '../../assets/fonts/Open_Sans/OpenSans-Light.ttf'
import openSansSemiBold from '../../assets/fonts/Open_Sans/OpenSans-SemiBold.ttf'
import openSansBold from '../../assets/fonts/Open_Sans/OpenSans-Bold.ttf'

import montserratRegular from '../../assets/fonts/Montserrat/Montserrat-Regular.ttf'
import montserratSemiBold from '../../assets/fonts/Montserrat/Montserrat-SemiBold.ttf'

export const GlobalStyle = createGlobalStyle`
  * {
    box-sizing: border-box;
    

    &::-webkit-scrollbar {
      width: 0;
    }

    &::-webkit-scrollbar-track {
      -webkit-box-shadow: none;
      border-radius: 0;
    }
    ::-moz-scrollbar{
      width: 0px;
      height: 0px;
    }

    &::-webkit-scrollbar-thumb {
      border: transparent 1px solid;
      border: transparent;
      background-color: transparent;
    }
    h1, h2, h3, h4, h5, h6 {
      margin: 0;
      margin-block-start: 0.5rem;
      margin-block-end: 0.5rem;
    }
    p {
      margin-block-start: 0em;
      margin-block-end: 0.5rem;
    }
  }


  @font-face {
    font-family: Open-Sans-Regular;
    src:
      url(${openSansRegular}) format('truetype');
    font-weight: 400;
    font-style: normal;
  }
  @font-face {
    font-family: Open-Sans-Light;
    src:
      url(${openSansLight}) format('truetype');
    font-weight: 300;
    font-style: normal;
  }
  @font-face {
    font-family: Open-Sans-SemiBold;
    src:
      url(${openSansSemiBold}) format('truetype');
    font-weight: 700;
    font-style: normal;
  }
  @font-face {
    font-family: Open-Sans-Bold;
    src:
      url(${openSansBold}) format('truetype');
    font-weight: 800;
    font-style: normal;
  }
  @font-face {
    font-family: Montserrat-Regular;
    src:
      url(${montserratRegular}) format('truetype');
    font-weight: 400;
    font-style: normal;
  }
  @font-face {
    font-family: Montserrat-SemiBold;
    src:
      url(${montserratSemiBold}) format('truetype');
    font-weight: 700;
    font-style: normal;
  }

  html {
    font-size: ${fontSize};
  }
  body {
    font-family: 'Monserrat-Regular', sans-serif;
    color: ${darkGrey};
    background-color: ${defaultBackground};
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
  }
  a {
    text-decoration: none;
  }
  input { border-style: none; background: transparent; outline: 0; font-family: 'Open-Sans-Regular', sans-serif; letter-spacing: 0.2px }
  button { padding: 0; background: none; border: none; outline: 0; }

  .title {
    font-family: 'Montserrat-SemiBold', sans-serif;
    color: ${darkGrey};
  }
  .title_light {
    color: white;
  }
  .title_h1 {
    font-size: 4.8rem;
    letter-spacing: 2px;
    @media (max-width: 1000px) {
      font-size: 3.8rem;
    }
    @media (max-width: 400px) {
      font-size: 2.5rem;
    }
  }
  .title_h2 {
    font-size: 2.7rem;
    @media (max-width: 400px) {
      font-size: 1.7rem;
    }
  }
  .title_h3 {
    font-size: 2.4rem;
  }
  .title_h4 {
    font-size: 1.73rem;
    @media (max-width: 400px) {
      font-size: 1.5rem;
    }
  }
  .title_h5 {
    font-size: 1.5rem;
  }
  .title_light {
    color: ${grey}
  }
  .text {
    font-family: 'Open-Sans-Regular', sans-serif;
    color: ${darkGrey};
  }
  .text_bold {
    font-family: 'Open-Sans-Bold', sans-serif;
  }
  .text_light {
    color: white;
  }
  .text_s1 {
    font-size: 1.12rem;
  }
  .text_s2 {
    font-size: 0.87rem;
  }
  .text_error {
    color: #E32636;
  }
  .text_subheading {
    font-size: 1.8rem;
    line-height: 3.4rem;
    @media (max-width: 1000px) {
      font-size: 1.5rem;
      line-height: 2.2rem;
    }
  }
  .text_p1 {
    font-size: 1.8rem;
    @media (max-width: 1000px) {
      font-size: 1.5rem;
    }
    @media (max-width: 400px) {
      font-size: 1rem;
    }
  }
  .text_p1-light {
    font-family: 'Open-Sans-Light', sans-serif;
    font-size: 1.4rem;
  }
  .text_p2-bold {
    font-family: 'Open-Sans-Bold', sans-serif;
    font-size: 1rem;
  }
  .text_p2 {
    font-size: 1.1rem;
  }
  .text_p3 {
    font-size: 1.5rem;
  }
  .text_subtitle {
    font-family: 'Open-Sans-Semibold', sans-serif;
    font-size: 1rem;
  }
  .text_input {
    font-size: 1.4rem;
  }
  .text_crossed {
    text-decoration: line-through;
    color: ${grey}
  }
  .text_italic {
    font-style: italic;
  }
  .text_link {
    transition: 0.3s all;
    &:hover {
      color: ${orange}
    }
  }
  .text_post {
    font-size: 1.3rem;
  }
  .text_location {
    margin-block-start: 0.5rem;
  }
  .ReactModal__Overlay--after-open {
    z-index: 200;
  }
  .ReactModal__Content--after-open {
    max-width: 44rem;
  }
`
