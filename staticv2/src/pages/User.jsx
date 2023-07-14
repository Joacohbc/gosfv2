import '../css/user.css';

export default function User() {
    return <>
    <div id="message"></div>

    <div className="container">

        <div id="user-info">
            <div id="user-username"></div>
            <div id="copy-id">Click to Copy ID</div>
            <img id="user-icon"/>
            
            <div className="icon-options">  
                <label htmlFor="icon-upload" className="btn-icon">Upload icon</label>
                <input id="icon-upload" type="file" style={{display: "none"}} required />
                <button id="btn-delete-icon" className="btn-icon">Remove icon</button>
            </div>
            
            <input type="text" id="username" placeholder="Username"/>
            <button id="btn-rename" className="btn-form">Update</button>

            <input type="password" id="current-password" placeholder="Current Password"/>
            <input type="password" id="new-password" placeholder="New Password"/>            
            <input type="password" id="confirm-password" placeholder="Confirm Password"/>

            <button id="btn-change-password" className="btn-form"> Changes Password</button>
            <button id="btn-delete-account" className="btn-form"> Delete Account</button>
        </div>

    </div>
    </>;
}
