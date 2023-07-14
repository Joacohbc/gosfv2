import { BrowserRouter, Routes, Route } from "react-router-dom";
import Layout from "./pages/Layout";
import NoPage from "./pages/NoPage";
import Files from "./pages/Files";
import User from "./pages/User";
import Login from "./pages/Login";
import Register from "./pages/Register";
import './css/index.css';

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<Login/>}/>
        <Route path="/register" element={ <Register/> } />
        <Route path="/" element={<Layout />}>
          <Route path="me" element={<User />} />
          <Route path="files" element={< Files/>} />
          <Route path="*" element={<NoPage />} />
        </Route>
      </Routes>
    </BrowserRouter>
  );
}

export default App
