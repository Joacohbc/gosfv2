import '../css/login.css';
import Logo from '../assets/gosf-icon.png';
import { Link } from 'react-router-dom';

export default function Login() {
    return <>
        <div className="logo-icon">
            <img src={Logo} alt="server"/>
        </div>

        <div className="login"> 
            
            <h1 className="login-title">LOGIN</h1>
            
            <form>
                <label className="form-label">Username:</label>
                <input type="text" className="form-control" id="username" placeholder="Enter username" required/>
                
                <br/>

                <label>Password:</label>
                <input type="password" className="form-control" id="password" placeholder="Enter password" required/>
                <br/>
                <button id="btn-login">Login</button>            
                <button id="btn-restore">Restore Tokens</button>
            </form>

            <div id="message"></div>

            <div className="link">
                <Link to="/register">{"Don't have an account yet? Sign up!"}</Link>
            </div>
        </div>
    </>;
}
