import "./SimpleModal.css";
import Modal from "react-bootstrap/Modal";
import Button from "./Button";
import "bootstrap/dist/css/bootstrap.min.css";
import PropTypes from "prop-types";
import { forwardRef, useImperativeHandle, useState } from "react";

const ConfirmDialog = forwardRef((props, ref) => {
    const [ showed, setShowed ] = useState(false);

    useImperativeHandle(ref, () => ({
        show: () => setShowed(true),
        hide: () => setShowed(false), 
        toogleShowed: () => setShowed(!showed),
        isShowed: () => showed,
    }));
    
    const handleClose = () => {
        setShowed(false);
    };

    const handleOk = () => {
        if(props.onOk) props.onOk()
        setShowed(false);
    };

    const handleCancel = () => {
        if(props.onCancel) props.onCancel();
        setShowed(false);
    };

    return (
        <Modal
            show={showed}
            onHide={handleClose}
            size={props.size || "sm"}
            enforceFocus
            autoFocus
        >
            <Modal.Header>
                <Modal.Title className="text-center">{props.title}</Modal.Title>
            </Modal.Header>
            <Modal.Body>{props.message}</Modal.Body>
            <Modal.Footer className="d-flex justify-content-center">
                <Button
                    text='Ok'
                    onClick={handleOk}
                    className={"p-2"}
                />
                <Button
                    text='Cancel'
                    onClick={handleCancel}
                    className={"p-2"}
                />
            </Modal.Footer>
        </Modal>
    );
});

ConfirmDialog.displayName = 'Modal';

ConfirmDialog.propTypes = {
    title: PropTypes.string,
    message: PropTypes.string,
    size: PropTypes.string,
    onOk: PropTypes.func,
    onCancel: PropTypes.func,
};

export default ConfirmDialog;
