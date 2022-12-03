import { Message } from "/static/modules/message.js";
import { createOverlay } from "/static/main/overlay.js";
import { File } from "/static/models/file.js";


const message = new Message("message");

class FileCustom {

    constructor(fileJson) {
        this.file = File.fromJSON(fileJson);
    }

    open() {
        window.open(`${this.file.link}`, '_blank');
    }

    download() {
        axios.get(`/api/files/${this.file.id}`, {
            responseType: 'blob'
        })
        .then(blob => {
            let url = window.URL.createObjectURL(blob.data);
            let a = document.createElement('a');
            a.href = url;
            a.download = this.file.filename;
            a.click();
        })
        .catch(err => {
            message.showError(err.response.data.message);
        });

    }

    delete() {
        axios.delete(`/api/files/${this.file.id}`)
        .then(res => {
            reloadTable();
            message.showSuccess(res.data.message);
        })
        .catch(err => {
            console.log(err);
            message.showError(err.response.data.message);
        });
    }

    update(filename, shared = this.file.shared) {
        axios.put(`/api/files/${this.file.id}`, {
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
        navigator.clipboard.writeText(`${window.location.origin}/api/files/share/${this.file.id}`)            
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
        id.innerText = this.file.id;
        part.appendChild(id);

        const filename = document.createElement('td');
        filename.classList.add('file-filename');
        filename.innerHTML = this.file.filename;
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

            let filename = prompt("Enter the new filename", this.file.filename);
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
            createOverlay(this.file.id);
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
            files.appendChild(new FileCustom(element).toTableRow());
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
