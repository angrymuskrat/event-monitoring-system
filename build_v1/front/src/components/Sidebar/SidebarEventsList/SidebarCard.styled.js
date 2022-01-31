import styled from 'styled-components'
import { lightGrey, orange } from '../../../config/styles'

export default styled.div`
  display: flex;
  flex-direction: row;
  justify-content: left;
  align-items: center;
  background-color: #fff;
  border: 1px solid ${lightGrey};
  box-sizing: border-box;
  box-shadow: 0px 3px 3px ${lightGrey};
  border-radius: 10px;
  padding: 0 1rem 0 0;
  margin: 1.1rem 0;
  width: 100%;
  transition: all 0.3s;
  overflow: hidden;
  &:hover {
    background-color: rgba(247, 248, 249, 0.95);
  }
  @media (max-width: 1130px) {
    flex-direction: column;
    padding: 0;
  }
  @media (max-width: 650px) {
    flex-direction: row;
  }
  @media (max-width: 450px) {
    flex-direction: column;
  }
  .event-card__image {
    position: relative;
    min-width: 14rem;
    min-height: 14rem;
    border-radius: 0.5rem;
    cursor: pointer;
    @media (max-width: 1130px) {
      flex-basis: 100%;
      min-width: 100%;
      min-height: 20rem;
    }
    @media (max-width: 650px) {
      flex-basis: 30%;
      min-width: 12rem;
      min-height: 12rem;
      max-width: 10rem;
      max-height: 10rem;
    }
    @media (max-width: 450px) {
      flex-basis: 100%;
      min-width: 100%;
      min-height: 15rem;
    }
  }
  .event-card__content {
    @media (max-width: 1130px) {
      align-self: flex-start;
      flex-basis: 100%;
      width: 100%;
    }
    margin-left: 1.5rem;
    flex-basis: 65%;
    width: 65%;
    overflow: hidden;
  }
  .event-card__title {
    cursor: pointer;
    text-overflow: ellipsis;
  }
  .event-card__tag {
    display: inline-block;
    background-color: ${lightGrey};
    border-radius: 5px;
    padding: 0.3rem;
    margin: 0.3rem;
    &:first-of-type {
      margin-left: 0;
    }
  }
  .event-card__button {
    color: ${orange};
    text-decoration: underline;
    margin: 1rem 0;
    cursor: pointer;
  }
  .text__events {
    margin-top: 1rem;
  }
`
