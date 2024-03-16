import Form from 'react-bootstrap/Form';
import { handleKeyUpWithTimeout } from '../utils/input-text';
import { useEffect } from 'react';
import { useNotes } from '../hooks/notes';
import { useState, useContext } from 'react';
import { MessageContext } from '../context/message-context';
import AuthContext from '../context/auth-context';
import Button from '../components/Button';

export default function Notes() {
    const [ text, setText ] = useState('');
    const { getNote, setNote } = useNotes({});
    const { isLogged } = useContext(AuthContext);
    const messageContext = useContext(MessageContext);

    useEffect(() => {
        if(!isLogged) return;

        getNote().then((note) => {
            setText(note.content);
            messageContext.showSuccess(note.message);
        }).catch((err) => {
            messageContext.showError(err.message)
        });
    }, [ getNote, messageContext, isLogged ]);

    const copyLink = async (e) => {
        e.preventDefault();
        try {
            const note = await getNote();
            await navigator.clipboard.writeText(note.content);
            messageContext.showSuccess('Copied to clipboard');
        } catch(err) {
            messageContext.showError(err.message);
        }
    }

    const onTextChange = handleKeyUpWithTimeout((e) => {
        setNote(e.target.value).then((res) => {
            messageContext.showSuccess(res.message);
        }).catch((err) => {
            messageContext.showError(err.message);
        });
    }, 500);

    return <div className='d-flex justify-content-center d-flex align-items-center flex-column m-3'>
        <Form.Control 
            as="textarea" 
            placeholder="Leave your note here"
            style={{ height: '25em', maxWidth: '50em'}} 
            onKeyUp={onTextChange}
            defaultValue={text}/>
        <Button onClick={copyLink} text={<span>{'Copy Link'} <i className='bi bi-clipboard-fill'/></span>} className="p-2 mt-2 rounded border border-white text-white"/>
    </div>
}
