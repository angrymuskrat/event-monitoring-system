import React from 'react'
import PropTypes from 'prop-types'

// images
import { Cross } from '../../../assets/svg/Cross'

// styled
import Container from './SidebarNav.styled'

function SidebarNav({ handleCheckboxChange, toggleSidebar }) {
  return (
    <>
      <Container>
        <button
          className="sidebar__close-button"
          onClick={() => toggleSidebar()}
        >
          <Cross />
        </button>
        <div className="sidebar__filters-toggle">
          <p className="text text_p2">Advanced Search</p>
          <label className="sidebar__toggle-active">
            <input type="checkbox" onChange={handleCheckboxChange} />
            <span className="slider round"></span>
          </label>
        </div>
      </Container>
    </>
  )
}
SidebarNav.propTypes = {
  handleCheckboxChange: PropTypes.func.isRequired,
  toggleSidebar: PropTypes.func.isRequired,
}

export default SidebarNav
