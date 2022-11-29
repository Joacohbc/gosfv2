const colors = {
    info: 'blue',
    error: 'red',
    success: 'green',
    warning: 'yellow'
};

export function showError(message, id = "#message") {
    showMessage(message, colors.error, id);
}

export function showSuccess(message, id = "#message") {
    showMessage(message, colors.success, id);
}

export function showWarning(message, id = "#message") {
    showMessage(message, colors.warning, id);
}

export function showInfo(message, id = "#message") {
    showMessage(message, colors.info, id);
}


function showMessage(message, color, id) {
    const element = document.querySelector(id);
    element.style.color = color;
    element.style.fontWeight = "bold";
    element.innerHTML = message[0].toUpperCase() + message.substring(1);
    setTimeout(() => {
        element.innerHTML = '';
    }, 5000);
}