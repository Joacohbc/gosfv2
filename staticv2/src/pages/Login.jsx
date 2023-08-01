import '../css/login.css';
import 'bootstrap/dist/css/bootstrap.min.css';
import Logo from '../assets/gosf-icon.png';
import Form from 'react-bootstrap/Form';
import { useContext, useRef } from 'react';
import AuthContext from '../context/auth-context';
import { MessageContext } from '../context/message-context';
import { Link } from 'react-router-dom';

export default function Login() {
    const username = useRef();
    const password = useRef();
    const auth = useContext(AuthContext);
    const messageContext = useContext(MessageContext);
    
    const loginHandler = (e) => {
        e.preventDefault();
        auth.onLogin(username.current.value, password.current.value)
        .catch(err => messageContext.showError(err.message));
    };
    
    return <>
        <div className="logo-icon">
            <img src={Logo} alt="server"/>
        </div>

        <div className="login"> 
            <h1 className="login-title">LOGIN</h1>
            <Form>
                <Form.Label className='login-label'>Username:</Form.Label>
                <Form.Control type="text" placeholder="Username" required ref={username}/>

                <Form.Label className='login-label'>Password:</Form.Label>
                <Form.Control type="password" placeholder="Password" required ref={password}/>

                <div className='mt-3 d-flex justify-content-center rounded'>
                    <button className='btn-login flex-fill' onClick={loginHandler}>Login</button>
                    <button className='btn-login flex-fill'>Restore Tokens</button>
                </div>
            </Form>
            
            <div className="link">
                <Link to="/register">{"Don't have an account yet? Sign up!"}</Link>
            </div>

        </div>
    </>
}
