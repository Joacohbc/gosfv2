import PropTypes from 'prop-types';
import { memo } from 'react';

const PreviewFile = memo(({ url, contentType, className }) => {
    return <>
        { !url && <h1>404</h1> }
        { contentType.includes('video') ? <video className={className} controls><source src={url} type={contentType}/></video> 
        : <iframe src={url} className={className}/> }
    </>;
});

PreviewFile.displayName = 'PreviewFile';

PreviewFile.propTypes = {
    url: PropTypes.string,
    contentType: PropTypes.string,
    className: PropTypes.string
}

export default PreviewFile;