import 'bootstrap/dist/css/bootstrap.min.css';
import Card from 'react-bootstrap/Card';
import '../components/FileItem.css';
import ToolTip from './ToolTip';
import PropTypes from 'prop-types';
import Button from './Button';
import { createPortal } from 'react-dom';
import { useRef } from 'react';
import Modal from './Modal';

const filesModal = document.getElementById('files-modals');

const FileItem = (props) => {
    const modalRef = useRef(null);

    const handleDownload = () => {
        console.log('Download');
    };
    
    const handleDelete = () => {
        console.log('Delete');
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

    return <>
        {createPortal(modal, filesModal) }
        <Card className='file'>
            <Card.Body>
                <Card.Title className='text-center'>File #{props.id}</Card.Title>
                <Card.Text>
                    <ToolTip toolTipMessage={props.filename} placement={'bottom'}>
                        <p className="text-center file-filename">{props.filename}</p>
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
    id: PropTypes.number.isRequired
};

export default FileItem;