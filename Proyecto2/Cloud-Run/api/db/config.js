const mongoose = require('mongoose');

const dbConnection = async() => {
    try {
        const stringConnect = process.env.DB_HOST + '/' + process.env.DB_NAME + '?authSource=admin';
        await mongoose.connect(stringConnect, {
            useNewUrlParser: true,
            useUnifiedTopology: true,
        });

        console.log('BD mongo online');
    } catch (error) {
        console.log('Error conectar Mongo: ' + error);
    }
}

module.exports = {
    dbConnection
}