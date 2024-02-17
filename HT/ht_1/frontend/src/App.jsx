import {useState, useEffect} from 'react';
import logo from './assets/images/logo-universal.png';
import './App.css';
import {Greet} from "../wailsjs/go/main/App";
import Chart from 'chart.js/auto';

const App = () => {
    const [uso, setUso] = useState(0);
    const [libre, setLibre] = useState(0);
    const [ramVacia, setRamVacia] = useState(0);
    const [ramLlena, setRamLlena] = useState(0);
    let graficaRef = null;

    useEffect(() => {
        const intervalId = setInterval(() => {
          Greet().then((data) => {
            
            const jsonObject = JSON.parse(data);
            const libre = jsonObject.libre;
            const uso = jsonObject.uso;
            const total = libre + uso;
    
            const ramVacia = +((libre / total) * 100).toFixed(2);
            const ramLlena = +((uso / total) * 100).toFixed(2);
    
            setUso(uso);
            setLibre(libre);
            setRamVacia(ramVacia);
            setRamLlena(ramLlena);
            updateGrafico();
          });
        }, 1000);
    
        return () => clearInterval(intervalId);
      }, []);

    useEffect(() => {
        const ctx = document.getElementById("grafica").getContext("2d");
        graficaRef = new Chart(ctx, {
          type: "pie",
          data: {
            labels: ["En Uso", "Libre"],
            datasets: [
              {
                data: [uso, libre],
                backgroundColor: ["rgb(255, 205, 86)", "rgb(54, 162, 235)"],
              },
            ],
          },
        });
    
        return () => {
            graficaRef.destroy()
        }
      }, [uso, libre]);

    const updateGrafico = () => {
        if (graficaRef) {
            graficaRef.data.datasets[0].data = [ramLlena, ramVacia];
            graficaRef.update();
        }
    };

    return (
        <>
            <h1>Monitoreo RAM</h1>
            <div className="result-container">
                <span>{ramLlena} % en uso</span>
            </div>
            <div className="result-container">
                <span>{ramVacia} % libre</span>
            </div>
            <div className="chart-container">
                <canvas id="grafica" width="200" height="200"></canvas>
            </div>
        </>
    )
}

export default App
