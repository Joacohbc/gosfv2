import './Header.css';
import Logo from '../assets/gosf-icon.png';
import { Link } from "react-router-dom";
import { useContext } from 'react';
import AuthContext from '../context/auth-context';


export default function Header() {
    const auth = useContext(AuthContext);

    const logoutHandler = () => {
        auth.onLogOut();
    };

    return <>   
    <div className="header-title">
        <img src={Logo} alt="server"/>
    </div>

    <header>
        <Link to="/files" className="header-btn">Files</Link>
        <Link to="/me" className="header-btn">User Info</Link>
        <button className="header-btn" onClick={logoutHandler}>Logout</button>
    </header>
    </>;
}
