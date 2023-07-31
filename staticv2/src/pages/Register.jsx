import '../css/register.css';
import Logo from '../assets/gosf-icon.png';
import { Link } from 'react-router-dom';
import Form from 'react-bootstrap/Form';
import { useRef, useContext } from 'react';
import AuthContext from '../context/auth-context';

export default function Register() {
    const username = useRef();
    const password = useRef();
    const { onRegister } = useContext(AuthContext);
    
    const registerHandler = (e) => {
        e.preventDefault();
        onRegister(username.current.value, password.current.value);
    };

    return <>
        <div className="logo-icon">
            <img src={Logo} alt="server" />
        </div>

        <div className="register"> 
        
            <h1 className="register-title">REGISTER</h1>
        
            <Form>
                <Form.Label className='login-label'>Username:</Form.Label>
                <Form.Control type="text" placeholder="Username" required ref={username}/>

                <Form.Label className='login-label'>Password:</Form.Label>
                <Form.Control type="password" placeholder="Password" required ref={password}/>

                <div className='mt-3 d-flex justify-content-center rounded'>
                    <button className='btn-login flex-fill' onClick={registerHandler}>Register</button>
                </div>
            </Form>

            <div className="link">
                <Link to="/login">{"Already have an account? Log in!"}</Link>
            </div>
        </div>
    </>;
}
