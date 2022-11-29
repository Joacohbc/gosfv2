import { getToken } from '/static/modules/request.js';
import { showError}  from '/static/modules/message.js';

if(getToken() != null) {
    // Si esta logueado, lo redirijo a la pagina principal
    window.location.href = '/static/main/files.html';
}

window.addEventListener('DOMContentLoaded', function() {
    // Agrego el evento click al botón de Login
    this.document.querySelector('#btn-login').addEventListener('click', function(e) {
        e.preventDefault();

        let url = window.location.origin + '/login';

        axios.post(url,{
                username: document.querySelector("#username").value,
                password: document.querySelector("#password").value
        })
        .then(req => {
            window.location.href = '/static/main/files.html';
        })
        .catch(err => {
            showError(err.response.data.message);
        });
    });
});
