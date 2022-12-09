
function refreshToken() {
    const expire = Date.parse(localStorage.getItem('expires'));
    const wait = (expire - 1000 * 60 * 1) - Date.now();

    if(wait <= 0) {
        window.location.href = '/static/login/login.html';
    }

    setTimeout(() => {
        axios.get('/auth/refresh?cookie=true')
        .then(res => {
            localStorage.setItem('token', res.data.token);
            localStorage.setItem('duration', res.data.duration);

            let expires = new Date(Date.now() + res.data.duration * 1000 * 60);
            localStorage.setItem('expires', expires);

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
        console.log(err);
        if(err.response.status === 401) {
            localStorage.removeItem('token');
            window.location.href = '/static/login/login.html';
        }
    });
} else {
    window.location.href = '/static/login/login.html';
}