import { createContext } from 'react';
import PropTypes from 'prop-types';
import useFilesIDB from '../hooks/useFilesIDB';

export const IDBContext = createContext({
    filesDB: null
});

const IDBContextProvider = (props) => {
    const filesDB = useFilesIDB();
    return (
        <IDBContext.Provider value={{ filesDB }}>
            {props.children}
        </IDBContext.Provider>
    );
};

IDBContextProvider.propTypes = {
    children: PropTypes.node.isRequired,
};


export default IDBContextProvider;