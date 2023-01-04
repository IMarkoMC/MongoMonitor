//Custom console colors can be defined here:
// To call a custom color use <color> example: Error in main <red> Error </red> (Not implemented) omegalul
const chalk = require('chalk');

module.exports = {
    info: chalk.rgb(3, 252, 48),
    debug: chalk.rgb(238, 59, 247),
    objects: chalk.rgb(255, 191, 0),
    arguments: chalk.rgb(255, 0, 242),
    green: chalk.rgb(0, 255, 0),
    red: chalk.rgb(255, 0, 0),
}