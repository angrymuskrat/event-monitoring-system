import React, { useState, useRef } from 'react'
import PropTypes from 'prop-types'

import Modal from 'react-modal'
import Slider from 'react-slick'
import 'slick-carousel/slick/slick.css'
import 'slick-carousel/slick/slick-theme.css'

// components
import { Cross } from '../../assets/svg/Cross'
import Post from '../Post/Post.jsx'
import SliderButton from '../../assets/svg/SliderButton'

// styled
import Container from './Popup.styled'
import ModalContainer from './Modal.styled'

// styles config
import { modalStyles } from '../../config/styles'

Modal.setAppElement('#root')

function Popup({ isPopupOpen, closePopup, data }) {
  const sliderRef = useRef()
  const [currentSlide, setCurrentSlide] = useState(0)
  const slickSettings = {
    dots: false,
    infinite: true,
    speed: 500,
    slidesToShow: 1,
    slidesToScroll: 1,
    beforeChange: (current, next) => setCurrentSlide(next),
  }
  const goNext = () => {
    if (sliderRef.current) {
      sliderRef.current.slickNext()
    }
  }
  const goPrev = () => {
    if (sliderRef.current) {
      sliderRef.current.slickPrev()
    }
  }
  return (
    <Modal
      isOpen={isPopupOpen}
      onRequestClose={closePopup}
      style={modalStyles}
      contentLabel="Modal"
    >
      <ModalContainer>
        <button className="modal__button" onClick={closePopup}>
          <Cross />
        </button>
        <Container>
          <Slider {...slickSettings} ref={sliderRef}>
            {data ? (
              data.map(p => {
                return p !== undefined && <Post key={p.id} post={p} />
              })
            ) : (
              <div className="popup__loading">
                <p>Loading...</p>
                <div className="spinner-container">
                  <div className="spinner"></div>
                </div>
              </div>
            )}
          </Slider>
          <div className="popup__counter">
            <p className="text text_subheading">
              {data && currentSlide + 1 + '/' + data.length}
            </p>
          </div>
          <button
            className="slider__button slider__button_next"
            onClick={goNext}
          >
            <span>
              <SliderButton />
            </span>
          </button>
          <button
            className="slider__button slider__button_prev"
            onClick={goPrev}
          >
            <span>
              <SliderButton />
            </span>
          </button>
        </Container>
      </ModalContainer>
    </Modal>
  )
}
Popup.propTypes = {
  isPopupOpen: PropTypes.bool.isRequired,
  closePopup: PropTypes.func.isRequired,
  postcodes: PropTypes.array,
}
export default Popup
