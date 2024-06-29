import './Header.css';
import Logo from '../assets/gosf-icon.png';
import { Link, useNavigate } from "react-router-dom";
import { useContext } from 'react';
import AuthContext from '../context/auth-context';
import { MessageContext } from '../context/message-context';
import 'bootstrap-icons/font/bootstrap-icons.css'

export default function Header() {
    const auth = useContext(AuthContext);
    const messageContext = useContext(MessageContext);
    const navigate = useNavigate();

    const logoutHandler = () => {
        auth.onLogOut()
        .then(() => navigate('/login'))
        .catch(err => messageContext.showError(err.message));
    };

    return <>   
        <div className="header-title">
            <img src={Logo} alt="server"/>
        </div>

        <header>
            <Link to="/me" className="header-btn"><i className="bi bi-person-fill"/> Me</Link>
            <Link to="/files" className="header-btn"><i className="bi bi-archive-fill"/> Files</Link>
            <Link to="/notes" className="header-btn"><i className="bi bi-pen-fill"/> Notes</Link>
            <button className="header-btn" onClick={logoutHandler}><i className='bi bi-door-open-fill'/></button>
        </header>
    </>;
}
