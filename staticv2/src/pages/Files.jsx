import 'bootstrap/dist/css/bootstrap.min.css';
import './Files.css';
import Container from 'react-bootstrap/Container';
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import FileItem from '../components/fileItem/FileItem';
import AuthContext from '../context/auth-context';
import Modal from 'react-bootstrap/Modal';
import { useCallback, useContext, useEffect,  useRef,  useState, useReducer } from 'react';
import { handleKeyUpWithTimeout } from '../utils/input-text';
import PreviewFile from './PreviewFile';
import { MessageContext } from '../context/message-context';
import { useGetInfo, useFiles } from '../hooks/files';
import SpinnerDiv from '../components/SpinnerDiv';
import useJobsQueue from '../hooks/jobsQueue';
import '../components/Message.css';

const emptyFile = Object.freeze({ id: null, filename: null, contentType: '', url: '', extesion: '', deleted: false });

export default function Files() {
    const messageContext = useContext(MessageContext);
    const { isLogged, cAxios } = useContext(AuthContext);

    const searching = useRef(false);
    
    const [ files, setFiles ] = useState([]);
    const [ progress, setProgress ] = useState(0);
    const [ fileLoader, setFileLoader ] = useState(() => {});
    const [ uploading, setUploading ] = useState(false);
    
    const { getFilenameInfo } = useGetInfo();
    const { getFiles, uploadFile } = useFiles();
    const { addJob, undoLastJob, jobsQueue } = useJobsQueue(5000);

    const previewReducer = (state, action) => {
        switch (action.type) {
            case 'SET_PREVIEW_FILE':
                return { ...state, previewFile: action.payload };
            case 'SHOW_PREVIEW':
                return { ...state, showPreview: true };
            case 'HIDE_PREVIEW':
                return { ...state, showPreview: false };
            default:
                return state;
        }
    };

    const [ state , setPreview ] = useReducer(previewReducer, { previewFile: emptyFile, showPreview: false });
    const { previewFile, showPreview } = state;

    // Genera una funcion custom que realiza una carga progresiva de los archivos mediante un closure
    const createFileLoader = useCallback(async (filterCb = (data) => data) => {
        try {
            const files = await getFiles();
            let data = filterCb(files.map(file => {
                const f = getFilenameInfo(file, true);
                f.deleted = false;
                return f;
            }));
            
            // Carga de 5 en 5 archivos o todos los archivos si son menos de 5
            const numberOfFilesPerLoad = data.length >= 5 ? 5 : data.length;
            setFiles(data);
            setProgress(0);
            return () => setProgress(prevProgress => prevProgress >= data.length ? data.length : prevProgress + numberOfFilesPerLoad);
        } catch(err) {
            messageContext.showError(err.message);
            return [];
        }
    }, [ messageContext, getFilenameInfo, getFiles ]);
    
    useEffect(() => {
        if(!cAxios || !isLogged) return;
    
        searching.current.loading();
        createFileLoader().then((loadInfo) => {
            setFileLoader(() => loadInfo);
            loadInfo();
        }).finally(() => searching.current.stopLoading());
    }, [ isLogged, cAxios, createFileLoader ]);

    const handleSearch = handleKeyUpWithTimeout((e) => {
        searching.current.loading();

        const filterCb = (data) => data.filter(file => file.filename.toLowerCase().includes(e.target.value.toLowerCase()) || file.id == e.target.value);
        createFileLoader(filterCb).then((loadInfo) => {
            setFileLoader(() => loadInfo);
            loadInfo();
        }).finally(() => searching.current.stopLoading());
    }, 500);

    const handleDelete = useCallback(async(deleteFunc, deletedFile) => {
        const undoDelete = () => {
            setFiles((files) => files.map(file => file.id == deletedFile.id ? { ...file, deleted: false } : file));
            messageContext.showSuccess(`File ${deletedFile.filename} (${deletedFile.id}) restored`);
        };
        addJob(deleteFunc, undoDelete);
        setFiles((files) => files.map(file => file.id == deletedFile.id ? { ...file, deleted: true } : file));
    }, [ addJob, messageContext ]);
    
    const handleUpdate = useCallback((openedFile) => setFiles((files) => files.map(file => file.id == openedFile.id ? openedFile : file)) , [ ]);

    const handleOpenPreview = useCallback((file) => {
        setPreview({ type: 'SET_PREVIEW_FILE', payload: file });
        setPreview({ type: 'SHOW_PREVIEW' });
    }, [ ]);

    const handleClosePreview = useCallback(() => {
        setPreview({ type: 'HIDE_PREVIEW' });
        setPreview({ type: 'SET_PREVIEW_FILE', payload: emptyFile });
    }, [ ]);
    
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

        uploadFile(form).then(res => {
            messageContext.showSuccess(res.message);
            searching.current.loading();
            createFileLoader().then((loadInfo) => {
                setFileLoader(() => loadInfo);
                loadInfo();
            }).finally(() => searching.current.stopLoading());
        })
        .catch(err => messageContext.showError(err.message))
        .finally(() => setUploading(false));
    }

    return <>
        <div className="loader file-loading" hidden={!uploading}> Uploading files </div> 

        {showPreview && 
        <Modal show={showPreview} onHide={handleClosePreview} className='d-flex modal-bg' fullscreen centered>
            <Modal.Header closeButton className='bg-modal' closeVariant='white'>{previewFile.filename}</Modal.Header>
            <div className='d-flex flex-fill'>
                <PreviewFile fileInfo={previewFile} className="flex-fill" />
            </div>
        </Modal>}

        <div className="d-flex justify-content-center align-items-center mb-4">
            <input type="text" placeholder="Enter Search" className='search-input' onKeyUp={handleSearch}/>
        </div>
        
        <SpinnerDiv ref={searching}>
        <Container fluid="md">
            <Col>
                <div className="d-flex justify-content-center align-items-center">
                    { files.length == 0 && <p className="text-center text-white">No files, start uploading files c:</p> }
                </div>
            </Col>
            
            <Row xs={1} sm={2}  md={3} lg={4} xl={5} className='row-gap-3 d-flex justify-content-center'>
            { files.slice(0, progress).map(file => 
                <Col key={file.id} hidden={file.deleted}>
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
                </Col>) }
            </Row>
            
            { files.length != 0 && <Row className='p-3'>
                <button className="btn btn-load" onClick={fileLoader}>
                    <i className="bi bi-arrow-down-square-fill" />
                </button>
            </Row> }
        </Container>
        </SpinnerDiv>
        
        <div className='d-flex justify-content-end mb-1'>
            { jobsQueue.length > 0 && <label className='undo-button' onClick={undoLastJob}><i className="bi bi-arrow-clockwise"/></label> }
            <label htmlFor="input-upload" className="btn-upload"><i className='bi bi-plus-square-dotted'/></label>
            <input id="input-upload" type="file" style={ {display: 'none'} } onChange={handleFileUpload} multiple/>
        </div>
    </>
}
