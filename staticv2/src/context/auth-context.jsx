import { createContext, useCallback, useEffect, useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import PropTypes from 'prop-types';
import axios from "axios";

const AuthContext = createContext({
    token: '',
    isLogged: false,
    cAxios: null,
    baseUrl: '',
    onLogOut: async () => {},
    onLogin: async (username, password) => {},
    onRegister: async (username, password) => {},
    onRestore: async (username, password) => {},
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
    // const BASE_URL = 'http://localhost:3000';
    const BASE_URL = window.location.origin;
    
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
                baseURL: BASE_URL,
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
    }, [ token, currentRoute, navigate, BASE_URL ]);

    const loginHandler = async (username, password) => {
        try {
            const req = await axios.post(BASE_URL + '/auth/login', {
                username: username,
                password: password
            });

            setAuthData(req.data.token, req.data.duration);
            setToken(req.data.token);
            setIsLogged(true);
            return Promise.resolve('User logged in successfully');
        } catch(err) {
            return Promise.reject(new Error(err.response.data.message));
        }
    };

    const restoreTokenHandler = async (username, password) => {
        try {
            await axios.post(BASE_URL + '/auth/restore', {
                username: username,
                password: password
            });
            resetAuthData();
            return Promise.resolve('All Tokens are restore successfully');
        } catch(err) {
            return Promise.reject(new Error(err.response.data.message));
        }
    };

    const logOutHandler = async () => {
        try {
            await axios.delete(BASE_URL + '/auth/logout', {
                headers: {
                    Authorization: `Bearer ${token}`
                }
            });
            return Promise.resolve('User logged out successfully');
        } catch(err) {
            return Promise.reject(new Error(err.response.data.message));
        } finally {
            resetAuthData();
            navigate('/login');
        }
    };

    const registerHandler = async (username, password) => {
        try {
            await axios.post(BASE_URL + '/auth/register', {
                username: username,
                password: password
            });
            return Promise.resolve('User created successfully');
        } catch(err) {
            return Promise.reject(new Error(err.response.data.message));
        }
    };
    
    const addTokenParam = (url) => {
        const urlObj = new URL(url);
        urlObj.searchParams.append('api-token', token);
        return urlObj.toString();
    };

    return <AuthContext.Provider value={{
        token: token,
        isLogged: isLogged,
        cAxios: cAxios,
        baseUrl: BASE_URL,
        onLogin: loginHandler,
        onLogOut: logOutHandler,
        onRegister: registerHandler,
        onRestore: restoreTokenHandler,    
        addTokenParam: addTokenParam
    }}>{props.children}</AuthContext.Provider>
};

AuthContextProvider.propTypes = {
    children: PropTypes.node.isRequired,
};

export default AuthContext; 