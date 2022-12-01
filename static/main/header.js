import { getToken } from '/static/modules/request.js';

window.addEventListener('DOMContentLoaded', function() {
    
    // Defino el URL base de las peticiones
    axios.defaults.baseURL = window.location.origin;

    // Agrego el evento click al botón de Login
    // Verifico que el Token este en la Cookie
    if (getToken() === null) {
        // Si no esta lo devuelvo al Login
        window.location.href = '/static/login/login.html';
    }

    document.body.innerHTML = `
    <div class="header-title">
        <img src="/static/images/gosf-icon.png" alt="server">
    </div>

    <header>
        <button id="btn-files" class="header-btn">Files</button>
        <button id="btn-user-info" class="header-btn">User Info</button>
        <button id="btn-share-info" class="header-btn">Share Info</button>
        <button id="btn-logout" class="header-btn">Logout</button>
    </header>
    ` + document.body.innerHTML;

    // Agrego funcionamiento a los botones
    document.querySelector("#btn-logout").addEventListener('click', (e) => {
        e.preventDefault();
        axios.post('/logout')
        .then(res => {
            window.location.href = '/static/index.html';
        })
        .catch(err => {
            alert(err.response.data.message);
        });
    });

    document.querySelector("#btn-files").addEventListener('click', (e) => {
        e.preventDefault();
        window.location.href = '/static/main/files.html';
    });

    document.querySelector("#btn-share-info").addEventListener('click', (e) => {
        e.preventDefault();
        window.location.href = '/static/main/share.html';
    });

    document.querySelector("#btn-user-info").addEventListener('click', (e) => {
        e.preventDefault();
        alert("Working in progress!");
    });
});