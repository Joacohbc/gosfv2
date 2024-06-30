import 'bootstrap/dist/css/bootstrap.min.css';
import './Files.css';
import AuthContext from '../context/auth-context';
import Modal from 'react-bootstrap/Modal';
import { useCallback, useContext, useEffect, useState, useReducer } from 'react';
import { handleKeyUpWithTimeout } from '../utils/input-text';
import PreviewFile from './PreviewFile';
import { MessageContext } from '../context/message-context';
import { useGetInfo, useFiles } from '../hooks/files';
import useJobsQueue from '../hooks/jobsQueue';
import useFilesIDB from '../hooks/useFilesIDB';
import FileContainer from '../components/FileContainer';
import { useParams, useNavigate } from 'react-router-dom';

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

function removeDefault(e) {
    e.preventDefault();
    e.stopPropagation();
}

const DELETE_UNDO_DURATION = 5000;

export default function Files() {
    const messageContext = useContext(MessageContext);
    const { isLogged } = useContext(AuthContext);
    const { sharedFileId } = useParams();
    const navigate = useNavigate();

    const [ files, setFiles ] = useState([]);
    const [ progress, setProgress ] = useState(0);
    const [ fileLoader, setFileLoader ] = useState(() => {});
    const [ searching, setSearching ]= useState(false);

    const [ state , setPreview ] = useReducer(previewReducer, { previewFile: emptyFile, showPreview: false });
    const { previewFile, showPreview } = state;

    const { getFilenameInfo } = useGetInfo();
    const { getFiles, uploadFile, getShareFileInfo } = useFiles();
    const { addJob, undoLastJob, undoJob, jobsQueue, executeAllJobs } = useJobsQueue(DELETE_UNDO_DURATION);
    const { getFileFromLocal } = useFilesIDB();

    const createFileLoader = useCallback(async (filterCb = (data) => data) => {
        try {
            setSearching(true)

            const files = await getFiles();

            // Filtra los archivos
            let data = filterCb(files.map(file => {
                const f = getFilenameInfo(file, false);
                f.deleted = false;
                return f;
            }))

            // Verifica si los archivos están guardados localmente
            await Promise.allSettled(data.map(async (file) => {
                const localFile = await getFileFromLocal(file.id);
                file.savedLocal = localFile != null;
            }));

            const x = 35;
            // Carga de X en X archivos o todos los archivos si son menos de X
            const numberOfFilesPerLoad = data.length >= x ? x : data.length;
            setFiles(data);
            setProgress(0);
            return () => {
                setProgress(prevProgress => {
                    const nextProgress = prevProgress + numberOfFilesPerLoad;
                    return nextProgress >= data.length ? data.length : nextProgress;
                });
            }
        } catch(err) {
            messageContext.showError(err.message);
            return [];
        } finally {
            setSearching(false);
        }
    }, [ messageContext, getFilenameInfo, getFiles, getFileFromLocal ]);
    
    useEffect(() => {
        if(!isLogged) return;
    
        if(sharedFileId) {
            getShareFileInfo(sharedFileId).then((file) => {
                setPreview({ type: 'SET_PREVIEW_FILE', payload: getFilenameInfo(file, false) });
                setPreview({ type: 'SHOW_PREVIEW' });
            }).catch(err => messageContext.showError(err.message));
        }

        createFileLoader().then((loadInfo) => {
            setFileLoader(() => loadInfo);
            loadInfo();
        })
    }, [ isLogged, createFileLoader, getShareFileInfo, sharedFileId, messageContext, getFilenameInfo]);

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

    const handleFileUpload = (form) => {
        const uploadPromise = async () => {
            
            try {
                const res = await uploadFile(form)
            
                const loadInfo = await createFileLoader()
                setFileLoader(() => loadInfo);
                loadInfo();
    
                return res.message;
            } catch(err) {
                throw new Error(err.message);
            }
        };
        messageContext.showPromise(uploadPromise, 
            'Uploading files...', 
            (message) => message, 
            'Error uploading files');
    }
    
    const handleFileUploadByClick = (e) => {
        e.preventDefault();

        const files = e.target.files;
        if(files.length == 0) {
            return;
        }

        const form = new FormData();
        for(let i = 0; i < files.length; i++) {
            form.append('files', files[i]);
        }
    
        handleFileUpload(form);
    }

    const handleFileUploadByDrop = (e) => {
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

        handleFileUpload(form);
    }

    const handleFilesDelete = useCallback(async(deleteFunc, deletedFile) => {
        
        // Id del mensaje de deshacer eliminación
        let removeMessageId = null;

        const job = addJob(() => {
            // Si la eliminación es exitosa, se elimina el mensaje de deshacer eliminación
            messageContext.dismiss(removeMessageId);
            deleteFunc();
        },
        () => {
            // Si se deshace la eliminación, se restaura el archivo
            setFiles((files) => files.map(file => file.id == deletedFile.id ? { ...file, deleted: false } : file));
            messageContext.dismiss(removeMessageId);
        });

        setFiles((files) => files.map(file => file.id == deletedFile.id ? { ...file, deleted: true } : file));

        removeMessageId = messageContext.showAction(
            `File ${deletedFile.filename} (${deletedFile.id}) deleted`, 
            'Undo',
            undoJob.bind(null, job), // bind retorna una nueva función que cuando se llama, llama a la función original con el job ya como parámetro
            job.deleteIn
        );        
    }, [addJob, undoJob, messageContext]);
    
    const handleFilesUpdate = useCallback((updatedFile) => 
        setFiles((files) => files.map(file => file.id == updatedFile.id ? updatedFile : file)),
    []);

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
        navigate('/files'); // To avoid that in the next re-load the file is opened again
    }, [ navigate ]);  

    return <>
        { showPreview && 
        <Modal show={showPreview} onHide={handleClosePreview} className='d-flex modal-bg' fullscreen centered>
            <Modal.Header closeButton className='bg-modal' closeVariant='white'>{previewFile.filename}</Modal.Header>
            <div className='d-flex flex-fill'>
                <PreviewFile contentType={previewFile.contentType} url={sharedFileId ? previewFile.sharedUrl : previewFile.url} className="flex-fill" />
            </div>
        </Modal> }

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
                onDrop={handleFileUploadByDrop}
                onDragLeaveCapture={removeDefault}
                onDragOverCapture={removeDefault}>
                <i className='bi bi-plus-square-dotted'/>
            </label>

            <input id="input-upload" type="file" style={{display: 'none'}} onChange={handleFileUploadByClick} multiple/>

            { jobsQueue.length > 0 && 
            <div className='undo-button'>
                <label onClick={undoLastJob}><i className="bi bi-arrow-clockwise fs-2"/></label>
                <span>{jobsQueue.length}</span>
            </div> }
            
            { jobsQueue.length > 0 && 
            <label className='undo-button' onClick={handleDeleteAllInQueue}><i className="bi bi-tornado fs-2"/></label> }
        </div>
    </>
}