import { showError, showSuccess, showInfo} from "/static/modules/message.js";

class User {
    constructor(id,  username) {
        this._id = id;
        this._username = username;
    }

    static fromJSON(json) {
        return new User(json.id, json.username);
    }

    get id() {
        return this._id;
    }

    get name() {
        return this._username;
    }

    getRow() {
        const part = document.createElement('tr');
        part.setAttribute("class", "file");

        const id = document.createElement('td');
        id.classList.add('user-id');
        id.innerText = this._id;
        part.appendChild(id);

        const username = document.createElement('td');
        username.classList.add('user-username');
        username.innerText = this._username;
        part.appendChild(username);

        return part;
    }

    getRowWithButton(fileId) {
        const part = this.getRow();

        const actions = document.createElement('td');
        actions.classList.add('file-actions');
        
        // La propiedad que tendrán los botones
        const btnAttribute = 'file-actions-item';
        const removeUser = document.createElement('button');
        removeUser.classList.add(btnAttribute);
        removeUser.classList.add('user-stop-share-btn');
        removeUser.innerHTML = 'Remove';
        removeUser.addEventListener('click', () => {
            axios.delete(`/api/files/share/${fileId}/user/${this._id}`)
            .then(res => {
                showSuccess(res.data.message);
                reloadTable(fileId);
            })
            .catch(err => {
                showError(err.response.data.message);
            });
        });
    
        actions.appendChild(removeUser);
        
        part.appendChild(actions);
        return part;
    }      
}


function reloadTable(fileId) { 
    const thead = document.querySelector('thead');
    thead.innerHTML = '';

    const tbody = document.querySelector('tbody');
    tbody.innerHTML = '';
    
    const addUser = document.querySelector("#add-user-btn");
    addUser.setAttribute("hidden", "true");

    axios.get(`/api/files/${fileId}/info`)
    .then(req => {
        addUser.removeAttribute("hidden");
        
        // Si no hay archivos que ingrese un mensaje personalizado
        // en el head de la tabla que indique que no hay archivos
        if(req.data.shared_with == null) {
            thead.innerHTML = "The file is not shared with anyone";
            return;
        }

        // Si hay archivos que agregue las columnas a la tabla
        thead.innerHTML = `
        <tr>
            <th>ID</th>
            <th>Username</th>
            <th>Actions</th>
        </tr>
        `;

        // Y cargue las filas de la tabla
        const users = document.createDocumentFragment();
        req.data.shared_with.forEach(user => {
            users.appendChild(User.fromJSON(user).getRowWithButton(fileId));
        });

        tbody.appendChild(users);
    }).catch(err => {
        showError(err.response.data.message);
    });
}

window.addEventListener('DOMContentLoaded', () => {

    document.querySelector("#input-file-id").addEventListener('change', (e) => { 
        if(e.target.value == "" || e.target.value <= e.target.min) {
            showError("Please enter a file id");
            return;
        }

        reloadTable(e.target.value);
    });


    document.querySelector("#btn-add-user").addEventListener('click', () => {
        const fileId = document.querySelector("#input-file-id").value;
        const userId = document.querySelector("#input-user-id").value;
        if(userId == "") {
            showError("Please enter a User ID");
            return;
        }
    
        axios.post(`/api/files/share/${fileId}/user/${userId}`)
        .then(res => {
            showSuccess(res.data.message);
            const userId = document.querySelector("#input-user-id").value = "";
            reloadTable(fileId);
        })
        .catch(err => {
            showError(err.response.data.message);
        });
    });
});
