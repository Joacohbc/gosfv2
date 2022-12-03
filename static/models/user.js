export class User {
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

    get username() {
        return this._username;
    }    
}