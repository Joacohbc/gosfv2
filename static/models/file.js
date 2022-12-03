export class File {
    constructor(id, filename, shared, shared_with) {
        this._id = id;
        this._filename = filename;
        this._shared = shared;
        this._shared_with = shared_with;
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


}