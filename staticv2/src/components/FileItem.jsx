import 'bootstrap/dist/css/bootstrap.min.css';
import Card from 'react-bootstrap/Card';
import './FileItem.css';
import ToolTip from './ToolTip';
import PropTypes from 'prop-types';
import Button from './Button';
import { createPortal } from 'react-dom';
import { forwardRef, memo, useContext, useRef, useState } from 'react';
import SimpleModal from './SimpleModal';
import AuthContext from '../context/auth-context';
import Form from 'react-bootstrap/Form';
import InputGroup from 'react-bootstrap/InputGroup';
import Container from 'react-bootstrap/Container';
import Row from 'react-bootstrap/Row';
import { useImperativeHandle } from 'react';
import Col from 'react-bootstrap/esm/Col';
import { MessageContext } from '../context/message-context';
import { useGetInfo } from '../hooks/files';

const filesModal = document.getElementById('files-modals');

const FileItem = memo((props) => {
    const download = useRef(null);
    const shareModal = useRef(null);
    const updateModal = useRef(null);
    const inputUpdate = useRef(null);
    const auth = useContext(AuthContext);
    const messageContext = useContext(MessageContext);
    const { getUserInfo } = useGetInfo();

    const file = Object.preventExtensions({
        id: props.id,
        filename: props.filename,
        contentType: props.contentType,
        url: props.url,
        extesion: props.extesion,
        name: props.name,
        shared: null,
        sharedWith: [],
    });

    const handleDownload = () => {
        download.current.click();
        props.onDownload(file);
    };

    const handleOpen = () => {
        props.onOpen(file);
    };

    const handleDelete = async() => {
        try {
            const res = await auth.cAxios.delete(`/api/files/${file.id}`);
            messageContext.showSuccess(res.data.message);
            props.onDelete(file);
        } catch(err) {
            messageContext.showError(err.response.data.message);
        }
    };

    const handleUpdate = async () => {
        if(inputUpdate.current.value.trim() == '') {
            messageContext.showError('Filename cannot be empty');
            return;
        }

        if(inputUpdate.current.value.trim() == file.name) {
            messageContext.showError('Filename cannot be the same');
            return;    
        }

        try {
            const res = await auth.cAxios.put(`/api/files/${file.id}`, { filename: inputUpdate.current.value + '.' + file.extesion });
            updateModal.current.hide();
            messageContext.showSuccess(res.data.message);
            props.onUpdate({
                ...file,
                filename: inputUpdate.current.value.trim() + '.' + file.extesion,
            });
        } catch(err) {
            messageContext.showError(err.response.data.message);
        }
    } 

    const openShareModel = async () => {
        try {
            const res = await auth.cAxios.get(`/api/files/${file.id}/info`);
            file.shared = res.data.shared;
            file.sharedWith = res.data.shared_with?.map(user => getUserInfo(user, true)) ?? [];
            shareModal.current.open(file);
        } catch(err) {
            messageContext.showError(err.data.message);
        }
    }

    const openUpdateNameModel = () => {
        updateModal.current.show();
    }

    const handleShare = () => { 
        props.onShare(file);
    };


    return <>
        {createPortal(<SharedWithModal
            ref={shareModal} 
            onUpdate={handleShare}
            />, filesModal) }
        
        {createPortal(
        <SimpleModal ref={updateModal} title={file.filename} buttonText={"Save changes"} onClick={handleUpdate}>
            <InputGroup className="mb-3">
                <Form.Control
                    placeholder="Filename"
                    defaultValue={file.name}
                    ref={inputUpdate}
                    onKeyDown={(e) => { if(e.key == 'Enter') handleUpdate() }}
                />
                <InputGroup.Text>.{file.extesion}</InputGroup.Text>
            </InputGroup>
        </SimpleModal>, filesModal) }

        <Card className='file'>
            <Card.Body>
                <Card.Title className='text-center'>File #{file.id}</Card.Title>
                <Card.Text>
                    <ToolTip toolTipMessage={file.filename} placement={'bottom'}>
                        <p className="text-center file-filename" onClick={handleOpen}>{file.filename}</p>
                    </ToolTip>
                </Card.Text>   

                <div className='text-center'>
                    <a href={auth.addTokenParam(file.url)} download={file.filename} ref={download} hidden/>
                    <Button text="Download" className="file-actions-item" onClick={handleDownload}/>
                    <Button text="Delete" className="file-actions-item" onClick={handleDelete}/>
                    <Button text="Update" className="file-actions-item" onClick={openUpdateNameModel}/>
                    <Button text="Share" className="file-actions-item" onClick={openShareModel}/>
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
    extesion: PropTypes.string.isRequired,
    name: PropTypes.string,
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

const SharedWithModal = forwardRef((props, ref) => {
    const shareModal = useRef(null);
    const messageContext = useContext(MessageContext);
    const { cAxios } = useContext(AuthContext);
    const userIdAdded = useRef(null);
    const [ file, setFile ] = useState(null);
    const { getUserInfo } = useGetInfo();

    useImperativeHandle(ref, () => ({
        open: (file) => {
            setFile(file);
            shareModal.current.show();
        },
    }), [ shareModal ]);

    const handleMarkAsPublic = async (e) => {
        try {
            const res = await cAxios.put(`/api/files/${file.id}`, { 
                shared: e.target.checked,
            });
            setFile((file) => ({ ...file, shared: e.target.checked }));
            messageContext.showSuccess(res.data.message);
            props.onUpdate();
        } catch(err) {
            messageContext.showError(err.response.data.message);
        }
    };

    const handleRemvoeUser = (userId) => {
        return async (e) => {
            e.preventDefault();
            try {
                const res = await cAxios.delete(`/api/files/share/${file.id}/user/${userId}`);
                setFile((file) => ({
                    ...file,
                    sharedWith: file.sharedWith.filter(user => user.id != userId)
                }));
                messageContext.showSuccess(res.data.message);
                props.onUpdate();
            } catch(err) {
                messageContext.showError(err.response.data.message);
            }
        }
    }

    const handleAddUser = async (e) => {
        e.preventDefault();
        props.onUpdate();

        const userId = userIdAdded.current.value;
        if(userId.trim() == "") {
            messageContext.showError("Please enter a User ID");
            return;
        }

        try {
            const id = userId.substring(userId.lastIndexOf('#') + 1);
            const res = await cAxios.post(`/api/files/share/${file.id}/user/${id}`);
            messageContext.showSuccess(res.data.message);

            const fileInfo = await cAxios.get(`/api/files/${file.id}/info`);
            setFile((prevFile) => ({
                ...prevFile,
                sharedWith: fileInfo.data.shared_with.map(user => getUserInfo(user, true)) 
            }));
            props.onUpdate();
        } catch(err) {
            messageContext.showError(err.response.data.message);
        }
    } 

    const handleCopyLink = async () => {
        try {
            await navigator.clipboard.writeText(`${window.location.origin}/shared/${file?.id}`);
            messageContext.showSuccess("Link copied to clipboard");
        } catch(err) {
            messageContext.showError(err);
        }
    }

    return <SimpleModal ref={shareModal} title={file?.filename}>
        <Form>
            <InputGroup>
                <Form.Control placeholder='Enter User ID' ref={userIdAdded}/>
                <Button text="Add" onClick={handleAddUser}/>
            </InputGroup>
        </Form>
        <hr className="hr" />

        <Container className='d-block'>
            { file?.sharedWith.map(user => {
                return <Row key={user.id} className='overlay-user-share align-items-center'>
                    <Col xs={9}>
                        <img src={user.icon} className='modal-user-share-icon'/>
                        <span className='ms-3'>{user.username} #{user.id}</span>
                    </Col>
                    <Col className='p-2 d-flex justify-content-end'>
                        <Button text={"Delete"} onClick={handleRemvoeUser(user.id)}/>
                    </Col>
                </Row>
            }) }

            { file?.sharedWith.length == 0 && <p className='text-center'>The file is not shared with anyone</p>}
        </Container>
        
        <hr className="hr" />
        <Form>
            <InputGroup>
                <Form.Control value={`${window.location.origin}/shared/${file?.id}`} onClick={handleCopyLink}/>
                <InputGroup.Checkbox label="Public" onChange={handleMarkAsPublic} defaultChecked={file?.shared}/>
            </InputGroup>
        </Form>
    </SimpleModal>
});

SharedWithModal.displayName = 'SharedWithModal';

SharedWithModal.propTypes = {
    onUpdate: PropTypes.func,
};