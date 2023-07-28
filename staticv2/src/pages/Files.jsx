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

const emptyFile = Object.freeze({ id: null, filename: null, contentType: '', url: ''});

export default function Files() {
    const [ files, setFiles ] = useState([]);
    const [ previewFile, setPreviewFile ] = useState(emptyFile);
    const [ showPreview, setShowPreview ] = useState(false); 
    const [ uploading, setUploading ] = useState(false);
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
        setShowPreview(true);
        setPreviewFile({ 
            id: id, 
            filename: filename, 
            contentType: getContentType(filename),
            url: `${window.location.origin}/api/files/${id}?api-token=${auth.token}`
        });
    };

    const handleClosePreview = () => {
        setShowPreview(false);
        setPreviewFile(emptyFile);
    };
    
    const handleDelete = async(id, message, error) => {
        if(error) {
            messageRef.current.showError(error);
            return;
        }
        messageRef.current.showInfo(message);
        setFiles((files) => files.filter(file => file.id != id));
    };

    const handleFileUpload = (e) => {
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

    const previewComponent = () => {
        if(previewFile.contentType.includes('video'))
            return <video className='flex-fill' controls><source src={previewFile.url} type={previewFile.contentType}/></video>;
        
        return <iframe src={previewFile.url} className='flex-fill'/>;
    }
    
    return <>
        <Message ref={messageRef} />  
        <div className="loader file-loading" hidden={!uploading}> Uploading files </div> 

        {showPreview && 
        <Modal show={showPreview} onHide={handleClosePreview} className='d-flex modal-bg' fullscreen centered>
            <Modal.Header closeButton className='bg-modal' closeVariant='white'>{previewFile.filename}</Modal.Header>
            <div className='flex-fill'>
                { previewComponent() }
            </div>
        </Modal>}

        <div className="d-flex justify-content-center align-items-center mb-4">
            <input type="text" placeholder="Enter Search" className='search-input' onKeyUp={handleSearch}/>
        </div>
        
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
            <input id="input-upload" type="file" style={ {display: 'none'} } onChange={handleFileUpload} multiple/>
        </div>
    </>;
}
