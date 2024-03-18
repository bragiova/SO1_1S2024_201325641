import { useState } from 'react'
import './App.css'
import { Route, Routes} from 'react-router-dom'
import { Home } from './components/Home';
import { Historical } from './components/Historical';
import NavBar from './components/NavBar';

function App() {

  return (
    <div className='App'>
      <header>
            <NavBar />
        </header>
      <Routes>
        <Route path='/' element={<Home />} />
        <Route path='/historical' element={<Historical />} />
      </Routes>
    </div>
    
  )
}

export default App
