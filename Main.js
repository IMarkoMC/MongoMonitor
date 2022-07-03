const { readFileSync } = require("fs"),
    { parse } = require('yaml'),
    { Info, Debug, Warn, Error, setDebug } = require('./src/Utils/Logger');

class Main {

    constructor() {
        this.config = parse(readFileSync('./Settings.yml', 'utf8'));

        Info('Starting!');

        setDebug(this.config.Debug)

        this.mongo = require('./src/Databases/Mongo');
        this.pushover = require('./src/Helpers/Pushover');

        this.__init__();
    }

    /** @private */
    __init__() {
        this.mongo._init_(this.config.Mongo.URL, (cont) => {

            if (!cont) {
                Warn('Mongo did not connect. Exiting in 5 seconds');

                setTimeout(() => {
                    process.exit(0)
                }, 5e3);
                return
            }

            this.pushover.setup(this.config.Pushover.User, this.config.Pushover.Token, this.config.Pushover.Sound);

            this.listen();
        })

    }

    listen() {



        this.mongo.getConn().on('serverDescriptionChanged', (data) => {
            //? Is the member down?

            if (data.previousDescription.type == 'RSSecondary' && data.newDescription.type == 'Unknown') {

                Warn('The node %s went down!', data.address)

                //! Send the pushover notificatrion

                //? Did we get a message?

                if ('error' in data.newDescription) {
                    Debug('Error code: %s', data.newDescription.error?.code)
                    Debug('Error codeName: %s', data.newDescription.error?.codeName)

                    this.pushover.sendNotification('A node went offline!', `The MongoDB replica member at ${data.address} went offline\nError: ${data.newDescription.error?.codeName} [${data.newDescription.error?.code}]`, (done) => {
                        Debug(done)
                        data = null
                    })

                    return
                }

                this.pushover.sendNotification('A node went offline!', `The MongoDB replica member at ${data.address} went offline}`, (done) => {
                    Debug(done)
                    data = null
                })

                return
            }


            if (data.previousDescription.type == 'Unknown' && data.newDescription.type == 'RSSecondary') {

                Info('The node %s is back up!', data.address);

                this.pushover.sendNotification('The node is back online!', `The MongoDB replica member at ${data.address} is online again\nCurrent round trip time ${data.newDescription.roundTripTime} ms`, (done) => {
                    Debug(done)
                })
                return
            }
        })
    }
}

new Main();
