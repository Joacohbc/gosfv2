import { forwardRef, useImperativeHandle, useState } from 'react';
import Spinner from 'react-bootstrap/Spinner';

const SpinnerDiv = forwardRef((props, ref) => {
    const [ isLoading, setLoading ] = useState(false);
    
    useImperativeHandle(ref, () => ({
        loading: () => setLoading(true),
        stopLoading: () => setLoading(false),
    }));

    return <div>
        <div className="d-flex justify-content-center align-items-center">
            <Spinner animation="border" role="status" hidden={!isLoading}/>
        </div>
        
        <div hidden={isLoading}>
            {props.children}
        </div>
    </div>;
});

SpinnerDiv.displayName = 'SpinnerDiv';
export default SpinnerDiv;