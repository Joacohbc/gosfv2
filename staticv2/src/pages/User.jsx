import './User.css';
import Form from 'react-bootstrap/Form';
import Button from '../components/Button';
import Stack from 'react-bootstrap/Stack';
import { useUsers } from '../hooks/users';
import { useEffect, useState, useContext, useRef } from 'react';
import { MessageContext } from '../context/message-context';
import AuthContext from '../context/auth-context';
import ConfirmDialog from '../components/ConfirmDialog';

export default function User() {
    const messageContext = useContext(MessageContext);
    const { onLogOut, isLogged } = useContext(AuthContext);
    const { getMyIconURL, getMyInfo, uploadIcon, updateUser, changePassword, deleteIcon, deleteAccount } = useUsers();
    const [ iconURL, setIconURL ] = useState('');
    const [ userInfo, setUserInfo ] = useState({});
    const newUsername = useRef(null);
    const currentPassword = useRef(null);
    const newPassword = useRef(null);
    const confirmPassword = useRef(null);
    const deleteAccountDialog = useRef(null);
    
    useEffect(() => {
        if(!isLogged) return;
        
        getMyInfo().then(data => setUserInfo(data))
        setIconURL(getMyIconURL());
    }, [ getMyInfo, messageContext, getMyIconURL, isLogged ])

    const handleCopyUserId = async (e) => {
        e.preventDefault();
        try {
            await navigator.clipboard.writeText(`${userInfo.username} #${userInfo.id}`);
            messageContext.showSuccess("Link copied to clipboard");
        } catch(err) {
            messageContext.showError(err.message);
        }
    }

    const handleUpdateUser = async (e) => {
        e.preventDefault();

        try {
            if(newUsername.current.value === userInfo.username) {
                messageContext.showWarning('Username is the same');
                return;
            }
            const res = await updateUser(newUsername.current.value);
            setIconURL(getMyIconURL(true));
            messageContext.showSuccess(res.message);
        } catch(err) {
            messageContext.showError(err.message);
        }
    }

    const handleUpdatePassword = async (e) => {
        e.preventDefault();

        if(newPassword.current.value !== confirmPassword.current.value) {
            messageContext.showError('Passwords do not match');
            return;
        }

        try {
            const res = await changePassword(currentPassword.current.value, newPassword.current.value);
            messageContext.showSuccess(res.message);
        } catch(err) {
            messageContext.showError(err.message);
        }
    }

    const handleUploadIcon = async (e) => {
        e.preventDefault();
        
        try {
            const res = await uploadIcon(e.target.files[0]);
            messageContext.showSuccess(res.message);

            // Call random to avoid cache issues
            setIconURL(getMyIconURL(true));
        } catch(err) {
            messageContext.showError(err.message);
        }
    }

    const handleDeleteIcon = async (e) => {
        e.preventDefault();

        try {
            const res = await deleteIcon();
            messageContext.showSuccess(res.message);
            setIconURL(getMyIconURL(true));
        } catch(err) {
            messageContext.showError(err.message);
        }
    }

    const handleDeleteAccount = async () => {
        try {
            await deleteAccount();
            onLogOut();
        } catch(err) {
            messageContext.showError(err.message);
        }
    };

    const handleDeleteAccountDialog = (e) => {
        e.preventDefault();
        deleteAccountDialog.current.show();
    }

    return <>
    <ConfirmDialog 
        ref={deleteAccountDialog} 
        onOk={handleDeleteAccount}
        title="Delete account"
        message="Are you sure you want to delete your account? This will be permanent"/>

    <div className='d-flex justify-content-center'>
        <div className="user-info">
            <div className="user-username">{userInfo.username}</div>
            <div className="copy-id" onClick={handleCopyUserId}>Click to Copy ID</div>
            <img className="user-icon" src={iconURL}/>

            <div className="icon-options">  
                <label htmlFor="icon-upload" className="btn-icon">Upload icon <i className='bi bi-camera'/></label>
                <input id="icon-upload" type="file" style={{display: "none"}} onChange={handleUploadIcon} required />
                <button className="btn-icon" onClick={handleDeleteIcon}>Remove icon <i className='bi bi-backspace'/></button>
            </div>

            <Stack gap={1}>
                <Form.Control type="text" placeholder="Username" defaultValue={userInfo.username} ref={newUsername}/>
                <Button className="btn-form" text={<span>{'Update'} <i className='bi bi-pencil-square'/></span>} onClick={handleUpdateUser}/>
            </Stack>

            <Stack gap={1}>
                <Form.Control type="password" placeholder="Current Password" ref={currentPassword}/>
                <Form.Control type="password" placeholder="Confirm Password" ref={confirmPassword}/>
                <Form.Control type="password" placeholder="New Password" ref={newPassword}/>
                <Button className="btn-form" text={<span>{'Change Password'} <i className='bi bi-house-lock-fill'/></span>} onClick={handleUpdatePassword}/>
                <Button className="btn-form" text={<span>{'Delete Account'} <i className='bi bi-trash3'/></span>} onClick={handleDeleteAccountDialog}/>
            </Stack>
        </div>
    </div>
    </>
}
