import { Message } from "/static/modules/message.js";
import { File } from "/static/models/file.js";
import { createOverlay } from "/static/main/overlay.js";

// Creo el objeto message para mostrar los mensajes
// y le digo que el ID sera "message"
const message = new Message("message");

// Una Clase FileCustom que tiene un objeto File que es convertido
// a elementos HTML para mostrar en la tabla y le agrega funcionalidades
// a los botones
class FileCustom {

    // Constructor que recibe un objeto JSON y 
    // lo convierte a un objeto File
    constructor(fileJson) {
        this.file = File.fromJSON(fileJson);

        const row = document.createElement('tr');
        row.setAttribute("class", "file");
        this.row = row;
    }

    // Función que abre el archivo en una nueva pestaña
    open() {
        window.open(`/api/files/${this.file.id}`, '_blank');
    }

    share() {
        window.open(`/api/files/share/${this.file.id}`, '_blank');
    }

    // Función que descarga el archivo
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

    // Función que elimina el archivo
    delete() {
        axios.delete(`/api/files/${this.file.id}`)
        .then(res => {
            this.row.remove();
    
            if(document.querySelector("tbody").childElementCount == 0) {
                document.querySelector("thead").innerHTML = "No files, start uploading files c:";
            }
            message.showSuccess(res.data.message);
        })
        .catch(err => {
            message.showError(err.response.data.message);
        });
    }

    // Función que actualiza el archivo (el nombre y si es publico o no)
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

    // Función que convierte el objeto File a elementos HTML
    toTableRow() {
        const id = document.createElement('td');
        id.classList.add('file-id');
        id.innerText = this.file.id;
        this.row.appendChild(id);

        const filename = document.createElement('td');
        filename.classList.add('file-filename');
        filename.innerHTML = this.file.filename;
        filename.addEventListener('click', (e) => {
            e.preventDefault();
            this.open();
        });
        this.row.appendChild(filename);
        
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
            createOverlay(this.file.id).then(show => show());
        });
        actions.appendChild(shareBtn);


        this.row.appendChild(actions);

        return this.row;
    }

    // Función que convierte el objeto File a elementos HTML
    toTableRowShare() {
        const id = document.createElement('td');
        id.classList.add('file-id');
        id.innerText = this.file.id;
        this.row.appendChild(id);

        const filename = document.createElement('td');
        filename.classList.add('file-filename');
        filename.innerHTML = this.file.filename;
        filename.addEventListener('click', (e) => {
            e.preventDefault();
            this.share();
        });
        this.row.appendChild(filename);
        
        const actions = document.createElement('td');
        actions.classList.add('file-actions');
        actions.innerHTML = "Shared file";
        this.row.appendChild(actions);

        return this.row;
    }
}

// Función que recarga la tabla de Archivos
export function reloadTable(cbFiltro = null) {
    const htmlFiles = [];
    document.querySelector('tbody').innerHTML = "";
    axios.get("/api/files/")
    .then(req => {
        
        let files = req.data || [];
        if(cbFiltro != null) {
            files = files.filter(cbFiltro);
        }

        files.forEach(element => {
            htmlFiles.push(new FileCustom(element).toTableRow());
        });

        return axios.get("/api/files/share");
    })
    .then(req => {
        let files = req.data || [];
        console.log(files);
        if(cbFiltro != null) {
            files = files.filter(cbFiltro);
        }
        
        files.forEach(file => {
            htmlFiles.push(new FileCustom(file).toTableRowShare());
        });

        return Promise.resolve(htmlFiles);
    })
    .then(files => {

        // Si no hay archivos que ingrese un mensaje personalizado
        // en el head de la tabla que indique que no hay archivos
        if(files.length == 0) {
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
        const part = document.createDocumentFragment();
        files.forEach(element => {
            part.appendChild(element);
        });
        document.querySelector('tbody').appendChild(part);
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
            reloadTable();
            message.showSuccess(req.data.message);
        })
        .catch(err => {
            message.showError(err.response.data.message); 
        });
    });

    document.querySelector('#search-input').addEventListener('keypress', (e) => {
        if(e.key != 'Enter') return;

        let search = document.querySelector('#search-input').value;
        if(search.trim() == "")  {
            reloadTable();
            return;
        }

        reloadTable((file) => {
            return file.filename.toLowerCase().includes(search.toLowerCase());
        });
    });
});