import '../css/header.css';
import Logo from '../assets/gosf-icon.png';
import { Link } from "react-router-dom";


export default function Header() {
    return <>
    <div className="header-title">
        <img src={Logo} alt="server"/>
    </div>

    <header>
        <Link to="/files" className="header-btn">Files</Link>
        <Link to="/me" className="header-btn">User Info</Link>
        <button className="header-btn">Logout</button>
    </header>
    </>;
}
