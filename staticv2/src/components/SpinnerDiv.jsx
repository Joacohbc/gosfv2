import { memo } from 'react';
import Spinner from 'react-bootstrap/Spinner';
import PropTypes from 'prop-types';

const SpinnerDiv = memo(({isLoading, children}) => {

    return <div>
        <div className="d-flex justify-content-center align-items-center">
            <Spinner animation="border" role="status" hidden={!isLoading}/>
        </div>
        
        <div hidden={isLoading}>
            {children}
        </div>
    </div>;
});

SpinnerDiv.propTypes = {
    isLoading: PropTypes.bool.isRequired,
    children: PropTypes.element.isRequired
};

SpinnerDiv.displayName = 'SpinnerDiv';
export default SpinnerDiv;