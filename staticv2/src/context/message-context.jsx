import { createContext } from 'react';
import { Toaster, toast } from 'sonner'
import PropTypes from 'prop-types';

export const MessageContext = createContext({
    showError: (message) => {},
    showInfo: (message) => {},
    showWarning: (message) => {},
    showSuccess: (message) => {},
    showPromise: (promise, loadingMessage, successMessage, errorMessage) => {},
    showAction: (message, label, onClick) => {},
    dismiss: (id) => {},
});

const MessageComponentProvider = (props) => {
    const showError = (message) => {
        return toast.error(message);
    };

    const showInfo = (message) => {
        return toast.info(message);
    };

    const showWarning = (message) => {
        return toast.warning(message);
    };

    const showSuccess = (message) => {
        return toast.success(message);
    };

    const showPromise = (promise, loadingMessage, successMessage, errorMessage) => {
        return toast.promise(promise, {
            loading: loadingMessage,
            success: successMessage,
            error: errorMessage,
        });
    }

    const showAction = (message, label, onClick) => {
        return toast(message, {
            action: {
                label: label,
                onClick: onClick
            }
        });
    }

    return (
        <MessageContext.Provider value={{
            showError,
            showInfo,
            showWarning,
            showSuccess,
            showPromise,
            showAction,
            dismiss: toast.dismiss,
        }}>
            <Toaster position="bottom-center" richColors/>
            {props.children}
        </MessageContext.Provider>
    );
};

MessageComponentProvider.propTypes = {
    children: PropTypes.node.isRequired,
};


export default MessageComponentProvider;