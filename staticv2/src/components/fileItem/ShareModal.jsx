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
import { useFiles, useGetInfo } from '../../hooks/files';

const SharedWithModal = forwardRef((props, ref) => {
    const shareModal = useRef(null);
    const userIdAdded = useRef(null);
    const messageContext = useContext(MessageContext);
    const [ file, setFile ] = useState(null);
    const { getUserInfo } = useGetInfo();
    const { updateFile, removeUserFromFile, addUserToFile } = useFiles();

    useImperativeHandle(ref, () => ({
        open: (file) => {
            setFile(file);
            shareModal.current.show();
        },
    }), [ shareModal ]);

    const handleMarkAsPublic = async (e) => {
        try {
            const res = await updateFile(file.id, { 
                shared: e.target.checked,
            });
            setFile((file) => ({ ...file, shared: e.target.checked }));
            messageContext.showSuccess(res.data.shared ? "File updated to: public" : "File updated to: restricted");
            props.onUpdate();
        } catch(err) {
            messageContext.showError(err.message);
        }
    };

    const handleRemoveUser = (userId) => {
        return async (e) => {
            e.preventDefault();
            try {
                const res = await removeUserFromFile(file.id, userId);
                setFile((file) => ({
                    ...file,
                    sharedWith: res.data.sharedWith.map(user => getUserInfo(user, true)),
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
                sharedWith: res.data.sharedWith.map(user => getUserInfo(user, true)),
            }));
            messageContext.showSuccess(res.message);
            props.onUpdate();
        } catch(err) {
            messageContext.showError(err.message);
        }
    } 

    const handleCopyLink = async () => {
        try {
            await navigator.clipboard.writeText(`${window.location.origin}/shared/${file?.id}`);
            messageContext.showSuccess("Link copied to clipboard");
        } catch(err) {
            messageContext.showError(err.message);
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
                        <Button text={"Delete"} onClick={handleRemoveUser(user.id)}/>
                    </Col>
                </Row>
            }) }

            { file?.sharedWith.length == 0 && <p className='text-center'>The file is not shared with anyone</p>}
        </Container>
        
        <hr className="hr" />
        <Form>
            <InputGroup>
                <Form.Control value={`${window.location.origin}/shared/${file?.id}`} onClick={handleCopyLink} readOnly/>
                <InputGroup.Checkbox label="Public" onChange={handleMarkAsPublic} defaultChecked={file?.shared}/>
            </InputGroup>
        </Form>
    </SimpleModal>
});

SharedWithModal.displayName = 'SharedWithModal';

SharedWithModal.propTypes = {
    onUpdate: PropTypes.func,
};

export default SharedWithModal;