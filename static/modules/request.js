
if (localStorage.getItem('token') != null) {
    // Defino el URL base de las peticiones
    axios.defaults.baseURL = window.location.origin;
    axios.defaults.headers.common.Authorization = 'Bearer ' + localStorage.getItem('token');
    // axios.defaults.params = {
    //     token: localStorage.getItem('token')
    // };
    
    axios.get(window.location.origin + '/auth/verify')
    .catch(err => {
        if(err.response.status === 401) {
            localStorage.removeItem('token');
            window.location.href = '/static/login/login.html';
        }
    });
} else {
    window.location.href = '/static/login/login.html';
}

