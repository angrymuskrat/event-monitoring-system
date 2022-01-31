import styled from 'styled-components'

export default styled.div`
  display: flex;
  justify-content: center;
  flex-direction: row;
  overflow: hidden;
  height: 100vh;
  #mapId {
    height: calc(100vh - 72px);
    margin-top: 72px;
    width: 100vw;
    z-index: 1;
  }
  .marker {
    width: 3rem !important;
    height: 3rem !important;
    will-change: transform, opacity, background;
    animation: marker 0.3s cubic-bezier(0.62, 0.2, 0.87, 1.47);
  }
  .marker-inner {
    position: relative;
    width: 5rem !important;
    height: 5rem !important;
    border: 3px solid white;
    border-radius: 50%;
  }
  .marker__description {
    position: absolute;
    top: 13%;
    left: 84%;
    z-index: -1;
    background-color: hsla(0, 0%, 100%, 0.95);
    box-shadow: 0 1px 4px rgba(0, 0, 0, 0.12), 0 2px 8px rgba(0, 0, 0, 0.06),
      0 0 16px rgba(0, 0, 0, 0.24);
    border-radius: 0px 5px 5px 0px;
    max-width: 15rem;
    min-width: 7.5rem;
    max-height: 3.7rem;
    overflow: hidden;
    padding: 0.8rem 0 0.8rem 1.3rem;
    display: -webkit-box;
    display: -webkit-flex;
    display: -ms-flexbox;
    display: flex;
    -webkit-align-items: center;
    -webkit-box-align: center;
    -ms-flex-align: center;
    align-items: center;
    -webkit-transition: 0.2s cubic-bezier(0.7, 0.2, 0.47, 0.2);
    transition: 0.2s cubic-bezier(0.7, 0.2, 0.47, 0.2);
    z-index: -10;
    font-size: 1rem;
    :hover {
      background-color: rgba(247, 248, 249, 0.95);
    }
  }
  @keyframes marker {
    from {
      opacity: 0;
    }
    to {
      opacity: 1;
    }
  }
  .leaflet-overlay-pane {
    z-index: 1 !important;
  }
  .leaflet-pane {
    z-index: 0 !important;
  }
  .leaflet-top,
  .leaflet-bottom {
    z-index: 2 !important;
  }
  .leaflet-control {
    z-index: 3 !important;
  }
  .leaflet-left {
    left: unset !important;
    right: 2rem !important;
  }
  .leaflet-top {
    top: 1rem !important;
  }
`
