import { Message } from "/static/modules/message.js";
import { createOverlay } from "/static/main/overlay.js";
const message = new Message("message");

class File {
    constructor(id, filename, shared, shared_with) {
        this._id = id;
        this._filename = filename;
        this._shared = shared;
        this._shared_with = shared_with;
    }

    static searchById(fileId) {
        axios.get(`/api/files/${fileId}/info`)
        .then(res => {
            return new File(fileId, res.data.filename, res.data.shared, res.data.shared_with);
        })
        .catch(err => {
            message.showError(err.response.data.message);
            return null;
        });
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
            message.showError(err.response.data.message);
        });

    }

    delete() {
        axios.delete(`/api/files/${this._id}`)
        .then(res => {
            reloadTable();
            message.showSuccess(res.data.message);
        })
        .catch(err => {
            console.log(err);
            message.showError(err.response.data.message);
        });
    }

    update(filename, shared = this._shared) {
        axios.put(`/api/files/${this._id}`, {
            filename: filename,
            shared: shared
        })
        .then(res => {
            message.showSuccess(res.data.message);
            reloadTable();
        })
        .catch(err => {
            message.showError(err.response.data.message);
        });
    }

    getShared() {
        navigator.clipboard.writeText(`${window.location.origin}/api/files/share/${this._id}`)            
        .then(() => {
            message.showInfo("The link has been copied to the clipboard");
        })
        .catch(err => {
            message.showError("Error copying the link");
            console.log(err);
        });
    }

    toTableRow() {
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

        const download = document.createElement('button');
        download.classList.add(btnAttribute);
        download.classList.add('file-download-btn');
        download.innerHTML = 'Download';
        download.addEventListener('click', (e) => {
            e.preventDefault();
            this.download();
        });
        actions.appendChild(download);

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
                message.showInfo("You have canceled the update");
                return;
            }
            this.update(filename);
        });
        actions.appendChild(updateBtn);
        
        const shareBtn = document.createElement('button');
        shareBtn.classList.add(btnAttribute);
        shareBtn.classList.add('file-share-btn');
        shareBtn.innerHTML = 'Share';
        shareBtn.addEventListener('click', (e) => {
            e.preventDefault();
            createOverlay(this);
        });
        actions.appendChild(shareBtn);


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
            files.appendChild(File.fromJSON(element).toTableRow());
        });
        document.querySelector('tbody').appendChild(files);
    })
    .catch(err => {
        message.showError(err.response.data.message);
    });
}

window.addEventListener('DOMContentLoaded', function() {
    
    // Cargo la tabla de archivos
    reloadTable();

    document.getElementById("input-upload").addEventListener('change', (e) => {
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
            message.showSuccess(req.data.message);
        })
        .catch(err => {
            message.showError(err.response.data.message); 
        });
    });
});
