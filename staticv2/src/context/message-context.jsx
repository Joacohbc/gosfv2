import { createContext, useRef } from 'react';
import Message from '../components/Message';
import PropTypes from 'prop-types';

export const MessageContext = createContext({
    showError: (message) => {},
    showInfo: (message) => {},
    showWarning: (message) => {},
    showSuccess: (message) => {}
});

const MessageComponentProvider = (props) => {
    const messageRef = useRef(null);

    const showError = (message) => {
        messageRef.current.showError(message);
    };

    const showInfo = (message) => {
        messageRef.current.showInfo(message);
    };

    const showWarning = (message) => {
        messageRef.current.showWarning(message);
    };

    const showSuccess = (message) => {
        messageRef.current.showSuccess(message);
    };

    return (
        <MessageContext.Provider value={{
            showError,
            showInfo,
            showWarning,
            showSuccess
        }}>
            <Message ref={messageRef}/>
            {props.children}
        </MessageContext.Provider>
    );
};

MessageComponentProvider.propTypes = {
    children: PropTypes.node.isRequired,
};


export default MessageComponentProvider;