import 'bootstrap/dist/css/bootstrap.min.css';
import './Message.css';
import { forwardRef, useImperativeHandle, useState } from 'react';

const colors = {
    info: 'blue',
    error: 'red',
    success: 'green',
    warning: 'yellow'
};

const Message = forwardRef((prop, ref) => {
    const [ message, setMessage ] = useState('');    
    const [ color, setColor ] = useState('');

    function showMessage(message, color) {
        setColor(color);
        setMessage(message[0].toUpperCase() + message.substring(1));
        setTimeout(() => {
            setMessage('');
        }, 1500);
    }

    useImperativeHandle(ref, () => ({
        showError(message) {
            showMessage(message, colors.error);
        },
        showSuccess(message) {
            showMessage(message, colors.success);
        },
        showWarning(message) {
            showMessage(message, colors.warning);
        },
        showInfo(message) {
            showMessage(message, colors.info);
        }
    }));

    return <div className="fixed-bottom p-3 message" style={{ color: color}}>{message}</div>
});

Message.displayName = 'Message';

export default Message;