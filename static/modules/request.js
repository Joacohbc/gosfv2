// Defino el URL base de las peticiones
axios.defaults.baseURL = window.location.origin;
axios.defaults.headers.common.Authorization = 'Bearer ' + localStorage.getItem('token');

if (localStorage.getItem('token') != null) {
    axios.get('/api/auth')
    .catch(err => {
        localStorage.removeItem('token');
        if(err.response.status === 401) window.location.href = '/static/login/login.html';
    });
}
