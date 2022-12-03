import { showError, showSuccess, showInfo, setMessageId} from "/static/modules/message.js";

class User {
    constructor(id,  username) {
        this._id = id;
        this._username = username;
    }

    static fromJSON(json) {
        return new User(json.id, json.username);
    }

    get id() {
        return this._id;
    }

    get username() {
        return this._username;
    }    
}

// Obtiene el div con los Datos del Usuario (el Username)
// y el botón para eliminar el acceso al archivo)
function getUserDiv(user, file) {
    const userShare = document.createElement('div');
    userShare.classList.add("overlay-user-share");

    const username = document.createElement("div");
    username.innerHTML = user.username;
    userShare.appendChild(username);

    const btnDelete = document.createElement('button');
    btnDelete.id = "btn-delete-overlay";
    btnDelete.innerText = "Delete";
    btnDelete.addEventListener('click', () => {
        axios.delete(`/api/files/share/${file.id}/user/${user.id}`)
        .then(res => {
            showSuccess(res.data.message);
            reloadUsers(file);
        })
        .catch(err => {
            showError(err.response.data.message);
        });
    });
    userShare.appendChild(btnDelete);
    return userShare;
}

// Obtiene el div con los Usuarios con acceso al archivo
function createOverlayShare(file) {
    // Si no existe el overlayShare lo crea
    let overlayShare = document.querySelector(".overlay-share");
    if(overlayShare == null) {
        overlayShare = document.createElement('div');
        overlayShare.classList.add("overlay-share");
    } else {
        overlayShare.innerHTML = "";
    }

    axios.get(`/api/files/${file.id}/info`)
    .then(req => {
        
        // Si no hay archivos que ingrese un mensaje personalizado
        // en el head de la tabla que indique que no hay archivos
        if(req.data.shared_with == null) {
            overlayShare.innerHTML = "No users have access to this file";
            return;
        }

        // Y cargue las filas de la tabla
        const users = document.createDocumentFragment();
        req.data.shared_with.forEach(user => {
            users.appendChild(getUserDiv(User.fromJSON(user), file));
        });
        overlayShare.appendChild(users);

    }).catch(err => {
        console.log(err);
        showError(err.response.data.message);
    });

    return overlayShare;
}

// Recarga el div con los Usuarios con acceso al archivo
function reloadUsers(file) {
    // Si no existe el overlayShare lo crea
    let overlayShare = document.querySelector(".overlay-share");
    overlayShare.innerHTML = "";
    
    axios.get(`/api/files/${file.id}/info`)
    .then(req => {
        
        // Si no hay archivos que ingrese un mensaje personalizado
        // en el head de la tabla que indique que no hay archivos
        if(req.data.shared_with == null) {
            overlayShare.innerHTML = "No users have access to this file";
            return;
        }

        // Y cargue las filas de la tabla
        const users = document.createDocumentFragment();
        req.data.shared_with.forEach(user => {
            users.appendChild(getUserDiv(User.fromJSON(user), file));
        });
        overlayShare.appendChild(users);

    }).catch(err => {
        console.log(err);
        showError(err.response.data.message);
    });
}

// Crea el overlay para agregar Usuarios con acceso al archivos
function createOverlayBody(file) {
    //* Creo el body del overlay 
    const overlayBody = document.createElement('div');  
    overlayBody.classList.add("overlay-body");

    const inputUserId = document.createElement('input');
    inputUserId.id = "user-id";
    overlayBody.appendChild(inputUserId);

    const btnAddUser = document.createElement('button');
    btnAddUser.id = "btn-add-user-overlay";
    btnAddUser.innerText = "Add User";
    btnAddUser.addEventListener('click', () => {
        const userId = document.getElementById('user-id').value;
        if(userId == "") {
            showError("Please enter a User ID");
            return;
        }
    
        axios.post(`/api/files/share/${file.id}/user/${userId}`)
        .then(res => {
            showSuccess(res.data.message);
            createOverlayShare(file);
        })
        .catch(err => {
            console.log(err);
            showError(err.response.data.message);
        });
    });
    overlayBody.appendChild(btnAddUser);

    return overlayBody;
}

// Retorna el link de compartir
function getShareLink(file) {
    return `${window.location.origin}/api/files/share/${file.id}`;
}

// Crea el footer del overlay (el link de compartir y el checkbox para compartir con todos)
// Ademas del label de los mensajes de error y éxito
function createOverlayFooter(file) {
    //* Creo el footer del overlay
    const overlayFooter = document.createElement('div');
    overlayFooter.classList.add("overlay-footer");

    const lblCheckbox = document.createElement('label');
    lblCheckbox.for = "share-anyone";
    lblCheckbox.innerText = "Share with Anyone";
    overlayFooter.appendChild(lblCheckbox);

    const checkbox = document.createElement('input');
    checkbox.type = "checkbox";
    checkbox.id = "share-anyone";
    checkbox.checked = file.shared;
    checkbox.addEventListener('change', () => {
        file.update(file.filename, checkbox.checked);
    });

    overlayFooter.appendChild(checkbox);

    const inputLink = document.createElement('input');
    inputLink.type = "text";
    inputLink.id = "share-link";
    inputLink.value = getShareLink(file);
    inputLink.readOnly = true;
    inputLink.addEventListener('click', () => {
        navigator.clipboard.writeText(getShareLink(file))            
        .then(() => {
            showInfo("The link has been copied to the clipboard");
        })
        .catch(err => {
            console.log(err);
            showError("Error copying the link");
        });
    });
    overlayFooter.appendChild(inputLink);
    
    
    const btnClose = document.createElement('button');
    btnClose.id = "btn-close-overlay";
    btnClose.innerText = "Close";
    btnClose.addEventListener('click', () => {
        document.querySelector(".overlay-background").setAttribute("hidden", "");
        setMessageId("message");
    });
    overlayFooter.appendChild(btnClose);
    
    const lblMessage = document.createElement('label');
    lblMessage.id = "overlay-message";
    overlayFooter.appendChild(lblMessage);

    return overlayFooter;
}

// Crea el overlay y lo agrega al DOM
export function createOverlay(file) {
    const overlay = document.querySelector(".overlay-background");
    overlay.removeAttribute("hidden");

    const overlayContent = document.querySelector(".overlay-content");
    overlayContent.innerHTML = "";
    overlayContent.appendChild(createOverlayBody(file));
    overlayContent.appendChild(createOverlayShare(file));
    overlayContent.appendChild(createOverlayFooter(file));

    overlay.appendChild(overlayContent);
    setMessageId("overlay-message");
}

