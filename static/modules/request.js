// Defino el URL base de las peticiones
axios.defaults.baseURL = window.location.origin;
axios.defaults.headers.common.Authorization = 'Bearer ' + localStorage.getItem('token');

if (localStorage.getItem('token') != null) {
    axios.get(window.location.origin + '/api/auth')
    .catch(err => {
        if(err.response.status === 401) {
            localStorage.removeItem('token');
            window.location.href = '/static/login/login.html';
        }
    });
}
