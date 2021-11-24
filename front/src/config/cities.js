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
    avaliable: true,
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
    avaliable: true,
    picture: moscowPic,
    lng: 37.61296270571807,
    lat: 55.73652849918221,
    topLeft: [56.07182, 36.5666568499],
    bottomRight: [55.0646688279, 38],
  },
  {
    id: 'london',
    city: 'London',
    country: 'United Kingdom',
    avaliable: true,
    picture: londonPic,
    lng: -0.127696,
    lat: 51.507351,
    topLeft: [51.6578218154, -0.3478717804],
    bottomRight: [51.37582, 0.24472],
  },
]
