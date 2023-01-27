import { Message } from "/static/modules/message.js";
import { User } from "/static/models/user.js";

const message = new Message("message");

const getFullUsername = (username, id) => `${username} #${id}`;

// Obtengo la información del usuario actual
async function getUserInfo() {
    try {
        const res = await axios.get("/api/users/me");
        return User.fromJSON(res.data);
    } catch (err) {
        throw new Error(message.showError(err.response.data.message));
    }
}

function updateIcon() {
    // Agrego el timestamp para que se actualice la imagen (no la saque del Cache)
    document.getElementById("user-icon").src = window.location.origin + "/api/users/icon/me?" + new Date().getTime();
}

window.addEventListener("DOMContentLoaded", () => {

    if(localStorage.getItem("username") == null || localStorage.getItem("userId") == null) {
        getUserInfo().then((user) => {
            localStorage.setItem("userId", user.id);
            localStorage.setItem("username", user.username);
        });
    }

    updateIcon();
    let lsUsername = localStorage.getItem("username");
    let lsUserId = localStorage.getItem("userId");

    document.getElementById("user-username").innerHTML = getFullUsername(lsUsername, lsUserId);
    document.getElementById("username").value = lsUsername;

    document.getElementById('copy-id').addEventListener("click", () => {
        // Copy content to clipboard
        navigator.clipboard.writeText(document.getElementById("user-username").innerHTML)
        .then(() => {
            message.showSuccess("ID Copied to clipboard");
        })
        .catch(err => {
            message.showError("Error copying ID to clipboard");
        });
    });

    document.getElementById("btn-rename").addEventListener("click", (e) => {
        e.preventDefault();

        const newUsername = document.getElementById("username").value;
        axios.put("/api/users/rename", {
            username: newUsername,
        })
        .then((res) => {
            document.getElementById("user-username").innerHTML = getFullUsername(newUsername, lsUserId);
            localStorage.setItem("username", newUsername);
            message.showSuccess(res.data.message);
        }).catch((err) => {
            console.log(err);
            message.showError(err.response.data.message);
        });
    });

    // Cambio de contraseña
    document.getElementById("btn-change-password").addEventListener("click", (e) => {
        e.preventDefault();

        const oldPassword = document.getElementById("current-password").value;
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
            document.getElementById("current-password").value = "";
            document.getElementById("new-password").value = "";
            document.getElementById("confirm-password").value = "";
        }).catch((err) => {
            console.log(err);
            message.showError(err.response.data.message);
        });
    });

    // Subida de icono
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

    // Eliminación de icono
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

    // Eliminación de cuenta
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