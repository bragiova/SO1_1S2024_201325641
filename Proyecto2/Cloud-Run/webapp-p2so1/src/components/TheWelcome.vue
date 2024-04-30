<script setup>

</script>

<template>
  <div class="container">
    <textarea name="console" id="console" cols="100" rows="30" readonly v-model="consoleOutput" style="font-family: monospace;"></textarea>
    <div>
      <button @click="fetchData" class="btn btn-primary">Actualizar logs</button>
    </div>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  data() {
    return {
      consoleOutput: '',
    };
  },
  methods: {
    async fetchData() {
      try {
        // Hacer la solicitud a la API
        const response = await axios.get('https://apilogsso1-o5auvpo3pq-uk.a.run.app/api/logs');

        this.consoleOutput = '';
        
        // Actualizar el contenido del textarea con los datos de la API
        response.data.forEach(element => {
          this.consoleOutput += `${ element.fecha } - ${ element.hora } - ${ element.voto }\n`
        });
        // this.consoleOutput = JSON.stringify(response.data, null, 2);
      } catch (error) {
        console.error('Error al obtener los datos de la API:', error);
      }
    },
    prueba () {
      this.consoleOutput += 'Prueba de consola\n'
    }
  },
};
</script>

<style scoped>
/* Estilos específicos para este componente Vue.js */
textarea {
  font-family: monospace;
  background-color: #181717;
  color: #2dc008;
  border: 1px solid #2edff7;
  padding: 5px;
  margin-bottom: 10px;
  border-radius: 5px;
  width: calc(100% - 10px); /* Ajuste para compensar el padding */
}

button {
  background-color: #55a5f0;
  color: #2c2a2a;
  border: none;
  padding: 10px 20px;
  cursor: pointer;
  border-radius: 5px;
}
</style>
