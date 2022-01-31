import styled from 'styled-components'
import { lightGrey, darkGrey } from '../../config/styles'

export default styled.div`
  flex-basis: 50%;
  margin-right: 1rem;
  svg {
    display: block;
    width: 1em;
    height: 1em;
    fill: currentColor;
  }
  .select {
    position: relative;
    display: inline-block;
    min-width: 11rem;
    width: 100%;

    &:focus {
      outline: 0;
      &__input {
        box-shadow: 0 0 1px 1px #00a9e0;
      }
    }
  }
  .select__input {
    position: relative;
    padding: 0.7rem;
    border: 1px solid rgba(143, 155, 179, 0.23);
    border-radius: 1rem;
    background: #fff;
  }
  .select__input_open {
    border-radius: 1rem 1rem 0px 0px;
  }
  .select__value {
    position: relative;
    display: inline-block;
    padding: 0.5rem 1rem;
  }
  .select__placeholder {
    padding: 0.5rem 1rem;
    color: ${darkGrey};
  }
  .select__arrow {
    position: absolute;
    top: 0.5rem;
    right: 0.5rem;
    display: block;
    padding: 1rem;
    font-size: 1rem;
    color: ${darkGrey};
  }
  .select__options {
    position: absolute;
    z-index: 200;
    top: 100%;
    left: 0;
    right: 0;
    border-width: 0 1px;
  }
  .select__option {
    padding: 1rem 1.5rem;
    border-bottom: 1px solid rgba(143, 155, 179, 0.23);
    border-left: 1px solid rgba(143, 155, 179, 0.23);
    border-right: 1px solid rgba(143, 155, 179, 0.23);
    background-color: #ffffff;
    cursor: pointer;
    &_selected {
      border: 1px solid #00a9e0;
      margin: -1px -1px 0;
      background: #d9f2fb;
      pointer-events: none;
    }
    &_focused {
      background: ${lightGrey};
    }
    &_last {
      border-radius: 0px 0px 10px 10px;
    }
  }
`
