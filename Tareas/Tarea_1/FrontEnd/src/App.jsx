import {useState} from 'react'

export const App = () => {
  const [info, setInfo] = useState()

  const getData = async() => {
    const response = await fetch('http://localhost:3000/data')
    const data = await response.json();

    setInfo(() => `${data.name}\n${data.carnet}`)
   }

  return (
    <>
       <h1>Tarea 1 - SO1</h1>
       <button onClick={getData}>Mostrar Datos</button>
       <p className='parrafo'>{info}</p>
    </>
  )
}
