// Objeto que contiene los colores de los mensajes
// Que se pueden cambiar de forma global
const colors = {
    info: 'blue',
    error: 'red',
    success: 'green',
    warning: 'yellow'
};

// Variable que determina el ID (Global) del elemento que mostrará el mensaje
const messageId = "message";

export class Message {
    constructor(id = messageId, color = colors) {
        this._id = id;
        this._color = color;
    }

    get id() {
        return this._id;
    }

    // Debe contener el ID del elemento que mostrará el mensaje
    set id(id) {
        this._id = id;
    }
    
    get color() {
        return this._color;
    }

    // Debe contendrá un objeto con los colores de los mensajes
    // - info
    // - error
    // - success
    // - warning
    set color(color) {
        if(this._color.info == undefined || this._color.error == undefined || this._color.success == undefined || this._color.warning == undefined) {
            throw new Error("The object must contain the this._color of the messages");
        }

        this._color = color;
    }

    showError(message) {
        this.showMessage(message, this._color.error, this._id);
    }
    
    showSuccess(message) {
        this.showMessage(message, this._color.success, this._id);
    }
    
    showWarning(message) {
        this.showMessage(message, this._color.warning, this._id);
    }
    
    showInfo(message) {
        this.showMessage(message, this._color.info, this._id);
    }
    
    showMessage(message, color) {
        const element = document.getElementById(this._id);
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
}