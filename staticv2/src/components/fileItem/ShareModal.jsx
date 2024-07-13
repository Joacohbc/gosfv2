import 'bootstrap/dist/css/bootstrap.min.css';
import './ShareModal.css';
import PropTypes from 'prop-types';
import Button from '../Button';
import { forwardRef, useContext, useRef, useState } from 'react';
import SimpleModal from '../SimpleModal';
import Form from 'react-bootstrap/Form';
import InputGroup from 'react-bootstrap/InputGroup';
import Container from 'react-bootstrap/Container';
import Row from 'react-bootstrap/Row';
import { useImperativeHandle } from 'react';
import Col from 'react-bootstrap/esm/Col';
import { MessageContext } from '../../context/message-context';
import { useFiles, useGetInfo } from '../../hooks/useFiles';
import Placeholder from 'react-bootstrap/Placeholder';
import { getDisplayFilename } from '../../services/files';

const SharedWithModal = forwardRef((props, ref) => {
    const shareModal = useRef(null);
    const userIdAdded = useRef(null);
    const messageContext = useContext(MessageContext);
    const [ completedInfo, setCompletedInfo ] = useState(false);
    const [ file, setFile ] = useState(null);
    const { getUserInfo } = useGetInfo();
    const { updateFile, removeUserFromFile, addUserToFile } = useFiles();

    useImperativeHandle(ref, () => ({
        open: shareModal.current.show,
        setFile: (file, completed) => {
            setCompletedInfo(completed);
            setFile(file);
        },
    }), [ shareModal ]);

    const handleMarkAsPublic = (value) => {
        return async (e) => {
            e.preventDefault();
            try {
                const res = await updateFile(file.id, { 
                    shared: value,
                });
                setFile((file) => ({ ...file, shared: value }));
                messageContext.showSuccess(res.data.shared ? "File updated to: public" : "File updated to: restricted");
                props.onUpdate();
            } catch(err) {
                messageContext.showError(err.message);
            }
        }
    };

    const handleRemoveUser = (userId) => {
        return async (e) => {
            e.preventDefault();
            try {
                const res = await removeUserFromFile(file.id, userId);
                setFile((file) => ({
                    ...file,
                    sharedWith: res.data.sharedWith.map(user => getUserInfo(user, false)),
                }));
                messageContext.showSuccess(res.message);
                props.onUpdate();
            } catch(err) {
                messageContext.showError(err.message);
            }
        }

    }

    const handleAddUser = async (e) => {
        e.preventDefault();
        props.onUpdate();

        const username = userIdAdded.current.value;
        if(username.trim() == "") {
            messageContext.showError("Please enter a User ID");
            return;
        }

        try {
            const userId = username.substring(username.lastIndexOf('#') + 1);
            const res = await addUserToFile(file.id, userId)
            setFile((file) => ({
                ...file,
                sharedWith: res.data.sharedWith.map(user => getUserInfo(user, false)),
            }));
            messageContext.showSuccess(res.message);
            props.onUpdate();
        } catch(err) {
            messageContext.showError(err.message);
        }
    } 

    const handleCopyLink = async (e) => {
        e.preventDefault();
        try {
            await navigator.clipboard.writeText(`${window.location.origin}/shared/${file?.id}`);
            messageContext.showSuccess("Link copied to clipboard");
        } catch(err) {
            messageContext.showError(err.message);
        }
    }

    return <SimpleModal ref={shareModal} title={getDisplayFilename(file?.filename)}>
        <Form>
            <InputGroup>
                <Form.Control placeholder='Enter User ID' ref={userIdAdded}/>
                <Button text={<i className='bi bi-person-fill-add'/>} onClick={handleAddUser}/>
            </InputGroup>
        </Form>
        <hr className="hr" />

        <Container className='d-block'>
            { completedInfo && file?.sharedWith.map(user => {
                return <Row key={user.id} className='overlay-user-share align-items-center'>
                    <Col xs={9}>
                        <img src={user.icon} className='modal-user-share-icon'/>
                        <span className='ms-3'>{user.username} #{user.id}</span>
                    </Col>
                    <Col className='p-2 d-flex justify-content-end'>
                        <Button text={<i className='bi bi-person-fill-dash'/>} onClick={handleRemoveUser(user.id)}/>
                    </Col>
                </Row>
            }) }

            { !completedInfo && Array.from({ length: 3 }).map((_, i) =>
                <Row key={i} className='overlay-user-share align-items-center'>
                    <Col xs={9}>
                        <Placeholder as="div" animation="wave">
                            <Placeholder as="img" xs={3} className='modal-user-share-icon'/>
                            <Placeholder as="span" xs={8} className='ms-3'/>
                        </Placeholder>
                    </Col>
                    <Col className='p-2 d-flex justify-content-end'>
                        <Placeholder as="button" xs={1} className='file-actions-item'><i className='bi bi-person-fill-dash'/></Placeholder>
                    </Col>
                </Row>
            )}


            { completedInfo && file?.sharedWith.length == 0 && <p className='text-center'>The file is not shared with anyone</p>}
        </Container>
        
        <hr className="hr" />
        <Form>
            <InputGroup>
                <Form.Control value={`${window.location.origin}/shared/${file?.id}`} onClick={handleCopyLink} readOnly/>
                
                <InputGroup.Text className='cursor-pointer'>
                    { completedInfo && <Button text={<i className='bi bi-clipboard'/>} onClick={handleCopyLink}/> }
                    { !completedInfo && <Placeholder as="button" xs={1} className='file-actions-item'><i className='bi bi-clipboard'/></Placeholder> }
                </InputGroup.Text>

                <InputGroup.Text className='cursor-pointer'>
                    { completedInfo && (file?.shared ? 
                        <Button text={<i className='bi bi-unlock-fill'/>} onClick={handleMarkAsPublic(false)}/> 
                        : <Button text={<i className='bi bi-lock-fill'/>} onClick={handleMarkAsPublic(true)}/>) }

                    { !completedInfo && <Placeholder as="button" xs={1} className='file-actions-item'><i className='bi bi-lock-fill'/></Placeholder> }
                </InputGroup.Text>
            </InputGroup>
        </Form>
    </SimpleModal>
});

SharedWithModal.displayName = 'SharedWithModal';

SharedWithModal.propTypes = {
    onUpdate: PropTypes.func,
};

export default SharedWithModal;