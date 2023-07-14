import Container from 'react-bootstrap/Container';
import Nav from 'react-bootstrap/Nav';
import Navbar from 'react-bootstrap/Navbar';
import NavDropdown from 'react-bootstrap/NavDropdown';
import Col from 'react-bootstrap/Col';
import Image from 'react-bootstrap/Image';
import Logo from '../assets/logo.png';

function Header() {
    return (
        <Navbar expand="lg" className="bg-dark">
            <Container className="gap-2">
                <Col xs={2} md={2}>
                    <Navbar.Brand href="#home">
                        <Image src={Logo} fluid/>
                    </Navbar.Brand>
                </Col>

                <Col className='justify-content-end'>
                </Col>  
            </Container>
        </Navbar>
    );
}

export default Header;
