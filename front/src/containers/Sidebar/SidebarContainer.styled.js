import styled from 'styled-components'

export default styled.div`
  position: relative;
  transition: all 0.3s;
  flex-basis: 30%;
  min-width: 30%;
  margin-top: 7rem;
  /* Sidebar container styles */
  &.sidebar-container_active {
    margin-left: 0;
  }
  &.sidebar-container_closed {
    margin-left: -30%;
  }
  @media (max-width: 1300px) {
    flex-basis: 40%;
    min-width: 40%;
    &.sidebar-container_active {
      margin-left: 0;
    }
    &.sidebar-container_closed {
      margin-left: -40%;
    }
  }
  @media (max-width: 880px) {
    flex-basis: 50%;
    min-width: 50%;
    &.sidebar-container_active {
      margin-left: 0;
    }
    &.sidebar-container_closed {
      margin-left: -50%;
    }
  }
  @media (max-width: 650px) {
    position: fixed;
    flex-basis: 100%;
    min-width: 100%;
    z-index: 2;
    &.sidebar-container_active {
      margin-left: 0;
    }
    &.sidebar-container_closed {
      margin-left: -100%;
    }
  }
`
