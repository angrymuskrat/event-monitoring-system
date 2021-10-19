import styled from 'styled-components'
import { lightGrey, orange, textError } from '../../config/styles'

export default styled.div`
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  height: 100vh;
  width: 100vw;
  .label__hidden {
    border: 0;
    clip: rect(0 0 0 0);
    height: 1px;
    margin: -1px;
    overflow: hidden;
    padding: 0;
    position: absolute;
    width: 1px;
  }
  .login-form {
    width: 40%;
    max-width: 400px;
    @media (max-width: 800px) {
      width: 70%;
    }
  }
  .form__field {
    width: 100%;
    margin-bottom: 1rem;
    min-height: 13px;
  }
  .form__field-error {
    font-size: 1.3rem;
    color: ${textError};
  }
  & input {
    background-image: none;
    border: 0;
    color: inherit;
    font: inherit;
    margin: 0;
    outline: 0;
    padding: 0;
    transition: background-color 0.3s;
    width: 100%;
    height: 5rem;
    border-radius: 0.5rem;
    padding: 0 1.3rem;
    font-size: 1.5rem;
    background-color: ${lightGrey};
    transition: all 0.3s;
    :hover {
      opacity: 0.7;
    }
    @media (max-width: 450px) {
      height: 4rem;
      font-size: 1.3rem;
    }
  }
  & button {
    background-color: ${orange};
    color: #ffffff;
    width: 100%;
    font-size: 1.5rem;
    font-weight: 700;
    cursor: pointer;
    height: 5rem;
    border-radius: 0.5rem;
    transition: all 0.3s;
    :hover {
      opacity: 0.7;
    }
    @media (max-width: 450px) {
      height: 4rem;
      font-size: 1.3rem;
    }
  }
`
