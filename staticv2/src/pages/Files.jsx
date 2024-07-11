import './Files.css';
import AuthContext from '../context/auth-context';
import Modal from 'react-bootstrap/Modal';
import { useCallback, useContext, useEffect, useReducer } from 'react';
import PreviewFile from './PreviewFile';
import { MessageContext } from '../context/message-context';
import { useGetInfo, useFiles } from '../hooks/useFiles';
import useJobsQueue from '../hooks/useJobsQueue';
import FileContainer from '../components/FileContainer';
import { useParams, useNavigate } from 'react-router-dom';
import { useCache } from '../hooks/cache';
import SearchBar from '../components/SearchBar';
import { getDisplayFilename } from '../services/files';

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

const initialFilesLoadingState = {
    files: [],
    progress: 0,
    loading: true,
    fileLoader: () => {},
    defaultNumberOfFilesPerLoad: 35
};

function filesReducer(state, action) {

    const calculateNextProgress = (currentProgress, filesLength) => {
        const nextProgress = currentProgress + state.defaultNumberOfFilesPerLoad;
        return nextProgress >= filesLength ? filesLength : nextProgress;
    }

    const resetProgress = (filesLength) => calculateNextProgress(0, filesLength);

    const calculateProgressOnDelete = (currentProgress, deletedFiles) => {
        const nextProgress = currentProgress - deletedFiles.length;
        return nextProgress < 0 ? 0 : nextProgress;
    }
    
    switch (action.type) {
        case 'SET_FILES':
            return { ...state, files: action.payload, progress: resetProgress(action.payload.length) };
        case 'ADD_FILES': {
            const newFiles = [ ...state.files, ...action.payload ];
            return { ...state, files: newFiles, progress: calculateNextProgress(state.progress, newFiles.length) };
        }
        case 'REMOVE_FILES': {
            const newFiles = state.files.filter(file => !action.payload.includes(file.id));
            return { ...state, files: newFiles, progress: calculateProgressOnDelete(state.progress, action.payload) };
        }
        case 'UPDATE_FILE':
            return { ...state, files: state.files.map(file => file.id === action.payload.id ? action.payload : file) };
        case 'MARK_DELETED':
            return { ...state, files: state.files.map(file => file.id === action.payload ? { ...file, deleted: true } : file) };
        case 'UNMARK_DELETED':
            return { ...state, files: state.files.map(file => file.id === action.payload ? { ...file, deleted: false } : file) };
        case 'NEXT_PROGRESS':
            return { ...state, progress: calculateNextProgress(state.progress, state.files.length) };
        case 'SET_LOADING':
            return { ...state, loading: action.payload };
        case 'SET_FILE_LOADER':
            return { ...state, fileLoader: action.payload };
        default:
            throw new Error(`Unhandled action type: ${action.type}`);
    }
}

const DELETE_UNDO_DURATION = 5000;

export default function Files() {
    const messageContext = useContext(MessageContext);
    const { isLogged } = useContext(AuthContext);
    const { sharedFileId } = useParams();
    const navigate = useNavigate();

    const { cacheService } = useCache();

    const [ filesLoadingState, dispatchFilesLoading ] = useReducer(filesReducer, initialFilesLoadingState);
    const { files, progress, loading, fileLoader } = filesLoadingState;

    const [ previewState, setPreview ] = useReducer(previewReducer, { previewFile: emptyFile, showPreview: false });
    const { previewFile, showPreview } = previewState;

    const { getFilenameInfo } = useGetInfo();
    const { getFiles, uploadFile, getShareFileInfo, deleteFiles } = useFiles();
    const { addJob, undoLastJob, undoJob, jobsQueue, clearAllJobs, undoAllJobs } = useJobsQueue(DELETE_UNDO_DURATION);
    
    const createFileLoader = useCallback(async (filterCb = (data) => data) => {
        try {
            dispatchFilesLoading({ type: 'SET_LOADING', payload: true });
            
            let files = [];
            const cacheFiles = cacheService.getCacheFiles();
            
            if (cacheFiles.value && cacheFiles.timestamp.getTime() > Date.now() - 5000) {
                files = cacheFiles.value;
            } else {
                files = await getFiles();
            }
            
            let data = filterCb(files.map(file => {
                const f = getFilenameInfo(file, false);
                f.deleted = false;
                return f;
            }));
            
            dispatchFilesLoading({ type: 'SET_FILES', payload: data });
            return () => { dispatchFilesLoading({ type: 'NEXT_PROGRESS' } ) };
        } catch(err) {
            messageContext.showError(err.message);
            return () => [];
        } finally {
            dispatchFilesLoading({ type: 'SET_LOADING', payload: false });
        }
    }, [ messageContext, getFilenameInfo, getFiles, cacheService ]);
    
    
    useEffect(() => {
        if(!isLogged) return;
    
        if(sharedFileId) {
            getShareFileInfo(sharedFileId).then((file) => {
                setPreview({ type: 'SET_PREVIEW_FILE', payload: getFilenameInfo(file, false) });
                setPreview({ type: 'SHOW_PREVIEW' });
            }).catch(err => messageContext.showError(err.message));
        }

        createFileLoader().then((loadInfo) => {
            dispatchFilesLoading({ type: 'SET_FILE_LOADER', payload: loadInfo });
        })
    }, [ isLogged, createFileLoader, getShareFileInfo, sharedFileId, messageContext, getFilenameInfo ]);

    const handleDeleteAllInQueue = () => {
        const fileIds = jobsQueue.map(job => job.info.fileId);
        deleteFiles(fileIds, true)
        .then((res) => {
            clearAllJobs();
            dispatchFilesLoading({ type: 'REMOVE_FILES', payload: fileIds });
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
                dispatchFilesLoading({ type: 'SET_FILE_LOADER', payload: loadInfo });
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
                dispatchFilesLoading({ type: 'UNMARK_DELETED', payload: deletedFile.id });
                messageContext.dismiss(removeMessageId);
            },
            clearCb: () => {
                // Si se limpia el job, se elimina el mensaje de deshacer eliminación
                messageContext.dismiss(removeMessageId);
            },
            actionInfo: { fileId: deletedFile.id }
        });

        dispatchFilesLoading({ type: 'MARK_DELETED', payload: deletedFile.id });

        removeMessageId = messageContext.showAction(
            `File ${deletedFile.filename} (${deletedFile.id}) deleted`, 
            'Undo',
            undoJob.bind(null, job), // bind retorna una nueva función que cuando se llama, llama a la función original con el job ya como parámetro
            job.deleteIn
        );

    }, [addJob, undoJob, messageContext]);
    
    const handleFilesUpdate = useCallback((updatedFile) => 
        dispatchFilesLoading({ type: 'UPDATE_FILE', payload: updatedFile }), 
    []);

    const handleOpenPreview = useCallback(async (file) => {        
        setPreview({ type: 'SET_PREVIEW_FILE', payload: file });
        setPreview({ type: 'SHOW_PREVIEW' });
    }, [ ]);

    const handleClosePreview = useCallback(() => {
        setPreview({ type: 'HIDE_PREVIEW' });
        setPreview({ type: 'SET_PREVIEW_FILE', payload: emptyFile });
        navigate('/files'); // To avoid that in the next re-load the file is opened again
    }, [ navigate ]);  

    return <>
        { showPreview && 
        <Modal show={showPreview} onHide={handleClosePreview} className='d-flex modal-bg' fullscreen centered>
            <Modal.Header closeButton className='bg-modal' closeVariant='white'>{getDisplayFilename(previewFile.filename)}</Modal.Header>
            <div className='d-flex flex-fill'>
                <PreviewFile contentType={previewFile.contentType} url={sharedFileId ? previewFile.sharedUrl : previewFile.url} className="flex-fill" />
            </div>
        </Modal> }

        <SearchBar createFileLoader={createFileLoader} setFileLoader={(fileLoader) => dispatchFilesLoading({ type: 'SET_FILE_LOADER', payload: fileLoader })}/>
        
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
