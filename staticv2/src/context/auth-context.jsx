import { createContext, useEffect, useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import PropTypes from 'prop-types';
import axios from "axios";

const AuthContext = createContext({
    token: '',
    isLogged: false,
    cAxios: null,
    onLogOut: async () => {},
    onLogin: async (username, password) => {},
    onRegister: async (username, password) => {},
    addTokenParam: (url) => {},
});

const setAuthData = (token, duration) => {
    localStorage.setItem('token', token);
    localStorage.setItem('duration', duration);
    let expires = new Date(Date.now() + duration * 1000 * 60);
    localStorage.setItem('expires', expires);
}

const resetAuthData = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('duration');
    localStorage.removeItem('expires');
}

export const AuthContextProvider = (props) => { 
    const navigate = useNavigate();
    const location = useLocation();

    const [ isLogged, setIsLogged ] = useState(false);
    const [ token, setToken ] = useState('');
    const [ cAxios, setCAxios ] = useState(null);

    const { pathname: currentRoute } = location;

    useEffect(() => {
        const token = localStorage.getItem('token');

        // Si el token esta seteado y la ruta actual es login o register redireccionar a files
        if (token) {
            setIsLogged(true);
            setToken(token);
            setCAxios(() => axios.create({
                // baseURL: window.location.origin,
                baseURL: 'http://localhost:3000',
                headers: {
                    Authorization: `Bearer ${token}`
                }
            }));
            
            if(currentRoute == '/login' || currentRoute == '/register') {
                navigate("/files");
            }
            
            // Si el token expiro redireccionar a login
            const expire = Date.parse(localStorage.getItem('expires'));
            if(expire - Date.now() <= 0) {
                navigate("/login");
                resetAuthData();
            }
            return;
        } 
        
        // Si el token no esta seteado y la ruta actual es login o register no redireccionar a login
        if(currentRoute == '/login' || currentRoute == '/register') return;
        
        // Si el token no esta seteado redireccionar a login
        navigate("/login");
        resetAuthData();
    }, [ token, currentRoute, navigate ]);

    const loginHandler = async (username, password) => {
        // const url = window.location.origin + '/auth/login';
        let url = 'http://localhost:3000' + '/auth/login';

        try {
            const req = await axios.post(url, {
                username: username,
                password: password
            });

            setAuthData(req.data.token, req.data.duration);
            setToken(req.data.token);
            setIsLogged(true);
        } catch(e) {
            console.log('Ocurrio un error');
            console.log(e);
        }
    };

    const logOutHandler = async () => {
        try {
            // axios.delete(window.location.origin + '/auth/logout');
            axios.delete('http://localhost:3000' + '/auth/logout');
        } catch(e) {
            console.log("Error on logout:" + e);
        }

        resetAuthData();
        navigate('/login');
    };

    const registerHandler = async (username, password) => {
        // const url = window.location.origin + '/auth/register';
        let url = 'http://localhost:3000' + '/auth/register';
        
        try {
            const req = await axios.post(url, {
                username: username,
                password: password
            });
        } catch(e) {
            console.log('Ocurrio un error');
            console.log(e);
        }
    };
    
    const addTokenParam = (url) => {
        return url + `?api-token=${token}`;
    };

    return <AuthContext.Provider value={{
        token: token,
        isLogged: isLogged,
        cAxios: cAxios,
        onLogin: loginHandler,
        onLogOut: logOutHandler,
        onRegister: registerHandler,    
        addTokenParam: addTokenParam
    }}>{props.children}</AuthContext.Provider>
};

AuthContextProvider.propTypes = {
    children: PropTypes.node.isRequired,
};

export default AuthContext; 