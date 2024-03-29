import 'bootstrap/dist/css/bootstrap.min.css';
import './Files.css';
import AuthContext from '../context/auth-context';
import Modal from 'react-bootstrap/Modal';
import { useCallback, useContext, useEffect, useState, useReducer, useRef } from 'react';
import { handleKeyUpWithTimeout } from '../utils/input-text';
import PreviewFile from './PreviewFile';
import { MessageContext } from '../context/message-context';
import { useGetInfo, useFiles } from '../hooks/files';
import useJobsQueue from '../hooks/jobsQueue';
import '../components/Message.css';
import useFilesIDB from '../hooks/useFilesIDB';
import FileContainer from '../components/FileContainer';

const emptyFile = Object.freeze({ id: null, filename: null, contentType: '', url: '', extension: '', deleted: false });

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

export default function Files() {
    const messageContext = useContext(MessageContext);
    const { isLogged, cAxios } = useContext(AuthContext);

    const uploadButton = useRef(null);
    
    const [ files, setFiles ] = useState([]);
    const [ progress, setProgress ] = useState(0);
    const [ fileLoader, setFileLoader ] = useState(() => {});
    const [ uploading, setUploading ] = useState(false);
    const [ searching, setSearching ]= useState(false);

    const [ state , setPreview ] = useReducer(previewReducer, { previewFile: emptyFile, showPreview: false });
    const { previewFile, showPreview } = state;

    const { getFilenameInfo } = useGetInfo();
    const { getFiles, uploadFile } = useFiles();
    const { addJob, undoLastJob, jobsQueue, executeAllJobs } = useJobsQueue(3000);
    const { getFileFromLocal } = useFilesIDB();

    const createFileLoader = useCallback(async (filterCb = (data) => data) => {
        try {
            const files = await getFiles();

            // Filtra los archivos
            let data = filterCb(files.map(file => {
                const f = getFilenameInfo(file, true);
                f.deleted = false;
                return f;
            }))

            // Verifica si los archivos estan guardados localmente
            await Promise.allSettled(data.map(async (file) => {
                const localFile = await getFileFromLocal(file.id);
                file.savedLocal = localFile != null;
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
    }, [ messageContext, getFilenameInfo, getFiles, getFileFromLocal ]);
    
    useEffect(() => {
        if(!cAxios || !isLogged) return;
    
        setSearching(true);
        createFileLoader().then((loadInfo) => {
            setFileLoader(() => loadInfo);
            loadInfo();
        }).finally(() => setSearching(false));
    }, [ isLogged, cAxios, createFileLoader ]);

    const handleDeleteAllInQueue = () => {
        executeAllJobs();
        messageContext.showSuccess('All files deleted');
    }

    const handleSearch = handleKeyUpWithTimeout((e) => {
        setSearching(true);

        const filterCb = (data) => data.filter(file => file.filename.toLowerCase().includes(e.target.value.toLowerCase()) || file.id == e.target.value);
        createFileLoader(filterCb).then((loadInfo) => {
            setFileLoader(() => loadInfo);
            loadInfo();
        }).finally(() => setSearching(false));
    }, 500);

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
            setSearching(true);
            createFileLoader().then((loadInfo) => {
                setFileLoader(() => loadInfo);
                loadInfo();
            }).finally(() => setSearching(false));
        })
        .catch(err => messageContext.showError(err.message))
        .finally(() => setUploading(false));
    }

    const handleFilesDelete = useCallback(async(deleteFunc, deletedFile) => {
        const undoDelete = () => {
            setFiles((files) => files.map(file => file.id == deletedFile.id ? { ...file, deleted: false } : file));
            messageContext.showSuccess(`File ${deletedFile.filename} (${deletedFile.id}) restored`);
        };

        addJob(deleteFunc, undoDelete);
        setFiles((files) => files.map(file => file.id == deletedFile.id ? { ...file, deleted: true } : file));
    }, [ addJob, messageContext ]);
    
    const handleFilesUpdate = useCallback((updatedFile) => setFiles((files) => files.map(file => file.id == updatedFile.id ? updatedFile : file)) , [ ]);

    const handleOpenPreview = useCallback(async (file) => {

        // Si el archivo esta en la base de datos local, se obtiene la URL del archivo
        const localFile = await getFileFromLocal(file.id);
        if (localFile != null) {
            file.url = window.URL.createObjectURL(localFile.blob);
        }
        
        setPreview({ type: 'SET_PREVIEW_FILE', payload: file });
        setPreview({ type: 'SHOW_PREVIEW' });
    }, [ getFileFromLocal ]);

    const handleClosePreview = useCallback(() => {
        setPreview({ type: 'HIDE_PREVIEW' });
        setPreview({ type: 'SET_PREVIEW_FILE', payload: emptyFile });
    }, [ ]);  

    const handleFileDropEnd = (e) => {
        e.preventDefault();
        e.stopPropagation();
    }

    const handleFileDropOver = (e) => {
        e.preventDefault();
        e.stopPropagation();
    }

    const handleFileDrop = (e) => {
        e.preventDefault();
        e.stopPropagation();

        const files = e.dataTransfer.files;
        if(files.length == 0) {
            return;
        }

        const form = new FormData();
        for(let i = 0; i < files.length; i++) {
            form.append('files', files[i]);
        }

        console.log(form);

        setUploading(true);

        uploadFile(form).then(res => {
            messageContext.showSuccess(res.message);
            setSearching(true);
            createFileLoader().then((loadInfo) => {
                setFileLoader(() => loadInfo);
                loadInfo();
            }).finally(() => setSearching(false));
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
        
        <FileContainer
            files={files}
            progress={progress}
            fileLoader={fileLoader}
            loading={searching}
            handleOpenPreview={handleOpenPreview}
            handleFilesDelete={handleFilesDelete}
            handleFilesUpdate={handleFilesUpdate}
        />
            
        <div className='d-flex flex-column fixed-top align-items-end mr-1'>
            <label htmlFor="input-upload" className="btn-upload"
                onDrop={handleFileDrop}
                onDragLeaveCapture={handleFileDropEnd}
                onDragOverCapture={handleFileDropOver} 
                ref={uploadButton} >
                <i className='bi bi-plus-square-dotted'/>
            </label>
            <input id="input-upload" type="file" style={{display: 'none'}} onChange={handleFileUpload} multiple/>

            { jobsQueue.length > 0 && 
                <div className='undo-button'>
                    <label onClick={undoLastJob}><i className="bi bi-arrow-clockwise fs-2"/></label>
                    <span>{jobsQueue.length}</span>
                </div> }
            { jobsQueue.length > 0 && <label className='undo-button' onClick={handleDeleteAllInQueue}><i className="bi bi-tornado fs-2"/></label> }
        </div>
    </>
}