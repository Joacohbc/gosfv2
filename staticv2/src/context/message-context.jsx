import { createContext } from 'react';
import { Toaster, toast } from 'sonner'
import PropTypes from 'prop-types';

export const MessageContext = createContext({
    showError: (message, duration) => {},
    showInfo: (message, duration) => {},
    showWarning: (message, duration) => {},
    showSuccess: (message, duration) => {},
    showPromise: (promise, loadingMessage, successMessage, errorMessage, duration) => {},
    showAction: (message, label, onClick, duration) => {},
    dismiss: (id) => {},
});

const DEFAULT_DURATION = 2000;

const MessageComponentProvider = (props) => {
    const showError = (message, duration = DEFAULT_DURATION) => {
        return toast.error(message, {
            duration: duration,
        });
    };

    const showInfo = (message, duration = DEFAULT_DURATION) => {
        return toast.info(message, {
            duration: duration,
        });
    };

    const showWarning = (message, duration = DEFAULT_DURATION) => {
        return toast.warning(message, {
            duration: duration,
        });
    };

    const showSuccess = (message, duration = DEFAULT_DURATION) => {
        return toast.success(message, {
            duration: duration,
        });
    };

    const showPromise = (promise, loadingMessage, successMessage, errorMessage, duration = DEFAULT_DURATION) => {
        return toast.promise(promise, {
            loading: loadingMessage,
            success: successMessage,
            error: errorMessage,
            duration: duration,
        });
    }

    const showAction = (message, label, onClick, duration = DEFAULT_DURATION) => {
        return toast(message, {
            action: {
                label: label,
                onClick: onClick
            },
            duration: duration,
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
            { window.screen.width < 768 && <Toaster position="top-center" richColors /> }
            { window.screen.width >= 768 && <Toaster position="bottom-left" richColors /> }
            {props.children}
        </MessageContext.Provider>
    );
};

MessageComponentProvider.propTypes = {
    children: PropTypes.node.isRequired,
};


export default MessageComponentProvider;