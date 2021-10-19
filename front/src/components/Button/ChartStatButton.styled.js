import { badgeOrange, orange, lightGrey } from '../../config/styles'

import styled from 'styled-components'

export default styled.button`
  display: block;
  margin-left: 1rem;
  background-color: ${orange};
  border-radius: 2rem;
  color: #fff;
  font-weight: 600;
  letter-spacing: 0.5px;
  padding: 0.4rem 0.8rem;
  box-shadow: 0px 4px 3px ${lightGrey};
  cursor: pointer;
  transition-duration: 0.25s;
  &:hover {
    transform: translateY(-3px);
    box-shadow: 0 1px 4px rgba(0, 0, 0, 0.08), 0 2px 8px rgba(0, 0, 0, 0.06),
      0 0 16px rgba(0, 0, 0, 0.18);
  }
  background-color: ${({ type, isDemoPlay }) =>
    type === 'period' && isDemoPlay ? `${badgeOrange}` : `${orange}`};
  pointer-events: ${({ type, isDemoPlay }) =>
    type === 'period' && isDemoPlay ? 'none' : `inherit`};
  & .text {
    color: ${({ type, isDemoPlay }) =>
      type === 'period' && isDemoPlay ? '#C5CDD9' : `#fff`};
  }
`
