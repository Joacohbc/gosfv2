import 'bootstrap/dist/css/bootstrap.min.css';
import Card from 'react-bootstrap/Card';
import './FileItem.css';
import ToolTip from '../ToolTip';
import PropTypes from 'prop-types';
import { createPortal } from 'react-dom';
import { memo, useContext, useRef  } from 'react';
import SimpleModal from '../SimpleModal';
import Form from 'react-bootstrap/Form';
import InputGroup from 'react-bootstrap/InputGroup';
import { MessageContext } from '../../context/message-context';
import { useGetInfo, useFiles } from '../../hooks/files';
import ConfirmDialog from '../ConfirmDialog';
import SharedWithModal from './ShareModal';
import 'bootstrap-icons/font/bootstrap-icons.css'
import { getDisplayFilename } from '../../services/files';

const filesModal = document.getElementById('files-modals');

const FileItem = memo((props) => {
    const download = useRef(null);
    const shareModal = useRef(null);
    const updateModal = useRef(null);
    const inputUpdate = useRef(null);
    const forceDeleteDialog = useRef(null);

    const { getFileInfo, deleteFile, updateFile } = useFiles();
    const messageContext = useContext(MessageContext);
    const { getUserInfo } = useGetInfo();

    const file = Object.preventExtensions({
        id: props.id,
        filename: props.filename,
        contentType: props.contentType,
        url: props.url,
        extension: props.extension,
        name: props.name,
        shared: props.shared,
        sharedWith: props.sharedWith ?? [],
        savedLocal: props.savedLocal ?? false,
    });

    const handleDownload = () => {
        download.current.click();
        props.onDownload(file);
    };

    const handleOpen = (e) => {
        e.preventDefault();
        e.stopPropagation();
        props.onOpen(file);
    };

    const handleDelete = async (e) => {
        e.preventDefault();
        e.stopPropagation();
        try {
            const res = await getFileInfo(file.id);
            file.shared = res.shared;
            file.sharedWith = res.sharedWith?.map(user => getUserInfo(user, false)) ?? [];

            // Si el archivo esta compartido, se muestra el dialogo de confirmación
            if(file.shared || file.sharedWith.length > 0) {
                forceDeleteDialog.current.show();
                return;
            }
            
            props.onDelete(async () => {
                try {
                    await deleteFile(file.id);
                } catch(err) {
                    messageContext.showError(err.message);
                }
            }, file);
        } catch(err) {
            messageContext.showError(err.message);
        }
    };

    const forceFileDelete = async(e) => {
        e.preventDefault();
        e.stopPropagation();
        props.onDelete(async () => {
            try {
                const res = await deleteFile(file.id, true);
                messageContext.showSuccess(res.message);
            } catch(err) {
                messageContext.showError(err.message);
            }
        }, file);
    };

    const handleUpdate = async (e) => {
        e.preventDefault();
        e.stopPropagation();
        if(inputUpdate.current.value.trim() == '') {
            messageContext.showError('Filename cannot be empty');
            return;
        }

        if(inputUpdate.current.value.trim() == file.name) {
            messageContext.showError('Filename cannot be the same');
            return;    
        }

        try {
            const newFilename = inputUpdate.current.value.trim() + '.' + file.extension;
            const res = await updateFile(file.id, {
                filename: newFilename,
            })
            updateModal.current.hide();
            messageContext.showSuccess(res.message);
            props.onUpdate({
                ...file,
                filename: newFilename,
            });
        } catch(err) {
            messageContext.showError(err.message);
        }
    } 

    const openShareModel = async (e) => {
        e.preventDefault();
        e.stopPropagation();
        try {
            shareModal.current.open();
            shareModal.current.setFile(file, false);

            const res = await getFileInfo(file.id);
            file.shared = res.shared;
            file.sharedWith = res.sharedWith?.map(user => getUserInfo(user, false)) ?? [];
            shareModal.current.setFile(file, true);
        } catch(err) {
            messageContext.showError(err.message);
        }
    }

    const openUpdateNameModel = (e) => {
        e.preventDefault();
        e.stopPropagation();
        updateModal.current.show();
    }

    const handleShare = (e) => { 
        e.preventDefault();
        e.stopPropagation();
        props.onShare(file);
    };

    return <>
        <ConfirmDialog 
            title="Are you sure you want to delete this file?"
            message="This file it shared with other users, if you delete it, it will be deleted for you and all users permanently."
            onOk={forceFileDelete} ref={forceDeleteDialog}/>

        {createPortal(<SharedWithModal
            ref={shareModal} 
            onUpdate={handleShare}
            />, filesModal) }
        
        {createPortal(
        <SimpleModal ref={updateModal} title={file.filename} buttonText={<span>Save changes <i className='bi bi-save'/></span>} onClick={handleUpdate}>
            <InputGroup className="mb-3">
                <Form.Control
                    placeholder="Filename"
                    defaultValue={file.name}
                    ref={inputUpdate}
                    onKeyDown={(e) => { if(e.key == 'Enter') handleUpdate() }}
                />
                <InputGroup.Text>.{file.extension}</InputGroup.Text>
            </InputGroup>
        </SimpleModal>, filesModal) }

        <Card className='file' onClick={handleOpen}>
            <Card.Body>
                <ToolTip toolTipMessage={file.filename} placement='bottom'>
                    <p className="text-center file-filename">{getDisplayFilename(file.filename)}</p>
                </ToolTip>

                <div className='text-center'>
                    <a href={file.url} download={file.filename} ref={download} hidden/>
                    <button className='file-actions-item' onClick={handleDownload}><i className='bi bi-file-arrow-down-fill'/></button>
                    <button className='file-actions-item' onClick={handleDelete}><i className='bi bi-trash3-fill'/></button>
                    <button className='file-actions-item' onClick={openUpdateNameModel}><i className='bi bi-pencil-square'/></button>
                    <button className='file-actions-item' onClick={openShareModel}><i className='bi bi-share-fill'/></button>
                </div>
            </Card.Body>
        </Card>
    </>
});

FileItem.propTypes = {
    id: PropTypes.number.isRequired,
    filename: PropTypes.string.isRequired,
    contentType: PropTypes.string.isRequired,
    url: PropTypes.string.isRequired,
    extension: PropTypes.string.isRequired,
    name: PropTypes.string,
    shared: PropTypes.bool,
    sharedWith: PropTypes.array,
    savedLocal: PropTypes.bool,
    onOpen: PropTypes.func,
    onDelete: PropTypes.func,
    onShare: PropTypes.func,
    onUpdate: PropTypes.func,
    onDownload: PropTypes.func
};

FileItem.defaultProps = {
    onOpen: () => {},
    onDelete: () => {},
    onShare: () => {},
    onUpdate: () => {},
    onDownload: () => {}
};

FileItem.displayName = 'FileItem';

export default FileItem;
