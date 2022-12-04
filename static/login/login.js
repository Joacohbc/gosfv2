import { getToken } from '/static/modules/request.js';
import { Message }  from '/static/modules/message.js';
const message = new Message("message");

if(getToken() != null) {
    // Si esta logueado, lo redirijo a la pagina principal
    window.location.href = '/static/main/files.html';
}

window.addEventListener('DOMContentLoaded', function() {

    const username = document.getElementById('username');
    username.addEventListener('keyup', function(e) {
        username.style = "color: white;";
    });

    const password = document.getElementById('password');
    password.addEventListener('keyup', function(e) {
        password.style = "color: white;";
    });

    // Agrego el evento click al botón de Login
    document.getElementById('btn-login').addEventListener('click', function(e) {
        e.preventDefault();

        let url = window.location.origin + '/login';

        axios.post(url,{
                username: username.value,
                password: password.value
        })
        .then(req => {
            window.location.href = '/static/main/files.html';
        })
        .catch(err => {
            message. showError(err.response.data.message);
        });
    });
});
