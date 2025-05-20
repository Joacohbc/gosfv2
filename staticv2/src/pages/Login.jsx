import './Login.css';
import 'bootstrap/dist/css/bootstrap.min.css';
import Logo from '../assets/gosf-icon.png';
import Form from 'react-bootstrap/Form';
import { useContext, useRef } from 'react';
import AuthContext from '../context/auth-context';
import { MessageContext } from '../context/message-context';
import { Link } from 'react-router-dom';
import ConfirmDialog from '../components/ConfirmDialog';

export default function Login() {
    const username = useRef();
    const password = useRef();
    const auth = useContext(AuthContext);
    const messageContext = useContext(MessageContext);
    const restoreTokenDialog = useRef();

    const loginHandler = (e) => {
        e?.preventDefault();
        auth.onLogin(username.current.value, password.current.value)
        .catch(err => messageContext.showError(err.message));
    };

    const restoreTokenHandler = () => {
        auth.onRestore(username.current.value, password.current.value)
        .then(message => messageContext.showSuccess(message))
        .catch(err => messageContext.showError(err.message));
    };
    
    const showRestoreTokenDialog = (e) => {
        e?.preventDefault();
        restoreTokenDialog.current.show();
    }

    return <>
        <ConfirmDialog
            title="Are you sure you want to close all session of your Account?"
            message="This will close all your session of your account. You will need to login again to use the application."
            onOk={restoreTokenHandler} ref={restoreTokenDialog}/>

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
                    <button className='btn-login flex-fill' onClick={showRestoreTokenDialog}>Restore Tokens</button>
                </div>
            </Form>
            
            {/* New Google Sign-In Button */}
            <div className="mt-3 d-flex justify-content-center">
                <a href="/auth/google/login" className="btn btn-outline-primary flex-fill" role="button">
                    Sign in with Google
                </a>
            </div>
            {/* End of New Google Sign-In Button */}

            <div className="link text-center mt-2"> {/* Added mt-2 for spacing */}
                <Link to="/register">{"Don't have an account yet? Sign up!"}</Link>
            </div>
        </div>
    </>
}
