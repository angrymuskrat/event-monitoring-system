import axios from 'axios'
const server = 'http://pp.onti.actcognitive.org:17112'
const config = {
  headers: {
    'Content-Type': 'application/x-www-form-urlencoded',
  },
  withCredentials: true,
  credentials: 'same-origin',
}
export const requestSessionCookie = data => {
  return axios(`${server}/login`, {
    method: 'post',
    data: data,
    withCredentials: true,
    headers: { 'Content-Type': 'application/json' },
  })
}




export const fetchHeatmap = ({
  city,
  topLeft,
  botRight,
  time,
  selectedHour,
}) => {
  return axios.get(
    `${server}/heatmap/${city}/${topLeft}/${botRight}/${
      selectedHour ? selectedHour : time
    }`,
    config
  )
}
export const fetchEvents = ({
  city,
  topLeft,
  botRight,
  time,
  selectedHour,
}) => {
  return axios.get(
    `${server}/events/${city}/${topLeft}/${botRight}/${
      selectedHour ? selectedHour : time
    }`,
    config
  )
}
export const fetchChartTimeData = (city, start, finish) => {
  return axios.get(`${server}/timeline/${city}/${start}/${finish}`, config)
}
export const fetchSearchEventsData = ({ city, tags, start, finish }) => {
  return axios.get(
    `${server}/search/${city}/${tags}/${start}/${finish}`,
    config
  )
}
export const fetchPostData = postcode => {
  return axios.get(`https://www.instagram.com/p/${postcode}/?__a=1`)
}
// export const requestSessionCookie = data => {
//   return axios(`${server}/login`, {
//     method: 'post',
//     data: data,
//     withCredentials: true,
//     headers: { 'Content-Type': 'application/json' },
//   })
// }
