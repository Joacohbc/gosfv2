import 'bootstrap/dist/css/bootstrap.min.css';
import Card from 'react-bootstrap/Card';
import './FileItem.css';
import ToolTip from './ToolTip';
import PropTypes from 'prop-types';
import Button from './Button';
import { createPortal } from 'react-dom';
import { useContext, useRef } from 'react';
import Modal from './Modal';
import AuthContext from '../context/auth-context';

const filesModal = document.getElementById('files-modals');

const FileItem = (props) => {
    const modalRef = useRef(null);
    const auth = useContext(AuthContext);

    const handleDownload = () => {
        console.log('Download');
    };
    
    const handleDelete = async() => {
        try {
            const res = await auth.cAxios.delete(`/api/files/${props.id}`);
            props.onDelete(props.id, res.data.message);
        } catch(err) {
            props.onDelete(null, err.data.message);
        }
    };

    const handleUpdate = () => {
        console.log('Update');
    };

    const handleShare = () => {
        modalRef.current.setShowed(true);
    };

    const modal = <Modal ref={modalRef} title={props.filename}>
        <p>{props.id}</p>
    </Modal>;

    const handleOpen = () => {
        props.onOpen(props.id, props.filename);
    };

    return <>
        {createPortal(modal, filesModal) }
        <Card className='file'>
            <Card.Body>
                <Card.Title className='text-center'>File #{props.id}</Card.Title>
                <Card.Text>
                    <ToolTip toolTipMessage={props.filename} placement={'bottom'}>
                        <p className="text-center file-filename" onClick={handleOpen}>{props.filename}</p>
                    </ToolTip>
                </Card.Text>

                <div className='text-center'>
                    <Button text="Download" className="file-actions-item" onClick={handleDownload }/>
                    <Button text="Delete" className="file-actions-item" onClick={handleDelete}/>
                    <Button text="Update" className="file-actions-item" onClick={handleUpdate}/>
                    <Button text="Share" className="file-actions-item" onClick={handleShare}/>
                </div>
            </Card.Body>
        </Card>
    </>
};

FileItem.propTypes = {
    filename: PropTypes.string.isRequired,
    id: PropTypes.number.isRequired,
    onOpen: PropTypes.func.isRequired,
    onDelete: PropTypes.func.isRequired,
    onShare: PropTypes.func.isRequired,
    onUpdate: PropTypes.func.isRequired,
    onDownload: PropTypes.func.isRequired
};

export default FileItem;