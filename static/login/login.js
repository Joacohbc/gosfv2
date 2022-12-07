import { Message }  from '/static/modules/message.js';
const message = new Message("message");

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

        let url = window.location.origin + '/auth/login?cookie=true';
        
        axios.post(url,{
                username: username.value,
                password: password.value
        })
        .then(req => {
            localStorage.setItem('token', req.data.token);
            localStorage.setItem('duration', req.data.duration);

            let expires = new Date(Date.now() + req.data.duration * 1000 * 60);
            localStorage.setItem('expires', expires);
            
            window.location.href = '/static/main/files.html';
        })
        .catch(err => {
            console.log(err);
            message. showError(err.response.data.message);
        });
    });
});
