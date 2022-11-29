import { showError, showSuccess} from "/static/modules/message.js";
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
        axios.get(`/api/files/${this._id}`, {
            responseType: 'blob'
        })
        .then(blob => {
            let url = window.URL.createObjectURL(blob.data);
            window.open(url, '_blank');
        })
        .catch(err => {
            showError(err.response.data.message);
        });
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
            shared: shared
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

        const open = document.createElement('button');
        open.classList.add(btnAttribute);
        open.classList.add('file-open-btn');
        open.innerHTML = 'Open';
        open.addEventListener('click', (e) => {
            e.preventDefault();
            this.open();
        });
        actions.appendChild(open);

        const getShared = document.createElement('button');
        getShared.classList.add(btnAttribute);
        getShared.classList.add('file-share-btn');
        getShared.innerHTML = 'Share';
        getShared.addEventListener('click', (e) => {
            e.preventDefault();
            this.getShared();
        });
        actions.appendChild(getShared);

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
            let shared = confirm("Do you want to share this file?");
            if(filename === null) {
                return;
            }

            this.update(filename, shared);
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

        const files = document.createDocumentFragment();
        req.data.forEach(element => {
            files.appendChild(File.fromJSON(element).getHTML());
        });
        document.querySelector('tbody').appendChild(files);
    })
    .catch(err => {
        console.log(err);
        if(err.request.status == 401) {
            alert(err.response.data.message);
            window.location.href = '/static/login/login.html';
        }
    });
}

window.addEventListener('DOMContentLoaded', function() {
    // Si tengo el Token, lo agrego al header de las peticiones
    axios.defaults.headers.Authorization = "Bearer " + getToken();
    axios.defaults.baseURL = window.location.origin;

    reloadTable();

    document.querySelector(`#btn-logout`).addEventListener('click', (e) => {
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
});
