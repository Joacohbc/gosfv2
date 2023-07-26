import 'bootstrap/dist/css/bootstrap.min.css';
import './Files.css';
import Container from 'react-bootstrap/Container';
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import FileItem from '../components/FileItem';
import AuthContext from '../context/auth-context';
import { useCallback, useContext, useEffect,  useState } from 'react';

export default function Files() {
    const [ files, setFiles ] = useState([]);
    
    const auth = useContext(AuthContext);
    const { isLogged, cAxios } = auth;

    const fetchDataFiles = useCallback(async (cb) => {
        try {
            const res = await cAxios.get('/api/files/');
            if(res.data?.length == 0) return [];
            cb(res.data)
        } catch(err) {
            console.log(err);
            return [];
        }
    }, [ cAxios ]);
    
    useEffect(() => {
        if(!cAxios || !isLogged) return;
        fetchDataFiles((data) => setFiles(data));
    }, [ isLogged, cAxios, fetchDataFiles ]);

    let searchTimeout = null;
    const handleSearch = (e) => {
        if(searchTimeout) clearTimeout(searchTimeout);
        searchTimeout = setTimeout(() => {
            fetchDataFiles((data) => {
                setFiles(() => data.filter(file => file.filename.includes(e.target.value)));
            });
        }, 500);
    };
    
    return <>
        <div className="d-flex justify-content-center align-items-center mb-4">
            <input type="text" placeholder="Enter Search" className='search-input' onKeyUp={handleSearch}/>
        </div>
        
        <div className="loader file-loading" hidden> Uploading files </div> 
        <div className="message"></div>

        <Container fluid="md">
            {files.length == 0 && <Col>
                    <div className="d-flex justify-content-center align-items-center">
                        <p className="text-center text-white">No files, start uploading files c:</p>
                    </div>
            </Col>}

            {files.length != 0 && <Row xs={1} sm={2}  md={3} lg={4} xl={5} className='row-gap-3 d-flex justify-content-center'>
                {files.map(file => <Col key={file.id}>
                    <FileItem filename={file.filename} id={file.id}/>
                </Col>)}
            </Row>}
        </Container>
        
        <div className="d-flex justify-content-center align-items-center mt-4">
            <label htmlFor="input-upload" id="btn-upload" className="btn-upload">Upload file/s</label>
            <input id="input-upload" type="file" style={ {display: 'none'} } multiple/>
        </div>
    </>;
}
