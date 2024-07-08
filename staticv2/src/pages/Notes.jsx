import './Notes.css';
import Form from 'react-bootstrap/Form';
import { useEffect } from 'react';
import { useNotes } from '../hooks/notes';
import { useState, useContext } from 'react';
import { MessageContext } from '../context/message-context';
import AuthContext from '../context/auth-context';
import Button from '../components/Button';

export default function Notes() {
    const [ text, setText ] = useState('');
    const [ saved, setSaved ] = useState({ saved: true, id: null});
    const { getNote, setNote } = useNotes({});
    const { isLogged } = useContext(AuthContext);
    const messageContext = useContext(MessageContext);

    useEffect(() => {
        if(!isLogged) return;

        getNote().then((note) => setText(note.content))
            .catch((err) => messageContext.showInfo(err));
    }, [ getNote, messageContext, isLogged ]);

    const handleCopy = async (e) => {
        e.preventDefault();
        try {
            await navigator.clipboard.writeText(text);
            messageContext.showInfo('Link copied to clipboard');
        } catch(err) {
            messageContext.showError(err.message);
        }
    }

    const onTextChange = (e) => {
        const id = setTimeout(() =>{
            setNote(e.target.value).then(() => setSaved({ saved: true, id: null }))
                .catch((err) => messageContext.showError(err.message));
        }, 2500);

        if(saved.id) clearTimeout(saved.id);
        setSaved({ saved: false, id});
    }
    return <div className='d-flex justify-content-center d-flex align-items-center flex-column m-3'>
        <Form.Control 
            as="textarea" 
            placeholder="Leave your note here"
            style={{ height: '25em', maxWidth: '50em'}} 
            onKeyUp={onTextChange}
            defaultValue={text}/>


        <div className='d-flex flex-row align-items-center'>
            <Button onClick={handleCopy} text={<span>{'Copy'} <i className='bi bi-clipboard-fill'/></span>} className="p-2 mt-2 rounded border border-white text-white"/>
            { !saved.saved ? <i className="bi bi-cloud-slash saving-cloud"></i> 
                : <i className="bi bi-cloud-check saved-cloud"></i> }
        </div>
    </div>
}
