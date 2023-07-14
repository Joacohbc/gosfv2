import '../css/files.css';

export default function Files() {
    return <>
        <div className="search">
            <input type="text" name="Buscar Archivo" placeholder="Enter Search" id="search-input"/>
        </div>
        
        <div id="file-loading" className="loader" hidden> Uploading files </div> 
        <div id="message"></div>

        <div className="table-container">
            <table>
                <thead>
                    
                </thead>
        
                <tbody>

                </tbody>
            </table>
        </div>
        
        <div className="table-buttons">
            <label htmlFor="input-upload" id="btn-upload" className="header-btn">Upload file/s</label>
            <input id="input-upload" type="file" style={ {display: 'none'} } multiple/>
        </div>
    </>;
}
