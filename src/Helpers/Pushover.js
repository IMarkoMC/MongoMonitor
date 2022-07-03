var Push = require('pushover-notification-client')
const { Warn, Debug, writeErr } = require('../Utils/Logger')

var p = null,
    sound = null;

exports.setup = function (user, token, customSound) {
    p = new Push({
        user: user,
        token: token,
        update_sounds: true // update the list of sounds every day - will
    })

    sound = customSound;

}


exports.sendNotification = function (title, message, callback) {

    if (process.env.NODE_ENV == 'development') {
        console.log({
            message: message,	// required
            title: title,
            sound: sound,
            priority: 2
        })
        return
    }

    var msg = {
        message: message,	// required
        title: title,
        sound: sound,
        priority: 2,
        retry: 60,
        expire: 1800
    }

    p.send(msg, function (err, result) {
        if (err) {
            // throw err
            Warn('Pushover error %s', err)
            writeErr(err)
            callback(err, null)
        }

        Debug('Pushover result %s', result)
        callback(null, result)
    })
}
