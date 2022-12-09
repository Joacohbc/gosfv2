import { Message } from "/static/modules/message.js";

const message = new Message("message");

window.addEventListener("DOMContentLoaded", () => {

    axios.get("/api/users/me")
    .then((res) => {
        document.getElementById("user-username").innerHTML = res.data.username;
        document.getElementById("username").value = res.data.username;
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
});