import Header from "../components/Header";
import { Outlet } from "react-router-dom";

export default function Layout() {
    return <div className='pb-5'>
        <Header></Header>
        <Outlet/>
    </div>
}
