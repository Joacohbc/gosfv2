// Genero un closure para que el timeoutId no se pierda entre llamadas
export function handleKeyUpWithTimeout(cb, time) {
    let timeoutId = null;
    return (keyUpEvent) => {
        if (timeoutId) clearTimeout(timeoutId);
        timeoutId = setTimeout(() => {
            cb(keyUpEvent);
        }, time);
    };
}
