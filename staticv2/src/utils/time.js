
const ONE_SECOND = 1000;
const ONE_MINUTE = 60 * ONE_SECOND;
const ONE_HOUR = 60 * ONE_MINUTE;
const ONE_DAY = 24 * ONE_HOUR;

const seconds = (s) => s * ONE_SECOND;
const minutes = (m) => m * ONE_MINUTE;
const hours = (h) => h * ONE_HOUR;
const days = (d) => d * ONE_DAY;

const checkTimeFrom = (date, ms) => {
    return date.getTime() > Date.now() - ms
}

export {
    ONE_SECOND,
    ONE_MINUTE,
    ONE_HOUR,
    ONE_DAY,
    seconds,
    minutes,
    hours,
    days,
    checkTimeFrom
}

