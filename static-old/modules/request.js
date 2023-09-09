import { logout } from "./Logout.js";

function refreshToken() {
    
    // Calculo el tiempo que falta para que expire el token
    const expire = Date.parse(localStorage.getItem('expires'));

    // Calculo el tiempo de espera que falta para que expire el token
    const wait = (expire - 1000 * 60 * 1) - Date.now();

    // Si el tiempo de espera es menor o igual a 0, significa que el token ya expiró
    if(wait <= 0) {
        window.location.href = '/static/login/login.html';
    }

    // Refresco el token
    setTimeout(() => {
        axios.get('/auth/refresh?cookie=true')
        .then(res => {
            // Guardo el token en el localStorage
            localStorage.setItem('token', res.data.token);
            localStorage.setItem('duration', res.data.duration);

            // Guardo la fecha de expiración del token (en el localStorage)
            let expires = new Date(Date.now() + res.data.duration * 1000 * 60);
            localStorage.setItem('expires', expires);

            // Agrego el token a las peticiones
            axios.defaults.headers.common.Authorization = 'Bearer ' + res.data.token;
            console.log('Token refreshed');
            refreshToken();
        })
        .catch(err => {
            console.log(err);
            if(err.response.status === 401) {
                localStorage.removeItem("token");
                localStorage.removeItem("duration");
                window.location.href = '/static/login/login.html';
            }
        });
    }, wait);    
}

if (localStorage.getItem('token') != null) {
    // Defino el URL base de las peticiones
    axios.defaults.baseURL = window.location.origin;
    axios.defaults.headers.common.Authorization = 'Bearer ' + localStorage.getItem('token');
    
    axios.get('/auth/verify')
    .then(res => {
        refreshToken();
    })
    .catch(err => {
        if(err.response.status === 401) {
            logout();
        }
    });

} else {
    window.location.href = '/static/login/login.html';
}