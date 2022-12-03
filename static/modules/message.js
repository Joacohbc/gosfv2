// Objeto que contiene los colores de los mensajes
// Que se pueden cambiar de forma global
let colors = {
    info: 'blue',
    error: 'red',
    success: 'green',
    warning: 'yellow'
};

// Función que permite cambiar los colores de los mensajes
// Que se pueden cambiar de forma global
// Debe contendrá un objeto con los colores de los mensajes
// - info
// - error
// - success
// - warning
export function setColors(newColors) {
    colors = newColors;
}

// Variable que determina el ID (Global) del elemento que mostrará el mensaje
let messageId = "message";

// Función que permite cambiar el ID (Global) del elemento que mostrará el mensaje
// Debe contener el ID del elemento que mostrará el mensaje
export function setMessageId(id) {
    messageId = id;
}

export function showError(message, id = messageId) {
    showMessage(message, colors.error, id);
}

export function showSuccess(message, id = messageId) {
    showMessage(message, colors.success, id);
}

export function showWarning(message, id = messageId) {
    showMessage(message, colors.warning, id);
}

export function showInfo(message, id = messageId) {
    showMessage(message, colors.info, id);
}

function showMessage(message, color, id) {
    const element = document.getElementById(id);
    element.style.color = color;
    element.style.fontWeight = "bold";
    element.innerHTML = message[0].toUpperCase() + message.substring(1);
    setTimeout(() => {
        // Si el mensaje que se muestra es el mismo que se quiere ocultar
        // Ya que si hay otro mensaje (porque se llamo en otra instancia al método
        // que este no borre el mensaje que se quiere mostrar)
        if(message == element.innerHTML){
            element.innerHTML = '';
        }
    }, 5000);
}