import styled from 'styled-components'

export default styled.div`
  cursor: pointer;
  height: 24vh;
  border-radius: 1rem;
  color: #ffffff;
  text-align: center;
  background: no-repeat center url(${props => props.picture});
  background-size: cover;
  filter: ${props => (props.avaliable ? 'grayscale(0%)' : 'grayscale(80%)')};
  align-self: end;
  box-shadow: 0px 5px 10px #bfc2c8;
  display: flex;
  align-items: center;
  justify-content: center;
  -webkit-transition: transform 0.3s;
  -moz-transition: transform 0.3s;
  -o-transition: transform 0.3s;
  transition: transform 0.3s;
  &:hover {
    transform: translateY(-5px);
  }
  .title {
    padding-top: 3rem;
    color: #ffffff;
  }
  .text {
    color: #ffffff;
  }
  @media (max-width: 490px) {
    max-width: 100%;
  }
`
