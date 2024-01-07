import 'bootstrap/dist/css/bootstrap.min.css';
import Container from 'react-bootstrap/Container';
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import FileItem from '../components/fileItem/FileItem';
import SpinnerDiv from '../components/SpinnerDiv';
import '../components/Message.css';
import { memo } from 'react';
import PropTypes from 'prop-types';

/**
 * Renders a container for displaying files.
 *
 * @component
 * @param {Object[]} files - The array of files to be displayed.
 * @param {number} progress - The number of files that have been loaded.
 * @param {function} fileLoader - The function to load more files.
 * @param {React.Ref} loading - Indicates if the files are loading.
 * @param {function} handleOpenPreview - The function to handle opening a file preview.
 * @param {function} handleFilesDelete - The function to handle deleting files.
 * @param {function} handleFilesUpdate - The function to handle updating files.
 * @returns {JSX.Element} The FileContainer component.
 */
const FileContainer = memo(({ files, progress = files.length, fileLoader = () => {}, loading, handleOpenPreview, handleFilesDelete, handleFilesUpdate }) => {

    return (
        <SpinnerDiv isLoading={loading}>
        <Container fluid="md">
            <Col>
                <div className="d-flex justify-content-center align-items-center">
                    { files.length == 0 && <p className="text-center text-white">No files, start uploading files c:</p> }
                </div>
            </Col>
            
            <Row xs={1} sm={2}  md={3} lg={4} xl={5} className='row-gap-3 d-flex justify-content-center'>
            { files.slice(0, progress).map(file => 
                <Col key={file.id} hidden={file.deleted}>
                    <FileItem 
                        id={file.id}
                        filename={file.filename}
                        name={file.name} 
                        contentType={file.contentType}
                        extesion={file.extesion}
                        url={file.url}
                        savedLocal={file.savedLocal}
                        onOpen={handleOpenPreview}
                        onDelete={handleFilesDelete}
                        onUpdate={handleFilesUpdate}
                    />
                </Col>) }
            </Row>
            
            { files.length != 0 && files.length != progress && 
            <Row className='p-3'>
                <button className="btn btn-load" onClick={fileLoader}>
                    <i className="bi bi-arrow-down-square-fill" />
                </button>
            </Row> }
        </Container>
        </SpinnerDiv>
    )
})

FileContainer.propTypes = {
    files: PropTypes.array.isRequired,
    progress: PropTypes.number,
    fileLoader: PropTypes.func,
    loading: PropTypes.bool.isRequired,
    handleOpenPreview: PropTypes.func.isRequired,
    handleFilesDelete: PropTypes.func.isRequired,
    handleFilesUpdate: PropTypes.func.isRequired
};

FileContainer.displayName = 'FileContainer';
export default FileContainer;