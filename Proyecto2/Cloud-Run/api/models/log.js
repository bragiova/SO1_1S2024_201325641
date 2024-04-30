const { Schema, model } = require('mongoose');

const LogSchema = Schema({
    fecha: String,
    hora: String,
    dato: String
});

module.exports = model('logs', LogSchema);

