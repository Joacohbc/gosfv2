import './Header.css';
import Logo from '../assets/gosf-icon.png';
import { Link } from "react-router-dom";
import { useContext } from 'react';
import AuthContext from '../context/auth-context';
import { MessageContext } from '../context/message-context';
import 'bootstrap-icons/font/bootstrap-icons.css'

export default function Header() {
    const auth = useContext(AuthContext);
    const messageContext = useContext(MessageContext);

    const logoutHandler = () => {
        auth.onLogOut()
        .catch(err => messageContext.showError(err.message));
    };

    return <>   
        <div className="header-title">
            <img src={Logo} alt="server"/>
        </div>

        <header>
            <Link to="/me" className="header-btn"><i className="bi bi-person-fill"/> Me</Link>
            <Link to="/files" className="header-btn"><i className="bi bi-archive-fill"/> My Files</Link>
            <button className="header-btn" onClick={logoutHandler}><i className='bi bi-door-open-fill'/> Bye!</button>
        </header>
    </>;
}
