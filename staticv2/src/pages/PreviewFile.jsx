import { useContext, useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import AuthContext from "../context/auth-context";
import PropTypes from 'prop-types';
import { getFileInfo } from "../utils/files";

const PreviewFile = (props) => {
    const { sharedFileId } = useParams();
    const [ previewFile, setPreviewFile ] = useState(props.fileInfo);
    const { cAxios, addTokenParam } = useContext(AuthContext);

    useEffect(() => {
        if(!sharedFileId) return;
        if(!cAxios) return;

        async function fetchDataFile() {
            try {
                const res = await cAxios.get(`/api/files/share/${sharedFileId}/info`);
                setPreviewFile(getFileInfo(res.data));
            } catch(err) {
                console.log(err);
            }
        }
        
        fetchDataFile();
    }, [ sharedFileId, cAxios ]);

    const previewComponent = () => {
        if(!previewFile) return <h1>404</h1>;

        const urlWithToken = !sharedFileId ? addTokenParam(previewFile?.url) : addTokenParam(previewFile?.sharedUrl);
        
        if(previewFile?.filename?.contentType?.includes('video'))
            return <video className={props.className} controls><source src={urlWithToken} type={previewFile?.contentType}/></video>;
        
        return <iframe src={urlWithToken} className={props.className}/>;
    }

    return previewComponent();
};

PreviewFile.propTypes = {
    file: PropTypes.object,
}

export default PreviewFile;