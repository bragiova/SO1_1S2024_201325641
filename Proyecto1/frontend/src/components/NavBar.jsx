import React from 'react'
import { Link } from 'react-router-dom'

const NavBar = () => {
  return (
    <nav className='navbar navbar-expand-lg navbar-dark bg-dark'>
        <Link className="navbar-brand" to="/">SO1</Link>
        <button className="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbarSupportedContent" aria-controls="navbarSupportedContent" aria-expanded="false" aria-label="Toggle navigation">
            <span className="navbar-toggler-icon"></span>
        </button>

        <div className="collapse navbar-collapse" id="navbarSupportedContent">
          <ul className="navbar-nav mr-auto">
            <li className="nav-item active">
              <Link className="nav-link" to="/">Home</Link>
            </li>
            <li className="nav-item">
              <Link className="nav-link" to="/historical">Histórico</Link>
            </li>
            <li className="nav-item">
              <Link className="nav-link disabled" to="#">Árbol de Procesos</Link>
            </li>
            <li className="nav-item">
              <Link className="nav-link disabled" to="#">Simulación de Estados</Link>
            </li>
          </ul>
        </div>
    </nav>
  )
}

export default NavBar