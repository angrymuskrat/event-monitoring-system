import styled from 'styled-components'

export default styled.div`
  transition: all 0.3s;
  min-height: 35%;
  height: 100vh;
  overflow: scroll;
  .sidebar-fixed {
    position: fixed;
    top: 1rem;
    left: 0;
    background-color: #ffffff;
    z-index: 1;
    width: 30%;
    padding: 3rem;
    padding-top: 0;
    transition: all 0.3s;

    &.sidebar-fixed_active {
      left: 0;
    }
    &.sidebar-fixed_closed {
      left: -30%;
    }
    @media (max-width: 1300px) {
      width: 40%;
      &.sidebar-fixed_active {
        left: 0;
      }
      &.sidebar-fixed_closed {
        left: -40%;
      }
    }
    @media (max-width: 880px) {
      width: 50%;
      &.sidebar-fixed_active {
        left: 0;
      }
      &.sidebar-fixed_closed {
        left: -50%;
      }
    }
    @media (max-width: 880px) {
      width: 50%;
      &.sidebar-fixed_active {
        left: 0;
      }
      &.sidebar-fixed_closed {
        left: -50%;
      }
    }
    @media (max-width: 650px) {
      width: 100%;
      &.sidebar-fixed_active {
        left: 0;
      }
      &.sidebar-fixed_closed {
        left: -100%;
      }
    }
  }
`
