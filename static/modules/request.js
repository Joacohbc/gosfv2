export function getToken(){
    // Verifico si tengo el Token en el LocalStorage
    let token = null;
    document.cookie.split(';').forEach((cookie) => {
        if(cookie.includes('token')) {
            token = cookie.split('=')[1];
            return;
        }
    });
    return token;
}
