import './Files.css';
import AuthContext from '../context/auth-context';
import Modal from 'react-bootstrap/Modal';
import { useCallback, useContext, useEffect, useState, useReducer } from 'react';
import PreviewFile from './PreviewFile';
import { MessageContext } from '../context/message-context';
import { useGetInfo, useFiles } from '../hooks/files';
import useJobsQueue from '../hooks/jobsQueue';
import useFilesIDB from '../hooks/useFilesIDB';
import FileContainer from '../components/FileContainer';
import { useParams, useNavigate } from 'react-router-dom';
import { useCache } from '../hooks/cache';
import SearchBar from '../components/SearchBar';

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

    const { cacheService } = useCache();
    const [ files, setFiles ] = useState([]);
    const [ progress, setProgress ] = useState(files.length);
    const [ fileLoader, setFileLoader ] = useState(() => {});
    const [ loading, setLoading ]= useState(true);

    const [ state , setPreview ] = useReducer(previewReducer, { previewFile: emptyFile, showPreview: false });
    const { previewFile, showPreview } = state;

    const { getFilenameInfo } = useGetInfo();
    const { getFiles, uploadFile, getShareFileInfo, deleteFiles } = useFiles();
    const { addJob, undoLastJob, undoJob, jobsQueue, clearAllJobs, undoAllJobs } = useJobsQueue(DELETE_UNDO_DURATION);
    const { getFileFromLocal } = useFilesIDB();
    
    const createFileLoader = useCallback(async (filterCb = (data) => data) => {
        try {
            setLoading(true);

            let files = [];

            // Check difference between the current files and the files in the cache (+1 second)
            const cacheFiles = cacheService.getCacheFiles();
            if(cacheFiles.timestamp.getTime() > Date.now() - 1000) {
                files = cacheFiles.value;
            } else {
                files = await getFiles();
            }

            // Filtra los archivos
            let data = filterCb(files.map(file => {
                const f = getFilenameInfo(file, false);
                f.deleted = false;
                return f;
            }));

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
            return () => [];
        } finally {
            setLoading(false);
        }
    }, [ messageContext, getFilenameInfo, getFiles, getFileFromLocal, cacheService ]);
    
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
        setLoading(false);
    }, [ isLogged, createFileLoader, getShareFileInfo, sharedFileId, messageContext, getFilenameInfo ]);

    const handleDeleteAllInQueue = () => {
        deleteFiles(jobsQueue.map(job => job.info.fileId), true)
        .then((res) => {
            clearAllJobs();
            messageContext.showSuccess(res.message);
        })
        .catch((err) => {
            undoAllJobs();
            messageContext.showError(err.message);
        });
    }

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

        const job = addJob({
            actionCb: () => {
                // Si la eliminación es exitosa, se elimina el mensaje de deshacer eliminación
                messageContext.dismiss(removeMessageId);
                deleteFunc();
            },
            undoCb: () => {
                // Si se deshace la eliminación, se restaura el archivo
                setFiles((files) => files.map(file => file.id == deletedFile.id ? { ...file, deleted: false } : file));
                messageContext.dismiss(removeMessageId);
            },
            clearCb: () => {
                // Si se limpia el job, se elimina el mensaje de deshacer eliminación
                messageContext.dismiss(removeMessageId);
            },
            actionInfo: { fileId: deletedFile.id }
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

        <SearchBar createFileLoader={createFileLoader} setFileLoader={setFileLoader}/>
        
        <FileContainer
            files={files}
            progress={progress}
            fileLoader={fileLoader}
            loading={loading}
            handleOpenPreview={handleOpenPreview}
            handleFilesDelete={handleFilesDelete}
            handleFilesUpdate={handleFilesUpdate}
        />
        
        <div className='d-flex flex-row sticky-bottom justify-content-center align-items-end gap-2'>
            { jobsQueue.length == 0 &&             
            <label htmlFor="input-upload" className="btn-upload"
                onDrop={handleFileUploadByDrop}
                onDragLeaveCapture={removeDefault}
                onDragOverCapture={removeDefault}>
                <i className='bi bi-plus-square-dotted'/>
            </label>}

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
