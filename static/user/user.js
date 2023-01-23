import { Message } from "/static/modules/message.js";

const message = new Message("message");


function updateIcon() {
    // Agrego el timestamp para que se actualice la imagen (no la saque del Cache)
    document.getElementById("user-icon").src = window.location.origin + "/api/users/icon/me?" + new Date().getTime();
}

window.addEventListener("DOMContentLoaded", () => {

    axios.get("/api/users/me")
    .then((res) => {
        const username = `${res.data.username} #${res.data.id}`;
        document.getElementById("user-username").innerHTML =username;
        
        document.getElementById('copy-id').addEventListener("click", () => {
            // Copy content to clipboard
            navigator.clipboard.writeText(username)
            .then(() => {
                message.showSuccess("ID Copied to clipboard");
            })
            .catch(err => {
                message.showError("Error copying ID to clipboard");
            });
        });

        document.getElementById("username").value = res.data.username;
        updateIcon();
    }).catch((err) => {
        message.showError(err.response.data.message);
    });

    document.getElementById("btn-update").addEventListener("click", (e) => {
        e.preventDefault();

        const newUsername =  document.getElementById("username").value;
        axios.put("/api/users/rename", {
            username: newUsername,
        })
        .then((res) => {
            document.getElementById("user-username").innerHTML = newUsername;
            message.showSuccess(res.data.message);
        }).catch((err) => {
            console.log(err);
            message.showError(err.response.data.message);
        });
    });

    document.getElementById("btn-change-password").addEventListener("click", (e) => {
        e.preventDefault();

        const oldPassword = document.getElementById("old-password").value;
        const newPassword = document.getElementById("new-password").value;
        const confirmPassword = document.getElementById("confirm-password").value;

        if (newPassword !== confirmPassword) {
            message.showError("New password and confirm password do not match");
            return;
        }

        axios.put("/api/users/password", {
            old_password: oldPassword,
            new_password: newPassword,
        })
        .then((res) => {
            message.showSuccess(res.data.message);
        }).catch((err) => {
            console.log(err);
            message.showError(err.response.data.message);
        });
    });

    document.getElementById("icon-upload").addEventListener('change', (e) => {
        e.preventDefault();

        if(e.target.files.length == 0) {
            return;
        }

        const form = new FormData();
        form.append("icon", e.target.files[0]);

        axios.post('/api/users/icon', form)
        .then(req => {
            updateIcon();
            message.showSuccess(req.data.message);
        })
        .catch(err => {
            message.showError(err.response.data.message); 
        });

        // Reseteo el input para que se pueda subir el mismo archivo
        e.target.files = [];
    });


    document.getElementById("btn-delete-icon").addEventListener("click", (e) => {
        e.preventDefault();

        axios.delete("/api/users/icon")
        .then((res) => {
            updateIcon();
            message.showSuccess(res.data.message);
        }).catch((err) => {
            message.showError(err.response.data.message);
        });
    });

    document.getElementById("btn-delete-account").addEventListener("click", (e) => {
        e.preventDefault();

        if(confirm("Are you sure you want to delete your account?") == false) return;

        axios.delete("/api/users/")
        .then((res) => {
            window.location.href = "/static/login/login.html";
            message.showSuccess(res.data.message);
        }).catch((err) => {
            message.showError(err.response.data.message);
        });
    });
});