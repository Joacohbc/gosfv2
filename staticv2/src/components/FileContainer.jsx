import 'bootstrap/dist/css/bootstrap.min.css';
import Container from 'react-bootstrap/Container';
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import FileItem from '../components/fileItem/FileItem';
import { memo } from 'react';
import PropTypes from 'prop-types';
import FileItemPlaceholder from './fileItem/FileItemPlaceholder';
import { useCache } from '../hooks/useCache';

/**
 * Renderiza un contenedor para mostrar archivos.
 *
 * @component
 * @param {Object[]} files - El array de archivos que se mostrarán.
 * @param {number} progress - El número de archivos que se han cargado.
 * @param {function} fileLoader - La función para cargar más archivos.
 * @param {React.Ref} loading - Indica si los archivos se están cargando.
 * @param {function} handleOpenPreview - La función para manejar la apertura de una vista previa de archivo.
 * @param {function} handleFilesDelete - La función para manejar la eliminación de archivos.
 * @param {function} handleFilesUpdate - La función para manejar la actualización de archivos.
 * @returns {JSX.Element} El componente FileContainer.
 */
const FileContainer = memo(({ files, progress = files.length, fileLoader = () => {}, loading, handleOpenPreview, handleFilesDelete, handleFilesUpdate }) => {
    const { cacheService } = useCache();
    const numLoadingFiles = cacheService.getCacheNumberOfFiles().value || 15;

    return (
        <Container fluid="md">
            <Col hidden={loading}>
                { files.length == 0 && 
                    <div className="d-flex justify-content-center align-items-center">
                        { files.length == 0 && <p className="text-center text-white">No files, start uploading files c:</p> }
                    </div>
                }
            </Col> 
            
            <Row xs={1} sm={2} md={3} lg={4} xl={5} className='row-gap-3 d-flex justify-content-center'>
                { !loading && files.slice(0, progress).map(file => 
                <Col key={file.id} hidden={file.deleted}>
                    <FileItem 
                        id={file.id}
                        filename={file.filename}
                        name={file.name} 
                        contentType={file.contentType}
                        extension={file.extension}
                        url={file.url}
                        shared={file.shared}
                        sharedWith={file.sharedWith}
                        savedLocal={file.savedLocal}
                        onOpen={handleOpenPreview}
                        onDelete={handleFilesDelete}
                        onUpdate={handleFilesUpdate}
                    />
                </Col>) }

                { loading && Array.from({ length: numLoadingFiles }).map((_, i) =>
                <Col key={i}>
                    <FileItemPlaceholder/>
                </Col>) }
            </Row>


            
            { files.length != 0 && files.length != progress && 
            <Row className='p-3'>
                <button className="btn btn-load" onClick={fileLoader}>
                    <i className="bi bi-arrow-down-square-fill" />
                </button>
            </Row> }
        </Container>
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