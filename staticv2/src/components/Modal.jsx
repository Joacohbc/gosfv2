import "./Modal.css";
import BoostrapModal from "react-bootstrap/Modal";
import Button from "./Button";
import "bootstrap/dist/css/bootstrap.min.css";
import PropTypes from "prop-types";
import { forwardRef, useImperativeHandle, useState } from "react";

const Modal = forwardRef((props, ref) => {
    const [ showed, setShowed ] = useState(false);

    useImperativeHandle(ref, () => ({
        setShowed: (value) => setShowed(value),
        toogleShowed: () => setShowed(!showed),
    }));
    
    const handleClose = () => {
        setShowed(false);
    };

    return (
        <BoostrapModal
            show={showed}
            onHide={handleClose}
            className="background"
            backdrop="static"
        >
            <div className="overlay-background">
                <BoostrapModal.Header closeButton closeVariant="white">
                    <BoostrapModal.Title>{props.title}</BoostrapModal.Title>
                </BoostrapModal.Header>
                <BoostrapModal.Body>{props.children}</BoostrapModal.Body>
                <BoostrapModal.Footer className="d-flex justify-content-center">
                    <Button
                        text={props.buttonText}
                        onClick={props.onClick}
                        className={"p-2"}
                    />
                </BoostrapModal.Footer>
            </div>
        </BoostrapModal>
    );
});

Modal.displayName = 'Modal';

Modal.propTypes = {
    title: PropTypes.string,
    children: PropTypes.node,
    buttonText: PropTypes.string,
    onClick: PropTypes.func,
};

export default Modal;
