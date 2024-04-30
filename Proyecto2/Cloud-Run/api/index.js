const express = require('express');
const cors = require('cors');
const dotenv = require('dotenv');
const { dbConnection } = require('./db/config');
const Log  = require('./models/log');

const app = express();
dotenv.config();
const router = express.Router();

app.use(express.json());
app.use(cors());

dbConnection();

router.get('/logs', async (req, res) => {
    try {
        const logs = await Log.find().sort({ _id: -1}).limit(20);

        res.status(200).json(logs);
    } catch (error) {
        res.status(500).json({ msg: error});
        console.log(error);
    }
});

app.use('/api', router);

const port = process.env.PORT || 5000
app.listen(port, () => {
    console.log(`API corriendo en puerto ${ port }`);
})
