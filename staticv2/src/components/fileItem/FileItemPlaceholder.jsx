import 'bootstrap/dist/css/bootstrap.min.css';
import Card from 'react-bootstrap/Card';
import './FileItem.css';
import Placeholder from 'react-bootstrap/Placeholder';
import 'bootstrap-icons/font/bootstrap-icons.css'
import { memo } from 'react';

function random(min, max) { 
    return Math.floor(Math.random() * (max - min + 1) + min);
}

const FileItemPlaceholder = memo(() => {
    return <>
        <Card className='file'>
            <Card.Body>
                <Placeholder as={Card.Title} animation="wave">
                    <Placeholder xs={6} />
                </Placeholder>
                <Placeholder as={Card.Text} animation="wave">
                    {Array.from({ length: random(3, 7) }).map((_, i) => (
                        <Placeholder key={i} xs={random(4, 8)} />
                    ))}
                </Placeholder>
                <Placeholder as="div" animation="wave" className='text-center'>
                    <Placeholder as="button" className='file-actions-item' xs={1}><i className='bi bi-file-arrow-down-fill'/></Placeholder>
                    <Placeholder as="button" className='file-actions-item' xs={1}><i className='bi bi-trash3-fill'/></Placeholder>
                    <Placeholder as="button" className='file-actions-item' xs={1}><i className='bi bi-pencil-square'/></Placeholder>
                    <Placeholder as="button" className='file-actions-item' xs={1}><i className='bi bi-share-fill'/></Placeholder>
                    <Placeholder as="button" className='file-actions-item' xs={1}><i className='bi bi-folder-check file-actions-item-no-hover'/></Placeholder>
                </Placeholder>
            </Card.Body>
        </Card>
    </>
});

FileItemPlaceholder.displayName = 'FileItemPlaceholder';

export default FileItemPlaceholder;
