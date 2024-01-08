import { useContext, useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { useGetInfo } from "../hooks/files";
import AuthContext from "../context/auth-context";
import PropTypes from 'prop-types';

const PreviewFile = (props) => {
    const { sharedFileId } = useParams();
    const [ previewFile, setPreviewFile ] = useState(props.fileInfo);
    const { cAxios } = useContext(AuthContext);
    const { getFilenameInfo } = useGetInfo();

    useEffect(() => {
        if(!sharedFileId) return;
        if(!cAxios) return;

        const fetchDataFile = async () => {
            try {
                const res = await cAxios.get(`/api/files/share/${sharedFileId}/info`);
                setPreviewFile(getFilenameInfo(res.data, true));
            } catch(err) {
                console.log(err);
            }
        }
        
        fetchDataFile();
    }, [ sharedFileId, cAxios, getFilenameInfo ]);

    const previewComponent = () => {
        if(!previewFile) return <h1>404</h1>;

        const url = !sharedFileId ? previewFile?.url : previewFile?.sharedUrl;

        console.log(url);
        if(previewFile?.contentType?.includes('video'))
            return <video className={props.className} controls><source src={url} type={previewFile?.contentType}/></video>;
        
        return <iframe src={url} className={props.className}/>;
    }

    return previewComponent();
};

PreviewFile.propTypes = {
    file: PropTypes.object,
}

export default PreviewFile;