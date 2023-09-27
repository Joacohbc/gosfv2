import 'bootstrap/dist/css/bootstrap.min.css';
import './Message.css';
import { forwardRef, useImperativeHandle, useState } from 'react';

const colors = {
    info: '#3498db',
    error: '#e74c3c',
    success: '#2ecc71',
    warning: '#f1c40f'
};

const Message = forwardRef((prop, ref) => {
    const [ message, setMessage ] = useState('');    
    const [ color, setColor ] = useState('');

    function showMessage(message, color) {
        if(!message || !color) return;
        setColor(color);
        setMessage(message[0].toUpperCase() + message.substring(1));
        setTimeout(() => {
            setMessage('');
        }, 2000);
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

    return message.length > 0 && <div className="fixed-bottom message m-2" style={{ color: color}}>{message}</div>
});

Message.displayName = 'Message';

export default Message;