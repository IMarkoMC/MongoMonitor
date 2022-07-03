
const MongoClient = require('mongodb').MongoClient,
    { Info, Error, writeErr } = require('../Utils/Logger');

exports._init_ = async function (url, cb) {
    let client = new MongoClient(url, { useNewUrlParser: true, useUnifiedTopology: true });

    try {
        await client.connect()

        Info('MongoDB connected')

        conn = client;

        //* When the client is connected. send a callback so the app loads
        cb(true);

    } catch (error) {
        Error('MongoDB error %s', error)
        console.log(error);
        writeErr(error)

        cb(false)
    }
}


exports.getConn = () => {
    return conn
}