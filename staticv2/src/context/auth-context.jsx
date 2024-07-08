import { createContext, useCallback, useEffect, useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import PropTypes from 'prop-types';
import axios from "axios";
import { useCache } from "../hooks/cache";

const AuthContext = createContext({
    token: '',
    isLogged: false,
    baseUrl: '',
    onLogOut: async () => {},
    onLogin: async (username, password) => {},
    onRegister: async (username, password) => {},
    onRestore: async (username, password) => {}
});

const deleteAllCookies = () => {
    document.cookie.split(';').forEach(cookie => {
        const eqPos = cookie.indexOf('=');
        const name = eqPos > -1 ? cookie.substring(0, eqPos) : cookie;
        document.cookie = name + '=;expires=Thu, 01 Jan 1970 00:00:00 GMT';
    });
}

export const AuthContextProvider = (props) => {
    const BASE_URL = window.location.origin;
    const { cacheService } = useCache();
    
    const navigate = useNavigate();
    const location = useLocation();

    const [ isLogged, setIsLogged ] = useState(false);
    const [ token, setToken ] = useState('');

    const { pathname: currentRoute } = location;
    
    const clearSession = useCallback(() => {
        setToken('');
        setIsLogged(false);
        cacheService.clean();
        deleteAllCookies();
    }, [ cacheService]);

    const refreshTokenHandler = useCallback(async () => {
        try {
            const req = await axios.get(BASE_URL + '/auth/refresh?cookie=yes');

            return {
                token: req.data.token,
                duration: req.data.duration
            };
        } catch(err) {
            throw new Error(err.response.data.message);
        }
    }, [ BASE_URL ]);

    const logOutHandler = useCallback(async () => {
        try {
            await axios.delete(BASE_URL + '/auth/logout?cookie=yes');
            return 'User logged out successfully';
        } catch(err) {
            throw new Error(err.response.data.message);
        } finally {
            clearSession();
            navigate('/login');
        }
    }, [BASE_URL, navigate, clearSession]);

    const verifyToken = useCallback(async () => {
        try {
            const req = await axios.get(BASE_URL + '/auth/verify');

            return {
                token: req.data.token,
                durationRemaining: req.data.durationRemaining
            }
        } catch(err) {
            throw new Error(err.response.data.message);
        }
    }, [ BASE_URL ]);

    useEffect(() => {
        if(currentRoute == '/register') return;

        if(isLogged) {
            if(currentRoute == '/login' || currentRoute == '/register' || currentRoute == '/') {
                navigate('/files');
            }
            return;
        }

        verifyToken()
        .then((currentTokenInfo) => {
            setIsLogged(true);
            setToken(currentTokenInfo.token);
            
            if(currentRoute == '/login' || currentRoute == '/register' || currentRoute == '/') {
                navigate('/files');
            }
            
            if(currentTokenInfo.durationRemaining <= 10) {
                refreshTokenHandler()
                .then((newTokenInfo) => {
                    setToken(newTokenInfo.token);
                    setIsLogged(true);
                }).catch(() => {
                    logOutHandler();
                    navigate('/login');
                });
            }
        })
        .catch(() => {
            clearSession();
            navigate('/login');
        });
    }, [token, currentRoute, navigate, BASE_URL, verifyToken, logOutHandler, isLogged, refreshTokenHandler, clearSession]);

    const loginHandler = async (username, password) => {
        try {
            const req = await axios.post(BASE_URL + '/auth/login?cookie=yes', {
                username: username,
                password: password
            });

            setToken(req.data.token);
            setIsLogged(true);
            cacheService.clean();
            navigate('/files');
            return 'User logged in successfully';
        } catch(err) {
            throw new Error(err.response.data.message);
        }
    };

    const restoreTokenHandler = async (username, password) => {
        try {
            await axios.post(BASE_URL + '/auth/restore?cookie=yes', {
                username: username,
                password: password
            });
            return 'All Tokens are restore successfully';
        } catch(err) {
            throw new Error(err.response.data.message);
        }
    };

    const registerHandler = async (username, password) => {
        try {
            await axios.post(BASE_URL + '/auth/register', {
                username: username,
                password: password
            });
            clearSession();
            return 'User created successfully';
        } catch(err) {
            throw new Error(err.response.data.message);
        }
    };
    
    return <AuthContext.Provider value={{
        token: token,
        isLogged: isLogged,
        baseUrl: BASE_URL,
        onLogin: loginHandler,
        onLogOut: logOutHandler,
        onRegister: registerHandler,
        onRestore: restoreTokenHandler
    }}>{props.children}</AuthContext.Provider>
};

AuthContextProvider.propTypes = {
    children: PropTypes.node.isRequired,
};

export default AuthContext; 