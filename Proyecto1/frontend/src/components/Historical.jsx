import { useEffect, useRef } from 'react'
import '../styles/Historical.css'
import Chart from 'chart.js/auto';

export const Historical = () => {
  let ramChartRef = useRef(null);
  let cpuChartRef = useRef(null);

  const getRamHistorical = () =>{
    fetch('/historicalram')
        .then((response) => response.json())
        .then((data) => {
            // console.log(data.map((info) => info.Percentage));
            const percentageVal = data.map((info) => info.Percentage);
            const dateVal = data.map((info) => info.Date);
            
            updateChartRamHistorical(percentageVal, dateVal);
            // updateChartCpu();
        })
        .catch((error) => {
            console.log("Error historical: ", error)
        })
  }

  const getCpuHistorical = () =>{
    fetch('/historicalcpu')
        .then((response) => response.json())
        .then((data) => {
            const percentageVal = data.map((info) => info.Percentage);
            const dateVal = data.map((info) => info.Date);
            
            updateChartCpuHistorical(percentageVal, dateVal);
        })
        .catch((error) => {
            console.log("Error historical: ", error)
        })
  }

  useEffect(() => {
    getRamHistorical();
    getCpuHistorical();
  }, []);

  useEffect(() => {
    const ctx = document.getElementById("chartRAM").getContext("2d");
    ramChartRef.current = new Chart(ctx, {
      type: "line",
      data: {
        labels: [],
        datasets: [
          {
            label: "Porcentaje de Uso",
            data: [],
            // backgroundColor: "rgb(255, 205, 86)", 
            borderColor: "rgb(255, 205, 86)",
            // fill: 'origin'
          },
        ],
        options: {
          scales: {
            x: {
              type: "category",
              title: {
                display: true,
                text: 'Fecha'
              }
            },
            y: {
              stacked: true,
              title: {
                display: true,
                text: 'Porcentaje'
              }
            }
          },
        },
      },
    });

    return () => {
        ramChartRef.current.destroy()
    }
  }, []);

  useEffect(() => {
    const ctx = document.getElementById("chartCPU").getContext("2d");
    cpuChartRef.current = new Chart(ctx, {
      type: "line",
      data: {
        labels: [],
        datasets: [
          {
            label: "Porcentaje de Uso",
            data: [],
            // backgroundColor: "rgb(255, 205, 86)", 
            borderColor: "rgb(133, 245, 159)",
            // fill: 'origin'
          },
        ],
        options: {
          scales: {
            x: {
              type: "category",
              title: {
                display: true,
                text: 'Fecha'
              }
            },
            y: {
              stacked: true,
              title: {
                display: true,
                text: 'Porcentaje'
              }
            }
          },
        },
      },
    });

    return () => {
        cpuChartRef.current.destroy()
    }
  }, []);

  const updateChartRamHistorical = (listPercentage, listDates) => {
    if (ramChartRef.current) {
      ramChartRef.current.data.labels = listDates;
      ramChartRef.current.data.datasets[0].data = listPercentage;
      ramChartRef.current.update();
    }
  };

  const updateChartCpuHistorical = (listPercentage, listDates) => {
    if (cpuChartRef.current) {
      cpuChartRef.current.data.labels = listDates;
      cpuChartRef.current.data.datasets[0].data = listPercentage;
      cpuChartRef.current.update();
    }
  };

  return (
    <>
      <div className='container'>
          <h1>Monitoreo Histórico</h1>
          <div className='flex-container'>
              <div className="chart-container flex-child">
                  <h2 style={{ textAlign: "center" }}>RAM</h2>
                  
                  <canvas id="chartRAM"></canvas>
              </div>
              <div className="chart-container flex-child">
                  <h2 style={{ textAlign: "center" }}>CPU</h2>
                  
                  <canvas id="chartCPU"></canvas>
              </div>
          </div>
      </div>
    </>
  )
}
