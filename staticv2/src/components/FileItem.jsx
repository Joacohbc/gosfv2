import 'bootstrap/dist/css/bootstrap.min.css';
import Card from 'react-bootstrap/Card';
import './FileItem.css';
import ToolTip from './ToolTip';
import PropTypes from 'prop-types';
import Button from './Button';
import { createPortal } from 'react-dom';
import { memo, useContext, useRef } from 'react';
import SimpleModal from './SimpleModal';
import AuthContext from '../context/auth-context';
import Form from 'react-bootstrap/Form';
import InputGroup from 'react-bootstrap/InputGroup';


const filesModal = document.getElementById('files-modals');

const FileItem = memo((props) => {
    const download = useRef(null);

    const file = Object.freeze({
        id: props.id,
        filename: props.filename,
        contentType: props.contentType,
        url: props.url,
        extesion: props.extesion,
        name: props.name,
    });

    const shareModal = useRef(null);
    const updateModal = useRef(null);
    const inputUpdate = useRef(null);
    const auth = useContext(AuthContext);

    const handleDownload = () => {
        download.current.click();
        props.onDownload(file, null, null);
    };
    
    const handleDelete = async() => {
        try {
            const res = await auth.cAxios.delete(`/api/files/${file.id}`);
            props.onDelete(file, res.data.message, null);
        } catch(err) {
            props.onDelete(null, null, err.data.message);
        }
    };
    
    const handleOpen = () => {
        props.onOpen(file);
    };

    const handleUpdate = async () => {
        try {
            const res = await auth.cAxios.put(`/api/files/${file.id}`, { filename: inputUpdate.current.value + '.' + file.extesion });
            updateModal.current.hide();
            props.onUpdate({
                ...file,
                filename: inputUpdate.current.value + '.' + file.extesion,
            }, res.data.message, null);
        } catch(err) {
            props.onUpdate(null, null, err.data.message);
        }
    } 

    const showUpdateModal = () => {
        updateModal.current.show();
    };

    const showShareModal = () => {
        shareModal.current.show();
    };


    return <>
        {createPortal(
        <SimpleModal ref={shareModal} title={file.filename} buttonText={"Save changes"} onClick={() => console.log("Confirm")}>
            <p>{file.id}</p>
        </SimpleModal>, filesModal) }
        
        {createPortal(
        <SimpleModal ref={updateModal} title={file.filename} buttonText={"Save changes"} onClick={handleUpdate}>
            <Form>
                <InputGroup className="mb-3">
                    <Form.Control
                        placeholder="Filename"
                        defaultValue={file.name}
                        ref={inputUpdate}
                    />
                    <InputGroup.Text>.{file.extesion}</InputGroup.Text>
                </InputGroup>
            </Form>
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
                    <Button text="Update" className="file-actions-item" onClick={showUpdateModal}/>
                    <Button text="Share" className="file-actions-item" onClick={showShareModal}/>
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