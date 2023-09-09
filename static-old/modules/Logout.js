export const logout = () => {
    // Aviso al servidor que el usuario se deslogueo
    axios.delete('/auth/logout?cookie=true');

    // Elimino los datos del usuario del localStorage
    localStorage.removeItem('token');
    localStorage.removeItem('duration');
    localStorage.removeItem('expires');
    localStorage.removeItem('user');
    localStorage.removeItem('userId');
    localStorage.removeItem('username');
    
    // Redirecciono al usuario al login
    window.location.href = '/static/login/login.html';
}