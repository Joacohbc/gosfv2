import 'bootstrap/dist/css/bootstrap.min.css';
import './Files.css';
import Container from 'react-bootstrap/Container';
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import FileItem from '../components/fileItem/FileItem';
import AuthContext from '../context/auth-context';
import Modal from 'react-bootstrap/Modal';
import { useCallback, useContext, useEffect,  useState } from 'react';
import { handleKeyUpWithTimeout } from '../utils/input-text';
import PreviewFile from './PreviewFile';
import { MessageContext } from '../context/message-context';
import { useGetInfo } from '../hooks/files';

const emptyFile = Object.freeze({ id: null, filename: null, contentType: '', url: '', extesion: ''});

export default function Files() {
    const [ files, setFiles ] = useState([]);
    const [ previewFile, setPreviewFile ] = useState(emptyFile);
    const [ showPreview, setShowPreview ] = useState(false); 
    const [ uploading, setUploading ] = useState(false);
    const { isLogged, cAxios } = useContext(AuthContext);
    const messageContext = useContext(MessageContext);
    const { getFileInfo } = useGetInfo();

    const fetchDataFiles = useCallback(async (cb) => {
        try {
            const res = await cAxios.get('/api/files/');
            if(!res.data) return [];
            const files = res.data.map(file => getFileInfo(file, true));
            cb(files);
        } catch(err) {
            messageContext.showError(err.response.data.message);
            return [];
        }
    }, [ cAxios, messageContext, getFileInfo ]);
    
    useEffect(() => {
        if(!cAxios || !isLogged) return;
        fetchDataFiles((data) => setFiles(data));
    }, [ isLogged, cAxios, fetchDataFiles ]);

    const handleSearch = handleKeyUpWithTimeout((e) => {
        fetchDataFiles((data) => setFiles(data.filter(file => file.filename.includes(e.target.value))));
    }, 200);

    const handleDelete = useCallback(async(openedFile) => {
        setFiles((files) => files.filter(file => file.id != openedFile.id));
    }, [ ]);
    
    const handleUpdate = useCallback((openedFile) => {
        setFiles((files) => files.map(file => file.id == openedFile.id ? openedFile : file));
    }, [ ]);

    const handleOpenPreview = useCallback((file) => {
        setPreviewFile(file);
        setShowPreview(true);
    }, [ ]);

    const handleClosePreview = () => {
        setShowPreview(false);
        setPreviewFile(emptyFile);
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
                messageContext.showSuccess(res.data.message);
            } catch(err) {
                messageContext.showError(err.response.data.message);
            } finally {
                setUploading(false);
            }
        }, 0);
    }

    return <>
        <div className="loader file-loading" hidden={!uploading}> Uploading files </div> 

        {showPreview && 
        <Modal show={showPreview} onHide={handleClosePreview} className='d-flex modal-bg' fullscreen centered>
            <Modal.Header closeButton className='bg-modal' closeVariant='white'>{previewFile.filename}</Modal.Header>
            <div className='flex-fill'>
                <PreviewFile fileInfo={previewFile} className={"flex-fill"} />
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

            <Row xs={1} sm={2}  md={3} lg={4} xl={5} className='row-gap-3 d-flex justify-content-center'>
                {files.map(file => 
                <Col key={file.id}>
                    <FileItem 
                        id={file.id}
                        filename={file.filename}
                        name={file.name} 
                        contentType={file.contentType}
                        extesion={file.extesion}
                        url={file.url}
                        onOpen={handleOpenPreview}
                        onDelete={handleDelete}
                        onUpdate={handleUpdate}
                        />
                </Col>)}
            </Row>
        </Container>
        
        <div className="d-flex justify-content-center align-items-center mt-4">
            <label htmlFor="input-upload" className="btn-upload">Upload file/s</label>
            <input id="input-upload" type="file" style={ {display: 'none'} } onChange={handleFileUpload} multiple/>
        </div>
    </>;
}
