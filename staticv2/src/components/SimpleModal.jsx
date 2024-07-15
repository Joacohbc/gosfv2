import "./SimpleModal.css";
import Modal from "react-bootstrap/Modal";
import Button from "./Button";
import "bootstrap/dist/css/bootstrap.min.css";
import PropTypes from "prop-types";
import { forwardRef, useImperativeHandle, useState } from "react";

const SimpleModal = forwardRef((props, ref) => {
    const [ showed, setShowed ] = useState(false);

    useImperativeHandle(ref, () => ({
        setShowed: (value) => setShowed(value),
        show: () => setShowed(true),
        hide: () => setShowed(false), 
        toggleShowed: () => setShowed(!showed),
        isShowed: () => showed,
    }));
    
    const handleClose = () => {
        setShowed(false);
    };

    return (
        <Modal
            show={showed}
            onHide={handleClose}
            size={props.size || "md"}
            enforceFocus
            scrollable
            autoFocus
        >
            <Modal.Header closeButton closeVariant="white">
                <Modal.Title className="text-center file-filename">{props.title}</Modal.Title>
            </Modal.Header>
            <Modal.Body>{props.children}</Modal.Body>
            {props.buttonText && props.onClick && <Modal.Footer className="d-flex justify-content-center">
                <Button
                    text={props.buttonText}
                    onClick={props.onClick}
                    className={"p-2"}
                />
            </Modal.Footer>}
        </Modal>
    );
});

SimpleModal.displayName = 'Modal';

SimpleModal.propTypes = {
    title: PropTypes.string,
    children: PropTypes.node,
    size: PropTypes.string,
    buttonText: PropTypes.any,
    onClick: PropTypes.func,
};

export default SimpleModal;
