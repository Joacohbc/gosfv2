import { showError, showSuccess, showInfo} from "/static/modules/message.js";
import { getToken } from '/static/modules/request.js';

class File {
    constructor(id, filename, shared) {
        this._id = id;
        this._filename = filename;
        this._shared = shared;
    }

    static fromJSON(json) {
        return new File(json.id, json.filename, json.shared);
    }
    
    get id() {
        return this._id;
    }

    get filename() {
        return this._filename;
    }

    get shared() {
        return this._shared;
    }

    get link() {
        return `/api/files/${this._id}`;
    }

    open() {
        window.open(`${this.link}`, '_blank');
    }

    download() {
        axios.get(`/api/files/${this._id}`, {
            responseType: 'blob'
        })
        .then(blob => {
            let url = window.URL.createObjectURL(blob.data);
            let a = document.createElement('a');
            a.href = url;
            a.download = this._filename;
            a.click();
        })
        .catch(err => {
            showError(err.response.data.message);
        });

    }

    delete() {
        axios.delete(`/api/files/${this._id}`)
        .then(res => {
            reloadTable();
            showSuccess(res.data.message);
        })
        .catch(err => {
            console.log(err);
            showError(err.response.data.message);
        });
    }

    update(filename, shared) {
        axios.put(`/api/files/${this._id}`, {
            filename: filename,
        })
        .then(res => {
            showSuccess(res.data.message);
            reloadTable();
        })
        .catch(err => {
            showError(err.response.data.message);
        });
    }

    getShared() {
        if(this._shared === false) {
            alert("This file not shared!");
            return;
        }

        alert(`${window.location.origin}/api/files/share/${this._id}`);
    }

    getHTML() {
        const part = document.createElement('tr');
        part.setAttribute("class", "file");

        const id = document.createElement('td');
        id.classList.add('file-id');
        id.innerText = this._id;
        part.appendChild(id);

        const filename = document.createElement('td');
        filename.classList.add('file-filename');
        filename.innerHTML = this._filename;
        filename.addEventListener('click', (e) => {
            e.preventDefault();
            this.open();
        });
        part.appendChild(filename);
        
        const actions = document.createElement('td');
        actions.classList.add('file-actions');

        // La propiedad que tendrán los botones
        const btnAttribute = 'file-actions-item';

        const link = document.createElement('button');
        link.classList.add(btnAttribute);
        link.classList.add('file-download-btn');
        link.innerHTML = 'Download';
        link.addEventListener('click', (e) => {
            e.preventDefault();
            this.download();
        });
        actions.appendChild(link);

        const deleteBtn = document.createElement('button');
        deleteBtn.classList.add(btnAttribute);
        deleteBtn.classList.add('file-delete-btn');
        deleteBtn.innerHTML = 'Delete';
        deleteBtn.addEventListener('click', (e) => {
            e.preventDefault();
            this.delete();
        });
        actions.appendChild(deleteBtn);

        const updateBtn = document.createElement('button');
        updateBtn.classList.add(btnAttribute);
        updateBtn.classList.add('file-update-btn');
        updateBtn.innerHTML = 'Update';
        updateBtn.addEventListener('click', (e) => {
            e.preventDefault();

            let filename = prompt("Enter the new filename", this._filename);
            if(filename === null || filename == "") {
                showInfo("You have canceled the update");
                return;
            }
            this.update(filename);
        });
        actions.appendChild(updateBtn);
        
        part.appendChild(actions);

        return part;
    }
}

function reloadTable() {
    axios.get("/api/files/")
    .then(req => {
        document.querySelector('tbody').innerHTML = '';

        // Si no hay archivos que ingrese un mensaje personalizado
        // en el head de la tabla que indique que no hay archivos
        if(req.data == null) {
            document.querySelector('thead').innerHTML = "No files, start uploading files c:";
            return;
        }

        // Si hay archivos que agregue las columnas a la tabla
        document.querySelector('thead').innerHTML = `
        <tr>
            <th>ID</th>
            <th>Filename</th>
            <th>Actions</th>
        </tr>
        `;

        // Y cargue las filas de la tabla
        const files = document.createDocumentFragment();
        req.data.forEach(element => {
            files.appendChild(File.fromJSON(element).getHTML());
        });
        document.querySelector('tbody').appendChild(files);
    })
    .catch(err => {
        showError(err.response.data.message);
    });
}

window.addEventListener('DOMContentLoaded', function() {

    // Verifico que el Token este en la Cookie
    if (getToken() === null) {
        // Si no esta lo devuelvo al Login
        window.location.href = '/static/login/login.html';
    }

    // Defino el URL base de las peticiones
    axios.defaults.baseURL = window.location.origin;

    // Cargo la tabla de archivos
    reloadTable();

    // Agrego funcionamiento a los botones
    document.querySelector("#btn-logout").addEventListener('click', (e) => {
        e.preventDefault();
        axios.post('/logout')
        .then(res => {
            window.location.href = '/static/index.html';
        })
        .catch(err => {
            alert(err.response.data.message);
        });
    });

    document.querySelector("#input-upload").addEventListener('change', (e) => {
        e.preventDefault();

        const files = e.target.files;
        if(files.length == 0) {
            return;
        }

        const form = new FormData();
        for(let i = 0; i < files.length; i++) {
            form.append('files', files[i]);
        }

        axios.post('/api/files/', form)
        .then(req => {
            setTimeout(() => {
                reloadTable();
            });
            showSuccess(req.data.message);
        })
        .catch(err => {
            showError(err.response.data.message); 
        });
    });

    document.querySelector("#btn-share-info").addEventListener('click', (e) => {
        e.preventDefault();
        alert("Working in progress!");
    });

    document.querySelector("#btn-user-info").addEventListener('click', (e) => {
        e.preventDefault();
        alert("Working in progress!");
    });
});
