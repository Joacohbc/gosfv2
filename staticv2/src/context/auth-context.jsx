import { createContext, useEffect, useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import PropTypes from 'prop-types';
import axios from "axios";

const AuthContext = createContext({
    token: '',
    isLogged: false,
    cAxios: null,
    onLogOut: async () => {},
    onLogin: async (username, password) => {}
});

const setAuthData = (token, duration) => {
    localStorage.setItem('token', token);
    localStorage.setItem('duration', duration);
    
    // Minutes to milliseconds
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

        // If token is set not redirect to login
        if (token) {
            setIsLogged(true);
            setToken(token);
            setCAxios(() => axios.create({
                baseURL: window.location.origin,
                headers: {
                    Authorization: `Bearer ${token}`
                }
            }));

            if(currentRoute == '/login' || currentRoute == '/register') {
                navigate("/files");
            }
            
            // Refresh token
            const expire = Date.parse(localStorage.getItem('expires'));
            if(expire - Date.now() <= 0) {
                navigate("/login");
                resetAuthData();
            }
            return;
        } 
        
        // If token is not set and current route is login or register not redirect to login
        if(currentRoute == '/login' || currentRoute == '/register') return;
        
        // If token is not set redirect to login
        navigate("/login");
        resetAuthData();
    }, [ token, currentRoute, navigate ]);

    const loginHandler = async (username, password) => {
        let url = window.location.origin + '/auth/login';

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
            axios.delete(window.location.origin + '/auth/logout');
        } catch(e) {
            console.log("Error on logout:" + e);
        }

        resetAuthData();
        navigate('/login');
    };

    return <AuthContext.Provider value={{
        token: token,
        isLogged: isLogged,
        onLogin: loginHandler,
        onLogOut: logOutHandler,
        cAxios: cAxios
    }}>{props.children}</AuthContext.Provider>
};

AuthContextProvider.propTypes = {
    children: PropTypes.node.isRequired,
};

export default AuthContext; 