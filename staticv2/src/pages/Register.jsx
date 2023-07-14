import '../css/register.css';
import Logo from '../assets/gosf-icon.png';
import { Link } from 'react-router-dom';

export default function Register() {
    return <>
        <div className="logo-icon">
            <img src={Logo} alt="server" />
        </div>

        <div className="register"> 
        
            <h1 className="register-title">REGISTER</h1>
        
            <form>
                <label className="form-label">Username:</label>
                <input type="text" className="form-control" id="username" placeholder="Enter username" required />
                
                <br/>

                <label>Password:</label>
                <input type="password" className="form-control" id="password" placeholder="Enter password" required />

                <br/>
                
                <button className="btn-menu" id="btn-register">Register</button>
            </form>        

            <div id="message"></div>

            <div className="link">
                <Link to="/login">Already have an account? Log in!</Link>
            </div>
        </div>
    </>;
}
