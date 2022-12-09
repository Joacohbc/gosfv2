window.addEventListener('DOMContentLoaded', function() {
    // Agrego el header a la página al Inicio del Body
    document.body.innerHTML = `
    <div class="header-title">
        <img src="/static/images/gosf-icon.png" alt="server">
    </div>

    <header>
        <button id="btn-files" class="header-btn">Files</button>
        <button id="btn-user-info" class="header-btn">User Info</button>
        <button id="btn-logout" class="header-btn">Logout</button>
    </header>
    ` + document.body.innerHTML;

    // Agrego funcionamiento a los botones
    document.querySelector("#btn-logout").addEventListener('click', (e) => {
        e.preventDefault();
        axios.delete('/auth/logout')
        .then(res => {
            window.location.href = '/static/login/login.html';
        })
        .catch(err => {
            localStorage.removeItem('token');
            window.location.href = '/static/login/login.html';
        });
    });

    document.querySelector("#btn-files").addEventListener('click', (e) => {
        e.preventDefault();
        window.location.href = '/static/main/files.html';
    });

    document.querySelector("#btn-user-info").addEventListener('click', (e) => {
        e.preventDefault();
        window.location.href = '/static/user/users.html';
    });
});