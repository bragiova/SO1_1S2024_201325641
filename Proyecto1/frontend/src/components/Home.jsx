import { useState, useEffect, useRef } from 'react'
import Chart from 'chart.js/auto';
import '../styles/Home.css'

export const Home = () => {
    const [cpuUsage, setCpuUsage] = useState(0);
    const [freeCpu, setFreeCpu] = useState(0);
    const [freeRam, setFreeRam] = useState(0);
    const [ramUsage, setRamUsage] = useState(0);

    let ramChartRef = useRef(null);
    let cpuChartRef = useRef(null);

    const getCpuRamInfo = () =>{
        fetch('/realtime')
            .then((response) => response.json())
            .then((data) => {
                // console.log(data);
                const usageRam = +data.ram_porcentaje.toFixed(2);
                const usageCpu = +data.cpu_porcentaje.toFixed(2);
                const ramFree = 100 - usageRam;
                const cpuFree = 100 - usageCpu;

                setCpuUsage(usageCpu);
                setFreeCpu(cpuFree);
                setRamUsage(usageRam);
                setFreeRam(ramFree);
                updateChartRam();
                updateChartCpu();
            })
            .catch((error) => {
                console.log("Error realtime: ", error)
            })
    }

    useEffect(() => {
        const intervalId = setInterval(() => {
            getCpuRamInfo();
        }, 2000);
    
        return () => clearInterval(intervalId);
      }, []);

    useEffect(() => {
        const ctx = document.getElementById("chartRAM").getContext("2d");
        ramChartRef.current = new Chart(ctx, {
          type: "pie",
          data: {
            labels: ["En Uso", "Libre"],
            datasets: [
              {
                data: [ramUsage, freeRam],
                backgroundColor: ["rgb(255, 205, 86)", "rgb(54, 162, 235)"], 
                // radius: '70%'
              },
            ]
          },
        });
    
        return () => {
            ramChartRef.current.destroy()
        }
    }, [ramUsage, freeRam]);

    const updateChartRam = () => {
        if (ramChartRef.current) {
            ramChartRef.current.data.datasets[0].data = [ramUsage, freeRam];
            ramChartRef.current.update();
        }
    };

    useEffect(() => {
        const ctx = document.getElementById("chartCPU").getContext("2d");
        cpuChartRef.current = new Chart(ctx, {
          type: "pie",
          data: {
            labels: ["En Uso", "Libre"],
            datasets: [
              {
                data: [cpuUsage, freeCpu],
                backgroundColor: ["rgb(235, 161, 60)", "rgb(168, 252, 219)"],
                // radius: '70%'
              },
            ],
          },
        });
    
        return () => {
            cpuChartRef.current.destroy()
        }
    }, [cpuUsage, freeCpu]);

    const updateChartCpu = () => {
        if (cpuChartRef.current) {
            cpuChartRef.current.data.datasets[0].data = [cpuUsage, freeCpu];
            cpuChartRef.current.update();
        }
    };


  return (
    <>
        <div className='container'>
            <h1>Monitoreo en Tiempo Real</h1>
            <div className='flex-container'>
              <div className="chart-container flex-child">
                <h2 style={{ textAlign: "center" }}>RAM</h2>
                <div className="result-container">
                    <span>{ramUsage} % en uso</span>
                </div>
                <div className="result-container">
                    <span>{freeRam} % libre</span>
                </div>
                <canvas id="chartRAM"></canvas>
              </div>
              <div className="chart-container flex-child">
                <h2 style={{ textAlign: "center" }}>CPU</h2>
                <div className="result-container">
                    <span>{cpuUsage} % en uso</span>
                </div>
                <div className="result-container">
                    <span>{freeCpu} % libre</span>
                </div>
                <canvas id="chartCPU"></canvas>
              </div>
            </div>
        </div>
    </>
  )
}
