import axios from 'axios'
//const server = 'http://pp.onti.actcognitive.org:17112'
const server = 'http://10.64.0.206:17112'

const config = {
    headers: {
    },
    withCredentials: true,
    credentials: 'same-origin',
}
export const requestSessionCookie = data => {
    return axios(`${server}/login`, {
        method: 'post',
        data: data,
        withCredentials: true,
        headers: {'Content-Type': 'application/json'},
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
export const fetchSearchEventsData = ({city, tags, start, finish}) => {
    return axios.get(
        `${server}/search/${city}/${tags}/${start}/${finish}`,
        config
    )
}
export const fetchPostData = ({city, postcode}) => {
   
    return axios.get(`${server}/singleShortPost/${city}/${postcode}`, config)
}
// export const requestSessionCookie = data => {
//   return axios(`${server}/login`, {
//     method: 'post',
//     data: data,
//     withCredentials: true,
//     headers: { 'Content-Type': 'application/json' },
//   })
// }
