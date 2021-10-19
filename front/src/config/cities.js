import londonPic from '../assets/img/cityPictures/london.jpg'
import moscowPic from '../assets/img/cityPictures/moscow.jpg'
import newyorkPic from '../assets/img/cityPictures/newyork.jpg'
import saintPetersburgPic from '../assets/img/cityPictures/spb.jpg'

export const cities = [
  {
    id: 'spb',
    city: 'Saint Petersburg',
    country: 'Russia',
    avaliable: true,
    picture: saintPetersburgPic,
    lng: 30.32315244950895,
    lat: 59.9271516041233,
    topLeft: [60.12, 30.11],
    bottomRight: [59.84, 30.69],
  },
  {
    id: 'nyc',
    city: 'New York',
    country: 'USA',
    avaliable: false,
    picture: newyorkPic,
    lng: -73.99152387944392,
    lat: 40.701733209232735,
    topLeft: [40.8482826, -73.9873646],
    bottomRight: [40.6185618, -74.0340492],
  },
  {
    id: 'moscow',
    city: 'Moscow',
    country: 'Russia',
    avaliable: false,
    picture: moscowPic,
    lng: 37.61296270571807,
    lat: 55.73652849918221,
  },
  {
    id: 'london',
    city: 'London',
    country: 'United Kingdom',
    avaliable: false,
    picture: londonPic,
    lng: 0,
    lat: 0,
  },
]
