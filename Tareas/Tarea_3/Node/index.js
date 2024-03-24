const Redis = require('ioredis');

const conexion = new Redis({
    host: '10.41.80.155',
    port: 6379,
    connectTimeout: 5000
});

const functionPub = () => {
    conexion.publish('test', JSON.stringify({ msg: 'Hola a todos' }))
        .then(() => {
            console.log('Mensaje publicado')
        })
        .catch(error => {
            console.log('Error al publicar')
        })
}

setInterval(functionPub, 5000);
