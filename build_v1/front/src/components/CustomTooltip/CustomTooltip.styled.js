import styled from 'styled-components'
import {
  lightGrey,
  chartLightGrey,
  badgeOrange,
  orange,
} from '../../config/styles'

export default styled.div`
  box-sizing: border-box;
  padding: 1.5rem;
  background-color: #fff;
  border: 1px solid ${lightGrey};
  border-radius: 0.5rem;
  & .chart__tooltip-text_main {
    font-family: 'Open-Sans-Bold', sans-serif;
    margin-bottom: 1rem;
  }
  & .chart__tooltip-label {
    display: inline-block;
    padding: 4px;
    border-radius: 0.5rem;
  }
  & .chart__tooltip-label_events {
    background-color: ${badgeOrange};
    color: ${orange};
  }
  & .chart__tooltip-label_posts {
    background-color: ${chartLightGrey};
  }
`
