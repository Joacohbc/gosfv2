import { getToken } from '/static/modules/request.js';
import { Message }  from '/static/modules/message.js';
const message = new Message("message");

if(getToken() != null) {
    // Si esta logueado, lo redirijo a la pagina principal
    window.location.href = '/static/main/files.html';
}

window.addEventListener('DOMContentLoaded', function() {
    // Agrego el evento click al botón de Login
    this.document.getElementById('btn-register').addEventListener('click', function(e) {
        e.preventDefault();

        let url = window.location.origin + '/register';

        axios.post(url,{
                username: document.getElementById("username").value,
                password: document.getElementById("password").value
        })
        .then(req => {
            console.log(req);
            message.showSuccess(req.data.message);
        })
        .catch(err => {
            console.log(err);
            message.showError(err.response.data.message);
        });
    });
});
