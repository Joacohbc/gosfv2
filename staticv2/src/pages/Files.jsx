import 'bootstrap/dist/css/bootstrap.min.css';
import './Files.css';
import Container from 'react-bootstrap/Container';
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import FileItem from '../components/FileItem';
import AuthContext from '../context/auth-context';
import Modal from 'react-bootstrap/Modal';
import { useCallback, useContext, useEffect,  useRef,  useState } from 'react';
import Message from '../components/Message';
import getContentType from '../utils/content-types';

export default function Files() {
    const [ files, setFiles ] = useState([]);
    const [ previewFile, setPreviewFile ] = useState({ id: null, filename: null});
    const [ showPreview, setShowPreview ] = useState(false);
    const [ uploading, setUploading ] = useState(false);
    const [ contentType, setContentType ] = useState(''); 
    const messageRef = useRef(null);

    const auth = useContext(AuthContext);
    const { isLogged, cAxios } = auth;

    const fetchDataFiles = useCallback(async (cb) => {
        try {
            const res = await cAxios.get('/api/files/');
            if(!res.data) return [];
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
    
    const handleOpenPreview = (id, filename) => {
        setContentType(getContentType(filename));
        setShowPreview(true);
        setPreviewFile(() => ({ id: id, filename: filename }));
    };

    const handleClosePreview = () => {
        setShowPreview(false);
        setPreviewFile(null);
    };
    
    const handleDelete = async(id, message) => {
        if(!id) {
            messageRef.current.showError(message);
            return;
        }
        messageRef.current.showInfo(message);
        fetchDataFiles((data) => setFiles(data));

    };

    const hanadleFileUpload = (e) => {
        e.preventDefault();

        const files = e.target.files;
        if(files.length == 0) {
            return;
        }

        const form = new FormData();
        for(let i = 0; i < files.length; i++) {
            form.append('files', files[i]);
        }

        setUploading(true);
        setTimeout(async () => {
            try {
                const res = await cAxios.post('/api/files/', form);
                fetchDataFiles((data) => setFiles(data));
                messageRef.current.showSuccess(res.data.message);
            } catch(err) {
                messageRef.current.showError(err.data.message);
            } finally {
                setUploading(false);
            }
        }, 0);
    }

    return <>
        {showPreview && 
        <Modal show={showPreview} onHide={handleClosePreview} className='d-flex modal-bg' fullscreen centered>
            <Modal.Header closeButton className='bg-modal' closeVariant='white'>{previewFile.filename}</Modal.Header>
            { !contentType.includes('video') && <iframe src={window.location.origin + '/api/files/' + previewFile.id +'?api-token='+auth.token} className='flex-fill'/>}
            
            { contentType.includes('video') && 
                <video className='flex-fill' controls>
                    <source src={window.location.origin + '/api/files/' + previewFile.id +'?api-token='+auth.token} type={contentType}/>
                </video>}
        </Modal>}

        <div className="d-flex justify-content-center align-items-center mb-4">
            <input type="text" placeholder="Enter Search" className='search-input' onKeyUp={handleSearch}/>
        </div>
        
        <div className="loader file-loading" hidden={!uploading}> Uploading files </div> 
        <Message ref={messageRef} />

        <Container fluid="md">
            {files.length == 0 && <Col>
                    <div className="d-flex justify-content-center align-items-center">
                        <p className="text-center text-white">No files, start uploading files c:</p>
                    </div>
            </Col>}

            {files.length != 0 && <Row xs={1} sm={2}  md={3} lg={4} xl={5} className='row-gap-3 d-flex justify-content-center'>
                {files.map(file => <Col key={file.id}>
                    <FileItem 
                        filename={file.filename} 
                        id={file.id} 
                        onOpen={handleOpenPreview}
                        onDelete={handleDelete}/>
                </Col>)}
            </Row>}
        </Container>
        
        <div className="d-flex justify-content-center align-items-center mt-4">
            <label htmlFor="input-upload" className="btn-upload">Upload file/s</label>
            <input id="input-upload" type="file" style={ {display: 'none'} } onChange={hanadleFileUpload} multiple/>
        </div>
    </>;
}
