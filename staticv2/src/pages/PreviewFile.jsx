import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { useFiles, useGetInfo } from "../hooks/files";
import PropTypes from 'prop-types';

const PreviewFile = (props) => {
    const { sharedFileId } = useParams();
    const [ previewFile, setPreviewFile ] = useState(props.fileInfo);
    const { getShareFileInfo } = useFiles();
    const { getFilenameInfo } = useGetInfo();

    useEffect(() => {
        if(!sharedFileId) return;
        getShareFileInfo(sharedFileId).then((file) => setPreviewFile(getFilenameInfo(file, true)));
    }, [ sharedFileId, getFilenameInfo, getShareFileInfo ]);

    const previewComponent = () => {
        if(!previewFile) return <h1>404</h1>;

        const url = !sharedFileId ? previewFile?.url : previewFile?.sharedUrl;

        if(previewFile?.contentType?.includes('video'))
            return <video className={props.className} controls><source src={url} type={previewFile?.contentType}/></video>;
        
        return <iframe src={url} className={props.className}/>
    }

    return previewComponent();
};

PreviewFile.propTypes = {
    file: PropTypes.object,
}

export default PreviewFile;