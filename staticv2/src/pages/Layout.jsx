import Header from "../components/Header";
import { Outlet } from "react-router-dom";
import MessageComponentProvider from "../context/message-context";

export default function Layout() {
    return <MessageComponentProvider>
        <Header></Header>
        <Outlet/>
    </MessageComponentProvider>;
}
